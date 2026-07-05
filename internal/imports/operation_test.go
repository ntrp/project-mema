package imports

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportFileHardlinkCopyAndMove(t *testing.T) {
	t.Run("hardlink", func(t *testing.T) {
		root := t.TempDir()
		source := writeImportSource(t, root, "source.mkv", "video")
		target := filepath.Join(root, "target.mkv")

		if err := importFile(source, target, ImportModeHardlink); err != nil {
			t.Fatal(err)
		}

		sourceInfo, err := os.Stat(source)
		if err != nil {
			t.Fatal(err)
		}
		targetInfo, err := os.Stat(target)
		if err != nil {
			t.Fatal(err)
		}
		if !os.SameFile(sourceInfo, targetInfo) {
			t.Fatalf("target is not linked to source")
		}
	})

	t.Run("copy", func(t *testing.T) {
		root := t.TempDir()
		source := writeImportSource(t, root, "source.mkv", "video")
		target := filepath.Join(root, "target.mkv")

		if err := importFile(source, target, ImportModeCopy); err != nil {
			t.Fatal(err)
		}

		sourceInfo, err := os.Stat(source)
		if err != nil {
			t.Fatal(err)
		}
		targetInfo, err := os.Stat(target)
		if err != nil {
			t.Fatal(err)
		}
		if os.SameFile(sourceInfo, targetInfo) {
			t.Fatalf("target should be a copy, not a hardlink")
		}
		if got := readImportFile(t, target); got != "video" {
			t.Fatalf("target content = %q", got)
		}
	})

	t.Run("move", func(t *testing.T) {
		root := t.TempDir()
		source := writeImportSource(t, root, "source.mkv", "video")
		target := filepath.Join(root, "target.mkv")

		if err := importFile(source, target, ImportModeMove); err != nil {
			t.Fatal(err)
		}

		if _, err := os.Stat(source); !os.IsNotExist(err) {
			t.Fatalf("source stat err = %v, want missing", err)
		}
		if got := readImportFile(t, target); got != "video" {
			t.Fatalf("target content = %q", got)
		}
	})
}

func TestImportFileMoveFailsBeforeRemovingSourceWhenTargetExists(t *testing.T) {
	root := t.TempDir()
	source := writeImportSource(t, root, "source.mkv", "video")
	target := writeImportSource(t, root, "target.mkv", "existing")

	err := importFile(source, target, ImportModeMove)
	if err == nil {
		t.Fatalf("expected target conflict")
	}
	if got := readImportFile(t, source); got != "video" {
		t.Fatalf("source content = %q", got)
	}
	if got := readImportFile(t, target); got != "existing" {
		t.Fatalf("target content = %q", got)
	}
}

func writeImportSource(t *testing.T, root string, name string, content string) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readImportFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
