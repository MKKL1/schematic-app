package main

import "github.com/MKKL1/schematic-app/server/internal/pkg/server"

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

type ApplicationConfig struct {
	Minio    MinioConfig           `koanf:"minio"`
	Database server.PostgresConfig `koanf:"database"`
	Kafka    KafkaConfig           `koanf:"kafka"`
	LogLevel string                `koanf:"log_level"`
}
