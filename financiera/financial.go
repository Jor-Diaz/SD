package main

import(
  //"os"
	//"encoding/csv"
	"fmt"
	// /"io"
  )


  // Los env´ıos completados.
  // La cantidad de veces que se intent´o entregar un paquete.
  // Los paquetes que no pudieron ser entregados.
  // P´erdidas o ganancias de cada paquete (en Dignipesos).

  type pack struct{
    packType int
    value int
    tries int
    income float64
  }

var deliveredPacks []*pack
var notDeliveredPacks []*pack
var total float64

func float(in int) float64 {
    return float64(in)
}

func financialBalance(packs []*pack)  {
  for _, pckt := range packs {
    if pckt.packType == 0{//retail
      if pckt.tries !=3{// entregado
        pckt.income = float(pckt.value-(pckt.tries)*10)
        total += pckt.income
        deliveredPacks = append(deliveredPacks, pckt)
      } else{
        pckt.income = float(pckt.value-(pckt.tries)*10)
        total += pckt.income
        notDeliveredPacks = append(deliveredPacks, pckt)

      }
    } else if pckt.packType == 1{//normal prioritario
      if pckt.tries !=3{// entregado
        pckt.income = float(pckt.value-(pckt.tries)*10)
        total += pckt.income
        deliveredPacks = append(deliveredPacks, pckt)
      } else{
        pckt.income = float((3*pckt.value/100)-(pckt.tries*10))
        total += pckt.income
        notDeliveredPacks = append(deliveredPacks, pckt)
      }
    } else { // paquete normal
      if pckt.tries !=3{// entregado
        pckt.income = float(pckt.value-(pckt.tries)*10)
        total += pckt.income
        deliveredPacks = append(deliveredPacks, pckt)
      } else{
        pckt.income = float(-(pckt.tries)*10)
        total += pckt.income
        notDeliveredPacks = append(deliveredPacks, pckt)
      }
    }
  }
}

func main()  {
  total = 0
  packs := []*pack{}
  p1 := pack{packType: 1, value: 100, tries: 1, income: -999 }
  p2 := pack{packType: 0, value: 122, tries: 2, income: -999 }
  p3 := pack{packType: 2, value: 30,  tries: 3, income: -999 }
  p4 := pack{packType: 0, value: 15,  tries: 1, income: -999 }
  p5 := pack{packType: 1, value: 5,   tries: 3, income: -999 }
  p6 := pack{packType: 0, value: 15,  tries: 3, income: -999 }
  packs = append(packs,&p1)
  packs = append(packs,&p2)
  packs = append(packs,&p3)
  packs = append(packs,&p4)
  packs = append(packs,&p5)
  packs = append(packs,&p6)

  financialBalance(packs)

  fmt.Println(total)
  

}
