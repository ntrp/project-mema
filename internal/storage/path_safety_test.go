package storage

import (
	"path/filepath"
	"testing"
)

func TestSafePathUnderRootRejectsTraversalAndRootTargets(t *testing.T) {
	root := t.TempDir()
	inside := filepath.Join(root, "Movie", "File.mkv")

	got, err := safePathUnderRoot(root, "Movie/File.mkv", false)
	if err != nil {
		t.Fatalf("relative path rejected: %v", err)
	}
	if got != inside {
		t.Fatalf("relative path = %q, want %q", got, inside)
	}

	got, err = safePathUnderRoot(root, inside, false)
	if err != nil {
		t.Fatalf("absolute path rejected: %v", err)
	}
	if got != inside {
		t.Fatalf("absolute path = %q, want %q", got, inside)
	}

	for _, value := range []string{
		"",
		"..",
		"../outside.mkv",
		filepath.Join(root, "..", "outside.mkv"),
		"%2e%2e/outside.mkv",
		"..%2foutside.mkv",
		root,
	} {
		if _, err := safePathUnderRoot(root, value, false); err == nil {
			t.Fatalf("expected %q to be rejected", value)
		}
	}

	if _, err := safePathUnderRoot(root, root, true); err != nil {
		t.Fatalf("root should be allowed when requested: %v", err)
	}
}

func TestMovedMediaFileTargetRejectsUnsafeSourceAndTarget(t *testing.T) {
	oldRoot := t.TempDir()
	newRoot := t.TempDir()
	source := filepath.Join(oldRoot, "Season 01", "Episode 01.mkv")

	target, err := movedMediaFileTarget(oldRoot, newRoot, source)
	if err != nil {
		t.Fatalf("movedMediaFileTarget returned error: %v", err)
	}
	if target != filepath.Join(newRoot, "Season 01", "Episode 01.mkv") {
		t.Fatalf("target = %q", target)
	}

	for _, source := range []string{
		filepath.Join(oldRoot, "..", "outside.mkv"),
		"../outside.mkv",
		"%2e%2e/outside.mkv",
		oldRoot,
	} {
		if _, err := movedMediaFileTarget(oldRoot, newRoot, source); err == nil {
			t.Fatalf("expected source %q to be rejected", source)
		}
	}

	if _, err := movedMediaFileTarget(oldRoot, string(filepath.Separator), source); err == nil {
		t.Fatal("expected filesystem root destination to be rejected")
	}
}
