package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres/db"
	"github.com/bwmarrin/snowflake"
)

func NewApplication(ctx context.Context) app.Application {
	dbPool, err := server.NewPostgreSQLClient(ctx, &server.PostgresConfig{
		Port:     "5432",
		Host:     "localhost",
		Username: "root",
		Password: "root",
		Database: "sh_post",
	})
	if err != nil {
		panic(err)
	}

	queries := db.New(dbPool)
	postRepo := postgres.NewPostPostgresRepository(queries)
	categoryRepo := postgres.NewCategoryPostgresRepository(queries)

	//clientRed := server.NewRedisClient()
	//TODO Move somewhere else
	//reuClient, err := rueidisaside.NewClient(rueidisaside.ClientOption{
	//	ClientBuilder: func(option rueidis.ClientOption) (rueidis.Client, error) {
	//		return clientRed, nil
	//	},
	//	ClientOption: rueidis.ClientOption{},
	//	ClientTTL:    time.Minute,
	//})
	//
	//postCacheRepo := redis.NewPostCacheRepository(postRepo, reuClient)

	idNode, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	userService := client.NewUsersClient(ctx, ":8001")

	schemaProvider := category.DefaultSchemaProvider{}

	return app.Application{
		Commands: app.Commands{
			CreatePost: command.NewCreatePostHandler(postRepo, categoryRepo, idNode, userService, schemaProvider),
		},
		Queries: app.Queries{
			GetPostById: query.NewGetPostByIdHandler(postRepo),
		},
	}
}
