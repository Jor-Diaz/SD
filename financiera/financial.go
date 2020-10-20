package main

import (
        "log"
        "encoding/json"
        "fmt"
        "github.com/streadway/amqp"
)


type pack struct{
  Pack_Type int32
  Value int32
  Tries int32
  Income float64
}
var deliveredPacks []*pack
var notDeliveredPacks []*pack
var total float64 = 0
var packs []*pack

/////////////////////////////////////////////
func failOnError(err error, msg string) {
        if err != nil {
                log.Fatalf("%s: %s", msg, err)
        }
}

/**------------------------------------------------------------------------------------------------*/
/*Pequeña función para castear un entero (int32) a flotante (float64)*/
func float(in int32) float64 {
    return float64(in)
}
/**------------------------------------------------------------------------------------------------*/



/**------------------------------------------------------------------------------------------------*/
/*Esta función ejecuta las condiciones de ganancia y perdida por paquete, trabajando con un arreglo
de punteros a paquete*/

func financialBalance(packs []*pack)  {
  for _, pckt := range packs {
    if pckt.Pack_Type == 0{//retail
      if pckt.Tries !=3{// entregado
        pckt.Income = float(pckt.Value-(pckt.Tries)*10)
        total += pckt.Income
        deliveredPacks = append(deliveredPacks, pckt)
      } else{
        pckt.Income = float(pckt.Value-(pckt.Tries)*10)
        total += pckt.Income
        notDeliveredPacks = append(deliveredPacks, pckt)
      }

    } else if pckt.Pack_Type == 1{//normal prioritario
      if pckt.Tries !=3{// entregado
        pckt.Income = float(pckt.Value-(pckt.Tries)*10)
        total += pckt.Income
        deliveredPacks = append(deliveredPacks, pckt)
      } else{
        pckt.Income = float((3*pckt.Value/100)-(pckt.Tries*10))
        total += pckt.Income
        notDeliveredPacks = append(deliveredPacks, pckt)
      }

    } else { // paquete normal
      if pckt.Tries !=3{// entregado
        pckt.Income = float(pckt.Value-(pckt.Tries)*10)
        total += pckt.Income
        deliveredPacks = append(deliveredPacks, pckt)
      } else{
        pckt.Income = float(-(pckt.Tries)*10)
        total += pckt.Income
        notDeliveredPacks = append(deliveredPacks, pckt)
      }

    }
  }
}

///// dist157

/*MAIN*/
func main() {
  /**Las siguiente LOC permiten hacer la connexió via RAbbitMQ
  entre la  maquina dist157 y dist158 (finanzas-logistico) */

        conn, err := amqp.Dial("amqp://test:holachao@localhost:5672/")
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()

        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()

        q, err := ch.QueueDeclare(
        	"hello-queue",  // name
        	false,         // durable
        	false,         // delete when usused
        	false,         // exclusive
        	false,         // no-wait
        	nil,           // arguments
        )
        	failOnError(err, "Failed to declare a queue")

        	msgs, err := ch.Consume(
        		q.Name,        // queue
        		"",            // consumer
        		true,          // auto-ack
        		false,         // exclusive
        		false,         // no-local
        		false,         // no-wait
        		nil,           // args
        	)
        	failOnError(err, "Failed to register a consumer")

        	forever := make(chan bool)
          var pck pack
        	go func() {
            for d := range msgs { // Unmarshall del json recibido desde logistico
              if err := json.Unmarshal(d.Body, &pck); err != nil {
                  panic(err)
              }
              packs = append(packs,&pck) //agregamos al arreglo
              financialBalance(packs) //balance
              /*Prints*/
              fmt.Println("[****TOTAL:",total, "Dignipesos **********]")

              fmt.Println("[*************************************]")
              if (pck.Pack_Type == 0){
                fmt.Println("[Tipo de Paquete]: Retail")
              } else if (pck.Pack_Type == 1){
                fmt.Println("[Tipo de Paquete]: Prioritario")
              }else{
                fmt.Println("[Tipo de Paquete]: Normal")
              }
              fmt.Println("[Valor del Paquete]: ", pck.Value)
              fmt.Println("[Entregas Fallidas]: ", pck.Tries)
              if (pck.Income >= 0){
                fmt.Println("[Ingresos generados]:", pck.Income, "Dignipesos")
              }else {
                fmt.Println("[Perdidas generadas]:", pck.Income, "Dignipesos")
              }
              fmt.Println("[*************************************]")
            }
          }()
      	   <-forever //Sigue corriendo hasta un Ctrl+C
}
