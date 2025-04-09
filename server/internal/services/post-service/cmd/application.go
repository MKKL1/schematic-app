package main

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/client/user"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	postgres2 "github.com/MKKL1/schematic-app/server/internal/services/post-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/postgres/db"
	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func setupApplication(logger zerolog.Logger, cfg *ApplicationConfig, dbPool *pgxpool.Pool, cqrsHandler kafka.CqrsHandler, metricsClient metrics.Client, grpcClient *grpc.ClientConn) (app.Application, error) {
	queries := db.New(dbPool)
	postRepo := postgres2.NewPostPostgresRepository(queries)
	categoryRepo := postgres2.NewCategoryPostgresRepository(queries)

	idNode, err := snowflake.NewNode(1) //Not sure where it should be
	if err != nil {
		return app.Application{}, err
	}

	userService := user.NewGrpcService(grpcClient)

	a := app.Application{
		Commands: app.Commands{
			CreatePost:     command.NewCreatePostHandler(postRepo, categoryRepo, idNode, userService, cqrsHandler.EventBus, logger, metricsClient),
			UpdateFileHash: command.NewUpdateAttachedFilesHandler(postRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetPostById: query.NewGetPostByIdHandler(postRepo, logger, metricsClient),
		},
	}

	return a, nil
}
