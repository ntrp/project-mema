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
	writeSparseImportFile(t, sourcePath, 800*1024*1024)
	extraPath := filepath.Join(appPath, "complete", "Toy.Story.5.2026.sample.mkv")
	writeSparseImportFile(t, extraPath, 2*1024*1024)

	folder, err := store.CreateLibraryFolder(ctx, libraryPath, "movie")
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
	history, err := store.ListMediaFileHistory(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].Operation != "imported" || history[0].Status != "succeeded" {
		t.Fatalf("history = %#v", history)
	}
	if history[0].SourcePath == nil || *history[0].SourcePath != sourcePath {
		t.Fatalf("history source = %#v", history[0].SourcePath)
	}
	if history[0].DestinationPath == nil || *history[0].DestinationPath != imported.FilePaths[0] {
		t.Fatalf("history destination = %#v", history[0].DestinationPath)
	}
	if _, err := store.DeleteMediaItemFile(ctx, item.ID, imported.FilePaths[0]); err != nil {
		t.Fatal(err)
	}
	history, err = store.ListMediaFileHistory(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 2 || history[0].Operation != "deleted" || history[0].Status != "succeeded" {
		t.Fatalf("delete history = %#v", history)
	}
}

func writeSparseImportFile(t *testing.T, path string, size int64) {
	t.Helper()
	if err := os.WriteFile(path, []byte{0}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Truncate(path, size); err != nil {
		t.Fatal(err)
	}
}
