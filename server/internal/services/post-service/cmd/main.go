package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	grpc2 "github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/ports"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"time"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		application := NewApplication(ctx)

		grpc2.RunGRPCServer(ctx, ":8002", ports.NewPostGrpcErrorMapper(), func(server *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterPostServiceServer(server, srv)
		})

	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
