package main

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"time"
)

type BucketsConfig struct {
	Files string `koanf:"files"`
	Temp  string `koanf:"temp"`
}

type MinioConfig struct {
	Endpoint  string        `koanf:"endpoint"`
	AccessKey string        `koanf:"access_key"`
	SecretKey string        `koanf:"secret_key"`
	UseSSL    bool          `koanf:"use_ssl"`
	Buckets   BucketsConfig `koanf:"buckets"`
}

type KafkaConfig struct {
	Brokers []string `koanf:"brokers"`
}

type UploadConfig struct {
	TmpExpire time.Duration `koanf:"expire_duration"`
}

type ServerConfig struct {
	Grpc GrpcConfig `koanf:"grpc"`
	Http HttpConfig `koanf:"http"`
}
type GrpcConfig struct {
	Host string `koanf:"host"`
}

type HttpConfig struct {
	Host string `koanf:"host"`
}

type ApplicationConfig struct {
	Minio    MinioConfig           `koanf:"minio"`
	Database server.PostgresConfig `koanf:"database"`
	Kafka    KafkaConfig           `koanf:"kafka"`
	LogLevel string                `koanf:"log_level"`
	Upload   UploadConfig          `koanf:"upload"`
	Server   ServerConfig          `koanf:"server"`
}
