package imports

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/storage"
	"media-manager/internal/testdb"
)

func TestAutomaticImportSkipsManualLibraryFolder(t *testing.T) {
	ctx := context.Background()
	databaseURL := testdb.Create(t)
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := storage.EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	store := storage.NewSettingsStore(pool)
	folder, err := store.CreateLibraryFolderWithManual(ctx, t.TempDir(), "movie", true)
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type: "movie", Title: "Manual Movie", Monitored: true, LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = NewService(store).ImportCompletedDownload(ctx, storage.DownloadActivity{
		ID: uuid.New(), MediaItemID: item.ID,
	}, nil)
	if err == nil || !strings.Contains(err.Error(), "only allows manual operations") {
		t.Fatalf("automatic import error = %v", err)
	}
}
