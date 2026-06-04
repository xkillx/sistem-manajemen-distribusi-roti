package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	JWTSecret string

	OwnerUsername string
	OwnerPassword string

	AppURL string
	port   string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "smdr"),
		DBUser:     getEnv("DB_USER", "smdr"),
		DBPassword: getEnv("DB_PASSWORD", "smdr"),

		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),

		OwnerUsername: getEnv("OWNER_USERNAME", "owner"),
		OwnerPassword: getEnv("OWNER_PASSWORD", "owner"),

		AppURL: getEnv("APP_URL", "http://localhost:8080"),
		port:   getEnv("PORT", "8080"),
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func (c *Config) Port() string {
	return c.port
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
