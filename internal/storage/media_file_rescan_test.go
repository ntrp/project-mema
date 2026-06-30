package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
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
