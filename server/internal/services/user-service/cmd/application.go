package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/rueidisaside"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/infra/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/infra/postgres/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/infra/redis"
	"github.com/bwmarrin/snowflake"
	"github.com/redis/rueidis"
	"time"
)

func NewApplication(ctx context.Context) app.Application {
	dbPool, err := server.NewPostgreSQLClient(ctx, &server.PostgresConfig{
		Port:     "5432",
		Host:     "localhost",
		Username: "root",
		Password: "root",
		Database: "sh_user",
	})
	if err != nil {
		panic(err)
	}

	queries := db.New(dbPool)
	userRepo := postgres.NewUserPostgresRepository(queries)

	clientRed := server.NewRedisClient()
	//TODO Move somewhere else
	client, err := rueidisaside.NewClient(rueidisaside.ClientOption{
		ClientBuilder: func(option rueidis.ClientOption) (rueidis.Client, error) {
			return clientRed, nil
		},
		ClientOption: rueidis.ClientOption{},
		ClientTTL:    time.Minute,
	})

	userCacheRepo := redis.NewCacheRepository(userRepo, client)

	idNode, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	return app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(userCacheRepo, idNode),
		},
		Queries: app.Queries{
			GetUserById:  query.NewGetUserByIdHandler(userCacheRepo),
			GetUserBySub: query.NewGetUserBySubHandler(userCacheRepo),
		},
	}
}
