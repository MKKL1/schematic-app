package minio

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/closer" // Import closer package
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
)

type Config struct {
	Endpoint  string        `koanf:"endpoint"`
	AccessKey string        `koanf:"access_key"`
	SecretKey string        `koanf:"secret_key"`
	UseSSL    bool          `koanf:"use_ssl"`
	Retry     closer.Config `koanf:"closer"`  // Add closer config
	Buckets   BucketConfig  `koanf:"buckets"` // Keep bucket names here
}

// BucketConfig holds Minio bucket names (could be nested in Minio Config)
type BucketConfig struct {
	Files string `koanf:"files"`
	Temp  string `koanf:"temp"`
}

// NewClient creates a new Minio client with closer logic.
func NewClient(ctx context.Context, conf Config, logger zerolog.Logger) (*minio.Client, error) {
	var minioClient *minio.Client
	var err error

	// Use the closer mechanism
	op := func(ctx context.Context) error {
		minioClient, err = minio.New(conf.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
			Secure: conf.UseSSL,
		})
		if err != nil {
			return fmt.Errorf("minio.New failed: %w", err)
		}

		// Ping check (ListBuckets is a common way to check connectivity)
		_, err = minioClient.ListBuckets(ctx) // Use the passed context
		if err != nil {
			return fmt.Errorf("minio connection check (ListBuckets) failed: %w", err)
		}
		return nil // Success
	}

	// Apply default closer config if not specified
	if conf.Retry.Attempts == 0 {
		conf.Retry = closer.DefaultConfig
		logger.Warn().Msg("Minio closer config not found, using defaults")
	}

	logger = logger.With().Str("component", "minio-client").Logger()
	err = closer.Do(ctx, logger, conf.Retry, op, "Connect to Minio")
	if err != nil {
		return nil, fmt.Errorf("connecting to minio after retries failed: %w", err)
	}

	logger.Info().Str("endpoint", conf.Endpoint).Msg("Successfully connected to Minio")
	// Note: Minio client doesn't have an explicit Close() method.
	// Connections are typically managed internally. No cleanup needed in Closer.
	return minioClient, nil
}
