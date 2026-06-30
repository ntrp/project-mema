package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveMediaFolderDeletesMediaRoot(t *testing.T) {
	root := t.TempDir()
	mediaPath := filepath.Join(root, "Movie")
	nestedFile := filepath.Join(mediaPath, "extras", "sample.txt")
	if err := os.MkdirAll(filepath.Dir(nestedFile), 0o755); err != nil {
		t.Fatalf("create nested folder: %v", err)
	}
	if err := os.WriteFile(nestedFile, []byte("sample"), 0o644); err != nil {
		t.Fatalf("write nested file: %v", err)
	}

	if err := removeMediaFolder(MediaItem{
		LibraryFolderPath: &root,
		MediaFolderPath:   &mediaPath,
	}); err != nil {
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
	if err := removeMediaFolder(MediaItem{}); err != nil {
		t.Fatalf("remove media folder without path: %v", err)
	}
}
