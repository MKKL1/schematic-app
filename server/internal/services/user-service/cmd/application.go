package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres/db"
	"github.com/bwmarrin/snowflake"
)

func NewApplication(ctx context.Context) app.Application {
	dbPool, err := server.NewPostgreSQL(ctx, &server.PostgresConfig{
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

	idNode, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	return app.Application{
		Commands: app.Commands{
			CreateUser: command.NewCreateUserHandler(userRepo, idNode),
		},
		Queries: app.Queries{
			GetUserById:  query.NewGetUserByIdHandler(userRepo),
			GetUserBySub: query.NewGetUserBySubHandler(userRepo),
		},
	}
}
