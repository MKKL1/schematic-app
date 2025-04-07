package config

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"strings"
)

func LoadConfig[T any](configFilePath string) (*T, error) {
	k := koanf.New(".")

	// Load YAML configuration from file.
	if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("loading yaml config file: %w", err)
	}

	// Optionally load environment variables with a prefix.
	// This example uses "MYAPP_" as the prefix.
	err := k.Load(env.Provider("MYAPP_", ".", func(s string) string {
		// Convert environment variable names to lower-case dot-delimited keys.
		// For example, MYAPP_MINIO_ENDPOINT becomes "minio.endpoint"
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "MYAPP_")), "_", ".", -1)
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("loading environment variable MYAPP_: %w", err)
	}

	// Unmarshal the configuration into our Config struct.
	var cfg T
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}
