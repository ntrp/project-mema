package imports

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/downloadclients"
	"media-manager/internal/storage"
	"media-manager/internal/testdb"
)

func TestImportCompletedDownloadLinksAndRecordsMediaFile(t *testing.T) {
	databaseURL := testdb.Create(t)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	store := storage.NewSettingsStore(pool)
	if err := storage.EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	if err := store.EnsureDefaultFileNamingSettings(ctx); err != nil {
		t.Fatal(err)
	}

	root := t.TempDir()
	clientPath := "/client-" + uuid.NewString()
	appPath := filepath.Join(root, "client")
	libraryPath := filepath.Join(root, "library")
	if err := os.MkdirAll(filepath.Join(appPath, "complete"), 0o755); err != nil {
		t.Fatal(err)
	}
	sourcePath := filepath.Join(appPath, "complete", "Toy.Story.5.2026.1080p.mkv")
	if err := os.WriteFile(sourcePath, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}
	extraPath := filepath.Join(appPath, "complete", "Toy.Story.5.2026.sample.mkv")
	if err := os.WriteFile(extraPath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	folder, err := store.CreateLibraryFolder(ctx, libraryPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = store.DeleteLibraryFolder(context.Background(), folder.ID)
	})
	mapping, err := store.CreatePathMapping(ctx, storage.PathMappingInput{
		ClientPath: clientPath,
		AppPath:    appPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = store.DeletePathMapping(context.Background(), mapping.ID)
	})

	year := int32(2026)
	item, err := store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type:            "movie",
		Title:           "Toy Story 5",
		Year:            &year,
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = store.DeleteMediaItem(context.Background(), item.ID, false)
	})

	activity := storage.DownloadActivity{
		ID:          uuid.New(),
		MediaItemID: item.ID,
	}
	err = NewService(store).ImportCompletedDownload(ctx, activity, []downloadclients.StatusFile{
		{
			Path:     clientPath + "/complete/Toy.Story.5.2026.1080p.mkv",
			Complete: true,
		},
		{
			Path:     clientPath + "/complete/Toy.Story.5.2026.sample.mkv",
			Complete: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	imported, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imported.Status != "downloaded" {
		t.Fatalf("status = %q", imported.Status)
	}
	if len(imported.FilePaths) != 1 {
		t.Fatalf("file paths = %#v", imported.FilePaths)
	}
	if filepath.Base(imported.FilePaths[0]) != "Toy.Story.5.2026.1080p.mkv" {
		t.Fatalf("file path = %q", imported.FilePaths[0])
	}
	targetInfo, err := os.Stat(imported.FilePaths[0])
	if err != nil {
		t.Fatal(err)
	}
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(sourceInfo, targetInfo) {
		t.Fatalf("imported file is not linked to source")
	}
	if _, err := os.Stat(filepath.Join(libraryPath, filepath.Base(extraPath))); !os.IsNotExist(err) {
		t.Fatalf("extra target stat err = %v, want missing", err)
	}
}
