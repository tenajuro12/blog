package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "host=localhost user=postgres password=murderpe dbname=blogdb port=5432 sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
