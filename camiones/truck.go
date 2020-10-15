// 3 camiones
// 2 retail
// 1 no retail
package main

import (
  //"os"
	"fmt"
	"strconv"
	"time"
	"math/rand"
)
/************paquetes tipos*************/
// 0: prioritario , 1:retail ,2: normal
/**************************************/

/************camiones tipos*************/
// 0: retail                   1: normal//
/**************************************/

type pack struct{
  id_pack string
  pack_type  int
  value int
  origin string
  destination string
  tries int
  delivery_date time.Time
}

type truck struct{
  type_t int
  pack1 *pack
	pack2 *pack
}


var pack404 pack = pack{id_pack: "not_found",pack_type: 3, value: -1, origin:  "not_found", destination:  "not_found", tries :-1, delivery_date:  time.Now()}
var camion404 truck = truck{type_t: -1 , pack1: &pack404, pack2: &pack404 }

/*************************************************************************************************************************************************/
/* Función apra crear un nuevo paquete
	Se espera que los valores sean entregados por otra maquina
	retorna un puntero a paquete
*/

func newPack(idPack string, typ int, val string, org string, dst string, trs int, date time.Time) *pack {
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
		pack1 : packA,
		pack2 : packB	}
	return &nTruck
}
/*************************************************************************************************************************************************/

/*************************************************************************************************************************************************/
/*	Función que retorna 1 siempre que sale un numero menor a 80, considerando numeros entre 0 y 100
*/
func chanceToDeliver() int{
	rand.Seed(time.Now().UnixNano())
	chance := rand.Intn(100)
	if chance <= 80{
		return 1
	} else{
		return 0
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
si se retorna 1 implica que se entrega .pack2
si se retorna 0 implica que se entrega .pack1

*/
func wichToDeliver(pack0 *pack, pack1 *pack) int {
	if pack0.value != -1 || pack1.value != -1 {
		if pack0.pack_type == pack1.pack_type {
			if pack0.value > pack1.value{
				return 0
			}else{
				return 1
			}
		}else if pack0.pack_type < pack1.pack_type{
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
	packToDeliver := wichToDeliver(deliver_truck.pack1, deliver_truck.pack2)
	if packToDeliver == 0{
		for deliver_truck.pack1.tries < 3{
			if chanceToDeliver() == 1{
				fmt.Println(".:| Entregando |:. -> ", deliver_truck.pack1.id_pack)
				deliver_truck.pack1 = deliver(deliver_truck.pack1)
				return deliver_truck

			} else{
				fmt.Println("Nuevo Intento de Entrega de ", deliver_truck.pack2.id_pack)

				deliver_truck.pack1.tries++
			}
		}
		return deliver_truck
	}else if packToDeliver == 1{
		for deliver_truck.pack2.tries < 3{
			if chanceToDeliver() == 1{
				fmt.Println(".:| Entregando |:. -> ", deliver_truck.pack2.id_pack)
				deliver_truck.pack2 = deliver(deliver_truck.pack2)
				return deliver_truck
			}else{
				fmt.Println("Nuevo Intento de Entrega de ", deliver_truck.pack2.id_pack)
				deliver_truck.pack2.tries++
			}
		}
		return deliver_truck
	}else{
		fmt.Println(".: Camión vacío :.")
		return deliver_truck
	}
}





func main()  {

  p1 := newPack("SA5897AS", 1,"95","tienda-A","casa-D",  0,  time.Now())
	p2 := newPack("SA6947GH", 0,"50","tienda-C","casa-C",  0,  time.Now())
	t1 := newTruck(1,p1,p2)
	fmt.Println("----------------------------------")

	fmt.Println("Antes de Entregar:")
	fmt.Println("_Camión_")
	fmt.Println("Paquete 1: *", t1.pack1.id_pack," *")
	fmt.Println("Paquete 2: *", t1.pack2.id_pack," *")
	fmt.Println("----------------------------------")



	// a := wichToDeliver(t1.pack1, t1.pack2)
	// b:= chanceToDeliver()
	// fmt.Println(a)
	// fmt.Println("chance:", b)
	// p1 = deliver(p1)
	//fmt.Println(p1.id_pack)
	t1 = delivery(t1)
	fmt.Println("----------------------------------")

	fmt.Println("Despues de 1er Entrega:")
	fmt.Println("_Camión_")
	fmt.Println("Paquete 1: *", t1.pack1.id_pack," *")
	fmt.Println("Paquete 2: *", t1.pack2.id_pack," *")
	fmt.Println("----------------------------------")

	t1 = delivery(t1)
	fmt.Println("----------------------------------")

	fmt.Println("Despues de 2da Entrega:")
	fmt.Println("_Camión_")
	fmt.Println("Paquete 1: *", t1.pack1.id_pack," *")
	fmt.Println("Paquete 2: *", t1.pack2.id_pack," *")
	fmt.Println("----------------------------------")



}
