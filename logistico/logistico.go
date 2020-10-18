package main

import(
  "fmt"
  //"os"
  "time"
  "log"
  "net"
  //"encoding/csv"
  "google.golang.org/grpc"
  "context"
  pb"Lab1/SD/pipeline"

  )

  const (
  	port = ":50051"
  )

  type Server struct {
      pb.UnimplementedGreeterServer
  }
  type Server1 struct {
      pb.UnimplementedEstadoServer
  }

  func (s *Server) SayHello(ctx context.Context, in *pb.Message) (*pb.Message, error) {
  	log.Printf("Orden recibida con datos:   %s %s %d %s %s %d", in.Id,in.Producto,in.Valor,in.Tienda,in.Destino, in.Prioridad )
    aux:=NewOrden(ordenes,in.Id,in.Producto,in.Valor,in.Tienda,in.Destino,in.Prioridad)
    ordenes=append(ordenes,aux)
  	return &pb.Message{Seguimiento: aux.seguimiento,}, nil
  }
  //func (s *Server1) ConEstado(ctx context.Context, in *pb.Consulta_estado) (*pb.Respuesta_consulta, error) {
  	//log.Printf("Cosulta recibida con datos:   %d", in.Seguimiento)
    //orden_aux:=searchOrder(in.Seguimiento)
  	//return &pb.Respuesta_consulta{Id: orden_aux.id_paquete,Producto:orden_aux.nombre,Valor:orden_aux.valor,Tienda:orden_aux.origen,Destino:orden_aux.destino,Prioridad:orden_aux.prioridad,Intentos:orden_aux.intentos,Estado:orden_aux.estado}, nil
  //}

  type orden struct {
      created_time time.Time
      id_paquete string
      nombre string
      valor  int32
      origen string
      destino string
      prioridad int32
      seguimiento int32
      intentos int32
      estado string
  }

func checkError(message string, err error) {
      if err != nil {
          log.Fatal(message, err)
      }
  }
func NewOrden(ordenes []*orden, id_paquete string, nombre string,
  valor  int32, origen string, destino string, prioridad int32 ) *orden {
    orden := orden{id_paquete: id_paquete,nombre:nombre,valor:valor,
    origen:origen,destino:destino,prioridad:prioridad,intentos:0,estado:"Pendiente Entrega"}
    orden.created_time = time.Now()
    orden.seguimiento = NewCodeSeguimiento(ordenes)
    //file, err := os.Open("data.csv")
    //checkError("Cannot create file", err)
    //defer file.Close()

    ///writer := csv.NewWriter(file)
    //defer writer.Flush()
    //err1 := writer.Write(["holi"])
    //checkError("Cannot write to file", err1)
    return &orden
}

func NewCodeSeguimiento(ordenes []*orden) int32{
    if len(ordenes)==0 {
      return 1
    }
    return ordenes[len(ordenes)-1].seguimiento+1
}

func searchOrder( codigo_seguimiento int32) *orden {
  for _, v := range ordenes {
    if v.seguimiento == codigo_seguimiento {
          return v
        }
      }
      //Item404 := Item{id: "not_found", producto: "not_found", valor:"not_found", tienda:"not_found",destino:"not_found"}
      return &Orden404
    }

func recepcion_ordenes(){
  grpcServer := grpc.NewServer()

  pb.RegisterGreeterServer(grpcServer, &Server{})

  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %s", err)
  }
}


var ordenes []*orden
var Orden404 orden = orden{id_paquete: "not_found", nombre: "not_found", valor:1, origen:"not_found",destino:"not_found", prioridad: -1,seguimiento:-1,estado:"No Existe"}

func main() {
    fmt.Println("Gracias por iniciar el receptor de ordenes de SD X-Wing Team")
  	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9000))
  	if err != nil {
  		log.Fatalf("failed to listen: %v", err)
  	}
    go recepcion_ordenes()
    fmt.Println("Wena profe")
    opcion:=0
    for opcion!=-1{
        fmt.Println("Ingrese el numero de seguimiento para consultar estado o -1 para salir : ")
        fmt.Scanf("%d", &opcion)
    }

    //aux=NewOrden(ordenes,"Paquete2","Bebida","Iñakikun",2000,"chilito","Corea")
    //ordenes=append(ordenes,aux)
    //for i := 0; i < len(ordenes); i++ {
    //  fmt.Println(ordenes[i])
    //  fmt.Println(ordenes[i].created_time.Format(time.ANSIC))
    //  fmt.Println("////")
    //}
    //fmt.Println(aux.created_time)
    //fmt.Println(aux.created_time.Format(time.ANSIC))
}
