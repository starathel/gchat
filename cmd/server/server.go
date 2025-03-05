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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey int

const (
	usernameKey contextKey = 0
)

type MessageStream = grpc.BidiStreamingServer[hello.HelloMessage, hello.HelloMessage]

type HelloServer struct {
	hello.UnimplementedHelloServiceServer
	mu                sync.RWMutex
	active_streams    []MessageStream
	incoming_messages chan *hello.HelloMessage
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s wrappedStream) Context() context.Context {
	return s.ctx
}

func newWrappedStream(s grpc.ServerStream) wrappedStream {
	return wrappedStream{ServerStream: s, ctx: s.Context()}
}

func AuthorizeUnary(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "missing metadata")
	}
	username := md["authorization"]
	if len(username) < 1 {
		return handler(ctx, req)
	}
	return handler(context.WithValue(ctx, usernameKey, username[0]), req)
}

func AuthorizaStream(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.InvalidArgument, "missing metadata")
	}
	username := md["authorization"]
	stream := newWrappedStream(ss)
	if len(username) >= 1 {
		stream.ctx = context.WithValue(stream.ctx, usernameKey, username[0])
	}
	return handler(srv, stream)
}

func authorize(ctx context.Context) error {
	username := ctx.Value(usernameKey)
	if username == nil {
		return status.Error(codes.Unauthenticated, "username required")
	}
	return nil
}

func (s *HelloServer) SayHello(ctx context.Context, msg *hello.HelloMessage) (*hello.HelloMessage, error) {
	err := authorize(ctx)
	if err != nil {
		return nil, err
	}
	return &hello.HelloMessage{
		Username: "server",
		Text:     fmt.Sprintf("Hello %s. From server", ctx.Value(usernameKey)),
	}, nil
}

func (s *HelloServer) Chat(stream MessageStream) error {
	if err := authorize(stream.Context()); err != nil {
		return status.Error(codes.Unauthenticated, "stream with no username")
	}
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
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(AuthorizeUnary), grpc.StreamInterceptor(AuthorizaStream))
	helloServer := NewHelloServer()
	hello.RegisterHelloServiceServer(grpcServer, helloServer)
	go helloServer.processMessages()
	err = grpcServer.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Started server")
}
