package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"slices"
	"sync"

	"github.com/starathel/gchat/gen/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessageStream = grpc.BidiStreamingServer[hello.HelloMessage, hello.HelloMessage]

type HelloServer struct {
	hello.UnimplementedHelloServiceServer
	mu                sync.RWMutex
	active_streams    []MessageStream
	incoming_messages chan *hello.HelloMessage
}

func (s *HelloServer) SayHello(_ context.Context, msg *hello.HelloMessage) (*hello.HelloMessage, error) {
	if len(msg.GetUsername()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Username should be not empty")
	}
	return &hello.HelloMessage{
		Username: "server",
		Text:     fmt.Sprintf("Hello %s. From server", msg.Username),
	}, nil
}

func (s *HelloServer) Chat(stream MessageStream) error {
	s.mu.Lock()
	s.active_streams = append(s.active_streams, stream)
	s.mu.Unlock()

	end := make(chan struct{})
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				s.mu.Lock()
				s.active_streams = slices.DeleteFunc(s.active_streams, func(s MessageStream) bool { return s == stream })
				s.mu.Unlock()
				end <- struct{}{}
				return
			}
			s.incoming_messages <- msg
		}
	}()
	<-end
	return nil
}

func (s *HelloServer) processMessages() {
	for msg := range s.incoming_messages {
		s.mu.RLock()
		for _, stream := range s.active_streams {
			err := stream.Send(msg)
			if err != nil {
				fmt.Println("WARNING: Attempt to write into closed stream")
				continue
			}
		}
		s.mu.RUnlock()
	}
}

func NewHelloServer() *HelloServer {
	server := HelloServer{}
	server.incoming_messages = make(chan *hello.HelloMessage, 10)
	return &server
}

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:6969")
	if err != nil {
		log.Fatalf("cannot start listening on 6969: %v", err)
	}
	fmt.Println("Started listener")
	grpcServer := grpc.NewServer()
	helloServer := NewHelloServer()
	hello.RegisterHelloServiceServer(grpcServer, helloServer)
	go helloServer.processMessages()
	err = grpcServer.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Started server")
}
