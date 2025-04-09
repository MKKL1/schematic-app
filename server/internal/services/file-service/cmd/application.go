package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	dMinio "github.com/MKKL1/schematic-app/server/internal/services/file-service/minio"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/ports"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
)

func setupApplication(
	ctx context.Context,
	logger zerolog.Logger,
	cfg *ApplicationConfig,
	dbPool *pgxpool.Pool,
	minioClient *minio.Client,
	cqrsHandler kafka.CqrsHandler,
	metricsClient metrics.Client,
) (app.Application, error) {
	logger.Info().Msg("Setting up File Service Application")

	queries := db.New(dbPool)
	fileRepo := postgres.NewFilePostgresRepository(queries)

	storageClient := dMinio.NewMinioStorageClient(minioClient, cfg.Minio.Buckets.Files, cfg.Minio.Buckets.Temp)

	application := app.Application{
		Commands: app.Commands{
			UploadTempFile:     command.NewUploadTempFileHandler(storageClient, fileRepo, cqrsHandler.EventBus, logger, metricsClient, cfg.Service.TmpExpire),
			DeleteExpiredFiles: command.NewDeleteExpiredFilesHandler(storageClient, fileRepo, logger, metricsClient),
			CommitTempFile:     command.NewCommitTempHandler(storageClient, fileRepo, cqrsHandler.EventBus, logger, metricsClient),
			PostCreatedHandler: command.NewPostCreatedHandler(cqrsHandler.CommandBus, logger), // Assuming CommandBus() is part of CqrsHandler interface
		},
		Queries: app.Queries{
			// GetImage: query.NewGetImageHandler(imageRepo, logger),
		},
	}

	logger.Info().Msg("Registering CQRS handlers with Kafka")
	// Register event/command handlers
	ports.NewEventHandlers(application, cqrsHandler)
	logger.Info().Msg("Application setup complete")
	return application, nil
}
