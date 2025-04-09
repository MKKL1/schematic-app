package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/common"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	grpcPkg "github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	httpPkg "github.com/MKKL1/schematic-app/server/internal/pkg/http"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/pkg/minio"
	"github.com/MKKL1/schematic-app/server/internal/pkg/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/ports"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"time"
)

const ServiceName = "file-service"

func main() {
	rt, err := common.NewRuntime[ApplicationConfig](ServiceName)
	if err != nil {
		fmt.Printf("Failed to initialize application runtime: %v\n", err)
		os.Exit(1)
	}

	logger := rt.Logger
	cfg := rt.Config
	ctx, _ := rt.InitCtx(0)
	appCloser := rt.Closer

	initCtx, initCancel := context.WithTimeout(ctx, 3*time.Minute)
	defer initCancel()

	dbPool, err := postgres.NewClient(initCtx, cfg.Database, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	postgres.AddToCloser(appCloser, dbPool, logger)

	minioClient, err := minio.NewClient(initCtx, cfg.Minio, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Minio")
	}

	cqrsHandler := kafka.NewCqrsHandler(cfg.Kafka, logger)
	appCloser.Add(func(ctx context.Context) error {
		return cqrsHandler.Close(ctx)
	})
	metricsClient := metrics.NewPrometheusMetrics()

	application, err := setupApplication(logger, cfg, dbPool, minioClient, cqrsHandler, metricsClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to setup application")
	}

	mux := http.NewServeMux()
	httpServer := ports.HttpServer{App: application}
	mux.HandleFunc("/upload-tmp", httpServer.UploadMultipartHandler)
	mux.Handle("/metrics", promhttp.Handler())

	rt.Go(func(ctx context.Context) error {
		logger.Info().Str("address", cfg.Server.Grpc.GetAddr()).Msg("Starting gRPC server")
		return grpcPkg.Run(ctx, cfg.Server.Grpc, logger, ports.NewFileGrpcErrorMapper(), func(s *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterFileServiceServer(s, srv)
			logger.Info().Msg("Registered FileService gRPC server")
		})
	})

	rt.Go(func(ctx context.Context) error {
		logger.Info().Str("address", cfg.Server.Http.GetAddr()).Msg("Starting HTTP server")
		return httpPkg.Run(ctx, cfg.Server.Http, logger, mux)
	})

	rt.Go(func(ctx context.Context) error {
		logger.Info().Msg("Starting Kafka handler")
		cqrsHandler.Run(ctx)
		logger.Info().Msg("Kafka handler finished")
		if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
			return fmt.Errorf("kafka handler stopped unexpectedly: %w", ctx.Err())
		}
		return nil
	})

	rt.Run()
}
