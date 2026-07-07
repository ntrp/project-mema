package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	defaultDatabaseURL = "postgres://media_manager:media_manager@localhost:15432/media_manager?sslmode=disable"
	gooseMigrationsDir = "internal/storage/migrations"
	gooseVersionTable  = "app.goose_db_version"
	defaultSeedPath    = "internal/storage/seeds/defaults.sql"
	languageSeedPath   = "internal/storage/seeds/languages.sql"
	devDefaultSeedPath = "scripts/seeds/dev.defaults.sql"
	devLocalSeedPath   = "scripts/seeds/dev.local.sql"
	remoteResetEnv     = "ALLOW_REMOTE_DEV_DB_RESET"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := run(ctx, os.Args[1:]); err != nil {
		slog.Error("dev database command failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: go run ./cmd/devdb [clean|reset|seed-local]")
	}
	databaseURL := envString("DATABASE_URL", defaultDatabaseURL)
	switch args[0] {
	case "clean":
		return withDB(ctx, databaseURL, func(db *sql.DB) error {
			if err := ensureLocalReset(databaseURL); err != nil {
				return err
			}
			return cleanSchema(ctx, db)
		})
	case "reset":
		return reset(ctx, databaseURL)
	case "seed-local":
		return withDB(ctx, databaseURL, func(db *sql.DB) error {
			return applyOptionalSeed(ctx, db, devLocalSeedPath)
		})
	default:
		return fmt.Errorf("unknown dev database command %q", args[0])
	}
}

func reset(ctx context.Context, databaseURL string) error {
	if err := ensureLocalReset(databaseURL); err != nil {
		return err
	}
	return withDB(ctx, databaseURL, func(db *sql.DB) error {
		if err := cleanSchema(ctx, db); err != nil {
			return err
		}
		if err := runMigrations(ctx, db); err != nil {
			return err
		}
		for _, path := range []string{defaultSeedPath, languageSeedPath, devDefaultSeedPath} {
			if err := applyRequiredSeed(ctx, db, path); err != nil {
				return err
			}
		}
		return applyOptionalSeed(ctx, db, devLocalSeedPath)
	})
}

func withDB(ctx context.Context, databaseURL string, fn func(*sql.DB) error) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	return fn(db)
}

func cleanSchema(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `drop schema if exists app cascade`); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, `create schema app`)
	return err
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	goose.SetLogger(goose.NopLogger())
	goose.SetTableName(gooseVersionTable)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.UpContext(ctx, db, gooseMigrationsDir); err != nil {
		return fmt.Errorf("goose migration failed: %w", err)
	}
	return nil
}

func applyRequiredSeed(ctx context.Context, db *sql.DB, path string) error {
	return applySeed(ctx, db, path)
}

func applyOptionalSeed(ctx context.Context, db *sql.DB, path string) error {
	err := applySeed(ctx, db, path)
	if err == nil {
		return nil
	}
	fmt.Fprintf(os.Stderr, "warning: skipped optional %s: %v\n", path, err)
	return nil
}

func applySeed(ctx context.Context, db *sql.DB, path string) error {
	seedSQL, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("seed read failed for %s: %w", path, err)
	}
	if strings.TrimSpace(string(seedSQL)) == "" {
		return nil
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("seed transaction failed for %s: %w", path, err)
	}
	if _, err := tx.ExecContext(ctx, string(seedSQL)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("seed apply failed for %s: %w", path, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seed commit failed for %s: %w", path, err)
	}
	return nil
}

func ensureLocalReset(databaseURL string) error {
	if os.Getenv(remoteResetEnv) == "true" {
		return nil
	}
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return fmt.Errorf("database URL parse failed: %w", err)
	}
	switch parsed.Hostname() {
	case "", "localhost", "127.0.0.1", "::1":
		return nil
	default:
		return fmt.Errorf("refusing to reset non-local database %q; set %s=true to override", parsed.Hostname(), remoteResetEnv)
	}
}

func envString(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
