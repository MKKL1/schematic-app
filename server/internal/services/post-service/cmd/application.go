package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/postgres/db"
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

	//clientRed := server.NewRedisClient()
	////TODO Move somewhere else
	//client, err := rueidisaside.NewClient(rueidisaside.ClientOption{
	//	ClientBuilder: func(option rueidis.ClientOption) (rueidis.Client, error) {
	//		return clientRed, nil
	//	},
	//	ClientOption: rueidis.ClientOption{},
	//	ClientTTL:    time.Minute,
	//})

	//userCacheRepo := redis.NewCacheRepository(userRepo, client)

	//idNode, err := snowflake.NewNode(1)
	//if err != nil {
	//	panic(err)
	//}

	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			GetPostById: query.NewGetPostByIdHandler(postRepo),
		},
	}
}
