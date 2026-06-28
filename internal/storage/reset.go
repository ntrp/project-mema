package storage

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/config"
)

var ErrDevResetNotAllowed = errors.New("dev reset is only allowed when APP_ENV=development and ALLOW_DEV_RESET=true")

//go:embed schema.sql
var schemaSQL string

func EnsureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, schemaSQL)
	return err
}

func ResetDevelopment(ctx context.Context, cfg config.Config) error {
	if !cfg.IsDevelopment() || !cfg.AllowDevReset {
		return ErrDevResetNotAllowed
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, `drop schema if exists app cascade`); err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, `create schema app`); err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		return err
	}
	return nil
}
