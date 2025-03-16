package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/infra/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/infra/postgres/db"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/ports"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	// Create a new MinIO client using the provided credentials and endpoint.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Optional: Check if you can list buckets or perform another operation to verify connection.
	return minioClient, nil
}

func NewApplication(ctx context.Context) app.Application {
	dbPool, err := server.NewPostgreSQLClient(ctx, &server.PostgresConfig{
		Port:     "5432",
		Host:     "localhost",
		Username: "root",
		Password: "root",
		Database: "sh_file",
	})
	if err != nil {
		panic(err)
	}

	queries := db.New(dbPool)
	repo := postgres.NewFilePostgresRepository(queries)

	endpoint := "localhost:9000"
	accessKeyID := "GdD14n1Oxz2U5hfQhdHo"
	secretAccessKey := "e1Peh4RLq7E4hgDW3GtV8nl4IaZjGrzDuS0WTPaB"

	minioClient, err := NewMinioClient(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	log.Printf("MinIO client initialized successfully at endpoint: %s", endpoint)

	cqrsHandler := kafka.NewCqrsHandler(kafka.KafkaConfig{Brokers: []string{"localhost:9092"}})

	a := app.Application{
		Commands: app.Commands{
			UploadTempFile:     command.NewUploadTempFileHandler(minioClient, repo, cqrsHandler.EventBus),
			DeleteExpiredFiles: command.NewDeleteExpiredFilesHandler(minioClient, repo),
			CommitTempFile:     command.NewCommitTempHandler(minioClient, repo, cqrsHandler.EventBus),
			PostCreatedHandler: command.NewPostCreatedHandler(cqrsHandler.CommandBus),
		},
		Queries: app.Queries{},
	}

	ports.NewEventHandlers(a, cqrsHandler)
	//Run has to be executed after handlers are registered
	cqrsHandler.Run(ctx)

	return a
}
