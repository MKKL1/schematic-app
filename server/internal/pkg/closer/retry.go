package closer

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// Config holds retry configuration.
type Config struct {
	Attempts    int           `koanf:"attempts"`
	InitialWait time.Duration `koanf:"initial_wait"`
	MaxWait     time.Duration `koanf:"max_wait"` // Optional: Cap the wait time
}

// DefaultConfig provides sensible defaults.
var DefaultConfig = Config{
	Attempts:    5,
	InitialWait: 1 * time.Second,
	MaxWait:     30 * time.Second,
}

// Func is a function that can be retried.
type Func func(ctx context.Context) error

// Do performs the operation with retries based on the config.
func Do(ctx context.Context, logger zerolog.Logger, cfg Config, op Func, description string) error {
	var err error
	wait := cfg.InitialWait

	for i := 0; i < cfg.Attempts; i++ {
		logger.Debug().Int("attempt", i+1).Str("operation", description).Msg("Attempting operation")
		err = op(ctx)
		if err == nil {
			logger.Info().Str("operation", description).Msg("Operation successful")
			return nil // Success
		}

		logger.Warn().Err(err).Int("attempt", i+1).Int("max_attempts", cfg.Attempts).Str("operation", description).Msg("Operation failed, will retry")

		// Check if context was cancelled
		if ctx.Err() != nil {
			logger.Error().Err(ctx.Err()).Str("operation", description).Msg("Context cancelled during retry, stopping")
			return ctx.Err()
		}

		// Don't wait after the last attempt
		if i == cfg.Attempts-1 {
			break
		}

		// Calculate next wait time (exponential backoff)
		sleep := wait
		wait *= 2
		if cfg.MaxWait > 0 && wait > cfg.MaxWait {
			wait = cfg.MaxWait // Cap the wait time
		}

		logger.Info().Dur("wait_duration", sleep).Str("operation", description).Msg("Waiting before next retry")
		select {
		case <-time.After(sleep):
			// Continue loop
		case <-ctx.Done():
			logger.Error().Err(ctx.Err()).Str("operation", description).Msg("Context cancelled during wait, stopping retries")
			return ctx.Err() // Context cancelled during wait
		}
	}

	logger.Error().Err(err).Int("attempts", cfg.Attempts).Str("operation", description).Msg("Operation failed after all attempts")
	return err // Return the last error
}
