package main

import(
  "fmt"
  "time"
  )

  type orden struct {
      created_time time.Time 'bson:"updated_at,omitempty" json:"updated_at,omitempty"
      id_paquete string
      tipo string
      nombre string
      valor  int
      origen string
      destino string
      seguimiento string
  }

func NewOrden( id_paquete string, tipo string, nombre string,
  valor  int, origen string, destino string ) *orden {
    orden := orden{id_paquete: id_paquete,tipo:tipo,nombre:nombre,valor:valor,
    origen:origen,destino:destino}
    orden.created_time = time.Now()
    orden.seguimiento = "000001"
    return &orden
}


func main() {
    fmt.Println("Wena profe")
    fmt.Println(NewOrden("1a1a1a1","gg","Jorgekun",1000,"chilito","membrillo"))
}
