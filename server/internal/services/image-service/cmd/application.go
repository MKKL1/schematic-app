package main

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/imgproxy"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/ports"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func setupApplication(
	logger zerolog.Logger,
	cfg *ApplicationConfig,
	dbPool *pgxpool.Pool,
	cqrsHandler kafka.CqrsHandler,
	urlGen *imgproxy.UrlGenerator,
	metricsClient metrics.Client,
) (app.Application, error) {
	logger.Info().Msg("Setting up Image Service Application")

	queries := db.New(dbPool)
	repo := postgres.NewImagePostgresRepository(queries)

	application := app.Application{
		Commands: app.Commands{
			ProcessUploadedImage: command.NewProcessUploadedImageHandler(repo, cqrsHandler.EventBus, logger, metricsClient),
		},
		Queries: app.Queries{
			GetImageSizes: query.NewGetImageHandler(repo, urlGen, logger, metricsClient),
		},
	}

	logger.Info().Msg("Registering CQRS handlers with Kafka")
	ports.NewEventHandlers(application, cqrsHandler)

	logger.Info().Msg("Application setup complete")
	return application, nil
}
