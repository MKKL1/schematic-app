package main

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/config"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
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
	_, err = minioClient.ListBuckets(context.Background()) // Basic check
	if err != nil {
		return nil, fmt.Errorf("minio connection check failed: %w", err)
	}
	return minioClient, nil
}

func NewApplication(ctx context.Context) (app.Application, error) {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	cfg, err := config.LoadConfig[ApplicationConfig]("config.yaml")
	if err != nil {
		logger.Fatal().Err(err).Msg("Loading config failed")
		return app.Application{}, err
	}
	logger.Info().Msg("Loaded config from config.yaml")

	logger.Info().Msg("Starting File Service Application Setup")

	// --- Database Setup ---
	logger.Info().Msg("Connecting to PostgreSQL")
	dbPool, err := server.NewPostgreSQLClient(ctx, &server.PostgresConfig{
		Port:     "5432",
		Host:     "localhost",
		Username: "root",
		Password: "root",
		Database: "sh_file",
	}) //TODO Use URL from config
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
		return app.Application{}, err
	}
	queries := db.New(dbPool)
	fileRepo := postgres.NewFilePostgresRepository(queries)

	// --- Minio Setup ---
	logger.Info().Str("endpoint", cfg.Minio.Endpoint).Msg("Connecting to Minio")
	minioClient, err := NewMinioClient(cfg.Minio.Endpoint, cfg.Minio.AccessKey, cfg.Minio.SecretKey, cfg.Minio.UseSSL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize MinIO client")
		return app.Application{}, err
	}
	// Pass all bucket names
	storageClient := dMinio.NewMinioStorageClient(minioClient, cfg.Minio.Buckets.Files, cfg.Minio.Buckets.Temp)
	logger.Info().Msg("MinIO client initialized successfully")
	// TODO: Add bucket existence checks/creation logic here if needed

	// --- CQRS Setup ---
	logger.Info().Strs("brokers", cfg.Kafka.Brokers).Msg("Setting up Kafka CQRS")
	cqrsHandler := kafka.NewCqrsHandler(kafka.KafkaConfig{Brokers: cfg.Kafka.Brokers}) // Pass logger

	// --- Application Layer ---
	application := app.Application{
		Commands: app.Commands{
			UploadTempFile:     command.NewUploadTempFileHandler(storageClient, fileRepo, cqrsHandler.EventBus),
			DeleteExpiredFiles: command.NewDeleteExpiredFilesHandler(storageClient, fileRepo),
			CommitTempFile:     command.NewCommitTempHandler(storageClient, fileRepo, cqrsHandler.EventBus, logger), // Pass logger
			PostCreatedHandler: command.NewPostCreatedHandler(cqrsHandler.CommandBus, logger),                       // Pass logger
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
	//cleanup := func() {
	//	logger.Info().Msg("Shutting down application components...")
	//	// Add cleanup for Kafka, DB pool, etc.
	//	cqrsHandler.Stop()
	//	dbPool.Close()
	//	logger.Info().Msg("Cleanup finished.")
	//}

	// Run CQRS handler (must be after handlers are registered)
	logger.Info().Msg("Starting CQRS handler")
	cqrsHandler.Run(ctx)
	logger.Info().Msg("Application setup complete")
	return application, nil
}
