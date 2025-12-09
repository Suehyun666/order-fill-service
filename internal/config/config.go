package config

import (
	"os"
	"strings"
)

type Config struct {
	KafkaBrokers []string
	LogLevel     string
}

func Load() *Config {
	return &Config{
		KafkaBrokers: getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
		LogLevel:     getEnv("LOG_LEVEL", "INFO"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsSlice(key string, fallback []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return fallback
}
