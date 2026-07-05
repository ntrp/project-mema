package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSCNMedia010MediaItemFileTargetValidatesPaths(t *testing.T) {
	root := t.TempDir()
	item := MediaItem{MediaFolderPath: stringPtr(root)}

	relative, err := mediaItemFileTarget(item, "Season 01/Episode 01.mkv")
	if err != nil {
		t.Fatalf("relative mediaItemFileTarget returned error: %v", err)
	}
	if relative != filepath.Join(root, "Season 01", "Episode 01.mkv") {
		t.Fatalf("relative mediaItemFileTarget = %q", relative)
	}

	absolute := filepath.Join(root, "Movie.mkv")
	target, err := mediaItemFileTarget(item, absolute)
	if err != nil {
		t.Fatalf("absolute mediaItemFileTarget returned error: %v", err)
	}
	if target != absolute {
		t.Fatalf("absolute mediaItemFileTarget = %q", target)
	}

	for _, value := range []string{"", "..", "../outside.mkv", filepath.Join(root, "..", "outside.mkv"), "%2e%2e/outside.mkv"} {
		if _, err := mediaItemFileTarget(item, value); err == nil {
			t.Fatalf("expected %q to be rejected", value)
		}
	}

	if _, err := mediaItemFileTarget(MediaItem{}, "Movie.mkv"); err == nil {
		t.Fatal("expected missing media folder to be rejected")
	}
}

func TestSCNMedia010MediaFilesInRootDiscoversOnlyVisibleMediaFiles(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "Movie.mkv"))
	writeTestFile(t, filepath.Join(root, "clip.txt"))
	writeTestFile(t, filepath.Join(root, ".hidden", "Hidden.mp4"))
	writeTestFile(t, filepath.Join(root, "Season 01", "Episode 01.MP4"))

	files, err := mediaFilesInRoot(root)
	if err != nil {
		t.Fatalf("mediaFilesInRoot returned error: %v", err)
	}
	expectStrings(t, files, []string{
		filepath.Join(root, "Movie.mkv"),
		filepath.Join(root, "Season 01", "Episode 01.MP4"),
	})

	fileRoot := filepath.Join(root, "Movie.mkv")
	if _, err := mediaFilesInRoot(fileRoot); err == nil {
		t.Fatal("expected a file root to be rejected")
	}
	if _, err := mediaFilesInRoot(filepath.Join(root, "missing")); err == nil {
		t.Fatal("expected a missing root to be rejected")
	}
}

func writeTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("media"), 0o644); err != nil {
		t.Fatal(err)
	}
}
