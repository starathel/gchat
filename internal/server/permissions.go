package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthorizationRequired(ctx context.Context) error {
	username, ok := ctx.Value(usernameKey).(string)
	if !ok || username == "" {
		return status.Error(codes.Unauthenticated, "username required")
	}
	return nil
}
