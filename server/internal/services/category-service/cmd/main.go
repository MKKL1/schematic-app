package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/category-service/ports"
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
			Port:     "1325",
			BasePath: "/",
			Timeout:  10000,
			Host:     "localhost",
		})

		application := NewApplication(ctx)

		//_, err := application.Commands.CreateCategoryVars.Handle(ctx, command.CreatePostCatValuesParams{
		//	PostId:     532,
		//	CategoryId: 3,
		//	Values:     []byte(`{"afkable": {"max": 433, "min": 212}, "mob_type": "spider", "spawn_rate": {"max": 433, "min": 212}}`),
		//})
		//if err != nil {
		//	log.Err(err).Msg("Error creating category vars")
		//}

		server.RunGRPCServer(ctx, ":8003", ports.ErrorMapper, func(server *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterCategoryServiceServer(server, srv)
		})

		//postController := http.NewPostController(application)
		//http.RegisterRoutes(e, postController)
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
