package cmd

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/http/middlewares"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/ports/http"
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
			Port:     "1325",
			BasePath: "/",
			Timeout:  10000,
			Host:     "localhost",
		})

		e.HTTPErrorHandler = middlewares.HTTPErrorHandler(http.MapAppError)

		application := NewApplication(ctx)

		//server.RunGRPCServer(ctx, ":8002", func(server *grpc.Server) {
		//	srv := ports.NewGrpcServer(application)
		//	genproto.RegisterUserServiceServer(server, srv)
		//})

		postController := http.NewPostController(application)
		http.RegisterRoutes(e, postController)
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
