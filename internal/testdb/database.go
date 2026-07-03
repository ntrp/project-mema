package testdb

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/url"
	"os"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Create returns a fresh PostgreSQL database URL derived from DATABASE_URL.
// The database is dropped during test cleanup, so integration tests never
// mutate the developer's configured database.
func Create(t testing.TB) string {
	t.Helper()
	baseURL := os.Getenv("DATABASE_URL")
	if baseURL == "" {
		t.Skip("DATABASE_URL is required for database integration tests")
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		t.Fatalf("parse DATABASE_URL: %v", err)
	}

	name := "project_mema_test_" + randomSuffix(t)
	ctx := context.Background()
	adminURL := withDatabase(parsed, "postgres")
	db, err := sql.Open("pgx", adminURL)
	if err != nil {
		t.Fatalf("open postgres admin database: %v", err)
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Fatalf("connect postgres admin database: %v", err)
	}
	if _, err := db.ExecContext(ctx, `create database `+quoteIdent(name)); err != nil {
		db.Close()
		t.Fatalf("create test database %s: %v", name, err)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), `drop database if exists `+quoteIdent(name)+` with (force)`)
		_ = db.Close()
	})

	return withDatabase(parsed, name)
}

func randomSuffix(t testing.TB) string {
	t.Helper()
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		t.Fatalf("read random database suffix: %v", err)
	}
	return hex.EncodeToString(bytes[:])
}

func withDatabase(parsed *url.URL, database string) string {
	copy := *parsed
	copy.Path = "/" + database
	return copy.String()
}

func quoteIdent(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
