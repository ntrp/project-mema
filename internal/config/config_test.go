package config

import (
	"testing"
	"time"
)

func TestScenarioSCNSystem001LoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("ADDR", "127.0.0.1:19090")
	t.Setenv("APP_ENV", "production")
	t.Setenv("ALLOW_DEV_RESET", "true")
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("MEDIA_DATA_DIR", "/tmp/media")
	t.Setenv("WEB_DIR", "/tmp/web")
	t.Setenv("APP_VERSION", "1.2.3")
	t.Setenv("APP_COMMIT", "abc123")
	t.Setenv("APP_SOURCE_URL", "https://example.invalid/source")
	t.Setenv("ADMIN_USERNAME", "root")
	t.Setenv("ADMIN_PASSWORD", "secret")
	t.Setenv("SESSION_TTL", "2h")

	cfg := Load()

	if cfg.Addr != "127.0.0.1:19090" || cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("config did not use env overrides: %#v", cfg)
	}
	if !cfg.AllowDevReset || cfg.IsDevelopment() {
		t.Fatalf("environment flags not applied: %#v", cfg)
	}
	if cfg.SessionTTL != 2*time.Hour {
		t.Fatalf("SessionTTL = %s", cfg.SessionTTL)
	}
	if cfg.AdminUsername != "root" || cfg.AdminPassword != "secret" {
		t.Fatalf("admin credentials not loaded: %#v", cfg)
	}
}

func TestScenarioSCNSystem001LoadFallsBackForInvalidOptionalValues(t *testing.T) {
	t.Setenv("ALLOW_DEV_RESET", "not-bool")
	t.Setenv("SESSION_TTL", "not-duration")

	cfg := Load()

	if cfg.AllowDevReset {
		t.Fatal("invalid ALLOW_DEV_RESET should fall back to false")
	}
	if cfg.SessionTTL != 24*time.Hour {
		t.Fatalf("SessionTTL = %s", cfg.SessionTTL)
	}
}
