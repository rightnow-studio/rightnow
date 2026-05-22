package config

import (
	"os"
)

type Config struct {
	ServerAddr   string
	DatabaseURL  string
	JWTSecret    string
	RefreshSecret string
}

func Load() *Config {
	return &Config{
		ServerAddr:    getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "cike.db"),
		JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
		RefreshSecret: getEnv("REFRESH_SECRET", "change-me-refresh"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
