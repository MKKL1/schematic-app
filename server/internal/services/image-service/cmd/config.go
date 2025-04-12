package main

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/MKKL1/schematic-app/server/internal/pkg/kafka"
	"github.com/MKKL1/schematic-app/server/internal/pkg/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/imgproxy"
	"time"
)

// ApplicationConfig defines the structure for your configuration file
type ApplicationConfig struct {
	Server struct {
		Grpc grpc.Config `koanf:"grpc"`
	} `koanf:"server"`
	Database postgres.Config   `koanf:"database"`
	Kafka    kafka.KafkaConfig `koanf:"kafka"`
	ImgProxy imgproxy.Config   `koanf:"imgproxy"`
	Service  struct {
		TmpExpire time.Duration `koanf:"tmp_expire"`
	} `koanf:"service"`
}
