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
  aux:=NewOrden(in.Id,in.Producto,in.Valor,in.Tienda,in.Destino,in.Prioridad)
  if(in.Prioridad==2){
    candados[0].mux.Lock()
    ordenes_retail=append(ordenes_retail,aux)
    candados[0].mux.Unlock()
  }else if( in.Prioridad==1){
    candados[1].mux.Lock()
    ordenes_prioridad_1=append(ordenes_prioridad_1,aux)
    candados[1].mux.Unlock()
  }else{
    candados[2].mux.Lock()
    ordenes_prioridad_0=append(ordenes_prioridad_0,aux)
    candados[2].mux.Unlock()
  }
	return &pb.Message{Seguimiento: aux.seguimiento,}, nil
}

func (s *Server) ConEstado(ctx context.Context, in *pb.ConsultaEstado) (*pb.RespuestaCon, error) {
	log.Printf("Cosulta recibida con datos:   %d", in.Seguimiento)
  orden_aux:=searchOrder(in.Seguimiento)
	return &pb.RespuestaCon{Id: orden_aux.id_paquete,Producto:orden_aux.nombre,Valor:orden_aux.valor,Tienda:orden_aux.origen,Destino:orden_aux.destino,Prioridad:orden_aux.prioridad,Intentos:orden_aux.intentos,Estado:orden_aux.estado,IdCamion:orden_aux.id_camion}, nil
}

func (s *Server) Solpedido(ctx context.Context, in *pb.Solcamion) (*pb.RespuestaCon, error) {
  	log.Printf("Peticion de orden por camion de id  %d", in.IdCamion)
    var orden_aux orden
    if (in.IdCamion == 1 ){
      orden_aux=searchOrder_retail(1)
    }
  	return &pb.RespuestaCon{Id: orden_aux.id_paquete,Producto:orden_aux.nombre,Valor:orden_aux.valor,Tienda:orden_aux.origen,Destino:orden_aux.destino,Prioridad:orden_aux.prioridad,Intentos:orden_aux.intentos,Estado:orden_aux.estado,Seguimiento:orden_aux.seguimiento,IdCamion:orden_aux.id_camion}, nil
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
    estado int32//0 en bodega; 1 en camino ; 2 recibido; 3 no recibido; -1 no existe
    id_camion int32

}

func checkError(message string, err error) {
      if err != nil {
          log.Fatal(message, err)
      }
  }
func NewOrden( id_paquete string, nombre string,
  valor  int32, origen string, destino string, prioridad int32 ) *orden {
    orden := orden{id_paquete: id_paquete,nombre:nombre,valor:valor,
    origen:origen,destino:destino,prioridad:prioridad,intentos:0,estado:0,id_camion:-1}
    orden.created_time = time.Now()
    orden.seguimiento = NewCodeSeguimiento()
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

func searchOrder_retail(id_camion int32) *orden {
  i:=0
  for _, v := range ordenes_retail {
    if v.estado == 0 && v.id_camion==-1  {
          candados[0].mux.Lock()
          ordenes_retail[i].id_camion=id_camion
          candados[0].mux.Unlock()
          return v
        }
    i=i+1
  }
  i=0
  for _, v := range ordenes_prioridad_1 {
    if v.estado == 0 && v.id_camion==-1  {
          candados[1].mux.Lock()
          ordenes_retail[i].id_camion=id_camion
          ordenes_retail[i].estado=1
          candados[1].mux.Unlock()
          return v
        }
    i=i+1
  }
  return &Orden404
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
var candados []*SafeCounter
var Orden404 orden = orden{id_paquete: "not_found", nombre: "not_found", valor:1, origen:"not_found",destino:"not_found", prioridad: -1,seguimiento:-1,estado:-1,id_camion:-1}

func main() {
    fmt.Println("Gracias por iniciar el receptor de ordenes de SD X-Wing Team")
    can:=new(SafeCounter)
    can.v= make(map[string]int)
    candados=append(candados,can)
    for i :=1; i<4 ;i++ {
        can=new(SafeCounter)
        can.v= make(map[string]int)
        candados=append(candados,can)
    }
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
