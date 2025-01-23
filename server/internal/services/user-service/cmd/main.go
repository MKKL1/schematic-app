package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/http/middlewares"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/http"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/ports"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"time"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		e := server.NewEchoServer()
		server.RunHttpServer(ctx, e, &server.EchoConfig{
			Port:     "1323",
			BasePath: "/",
			Timeout:  10000,
			Host:     "localhost",
		})

		e.HTTPErrorHandler = middlewares.HTTPErrorHandler(http.MapAppError)

		application := NewApplication(ctx)

		server.RunGRPCServer(ctx, ":8001", func(server *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterUserServiceServer(server, srv)
		})

		userController := ports.NewUserController(application)
		ports.RegisterRoutes(e, userController)
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
