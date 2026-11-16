package config

import (
	"fmt"
	"os"
)

type Config struct {
	GRPCAddress string
	SQLitePath  string
	LogLevel    string
	LogFormat   string
}

func Load() (Config, error) {
	cfg := Config{
		GRPCAddress: getEnv("GLACIER_GRPC_ADDRESS", ":50051"),
		SQLitePath:  getEnv("GLACIER_SQLITE_PATH", "glacier.db"),
		LogLevel:    getEnv("GLACIER_LOG_LEVEL", "info"),
		LogFormat:   getEnv("GLACIER_LOG_FORMAT", "text"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid log level %q", c.LogLevel)
	}

	switch c.LogFormat {
	case "text", "json":
	default:
		return fmt.Errorf("invalid log format %q", c.LogFormat)
	}

	if c.GRPCAddress == "" {
		return fmt.Errorf("gRPC address cannot be empty")
	}

	if c.SQLitePath == "" {
		return fmt.Errorf("SQLite path cannot be empty")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
