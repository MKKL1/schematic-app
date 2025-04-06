package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/ports"
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

		server.RunGRPCServer(ctx, ":8010", ports.NewImageGrpcErrorMapper(), func(server *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterImageServiceServer(server, srv)
		})

	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
