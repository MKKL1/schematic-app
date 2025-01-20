package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresConfig struct {
	Port     string `mapstructure:"port" validate:"required"`
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func NewPostgreSQLClient(ctx context.Context, conf *PostgresConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	dbPool, err := pgxpool.New(ctx, connString)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Str("port", conf.Port).Msg("shutting down postgres pool")
				dbPool.Close()
				log.Info().Msg("postgres pool shut down")
				return
			}
		}
	}()

	return dbPool, err
}
