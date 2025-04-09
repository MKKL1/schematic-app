package main

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/common"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	grpcPkg "github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/pkg/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/ports"
	"google.golang.org/grpc"
	"os"
	"time"
)

const ServiceName = "user-service"

func main() {

	//ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	//defer stop()
	//
	//go func() {
	//	application := NewApplication(ctx)
	//
	//	grpc2.RunGRPCServer(ctx, ":8001", ports.NewUserGrpcErrorMapper(), func(server *grpc.Server) {
	//		srv := ports.NewGrpcServer(application)
	//		genproto.RegisterUserServiceServer(server, srv)
	//	})
	//
	//}()
	//
	//<-ctx.Done()
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

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

	metricsClient := metrics.LogClient{Logger: logger}

	application, err := setupApplication(logger, cfg, dbPool, metricsClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to setup application")
	}

	rt.Go(func(ctx context.Context) error {
		logger.Info().Str("address", cfg.Server.Grpc.GetAddr()).Msg("Starting gRPC server")
		return grpcPkg.Run(ctx, cfg.Server.Grpc, logger, ports.NewUserGrpcErrorMapper(), func(s *grpc.Server) {
			srv := ports.NewGrpcServer(application)
			genproto.RegisterUserServiceServer(s, srv)
			logger.Info().Msg("Registered UserService gRPC server")
		})
	})

	rt.Run()
}
