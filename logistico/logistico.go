package main

import(
  "fmt"
  "time"
  "log"
  "net"
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

  func (s *Server) SayHello(ctx context.Context, in *pb.Message) (*pb.Message, error) {
  	log.Printf("Orden recibida con datos:   %s %s %d %s %s %d", in.Id,in.Producto,in.Valor,in.Tienda,in.Destino, in.Prioridad )
    aux:=NewOrden(ordenes,in.Id,in.Producto,in.Valor,in.Tienda,in.Destino,in.Prioridad)
    ordenes=append(ordenes,aux)
    //fmt.Println("#################")
    //fmt.Println("_Data en memoria__")
    //for i := 0; i < len(ordenes); i++ {
      //fmt.Println(ordenes[i])
      //fmt.Println(ordenes[i].created_time.Format(time.ANSIC))
      //fmt.Println("_________________")
    //}
  	return &pb.Message{Seguimiento: aux.seguimiento,}, nil
  }

  type orden struct {
      created_time time.Time
      id_paquete string
      nombre string
      valor  int32
      origen string
      destino string
      prioridad int32
      seguimiento int32
  }

func NewOrden(ordenes []*orden, id_paquete string, nombre string,
  valor  int32, origen string, destino string, prioridad int32 ) *orden {
    orden := orden{id_paquete: id_paquete,nombre:nombre,valor:valor,
    origen:origen,destino:destino,prioridad:prioridad}
    orden.created_time = time.Now()
    orden.seguimiento = NewCodeSeguimiento(ordenes)
    return &orden
}

func NewCodeSeguimiento(ordenes []*orden) int32{
    if len(ordenes)==0 {
      return 1
    }
    return ordenes[len(ordenes)-1].seguimiento+1
}

var ordenes []*orden
func main() {
    fmt.Println("Gracias por iniciar el receptor de ordenes de SD X-Wing Team")
  	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9000))
  	if err != nil {
  		log.Fatalf("failed to listen: %v", err)
  	}

  	grpcServer := grpc.NewServer()

  	pb.RegisterGreeterServer(grpcServer, &Server{})

  	if err := grpcServer.Serve(lis); err != nil {
  		log.Fatalf("failed to serve: %s", err)
  	}
    fmt.Println("Wena profe")

    //aux=NewOrden(ordenes,"Paquete2","Bebida","Iñakikun",2000,"chilito","Corea")
    //ordenes=append(ordenes,aux)
    for i := 0; i < len(ordenes); i++ {
      fmt.Println(ordenes[i])
      fmt.Println(ordenes[i].created_time.Format(time.ANSIC))
      fmt.Println("////")
    }
    //fmt.Println(aux.created_time)
    //fmt.Println(aux.created_time.Format(time.ANSIC))
}
