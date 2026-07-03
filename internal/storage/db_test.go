package storage

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/acceptance"
)

func testDBStore(t *testing.T) (context.Context, *SettingsStore) {
	t.Helper()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is required for storage integration tests")
	}
	ctx := context.Background()
	if err := EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return ctx, NewSettingsStore(pool)
}

func stringPtr(value string) *string {
	return &value
}

func int32Ptr(value int32) *int32 {
	return &value
}

func requireStorageScenario(t *testing.T, id string) {
	t.Helper()
	scenario, err := acceptance.RequireScenario("features/behavior", id)
	if err != nil {
		t.Fatal(err)
	}
	if !scenario.HasTag("integration") {
		t.Fatalf("%s missing @integration tag", id)
	}
}
