package storage

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateUser = errors.New("user already exists")
	ErrLastAdmin     = errors.New("at least one admin user is required")
	ErrNotFound      = errors.New("resource not found")
	ErrRequestClosed = errors.New("media request is not pending")
)

type SettingsStore struct {
	pool *pgxpool.Pool
}

func NewSettingsStore(pool *pgxpool.Pool) *SettingsStore {
	return &SettingsStore{pool: pool}
}
