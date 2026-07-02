package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"media-manager/internal/config"
)

var ErrDevResetNotAllowed = errors.New("dev reset is only allowed when APP_ENV=development and ALLOW_DEV_RESET=true")

const (
	gooseMigrationsDir = "migrations"
	gooseVersionTable  = "app.goose_db_version"
	defaultSeedPath    = "seeds/defaults.sql"
	languageSeedPath   = "seeds/languages.sql"
	devDefaultSeedPath = "seeds/dev.defaults.sql"
	devSeedPath        = "internal/storage/seeds/dev.local.sql"
)

//go:embed migrations/*.sql seeds/defaults.sql seeds/languages.sql seeds/dev.defaults.sql
var storageFS embed.FS

func init() {
	goose.SetBaseFS(storageFS)
	goose.SetLogger(goose.NopLogger())
	goose.SetTableName(gooseVersionTable)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
}

func EnsureSchema(ctx context.Context, databaseURL string) error {
	db, err := openMigrationDB(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := runMigrations(ctx, db); err != nil {
		return err
	}
	return applyDefaultSeed(ctx, db)
}

func ResetDevelopment(ctx context.Context, cfg config.Config) error {
	if !cfg.IsDevelopment() || !cfg.AllowDevReset {
		return ErrDevResetNotAllowed
	}

	db, err := openMigrationDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, `drop schema if exists app cascade`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `create schema app`); err != nil {
		return err
	}
	if err := runMigrations(ctx, db); err != nil {
		return err
	}
	if err := applyDefaultSeed(ctx, db); err != nil {
		return err
	}
	if err := applyDevDefaultSeed(ctx, db); err != nil {
		return err
	}
	return applyDevSeed(ctx, db)
}

func openMigrationDB(ctx context.Context, databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `create schema if not exists app`); err != nil {
		return fmt.Errorf("migration schema setup failed: %w", err)
	}
	if err := goose.UpContext(ctx, db, gooseMigrationsDir); err != nil {
		return fmt.Errorf("goose migration failed: %w", err)
	}
	return nil
}

func applyDefaultSeed(ctx context.Context, db *sql.DB) error {
	if err := applySeed(ctx, db, defaultSeedPath, storageFS.ReadFile); err != nil {
		return err
	}
	return applySeed(ctx, db, languageSeedPath, storageFS.ReadFile)
}

func applyDevDefaultSeed(ctx context.Context, db *sql.DB) error {
	return applySeed(ctx, db, devDefaultSeedPath, storageFS.ReadFile)
}

func applyDevSeed(ctx context.Context, db *sql.DB) error {
	return applySeed(ctx, db, devSeedPath, os.ReadFile)
}

func applySeed(
	ctx context.Context,
	db *sql.DB,
	path string,
	read func(string) ([]byte, error),
) error {
	seedSQL, err := read(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("seed read failed: %w", err)
	}
	if strings.TrimSpace(string(seedSQL)) == "" {
		return nil
	}
	if _, err := db.ExecContext(ctx, string(seedSQL)); err != nil {
		return fmt.Errorf("seed apply failed: %w", err)
	}
	return nil
}
