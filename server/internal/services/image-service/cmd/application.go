package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/imgproxy"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/ports"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/postgres/db"
	"github.com/rs/zerolog"
	"os"
)

func NewApplication(ctx context.Context) app.Application {
	dbPool, err := server.NewPostgreSQLClient(ctx, &server.PostgresConfig{
		Port:     "5432",
		Host:     "localhost",
		Username: "root",
		Password: "root",
		Database: "sh_images",
	})
	if err != nil {
		panic(err)
	}

	queries := db.New(dbPool)
	repo := postgres.NewImagePostgresRepository(queries)
	cqrsHandler := kafka.NewCqrsHandler(kafka.KafkaConfig{Brokers: []string{"localhost:9092"}})
	//TODO get from env var (not prod safe)
	urlGen := imgproxy.NewUrlGenerator("s3://files/", "http://localhost:9004", "2b33e772e9d2c041f289a5babf817bd0f4247be9d9e028cf9bf1b359d5cd6641", "087c98d789b3d79f145c1ca33cfcd2456daad64945816881de493fe875045d10")
	logger := zerolog.New(os.Stdout)

	a := app.Application{
		Commands: app.Commands{
			ProcessUploadedImage: command.NewProcessUploadedImageHandler(repo, cqrsHandler.EventBus, logger),
		},
		Queries: app.Queries{
			GetImageSizes: query.NewGetImageHandler(repo, urlGen, logger),
		},
	}

	ports.NewEventHandlers(a, cqrsHandler)
	cqrsHandler.Run(ctx)

	return a
}
