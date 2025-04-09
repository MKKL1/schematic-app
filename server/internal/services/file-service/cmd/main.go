package main

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/config"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/ports"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		cfg, err := config.LoadConfig[ApplicationConfig]("config.yaml")
		if err != nil {
			panic(fmt.Errorf("loading config: %v", err))
		}

		application, _ := NewApplication(ctx, cfg)

		server.RunGRPCServer(ctx, cfg.Server.Grpc.Host, ports.NewFileGrpcErrorMapper(), func(server *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterFileServiceServer(server, srv)
		})

		httpServer := ports.HttpServer{
			App: application,
		}
		http.HandleFunc("/upload-tmp", httpServer.UploadMultipartHandler)
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(cfg.Server.Http.Host, nil)
		if err != nil {
			return
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
