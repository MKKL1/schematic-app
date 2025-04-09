package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ClientConfig struct {
	Address           string        `koanf:"address"`
	Backoff           BackoffConfig `koanf:"backoff"`
	MinConnectTimeout time.Duration `koanf:"connect_timeout"`
}

type BackoffConfig struct {
	// BaseDelay is the amount of time to backoff after the first failure.
	BaseDelay time.Duration `koanf:"base_delay"`
	// Multiplier is the factor with which to multiply backoffs after a
	// failed retry. Should ideally be greater than 1.
	Multiplier float64 `koanf:"multiplier"`
	// Jitter is the factor with which backoffs are randomized.
	Jitter float64 `koanf:"jitter"`
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration `koanf:"max_delay"`
}

func NewClient(config ClientConfig) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  config.Backoff.BaseDelay,
				Multiplier: config.Backoff.Multiplier,
				Jitter:     config.Backoff.Jitter,
				MaxDelay:   config.Backoff.MaxDelay,
			},
			MinConnectTimeout: config.MinConnectTimeout,
		}),
	)

	return conn, err
}
