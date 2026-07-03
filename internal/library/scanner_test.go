package library

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSCNMedia010DiscoverClassifiesLibraryFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "Movies/Scenario.Movie.2026.mkv")
	writeFile(t, root, "Shows/Scenario Series/Season 1/Scenario.Series.S01E02.mkv")
	writeFile(t, root, "Unsafe/Other.Movie.2026.1080p.WEBDL.mkv")
	writeFile(t, root, ".hidden/Hidden.Movie.2026.mkv")
	writeFile(t, root, "Notes/readme.txt")

	files, err := Discover(root)
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if len(files) != 3 {
		t.Fatalf("discovered %d files, want 3: %#v", len(files), files)
	}

	byPath := map[string]DiscoveredFile{}
	for _, file := range files {
		byPath[file.Path] = file
	}

	movie := byPath["Movies/Scenario.Movie.2026.mkv"]
	if movie.DetectedKind != MediaKindMovie || movie.DetectedTitle != "Scenario Movie" {
		t.Fatalf("movie discovery = %#v", movie)
	}
	if movie.DetectedYear == nil || *movie.DetectedYear != 2026 {
		t.Fatalf("movie year = %#v, want 2026", movie.DetectedYear)
	}
	if !movie.SafeMatch {
		t.Fatalf("movie should be a safe match: %#v", movie)
	}

	episode := byPath["Shows/Scenario Series/Season 1/Scenario.Series.S01E02.mkv"]
	if episode.DetectedKind != MediaKindSeries || episode.DetectedTitle != "Scenario Series" {
		t.Fatalf("episode discovery = %#v", episode)
	}
	if !episode.SafeMatch {
		t.Fatalf("episode should be a safe match: %#v", episode)
	}

	unsafe := byPath["Unsafe/Other.Movie.2026.1080p.WEBDL.mkv"]
	if unsafe.SafeMatch {
		t.Fatalf("release-token file should not be a safe match: %#v", unsafe)
	}
}

func writeFile(t *testing.T, root string, relativePath string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %q: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte("fixture"), 0o644); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}
