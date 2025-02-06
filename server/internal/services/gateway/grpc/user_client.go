package grpc

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
)

func NewUserServiceClient(ctx context.Context) client.PostApplication {
	userService := client.NewUsersClient(ctx, ":8001")
	return userService
}
