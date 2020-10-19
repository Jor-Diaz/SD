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
	address     = "dist158:50051"
)


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

/************************************************************************************************************************/
func OrderReader(tipo int32){
	//Abir archivo
  archivo:="retail.csv"
  if tipo == 2{
    archivo="pymes.csv"
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
    ord := Orden{id:record[0],producto:record[1],valor:int32(numero),tienda:record[3],destino:record[4],prioritario:2}
    if tipo==2 {
      prio,_:=strconv.Atoi(record[5])
      ord = Orden{id:record[0],producto:record[1],valor:int32(numero),tienda:record[3],destino:record[4],prioritario:int32(prio)}
    }
    ordenes = append(ordenes, &ord)
  	}

}

/************************************************************************************************************************/
func searchOrder( _id string) *Orden {
  for _, v := range ordenes {
    if v.id == _id {
          return v
        }
      }
  return &Orden404
}

func enviar_ordenes( delta_tiempo float64){
  var conn *grpc.ClientConn
  conn, err := grpc.Dial("dist158:9000", grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %s", err)
  }
  defer conn.Close()
  c := pb.NewGreeterClient(conn)
  i:=0
  update_time:=time.Now()
  time2:=time.Now()
  for  i < len(ordenes){
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
}



func main() {
  // Set up a connection to the server.
    var delta_tiempo float64
    var tipo_cliente int32
    tipo_cliente=0
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


    //Función para crear el array de estructuras Item
    //RetailReader() //working!
    OrderReader(tipo_cliente) //working!
    go enviar_ordenes( delta_tiempo )


    var conn *grpc.ClientConn
    conn, err := grpc.Dial("dist158:9000", grpc.WithInsecure())
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
        if (opcion!=-1){
          response, err := c.ConEstado(context.Background(), &pb.ConsultaEstado{Seguimiento:opcion})
          if err != nil {
            log.Fatalf("Error when calling SayHello: %s", err)
          }
          if response.Estado==0{
            log.Printf("El Estado de la orden es : En bodega")
            log.Printf("Orden numero seguimiento: %d ; Numero de intentos %d;  ",response.Seguimiento,response.Intentos)
          }else if response.Estado==1{
            log.Printf("El Estado de la orden es : En camino")
            log.Printf("Orden numero seguimiento: %d ; Numero de intentos %d; Id Camion %d ",response.Seguimiento,response.Intentos,response.IdCamion)
          }else if response.Estado==2{
            log.Printf("El Estado de la orden es : Recibido")
            log.Printf("Orden numero seguimiento: %d ; Numero de intentos %d; Id Camion %d ",response.Seguimiento,response.Intentos,response.IdCamion)
            fmt.Println("Hora de entrega: ",response.TiempoEntrega.Format(time.ANSIC))            
          }else if response.Estado==3{
            log.Printf("El Estado de la orden es : No Recibido")
            log.Printf("Orden numero seguimiento: %d ; Numero de intentos %d; Id Camion %d ",response.Seguimiento,response.Intentos,response.IdCamion)
          }else{
            log.Printf("El Estado de la orden es : No existe")
          }
        }
    }
}
