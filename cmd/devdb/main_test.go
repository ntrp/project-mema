package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"os"
	"path/filepath"
	"testing"

	"media-manager/internal/storage"
)

type recordingExecer struct {
	query string
	args  []any
}

func (db *recordingExecer) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	db.query, db.args = query, args
	return driver.RowsAffected(1), nil
}

func TestResetSeedsConfiguredDefaultAdmin(t *testing.T) {
	t.Setenv("ADMIN_USERNAME", "local-admin")
	t.Setenv("ADMIN_PASSWORD", "local-password")
	db := &recordingExecer{}
	if err := ensureDefaultAdminUser(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	if len(db.args) != 3 || db.args[1] != "local-admin" {
		t.Fatalf("admin seed args = %#v", db.args)
	}
	hash, ok := db.args[2].(string)
	if !ok || !storage.VerifyPassword("local-password", hash) {
		t.Fatal("admin password was not seeded with the configured value")
	}
}

func TestOptionalSeedOnlyIgnoresMissingFile(t *testing.T) {
	db, err := sql.Open("pgx", defaultDatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	if err := applyOptionalSeed(context.Background(), db, filepath.Join(t.TempDir(), "missing.sql")); err != nil {
		t.Fatalf("missing optional seed: %v", err)
	}

	path := filepath.Join(t.TempDir(), "dev.local.sql")
	if err := os.WriteFile(path, []byte("select 1;"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := applyOptionalSeed(context.Background(), db, path); err == nil {
		t.Fatal("existing local seed failure must stop the reset")
	}
}
