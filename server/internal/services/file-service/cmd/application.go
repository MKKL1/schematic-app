package main

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/image"
	dMinio "github.com/MKKL1/schematic-app/server/internal/services/file-service/minio"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/ports"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres/db"
	"github.com/bwmarrin/snowflake"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	// Add a Ping check
	_, err = minioClient.ListBuckets(context.Background()) // Basic check
	if err != nil {
		return nil, fmt.Errorf("minio connection check failed: %w", err)
	}
	return minioClient, nil
}

// Consider moving config loading to a dedicated function/package
type ServiceConfig struct {
	PostgresURL      string   `env:"POSTGRES_URL" envDefault:"postgres://root:root@localhost:5432/sh_file?sslmode=disable"`
	KafkaBrokers     []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"localhost:9092"`
	MinioEndpoint    string   `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	MinioAccessKey   string   `env:"MINIO_ACCESS_KEY" envDefault:"GdD14n1Oxz2U5hfQhdHo"`
	MinioSecretKey   string   `env:"MINIO_SECRET_KEY" envDefault:"e1Peh4RLq7E4hgDW3GtV8nl4IaZjGrzDuS0WTPaB"`
	MinioUseSSL      bool     `env:"MINIO_USE_SSL" envDefault:"false"`
	MinioFileBucket  string   `env:"MINIO_FILE_BUCKET" envDefault:"files"`
	MinioTempBucket  string   `env:"MINIO_TEMP_BUCKET" envDefault:"temp-bucket"`
	MinioImageBucket string   `env:"MINIO_IMAGE_BUCKET" envDefault:"images"` // New config
	SnowflakeNodeID  int64    `env:"SNOWFLAKE_NODE_ID" envDefault:"1"`       // New config
	LogLevel         string   `env:"LOG_LEVEL" envDefault:"info"`
}

func NewApplication(ctx context.Context) app.Application {
	cfg := ServiceConfig{}
	if err := config.Load(&cfg); err != nil { // Use your config loading mechanism
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// --- Logger Setup ---
	logger := logging.NewLogger(cfg.LogLevel) // Use your logger setup

	logger.Info().Msg("Starting File Service Application Setup")

	// --- Database Setup ---
	logger.Info().Msg("Connecting to PostgreSQL")
	dbPool, err := server.NewPostgreSQLClient(ctx, cfg.PostgresURL) // Use URL from config
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	queries := db.New(dbPool)
	fileRepo := postgres.NewFilePostgresRepository(queries)
	imageRepo := postgres.NewImagePostgresRepository(queries) // Create image repo

	// --- Minio Setup ---
	logger.Info().Str("endpoint", cfg.MinioEndpoint).Msg("Connecting to Minio")
	minioClient, err := NewMinioClient(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioUseSSL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize MinIO client")
	}
	// Pass all bucket names
	storageClient := dMinio.NewMinioStorageClient(minioClient, cfg.MinioFileBucket, cfg.MinioTempBucket, cfg.MinioImageBucket)
	logger.Info().Msg("MinIO client initialized successfully")
	// TODO: Add bucket existence checks/creation logic here if needed

	// --- Snowflake ID Generator ---
	logger.Info().Int64("nodeId", cfg.SnowflakeNodeID).Msg("Initializing Snowflake node")
	idNode, err := snowflake.NewNode(cfg.SnowflakeNodeID)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create Snowflake node")
	}
	imageIDGenerator := image.NewSnowflakeIDGenerator(idNode)

	// --- CQRS Setup ---
	logger.Info().Strs("brokers", cfg.KafkaBrokers).Msg("Setting up Kafka CQRS")
	cqrsHandler := kafka.NewCqrsHandler(kafka.KafkaConfig{Brokers: cfg.KafkaBrokers}, logger) // Pass logger

	// --- Application Layer ---
	application := app.Application{
		Commands: app.Commands{
			UploadTempFile:     command.NewUploadTempFileHandler(storageClient, fileRepo, cqrsHandler.EventBus),
			DeleteExpiredFiles: command.NewDeleteExpiredFilesHandler(storageClient, fileRepo),
			CommitTempFile:     command.NewCommitTempHandler(storageClient, fileRepo, cqrsHandler.EventBus, logger),          // Pass logger
			PostCreatedHandler: command.NewPostCreatedHandler(cqrsHandler.CommandBus, logger),                                // Pass logger
			ProcessImage:       command.NewProcessImageHandler(fileRepo, imageRepo, storageClient, imageIDGenerator, logger), // Instantiate new handler
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
		cqrsHandler.Stop()
		dbPool.Close()
		logger.Info().Msg("Cleanup finished.")
	}

	// Run CQRS handler (must be after handlers are registered)
	logger.Info().Msg("Starting CQRS handler")
	runCtx, cancelRun := context.WithCancel(ctx) // Create a context for the run loop
	go func() {
		defer cancelRun() // Ensure cancel is called if Run exits
		if err := cqrsHandler.Run(runCtx); err != nil {
			// Log error unless it's context cancellation during shutdown
			if !errors.Is(err, context.Canceled) {
				logger.Error().Err(err).Msg("CQRS handler stopped with error")
			} else {
				logger.Info().Msg("CQRS handler stopped gracefully.")
			}
		}
	}()

	logger.Info().Msg("Application setup complete")
	return application, cleanup // Return app and cleanup function
}
