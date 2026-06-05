package config

import (
	"os"
)

type Config struct {
	ServerAddr    string
	DatabaseURL   string
	JWTSecret     string
	RefreshSecret string
	AdminSecret   string
	AdminUsername string
	AdminPassword string
}

func Load() *Config {
	return &Config{
		ServerAddr:    getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "cike.db"),
		JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
		RefreshSecret: getEnv("REFRESH_SECRET", "change-me-refresh"),
		AdminSecret:   getEnv("ADMIN_SECRET", "change-me-admin"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
