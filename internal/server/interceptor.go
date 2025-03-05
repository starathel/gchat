package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey int

const (
	usernameKey contextKey = 0
)

func CurrentUserOnStream(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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
