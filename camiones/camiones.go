// 3 camiones
// 2 retail
// 1 no retail
package main

import (
  //"os"
	"fmt"
	//"strconv"
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
  value int32
  origin string
  destination string
  tries int
  delivery_date time.Time
	seguimiento int32
}

type truck struct{
  type_t int32
  pack0 *pack
	pack1 *pack
}


var pack404 pack = pack{id_pack: "empty",pack_type: -1, value: -1, origin:  "empty", destination:  "empty", tries :-1, delivery_date:  time.Now(),seguimiento:-1}
var camion404 truck = truck{type_t: -1 , pack0: &pack404, pack1: &pack404 }

/*************************************************************************************************************************************************/
/* Función apra crear un nuevo paquete
	Se espera que los valores sean entregados por otra maquina
	retorna un puntero a paquete
*/

func newPack(idPack string, typ int32, val int32, org string, dst string, trs int, date time.Time, seguimiento int32) *pack {
		Npack := pack{
			id_pack : idPack,
		  pack_type: typ,
		  value : val,
		  origin: org,
		  destination: dst,
		  tries: trs,
		  delivery_date: date,
			seguimiento: seguimiento}
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
	if pack0.value != -1 && pack1.value != -1 {
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
	}else if pack0.value == -1 && pack1.value != -1{
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
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist158:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	packToDeliver := wichToDeliver(deliver_truck.pack0, deliver_truck.pack1)
	if packToDeliver == 0{
		for deliver_truck.pack0.tries < 3{
			if chanceToDeliver() == 1{
				fmt.Println(".:| Entregando |:. -> ", deliver_truck.pack0.id_pack)
				response, err := c.ActEntrega(context.Background(), &pb.ActCamion{Seguimiento:deliver_truck.pack0.seguimiento,Exito:1})
				if err != nil {
					log.Fatalf("Error when calling SayHello: %s", err)
				}
				log.Printf("Orden actualizada en logistica, orden numero %d",response.Seguimiento)
				deliver_truck.pack0 = deliver(deliver_truck.pack0)
				return deliver_truck

			} else{
				fmt.Println("Nuevo Intento de Entrega de ", deliver_truck.pack1.id_pack)
				response, err := c.ActEntrega(context.Background(), &pb.ActCamion{Seguimiento:deliver_truck.pack0.seguimiento,Exito:0})
				if err != nil {
					log.Fatalf("Error when calling SayHello: %s", err)
				}
				log.Printf("Orden actualizada en logistica, orden numero %d",response.Seguimiento)
				deliver_truck.pack0.tries++
			}
		}
		response, err := c.ActEntrega(context.Background(), &pb.ActCamion{Seguimiento:deliver_truck.pack0.seguimiento,Exito:-1})
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}
		log.Printf("Orden actualizada en logistica, orden numero %d",response.Seguimiento)
		deliver_truck.pack0=&pack404
		return deliver_truck
	}else if packToDeliver == 1{
		for deliver_truck.pack1.tries < 3{
			if chanceToDeliver() == 1{
				fmt.Println(".:| Entregando |:. -> ", deliver_truck.pack1.id_pack)
				response, err := c.ActEntrega(context.Background(), &pb.ActCamion{Seguimiento:deliver_truck.pack1.seguimiento,Exito:1})
				if err != nil {
					log.Fatalf("Error when calling SayHello: %s", err)
				}
				log.Printf("Orden actualizada en logistica, orden numero %d",response.Seguimiento)
				deliver_truck.pack1 = deliver(deliver_truck.pack1)
				return deliver_truck
			}else{
				fmt.Println("Nuevo Intento de Entrega de ", deliver_truck.pack1.id_pack)
				response, err := c.ActEntrega(context.Background(), &pb.ActCamion{Seguimiento:deliver_truck.pack1.seguimiento,Exito:0})
				if err != nil {
					log.Fatalf("Error when calling SayHello: %s", err)
				}
				log.Printf("Orden actualizada en logistica, orden numero %d",response.Seguimiento)
				deliver_truck.pack1.tries++
			}
		}
		response, err := c.ActEntrega(context.Background(), &pb.ActCamion{Seguimiento:deliver_truck.pack1.seguimiento,Exito:-1})
		if err != nil {
			log.Fatalf("Error when calling SayHello: %s", err)
		}
		log.Printf("Orden actualizada en logistica, orden numero %d",response.Seguimiento)
		deliver_truck.pack1=&pack404
		return deliver_truck
	}else{
		fmt.Println(".: Camión vacío :.")
		return deliver_truck
	}
}

var response pb.RespuestaCon
var tiempo_espera float64

func ejecucion_camion(Id_camion int32, tiempo_espera float64){
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist158:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	opcion:=0
	no_paquetes:=0
	camion1 := newTruck(Id_camion,&pack404,&pack404)
	update_time:=time.Now()
  time2:=time.Now()
	status:=0
	for  opcion!=-1{
			no_paquetes=0
			time2=time.Now()
			if ( time2.Sub(update_time).Seconds() > tiempo_espera){
			//fmt.Println("Ingrese el numero de seguimiento para consultar estado o -1 para salir : ")
			//fmt.Scanf("%d", &opcion)
				fmt.Println("Camion %d solicitando paquetes",Id_camion)
				response, err := c.Solpedido(context.Background(), &pb.Solcamion{IdCamion:Id_camion})
				if err != nil {
					log.Fatalf("Error when calling SayHello: %s", err)
				}
				///Verificar que es paquete valido
				no_paquetes=1
				if (no_paquetes==1){
						update_time=time.Now()
						log.Printf("Orden asignada con codigo seguimiento %d al camion %d",response.Seguimiento,Id_camion)
						paquete_1 := newPack(response.Id, 2, response.Valor, response.Tienda,response.Destino, 0,  time.Now(),response.Seguimiento)
						camion1.pack0
						fmt.Println("Esperando segundo paquete para entregar")
						time2:=time.Now()
						for (status!=0) {
							if ( time2.Sub(update_time).Seconds() > tiempo_espera){
									status=1
							}
						}
						fmt.Println("Camion %d solicitando segundo paquete",Id_camion)
						response, err := c.Solpedido(context.Background(), &pb.Solcamion{IdCamion:Id_camion})
						if err != nil {
							log.Fatalf("Error when calling SayHello: %s", err)
						}
						///Verificar que es paquete valido
						no_paquetes=2
						if (no_paquetes==2){
							log.Printf("Orden asignada con codigo seguimiento %d al camion %d",response.Seguimiento,Id_camion)
							paquete_1 := newPack(response.Id, 2, response.Valor, response.Tienda,response.Destino, 0,  time.Now(),response.Seguimiento)
							camion1.pack1=paquete_1
						}else{
							fmt.Println("No hay paquetes disponibles para repartir para el camion %d", Id_camion)
						}
				}else{
					fmt.Println("No hay paquetes disponibles para repartir para el camion %d", Id_camion)
				}
				if (no_paquetes>0){
					state := truckState(camion1)
					fmt.Println("-----------------------------------")
					fmt.Println("Antes de Entregar:")
					fmt.Println("_Camión_")
					fmt.Println("Capacidad del Camión:", state, "espacios")
					fmt.Println("Paquete 1: *", camion1.pack0.id_pack," *")
					fmt.Println("Paquete 2: *", camion1.pack1.id_pack," *")
					fmt.Println("------------------------------------")
					camion1 = delivery(camion1)
					fmt.Println("-----------------------------------")
					state = truckState(camion1)
					fmt.Println("Despues de 1er Entrega:")
					fmt.Println("_Camión_")
					fmt.Println("Capacidad del Camión:", state, "espacios")
					fmt.Println("Paquete 1: *", camion1.pack0.id_pack," *")
					fmt.Println("Paquete 2: *", camion1.pack1.id_pack," *")
					fmt.Println("-----------------------------------")
					camion1 = delivery(camion1)
					fmt.Println("-----------------------------------")
					state = truckState(camion1)
					fmt.Println("Despues de 2da Entrega:")
					fmt.Println("_Camión_")
					fmt.Println("Capacidad del Camión:", state, "espacios")
					fmt.Println("Paquete 1: *", camion1.pack0.id_pack," *")
					fmt.Println("Paquete 2: *", camion1.pack1.id_pack," *")
					fmt.Println("-----------------------------------")
				}
			}
	}
}

func main()  {
	fmt.Println("Ingrese el tiempo (en segundos) a esperar por parte de los camiones")
	fmt.Scanf("%f", &tiempo_espera)
	go ejecucion_camion(1,tiempo_espera)
	var jij int32
	fmt.Scanf("%d", &jij)
}
	//fmt.Println(aux)
  //p1 := newPack(aux, 2, response.Valor, response.Tienda,response.Destino, 0,  time.Now())
	//p2 := newPack("SA6947GH", 0, "50", "_","_",  0,  time.Now())
	// p3 := newPack("SA2589TR", 2, "5",  "_", "_", 0,  time.Now())
	// p4 := newPack("SA1597EF", 0, "20", "_", "_", 0,  time.Now())
	// p5 := newPack("SA6947GH", 1, "50", "_", "_", 0,  time.Now())
	// p6 := newPack("SA2596NH", 2, "90", "_", "_", 0,  time.Now())

	//t1 := newTruck(1,p1,p1)
