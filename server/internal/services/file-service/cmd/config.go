package main

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/pkg/http"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/minio"
	"github.com/MKKL1/schematic-app/server/internal/pkg/postgres"
	"time"
)

// ApplicationConfig defines the structure for your configuration file
type ApplicationConfig struct {
	Server struct {
		Grpc grpc.Config `koanf:"grpc"`
		Http http.Config `koanf:"http"`
	} `koanf:"server"`
	Database postgres.Config   `koanf:"database"`
	Minio    minio.Config      `koanf:"minio"`
	Kafka    kafka.KafkaConfig `koanf:"kafka"`
	Service  struct {
		TmpExpire time.Duration `koanf:"tmp_expire"`
	} `koanf:"service"`
}
