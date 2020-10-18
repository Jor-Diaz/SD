// 3 camiones
// 2 retail
// 1 no retail
package main

import (
  //"os"
	"fmt"
	"strconv"
	"time"
	"log"
	"math/rand"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb"Lab1/SD/pipeline"
)
/************paquetes tipos*************/
// 2: retail , 1: prioritario ,0: normal
/**************************************/

/************camiones tipos*************/
// 0: retail                   1: normal//
/**************************************/

type pack struct{
  id_pack string
  pack_type  int32
  value int
  origin string
  destination string
  tries int
  delivery_date time.Time
}

type truck struct{
  type_t int
  pack0 *pack
	pack1 *pack
}


var pack404 pack = pack{id_pack: "empty",pack_type: -1, value: -1, origin:  "empty", destination:  "empty", tries :-1, delivery_date:  time.Now()}
var camion404 truck = truck{type_t: -1 , pack0: &pack404, pack1: &pack404 }

/*************************************************************************************************************************************************/
/* Función apra crear un nuevo paquete
	Se espera que los valores sean entregados por otra maquina
	retorna un puntero a paquete
*/

func newPack(idPack string, typ int32, val string, org string, dst string, trs int, date time.Time) *pack {
	pVal, err := strconv.Atoi(val)
		if err == nil {
			//fmt.Println(pVal)
		}

		Npack := pack{
			id_pack : idPack,
		  pack_type: typ,
		  value : pVal,
		  origin: org,
		  destination: dst,
		  tries: trs,
		  delivery_date: date		}
		return &Npack
}
/*************************************************************************************************************************************************/

/*************************************************************************************************************************************************/

/* Función apra crear un nuevo camión
	Se espera que los valores sean entregados por otra maquina
	retorna un puntero a camión
*/
func newTruck(typ int, packA  *pack, packB *pack ) *truck  {
	nTruck := truck{
		type_t:  typ,
		pack0 : packA,
		pack1 : packB	}
	return &nTruck
}
/*************************************************************************************************************************************************/

/*************************************************************************************************************************************************/
/*	Función que retorna 1 siempre que sale un numero menor a 80, considerando numeros entre 0 y 100
*/
func chanceToDeliver() int{
	rand.Seed(time.Now().UTC().UnixNano())
	chance := rand.Intn(100)
	fmt.Println("-----",chance,"----")
	if chance < 81{
		return 1
	} else{
		return 0
	}
}
/*************************************************************************************************************************************************/

/*************************************************************************************************************************************************/
/**/

func truckState(trucky 	*truck) int{
	if trucky.pack0.id_pack == "empty" && trucky.pack1.id_pack == "empty"{
		return 0
	} else if (trucky.pack0.id_pack != "empty" && trucky.pack1.id_pack == "empty") || (trucky.pack0.id_pack == "empty" && trucky.pack1.id_pack != "empty"){
		return 1
	} else{
		return 2
	}

}
/*************************************************************************************************************************************************/
/*************************************************************************************************************************************************/
/* Función que indica que paquete entregar en esta iteración
Primero considera su prioridad
SI son iguales prioriza por su valor
de no existir un paquete este ya habra sido reemplazado por el paquete404
de estar vacío retorna -1
en la estructura de camión:
si se retorna 1 implica que se entrega .pack1
si se retorna 0 implica que se entrega .pack0

*/
func wichToDeliver(pack0 *pack, pack1 *pack) int {
	if pack0.value != -1 || pack1.value != -1 {
		if pack0.pack_type == pack1.pack_type {
			if pack0.value > pack1.value{
				return 0
			}else{
				return 1
			}
		}else if pack0.pack_type > pack1.pack_type{
			return 0
		}else {
			return 1
		}
	}else if pack0.value != -1 && pack1.value == -1{
		return 0
	}else if pack0.value != -1 && pack1.value == -1{
		return 1
	}else{
		return -1
	}
}
/*************************************************************************************************************************************************/

/*************************************************************************************************************************************************/
/*
Función que "borra" el paquete del camión*/
func deliver (pck *pack) *pack {
	pck = &pack404
	return pck

}
/*************************************************************************************************************************************************/

/*************************************************************************************************************************************************/
/* Función que hace la entrega de los paquetes, esta debe ser usada en un loop o directamente 2 veces

*/

func delivery(deliver_truck *truck) *truck {
	packToDeliver := wichToDeliver(deliver_truck.pack0, deliver_truck.pack1)
	if packToDeliver == 0{
		for deliver_truck.pack0.tries < 3{
			if chanceToDeliver() == 1{
				fmt.Println(".:| Entregando |:. -> ", deliver_truck.pack0.id_pack)
				deliver_truck.pack0 = deliver(deliver_truck.pack0)
				deliver_truck.pack0.tries++
				return deliver_truck

			} else{
				fmt.Println("Nuevo Intento de Entrega de ", deliver_truck.pack1.id_pack)

				deliver_truck.pack0.tries++
			}
		}
		return deliver_truck
	}else if packToDeliver == 1{
		for deliver_truck.pack1.tries < 3{
			if chanceToDeliver() == 1{
				fmt.Println(".:| Entregando |:. -> ", deliver_truck.pack1.id_pack)
				deliver_truck.pack1 = deliver(deliver_truck.pack1)
				deliver_truck.pack0.tries++
				return deliver_truck
			}else{
				fmt.Println("Nuevo Intento de Entrega de ", deliver_truck.pack1.id_pack)
				deliver_truck.pack1.tries++
			}
		}
		return deliver_truck
	}else{
		fmt.Println(".: Camión vacío :.")
		return deliver_truck
	}
}


var tiempo_espera float64
func main()  {
	tiempo_espera:=0
	for tiempo_espera<=0 {
			fmt.Println("Ingrese el tiempo (en segundos) a esperar por parte de los camiones")
			fmt.Scanf("%f", &tiempo_espera)
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist159:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	var opcion int32
	opcion=0
	for  opcion!=-1{
			fmt.Println("Ingrese el numero de seguimiento para consultar estado o -1 para salir : ")
			fmt.Scanf("%d", &opcion)
			response, err := c.ConEstado(context.Background(), &pb.ConsultaEstado{Seguimiento:opcion})
			if err != nil {
				log.Fatalf("Error when calling SayHello: %s", err)
			}
			log.Printf("El Estado de la orden es : %s", response.Estado)
	}

  p1 := newPack("SA5897AS", 2, "95", "_","_",  0,  time.Now())
	p2 := newPack("SA6947GH", 0, "50", "_","_",  0,  time.Now())
	// p3 := newPack("SA2589TR", 2, "5",  "_", "_", 0,  time.Now())
	// p4 := newPack("SA1597EF", 0, "20", "_", "_", 0,  time.Now())
	// p5 := newPack("SA6947GH", 1, "50", "_", "_", 0,  time.Now())
	// p6 := newPack("SA2596NH", 2, "90", "_", "_", 0,  time.Now())

	t1 := newTruck(1,p1,p2)
	state := truckState(t1)
	// t2 := newTruck(0,p3,p4)
	// t3 := newTruck(0,p5,p6)
	fmt.Println("-----------------------------------")

	fmt.Println("Antes de Entregar:")
	fmt.Println("_Camión_")
	fmt.Println("Capacidad del Camión:", state, "espacios")

	fmt.Println("Paquete 1: *", t1.pack0.id_pack," *")
	fmt.Println("Paquete 2: *", t1.pack1.id_pack," *")
	fmt.Println("------------------------------------")

	// a := wichToDeliver(t1.pack0, t1.pack1)
	// b:= chanceToDeliver()
	// fmt.Println(a)
	// fmt.Println("chance:", b)
	// p1 = deliver(p1)
	//fmt.Println(p1.id_pack)
	t1 = delivery(t1)
	fmt.Println("-----------------------------------")
	state = truckState(t1)
	fmt.Println("Despues de 1er Entrega:")
	fmt.Println("_Camión_")
	fmt.Println("Capacidad del Camión:", state, "espacios")
	fmt.Println("Paquete 1: *", t1.pack0.id_pack," *")
	fmt.Println("Paquete 2: *", t1.pack1.id_pack," *")
	fmt.Println("-----------------------------------")
	t1 = delivery(t1)
	fmt.Println("-----------------------------------")
	state = truckState(t1)
	fmt.Println("Despues de 2da Entrega:")
	fmt.Println("_Camión_")
	fmt.Println("Capacidad del Camión:", state, "espacios")
	fmt.Println("Paquete 1: *", t1.pack0.id_pack," *")
	fmt.Println("Paquete 2: *", t1.pack1.id_pack," *")
	fmt.Println("-----------------------------------")
}
