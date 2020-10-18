package main

import (
  "os"
  "log"
	"encoding/csv"
	"fmt"
	"io"
  "strconv"
  "time"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  pb"Lab1/SD/pipeline"
)

const (
	address     = "dist159:50051"
	defaultName = "world"
)


// Creación de Item struc
//type Item struct{
//  id string
//  producto string
//  valor int32
//  tienda string
//  destino string
//}

type Orden struct{
  id string
  producto string
  valor int32
  tienda string
  destino string
  prioritario int32
}

///retail prioridad 2
///pymes prioritario 1
///pymes normal 0


var ordenes []*Orden // Arreglo de pymes (pymes.csv)
var Orden404 Orden = Orden{id: "not_found", producto: "not_found", valor:1, tienda:"not_found",destino:"not_found", prioritario: -1}

/****************
func RetailReader(){
	//Abir archivo
  recordFile, err := os.Open("retail.csv")
	if err != nil {
		fmt.Println("An error encountered ::", err)
		os.Exit(0)
	}
	reader := csv.NewReader(recordFile)
	reader.Read() //saltar primera linea
	for i:= 0 ;; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // final del archivo
		} else if err != nil {
			fmt.Println("Error ::", err)
			break
		}
    numero,_:=strconv.Atoi(record[2])
    prod := Item{id: record[0], producto: record[1], valor:int32(numero), tienda:record[3],destino:record[4]}
    productos = append(productos, &prod)
  	}
}********/
/************************************************************************************************************************/

/************************************************************************************************************************/
func OrderReader(tipo int32){
	//Abir archivo
  if(tipo==1){
    archivo:="retail.csv"
  }else{
    archivo:="pymes.csv"
  }
  recordFile, err := os.Open(archivo)
	if err != nil {
		fmt.Println("An error encountered ::", err)
		os.Exit(0)
	}
	reader := csv.NewReader(recordFile)
	reader.Read() //saltar primera linea
	for i:= 0 ;; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // final del archivo
		} else if err != nil {
			fmt.Println("Error ::", err)
			break
		}
    numero,_:=strconv.Atoi(record[2])
    if(tipo==1){
        ord := Orden{id:record[0],producto:record[1],valor:record[2],tienda:record[3],destino:record[4],prioritario:2}
    }else{
      ord := Orden{id:record[0],producto:record[1],valor:record[2],tienda:record[3],destino:record[4],prioritario:record[5]}
    }
    ordenes = append(ordenes, &ord)
  	}

}
/************************************************************************************************************************/

/****************************************************************************************************************
func searchItem( _id string) *Item {
  for _, v := range productos {
    if v.id == _id {
          return v
        }
      }
      //Item404 := Item{id: "not_found", producto: "not_found", valor:"not_found", tienda:"not_found",destino:"not_found"}
      return &Item404
    }
    ********/
/************************************************************************************************************************/
func searchOrder( _id string) *Orden {
  for _, v := range ordenes {
    if v.id == _id {
          return v
        }
      }
      //Item404 := Item{id: "not_found", producto: "not_found", valor:"not_found", tienda:"not_found",destino:"not_found"}
      return &Orden404
    }

func main() {
  // Set up a connection to the server.
    var delta_tiempo float64
    var tipo_cliente=0
    fmt.Println("Gracias por iniciar el cliente de ordenes de SD X-Wing Team")
    fmt.Println("----------Configuracion------")
    fmt.Println("Ingrese el tiempo(en segundos) entre envio de ordenes")
    fmt.Scanf("%f", &delta_tiempo)
    for (tipo_cliente-1)*(tipo_cliente-2)!=0{
        fmt.Println("Ingrese 1 si es cliente retail o 2 en caso de ser cliente pymes")
        fmt.Scanf("%d", &tipo_cliente)
    }
    if (tipo_cliente==1){
      fmt.Println("A ingresado 1. Se cargara el archivo retail.csv")
    }else{
      fmt.Println("A ingresado 2. Se cargara el archivo pymes.csv")
    }
    var conn *grpc.ClientConn
  	conn, err := grpc.Dial("dist159:9000", grpc.WithInsecure())
  	if err != nil {
  		log.Fatalf("did not connect: %s", err)
  	}
  	defer conn.Close()

    //Función para crear el array de estructuras Item
    //RetailReader() //working!
    OrderReader() //working!

  	c := pb.NewGreeterClient(conn)
    i:=0
    update_time:=time.Now()
    time2:=time.Now()
    for  i < len(productos){
      time2=time.Now()
      if ( time2.Sub(update_time).Seconds() > delta_tiempo){
    	  response, err := c.SayHello(context.Background(), &pb.Message{Id:ordenes[i].id,Producto:ordenes[i].producto,Valor:ordenes[i].valor,Tienda:ordenes[i].tienda,Destino: ordenes[i].destino,Prioridad:ordenes[i].prioritario})
      	if err != nil {
      		log.Fatalf("Error when calling SayHello: %s", err)
      	}
      	log.Printf("El codigo de seguimiento del pedido es: %d", response.Seguimiento)
        i=i+1
        update_time=time.Now()
      }
    }

  id_ex := "CsC147"
  //item_ex := searchItem(id_ex)
  order_ex := searchOrder(id_ex)

  fmt.Println(id_ex,"->", order_ex.valor)
}
