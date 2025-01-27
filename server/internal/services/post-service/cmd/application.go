package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/pkg/rueidisaside"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres/db"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/redis"
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
		Database: "sh_schematic",
	})
	if err != nil {
		panic(err)
	}

	queries := db.New(dbPool)
	postRepo := postgres.NewPostPostgresRepository(queries)

	clientRed := server.NewRedisClient()
	//TODO Move somewhere else
	reuClient, err := rueidisaside.NewClient(rueidisaside.ClientOption{
		ClientBuilder: func(option rueidis.ClientOption) (rueidis.Client, error) {
			return clientRed, nil
		},
		ClientOption: rueidis.ClientOption{},
		ClientTTL:    time.Minute,
	})

	postCacheRepo := redis.NewPostCacheRepository(postRepo, reuClient)

	idNode, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	userService := client.NewUsersClient(ctx, ":8001")

	return app.Application{
		Commands: app.Commands{
			CreatePost: command.NewCreatePostHandler(postCacheRepo, idNode, userService),
		},
		Queries: app.Queries{
			GetPostById: query.NewGetPostByIdHandler(postCacheRepo),
		},
	}
}
