package main

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/infra/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/infra/postgres/db"
	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

//func NewApplication(ctx context.Context) app.Application {
//	dbPool, err := server.NewPostgreSQLClient(ctx, &server.PostgresConfig{
//		Port:     "5432",
//		Host:     "localhost",
//		Username: "root",
//		Password: "root",
//		Database: "sh_user",
//	})
//	if err != nil {
//		panic(err)
//	}
//
//}

func setupApplication(
	logger zerolog.Logger,
	cfg *ApplicationConfig,
	dbPool *pgxpool.Pool,
	metricsClient metrics.Client,
) (app.Application, error) {

	idNode, err := snowflake.NewNode(1)
	if err != nil {
		return app.Application{}, err
	}
	queries := db.New(dbPool)
	userRepo := postgres.NewUserPostgresRepository(queries)
	//userCacheRepo := redis.NewCacheRepository(userRepo, clientRedis)

	return app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(userRepo, idNode),
		},
		Queries: app.Queries{
			GetUserById:  query.NewGetUserByIdHandler(userRepo),
			GetUserBySub: query.NewGetUserBySubHandler(userRepo),
		},
	}, nil
}
