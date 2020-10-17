package pipeline

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Orden recibida con datos:  %s %s %s %d %s %s", in.Tipo,in.Id,in.Producto,in.Valor,in.Tienda,in.Destino )
	return &Message{Tipo: " Datos recibidos",}, nil
}
