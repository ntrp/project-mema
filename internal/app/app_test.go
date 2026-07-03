package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSCNSystem006EnsureMediaDataDirCreatesNestedPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "media", "movies")

	if err := ensureMediaDataDir(path); err != nil {
		t.Fatalf("ensureMediaDataDir returned error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected media data directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("%s is not a directory", path)
	}
}

func TestSCNSystem006EnsureMediaDataDirReportsFilesystemFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "media")
	if err := os.WriteFile(path, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write fixture file: %v", err)
	}

	err := ensureMediaDataDir(filepath.Join(path, "movies"))
	if err == nil {
		t.Fatal("expected setup error")
	}
	if !strings.Contains(err.Error(), "media data directory setup failed") {
		t.Fatalf("error = %q", err.Error())
	}
}
