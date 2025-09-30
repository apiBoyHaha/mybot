package config

import (
	"os"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "user:pass@tcp(localhost:3306)/dbname"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
