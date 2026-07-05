package storage

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

func TestMediaFilesInRootFindsVideoFiles(t *testing.T) {
	root := t.TempDir()
	video := filepath.Join(root, "Movie.2026.mkv")
	nested := filepath.Join(root, "Extras", "Sample.mp4")
	sidecar := filepath.Join(root, "Movie.2026.nfo")
	for _, path := range []string{video, nested, sidecar} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create folder: %v", err)
		}
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	files, err := mediaFilesInRoot(root)
	if err != nil {
		t.Fatalf("scan media root: %v", err)
	}
	expected := []string{nested, video}
	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("expected media files %v, got %v", expected, files)
	}
}

func TestCollectMetadataFilePathsFindsRelatedSidecars(t *testing.T) {
	root := t.TempDir()
	video := filepath.Join(root, "Movie.2026.mkv")
	metadata := filepath.Join(root, "Movie.2026.nfo")
	unrelated := filepath.Join(root, "Other.nfo")
	for _, path := range []string{video, metadata, unrelated} {
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	files := collectMetadataFilePaths([]string{video})
	expected := []string{metadata}
	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("expected metadata files %v, got %v", expected, files)
	}
}

func TestRescanMediaItemFilesMarksMissingAndRestoredFiles(t *testing.T) {
	ctx, store := testDBStore(t)
	item := rescanMediaItem(t, ctx, store)
	filePath := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.mkv")
	writeRescanFile(t, filePath)
	if err := store.RecordImportedMediaFile(ctx, item, filePath); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filePath); err != nil {
		t.Fatal(err)
	}

	updated, err := store.RescanMediaItemFiles(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.FilePaths) != 1 || updated.FilePaths[0] != filePath {
		t.Fatalf("missing file path should remain visible, got %#v", updated.FilePaths)
	}
	requireLibraryScanItemStatus(t, ctx, store, item.ID, filePath, "missing")
	requireHistoryOperation(t, ctx, store, item.ID, "missing")

	writeRescanFile(t, filePath)
	if _, err := store.RescanMediaItemFiles(ctx, item.ID); err != nil {
		t.Fatal(err)
	}
	requireLibraryScanItemStatus(t, ctx, store, item.ID, filePath, "restored")
	requireHistoryOperation(t, ctx, store, item.ID, "restored")
}

func TestRescanMediaItemFilesCreatesMovedCandidate(t *testing.T) {
	ctx, store := testDBStore(t)
	item := rescanMediaItem(t, ctx, store)
	oldPath := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.mkv")
	newPath := filepath.Join(*item.MediaFolderPath, "Subfolder", "Scenario.Movie.2026.mkv")
	writeRescanFile(t, oldPath)
	if err := store.RecordImportedMediaFile(ctx, item, oldPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(oldPath); err != nil {
		t.Fatal(err)
	}
	writeRescanFile(t, newPath)

	if _, err := store.RescanMediaItemFiles(ctx, item.ID); err != nil {
		t.Fatal(err)
	}

	requireLibraryScanItemStatus(t, ctx, store, item.ID, oldPath, "missing")
	requireLibraryScanItemStatus(t, ctx, store, item.ID, newPath, "moved_candidate")
	requireHistoryOperation(t, ctx, store, item.ID, "moved_candidate")
}

func rescanMediaItem(t *testing.T, ctx context.Context, store *SettingsStore) MediaItem {
	t.Helper()
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Scenario Movie",
		Year:            int32Ptr(2026),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func requireLibraryScanItemStatus(t *testing.T, ctx context.Context, store *SettingsStore, itemID uuid.UUID, path string, status string) {
	t.Helper()
	rows, err := storagegen.New(store.pool).ListMediaFileRecordsForItem(ctx, &itemID)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if row.Path == path && row.Status == status {
			return
		}
	}
	t.Fatalf("status %q for %s not found in %#v", status, path, rows)
}

func requireHistoryOperation(t *testing.T, ctx context.Context, store *SettingsStore, itemID uuid.UUID, operation string) {
	t.Helper()
	entries, err := store.ListMediaFileHistory(ctx, itemID)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Operation == operation {
			return
		}
	}
	t.Fatalf("history operation %q not found in %#v", operation, entries)
}

func writeRescanFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("movie"), 0o644); err != nil {
		t.Fatal(err)
	}
}
