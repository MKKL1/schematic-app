package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"
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
	endpoint := "localhost:9000"
	accessKeyID := "GdD14n1Oxz2U5hfQhdHo"
	secretAccessKey := "e1Peh4RLq7E4hgDW3GtV8nl4IaZjGrzDuS0WTPaB"

	minioClient, err := NewMinioClient(endpoint, accessKeyID, secretAccessKey, false)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	log.Printf("MinIO client initialized successfully at endpoint: %s", endpoint)

	return app.Application{
		Commands: app.Commands{
			UploadTempFile: command.NewUploadTempFileHandler(minioClient),
		},
		Queries: app.Queries{},
	}
}
