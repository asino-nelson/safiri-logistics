package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort          = "8080"
	defaultTokenLifetime = 24 * time.Hour
)

type Config struct {
	Port             string
	DatabaseURL      string
	JWTSecret        string
	JWTTokenLifetime time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:             getEnv("PORT", defaultPort),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTTokenLifetime: defaultTokenLifetime,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	if value := os.Getenv("JWT_TOKEN_TTL_HOURS"); value != "" {
		hours, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, fmt.Errorf("parse JWT_TOKEN_TTL_HOURS: %w", err)
		}

		cfg.JWTTokenLifetime = time.Duration(hours) * time.Hour
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
