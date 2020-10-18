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
type Item struct{
  id string
  producto string
  valor int32
  tienda string
  destino string
}

type Orden struct{
  id string
  producto string
  valor string
  tienda string
  destino string
  prioritario string
}

var productos []*Item // Arreglo de Items (retail.csv)
var ordenes []*Orden // Arreglo de pymes (pymes.csv)

var Item404 Item = Item{id: "not_found", producto: "not_found", valor:1, tienda:"not_found",destino:"not_found"}
var Orden404 Orden = Orden{id: "not_found", producto: "not_found", valor:"bla", tienda:"not_found",destino:"not_found", prioritario: "not_found"}

/************************************************************************************************************************/
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
}
/************************************************************************************************************************/

/************************************************************************************************************************/
func OrderReader(){
	//Abir archivo
  recordFile, err := os.Open("pymes.csv")
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
    ord := Orden{id:record[0],producto:record[1],valor:record[2],tienda:record[3],destino:record[4],prioritario:record[5]}
    ordenes = append(ordenes, &ord)
  	}

}
/************************************************************************************************************************/

/************************************************************************************************************************/
func searchItem( _id string) *Item {
  for _, v := range productos {
    if v.id == _id {
          return v
        }
      }
      //Item404 := Item{id: "not_found", producto: "not_found", valor:"not_found", tienda:"not_found",destino:"not_found"}
      return &Item404
    }
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
    delta_tiempo:=60.0
    var conn *grpc.ClientConn
  	conn, err := grpc.Dial("dist159:9000", grpc.WithInsecure())
  	if err != nil {
  		log.Fatalf("did not connect: %s", err)
  	}
  	defer conn.Close()

    //Función para crear el array de estructuras Item
    RetailReader() //working!
    OrderReader() //working!

  	c := pb.NewGreeterClient(conn)
    i:=0
    update_time:=time.Now()
    time2:=time.Now()
    for  i < len(productos){
      time2=time.Now()
      if ( time2.Sub(update_time).Seconds() > delta_tiempo){
    	  response, err := c.SayHello(context.Background(), &pb.Message{Tipo:"1",Id:productos[i].id,Producto:productos[i].producto,Valor:productos[i].valor,Tienda:productos[i].tienda,Destino: productos[i].destino})
      	if err != nil {
      		log.Fatalf("Error when calling SayHello: %s", err)
      	}
      	log.Printf("Response from server: %s", response.Producto)
        i=i+1
        update_time=time.Now()
      }
    }

  id_ex := "CsC147"
  //item_ex := searchItem(id_ex)
  order_ex := searchOrder(id_ex)

  fmt.Println(id_ex,"->", order_ex.valor)
}
