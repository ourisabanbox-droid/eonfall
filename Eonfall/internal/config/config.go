package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort    string
	DatabaseURL string
	RedisURL    string
	TickRateMs  int
}

func Load() (Config, error) {
	tickRate := 1000
	if raw := os.Getenv("TICK_RATE_MS"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("invalid TICK_RATE_MS: %w", err)
		}
		tickRate = v
	}

	cfg := Config{
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		TickRateMs:  tickRate,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
