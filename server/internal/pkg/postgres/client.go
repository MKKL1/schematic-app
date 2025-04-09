package postgres

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/closer" // Import retry package
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Config struct {
	Port     string        `koanf:"port"`
	Host     string        `koanf:"host"`
	Username string        `koanf:"username"`
	Password string        `koanf:"password"`
	Database string        `koanf:"database"`
	Retry    closer.Config `koanf:"retry"`
}

// NewClient creates a new PostgreSQL client pool with retry logic.
func NewClient(ctx context.Context, conf Config, logger zerolog.Logger) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=10", // Example: add pool settings
		conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)

	var dbPool *pgxpool.Pool
	var err error

	op := func(ctx context.Context) error {
		dbPool, err = pgxpool.New(ctx, connString)
		if err != nil {
			return fmt.Errorf("pgxpool.New failed: %w", err)
		}
		err = dbPool.Ping(ctx)
		if err != nil {
			dbPool.Close()
			return fmt.Errorf("dbPool.Ping failed: %w", err)
		}
		return nil
	}

	if conf.Retry.Attempts == 0 {
		conf.Retry = closer.DefaultConfig
		logger.Warn().Msg("Postgres retry config not found, using defaults")
	}

	logger = logger.With().Str("component", "postgres-client").Logger()
	err = closer.Do(ctx, logger, conf.Retry, op, "Connect to PostgreSQL")
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres after retries failed: %w", err)
	}

	logger.Info().Str("host", conf.Host).Str("db", conf.Database).Msg("Successfully connected to PostgreSQL")
	return dbPool, nil
}

func AddToCloser(c *closer.Closer, pool *pgxpool.Pool, logger zerolog.Logger) {
	c.Add(func(_ context.Context) error {
		logger.Info().Msg("Closing Postgres connection...")
		pool.Close()
		logger.Info().Msg("Closed Postgres connection")
		return nil
	})
}
