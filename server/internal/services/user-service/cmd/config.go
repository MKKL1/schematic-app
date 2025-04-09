package main

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/pkg/postgres"
)

// ApplicationConfig defines the structure for your configuration file
type ApplicationConfig struct {
	Server struct {
		Grpc grpc.Config `koanf:"grpc"`
	} `koanf:"server"`
	Database postgres.Config `koanf:"database"`
}
