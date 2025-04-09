package main

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	dMinio "github.com/MKKL1/schematic-app/server/internal/services/file-service/minio"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/ports"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres/db"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func NewMinioClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("creating minio client: %v", err)
	}
	// Add a Ping check
	_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		return nil, fmt.Errorf("minio connection check: %w", err)
	}
	return minioClient, nil
}

func NewApplication(ctx context.Context, cfg *ApplicationConfig) (app.Application, func()) {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	logger.Info().Msg("Starting File Service Application Setup")

	// --- Database Setup ---
	logger.Info().Msg("Connecting to PostgreSQL")
	dbPool, err := server.NewPostgreSQLClient(ctx, cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
		panic(err)
	}
	queries := db.New(dbPool)
	fileRepo := postgres.NewFilePostgresRepository(queries)

	// --- Minio Setup ---
	logger.Info().Str("endpoint", cfg.Minio.Endpoint).Msg("Connecting to Minio")
	minioClient, err := NewMinioClient(cfg.Minio.Endpoint, cfg.Minio.AccessKey, cfg.Minio.SecretKey, cfg.Minio.UseSSL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize MinIO client")
		panic(err)
	}

	storageClient := dMinio.NewMinioStorageClient(minioClient, cfg.Minio.Buckets.Files, cfg.Minio.Buckets.Temp)
	logger.Info().Msg("MinIO client initialized successfully")
	// TODO: Add bucket existence checks/creation logic here if needed

	// --- CQRS Setup ---
	logger.Info().Strs("brokers", cfg.Kafka.Brokers).Msg("Setting up Kafka CQRS")
	cqrsHandler := kafka.NewCqrsHandler(kafka.KafkaConfig{Brokers: cfg.Kafka.Brokers}, logger)

	metricsClient := metrics.NewPrometheusMetrics()

	// --- Application Layer ---
	application := app.Application{
		Commands: app.Commands{
			UploadTempFile:     command.NewUploadTempFileHandler(storageClient, fileRepo, cqrsHandler.EventBus, logger, metricsClient, cfg.Upload.TmpExpire),
			DeleteExpiredFiles: command.NewDeleteExpiredFilesHandler(storageClient, fileRepo, logger, metricsClient),
			CommitTempFile:     command.NewCommitTempHandler(storageClient, fileRepo, cqrsHandler.EventBus, logger, metricsClient),
			PostCreatedHandler: command.NewPostCreatedHandler(cqrsHandler.CommandBus, logger),
		},
		Queries: app.Queries{
			// Instantiate image query handler later
			// GetImage: query.NewGetImageHandler(imageRepo, logger /*, &someConfig */),
		},
	}

	// --- Ports Setup (Event Handlers, etc.) ---
	logger.Info().Msg("Registering CQRS handlers")
	ports.NewEventHandlers(application, cqrsHandler) // Register event/command handlers

	// Create cleanup function
	cleanup := func() {
		logger.Info().Msg("Shutting down application components...")
		// Add cleanup for Kafka, DB pool, etc.
		dbPool.Close()
		logger.Info().Msg("Cleanup finished.")
	}

	// Run CQRS handler (must be after handlers are registered)
	logger.Info().Msg("Starting CQRS handler")
	cqrsHandler.Run(ctx)
	logger.Info().Msg("Application setup complete")
	return application, cleanup
}
