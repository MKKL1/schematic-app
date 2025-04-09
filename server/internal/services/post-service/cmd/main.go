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
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/ports"
	"google.golang.org/grpc"
	"os"
	"time"
)

/*
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
*/

const ServiceName = "post-service"

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
	metricsClient := metrics.LogClient{Logger: logger}

	userGrpcClient, err := grpcPkg.NewClient(cfg.User)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create user grpc client")
	}

	application, err := setupApplication(logger, cfg, dbPool, cqrsHandler, metricsClient, userGrpcClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to setup application")
	}

	rt.Go(func(ctx context.Context) error {
		logger.Info().Str("address", cfg.Server.Grpc.GetAddr()).Msg("Starting gRPC server")
		return grpcPkg.Run(ctx, cfg.Server.Grpc, logger, ports.NewPostGrpcErrorMapper(), func(s *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterPostServiceServer(s, srv)
			logger.Info().Msg("Registered PostService gRPC server")
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
