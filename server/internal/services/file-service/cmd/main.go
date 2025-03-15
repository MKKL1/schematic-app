package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/infra/kafka"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/ports"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		application := NewApplication(ctx)

		server.RunGRPCServer(ctx, ":8005", post.ErrorMapper, func(server *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterFileServiceServer(server, srv)
		})

		brokers := []string{"localhost:9092"}
		sub, err := kafka.NewKafkaSubscriber(brokers)
		ports.NewEventListener(ctx, sub, application)

		httpServer := ports.HttpServer{
			App: application,
		}
		http.HandleFunc("/upload-tmp", httpServer.UploadMultipartHandler)
		err = http.ListenAndServe(":8006", nil)
		if err != nil {
			return
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
