package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveMediaFolderDeletesMediaRoot(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Movie",
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	mediaPath := *item.MediaFolderPath
	nestedFile := filepath.Join(mediaPath, "extras", "sample.txt")
	if err := os.MkdirAll(filepath.Dir(nestedFile), 0o755); err != nil {
		t.Fatalf("create nested folder: %v", err)
	}
	if err := os.WriteFile(nestedFile, []byte("sample"), 0o644); err != nil {
		t.Fatalf("write nested file: %v", err)
	}

	if err := store.removeMediaFolder(ctx, item); err != nil {
		t.Fatalf("remove media folder: %v", err)
	}
	if _, err := os.Stat(mediaPath); !os.IsNotExist(err) {
		t.Fatalf("expected media folder to be deleted, got err=%v", err)
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("expected library root to remain: %v", err)
	}
}

func TestMediaFolderDeletePathRejectsLibraryRoot(t *testing.T) {
	root := t.TempDir()

	_, _, err := mediaFolderDeletePath(MediaItem{
		LibraryFolderPath: &root,
		MediaFolderPath:   &root,
	})
	if err == nil {
		t.Fatal("expected library root to be rejected")
	}
}

func TestMediaFolderDeletePathRejectsOutsideLibraryRoot(t *testing.T) {
	root := t.TempDir()
	other := filepath.Join(t.TempDir(), "Movie")

	_, _, err := mediaFolderDeletePath(MediaItem{
		LibraryFolderPath: &root,
		MediaFolderPath:   &other,
	})
	if err == nil {
		t.Fatal("expected outside media folder to be rejected")
	}
}

func TestRemoveMediaFolderSkipsMissingPath(t *testing.T) {
	ctx, store := testDBStore(t)
	if err := store.removeMediaFolder(ctx, MediaItem{}); err != nil {
		t.Fatalf("remove media folder without path: %v", err)
	}
}

func deletePolicyItem(t *testing.T, ctx context.Context, store *SettingsStore) (MediaItem, string) {
	t.Helper()
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Policy Movie",
		Year:            int32Ptr(2026),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(*item.MediaFolderPath, "Policy.Movie.mkv")
	writeTestFile(t, path)
	if err := store.RecordImportedMediaFile(ctx, item, path); err != nil {
		t.Fatal(err)
	}
	return item, path
}
