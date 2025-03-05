package main

import (
	"fmt"
	"log"
	"net"

	"github.com/starathel/gchat/gen/chat"
	"github.com/starathel/gchat/internal/server"
	"google.golang.org/grpc"
)

type Server struct {
	host string
	port string
}

func (s *Server) StartAndListen() error {
	l, err := net.Listen("tcp", "127.0.0.1:6969")
	if err != nil {
		return err
	}
	fmt.Println("Started Listener")
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(server.CurrentUserOnStream))
	chatServer := server.NewChatServer()
	chat.RegisterChatServiceServer(grpcServer, chatServer)

	fmt.Println("Starting Server")
	err = grpcServer.Serve(l)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	server := Server{
		host: "127.0.0.1",
		port: "6969",
	}
	err := server.StartAndListen()
	if err != nil {
		log.Fatal(err)
	}
}
