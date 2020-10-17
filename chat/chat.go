package chat

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Orden recibida con datos:  %s %s %s %d %s %s", in.tipo,in.id,in.producto,in.valor,in.tienda,in.destino )
	return &Message{producto: " Datos recibidos",}, nil
}
