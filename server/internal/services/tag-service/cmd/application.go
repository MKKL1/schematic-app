package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/domain/category"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/infra/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/infra/postgres/db"
)

func NewApplication(ctx context.Context) app.Application {
	dbPool, err := server.NewPostgreSQLClient(ctx, &server.PostgresConfig{
		Port:     "5432",
		Host:     "localhost",
		Username: "root",
		Password: "root",
		Database: "sh_tag",
	})
	if err != nil {
		panic(err)
	}

	queries := db.New(dbPool)
	repo := postgres.NewCategoryPostgresRepository(queries)
	provider := category.DefaultSchemaProvider{}

	return app.Application{
		Commands: app.Commands{
			CreateCategoryVars: command.NewCreatePostCatValuesHandler(repo, provider),
		},
		Queries: app.Queries{},
	}
}
