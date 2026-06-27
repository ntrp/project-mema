package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr          string
	AppEnv        string
	AllowDevReset bool
	DatabaseURL   string
	WebDir        string
	Version       string
	Commit        string
	AdminUsername string
	AdminPassword string
	SessionTTL    time.Duration
}

func Load() Config {
	return Config{
		Addr:          envString("ADDR", ":8080"),
		AppEnv:        envString("APP_ENV", "development"),
		AllowDevReset: envBool("ALLOW_DEV_RESET", false),
		DatabaseURL:   envString("DATABASE_URL", "postgres://media_manager:media_manager@localhost:5432/media_manager?sslmode=disable"),
		WebDir:        envString("WEB_DIR", "web/build"),
		Version:       envString("APP_VERSION", "0.0.0-dev"),
		Commit:        envString("APP_COMMIT", "dev"),
		AdminUsername: envString("ADMIN_USERNAME", "admin"),
		AdminPassword: envString("ADMIN_PASSWORD", "admin"),
		SessionTTL:    envDuration("SESSION_TTL", 24*time.Hour),
	}
}

func (c Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func envString(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func envBool(name string, fallback bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(name string, fallback time.Duration) time.Duration {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
