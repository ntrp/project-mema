package httpapi

import (
	"testing"

	"media-manager/internal/storage"
)

func TestShouldBlockReleaseMismatchHonorsOverride(t *testing.T) {
	item := storage.MediaItem{Type: "series", Title: "Friends"}
	release := storage.ReleaseCandidate{Title: "Graceful.Friends.S01E01.1080p.WEB-DL"}

	if !shouldBlockReleaseMismatch(item, release, false) {
		t.Fatal("expected mismatched release to be blocked without override")
	}
	if shouldBlockReleaseMismatch(item, release, true) {
		t.Fatal("expected mismatched release to be allowed with override")
	}
}
