package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port:        envOrDefault("PORT", "8080"),
		DatabaseURL: envOrDefault("DATABASE_URL", "postgres://inventory_user:inventory_password@localhost:5432/inventory_reservation?sslmode=disable"),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
