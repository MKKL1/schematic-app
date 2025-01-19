package config

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func ConfigDB(ctx context.Context) *db.Queries {
	dbPool, err := pgxpool.New(ctx, "postgres://root:root@localhost:5432/sh_user")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to database")
	//defer dbPool.Close() //TODO clean shutdown
	return db.New(dbPool)
}
