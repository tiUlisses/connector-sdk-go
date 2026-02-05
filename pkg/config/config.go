package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains runtime settings sourced from environment variables.
type Config struct {
	LogLevel        string
	KafkaBrokers    []string
	KafkaTopic      string
	HTTPAddr        string
	KafkaMaxRetries int
	KafkaBackoff    time.Duration
}

// FromEnv builds a Config from process environment.
func FromEnv() (Config, error) {
	cfg := Config{
		LogLevel:        envOrDefault("CONNECTOR_LOG_LEVEL", "INFO"),
		KafkaBrokers:    splitCSV(envOrDefault("CONNECTOR_KAFKA_BROKERS", "localhost:9092")),
		KafkaTopic:      envOrDefault("CONNECTOR_KAFKA_TOPIC", "connector.events"),
		HTTPAddr:        envOrDefault("CONNECTOR_HTTP_ADDR", ":8080"),
		KafkaMaxRetries: envOrDefaultInt("CONNECTOR_KAFKA_MAX_RETRIES", 3),
		KafkaBackoff:    envOrDefaultDuration("CONNECTOR_KAFKA_BACKOFF", 500*time.Millisecond),
	}

	if len(cfg.KafkaBrokers) == 0 {
		return Config{}, fmt.Errorf("CONNECTOR_KAFKA_BROKERS must not be empty")
	}
	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		parsed, err := strconv.Atoi(v)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func envOrDefaultDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		parsed, err := time.ParseDuration(v)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

func splitCSV(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
