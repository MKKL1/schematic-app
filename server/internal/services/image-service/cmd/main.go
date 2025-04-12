package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/common"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	grpcPkg "github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/pkg/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/imgproxy"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/ports"
	"google.golang.org/grpc"
	"os"
	"time"
)

const ServiceName = "image-service"

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

	cqrsHandler := kafka.NewCqrsHandler(cfg.Kafka, logger)
	appCloser.Add(func(ctx context.Context) error {
		return cqrsHandler.Close(ctx)
	})

	metricsClient := metrics.NewPrometheusMetrics()
	urlGen := imgproxy.NewUrlGeneratorFromConfig(cfg.ImgProxy)

	application, err := setupApplication(logger, cfg, dbPool, cqrsHandler, urlGen, metricsClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to setup application")
	}

	rt.Go(func(ctx context.Context) error {
		logger.Info().Str("address", cfg.Server.Grpc.GetAddr()).Msg("Starting gRPC server")
		return grpcPkg.Run(ctx, cfg.Server.Grpc, logger, ports.NewImageGrpcErrorMapper(), func(s *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterImageServiceServer(s, srv)
			logger.Info().Msg("Registered ImageService gRPC server")
		})
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
