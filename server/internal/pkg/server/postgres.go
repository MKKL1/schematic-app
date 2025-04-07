package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresConfig struct {
	Port     string `koanf:"port"`
	Host     string `koanf:"host"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
}

// TODO no reason to pass by pointer

func NewPostgreSQLClient(ctx context.Context, conf *PostgresConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	return NewPostgreSQLClientByUrl(ctx, connString)
}

func NewPostgreSQLClientByUrl(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(ctx, connString)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("shutting down postgres pool")
				dbPool.Close()
				log.Info().Msg("postgres pool shut down")
				return
			}
		}
	}()

	return dbPool, err
}
