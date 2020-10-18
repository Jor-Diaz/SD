package main

import(
  "fmt"
  //"os"
  "sync"
  "time"
  "log"
  "net"
  //"encoding/csv"
  "google.golang.org/grpc"
  "context"
  pb"Lab1/SD/pipeline"

  )

  type Server struct {
      pb.UnimplementedGreeterServer
  }

  // SafeCounter is safe to use concurrently.
type SafeCounter struct {
	v   map[string]int
	mux sync.Mutex
}

func (s *Server) SayHello(ctx context.Context, in *pb.Message) (*pb.Message, error) {
  	log.Printf("Orden recibida con datos:   %s %s %d %s %s %d", in.Id,in.Producto,in.Valor,in.Tienda,in.Destino, in.Prioridad )
    aux:=NewOrden(ordenes,in.Id,in.Producto,in.Valor,in.Tienda,in.Destino,in.Prioridad)
    if(in.Prioridad==2){
      ordenes_retail=append(ordenes_retail,aux)
    }else if( in.Prioridad==0){
      ordenes_prioridad_0=append(ordenes_prioridad_0,aux)
    }else{
      ordenes_prioridad_1=append(ordenes_prioridad_1,aux)
    }
  	return &pb.Message{Seguimiento: aux.seguimiento,}, nil
  }

func (s *Server) ConEstado(ctx context.Context, in *pb.ConsultaEstado) (*pb.RespuestaCon, error) {
  	log.Printf("Cosulta recibida con datos:   %d", in.Seguimiento)
    orden_aux:=searchOrder(in.Seguimiento)
  	return &pb.RespuestaCon{Id: orden_aux.id_paquete,Producto:orden_aux.nombre,Valor:orden_aux.valor,Tienda:orden_aux.origen,Destino:orden_aux.destino,Prioridad:orden_aux.prioridad,Intentos:orden_aux.intentos,Estado:orden_aux.estado}, nil
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
    origen:origen,destino:destino,prioridad:prioridad,intentos:0,estado:"En Bodega"}
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

func NewCodeSeguimiento() int32{
    candados[3].mux.Lock()
    aux:=numero_seguimiento+1
    numero_seguimiento=numero_seguimiento+1
    candados[3].mux.Unlock()
    return aux
}

func searchOrder(codigo_seguimiento int32) *orden {
  for _, v := range ordenes_retail {
    if v.seguimiento == codigo_seguimiento {
          return v
        }
  }
  for _, v := range ordenes_prioridad_1 {
    if v.seguimiento == codigo_seguimiento {
          return v
        }
  }
  for _, v := range ordenes_prioridad_0 {
    if v.seguimiento == codigo_seguimiento {
          return v
        }
  }
  return &Orden404
}

func  recepcion_clientes(){
  lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 9000))
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  grpcServer := grpc.NewServer()

  pb.RegisterGreeterServer(grpcServer, &Server{})

  if err := grpcServer.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %s", err)
  }
}

var ordenes_retail []*orden
var ordenes_prioridad_0 []*orden
var ordenes_prioridad_1 []*orden
var numero_seguimiento int32
var candados []*SafeCounte
var Orden404 orden = orden{id_paquete: "not_found", nombre: "not_found", valor:1, origen:"not_found",destino:"not_found", prioridad: -1,seguimiento:-1,estado:"No Existe"}

func main() {
    fmt.Println("Gracias por iniciar el receptor de ordenes de SD X-Wing Team")
    candados= {SafeCounter{v: make(map[string]int)},SafeCounter{v: make(map[string]int)},SafeCounter{v: make(map[string]int)},SafeCounter{v: make(map[string]int)} }
    numero_seguimiento=0
    go recepcion_clientes()
    opcion:=0
    for opcion!=-1{
        fmt.Println("Ingrese -1 para cerrar el programa ")
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
