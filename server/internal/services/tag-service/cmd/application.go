package cmd

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/redis"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/tag-service/infra/postgres/db"
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
		Database: "sh_tag",
	})
	if err != nil {
		panic(err)
	}

	queries := db.New(dbPool)

	return app.Application{}
}
