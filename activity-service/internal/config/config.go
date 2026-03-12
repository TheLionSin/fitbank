package config

import (
	"os"
)

type Config struct {
	AppPort     string
	DBURL       string
	Environment string // "dev" или "prod"
}

func New() *Config {
	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DBURL:       getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/fitbank?sslmode=disable"),
		Environment: getEnv("APP_ENV", "dev"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
