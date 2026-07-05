package jobs

import (
	"testing"

	"media-manager/internal/storage"
)

func TestMovieYearFallbackQueryRemovesMovieYear(t *testing.T) {
	year := int32(2026)
	item := storage.MediaItem{Type: "movie", Title: "Obsession", Year: &year}
	if got := movieYearFallbackQuery(item, "Obsession 2026"); got != "Obsession" {
		t.Fatalf("fallback = %q, want Obsession", got)
	}
}

func TestMovieYearFallbackQueryRemovesSeparatedMovieYear(t *testing.T) {
	year := int32(2026)
	item := storage.MediaItem{Type: "movie", Title: "Obsession", Year: &year}
	if got := movieYearFallbackQuery(item, "Obsession.2026.2160p"); got != "Obsession 2160p" {
		t.Fatalf("fallback = %q, want Obsession 2160p", got)
	}
}

func TestMovieYearFallbackQuerySkipsSeries(t *testing.T) {
	year := int32(2026)
	item := storage.MediaItem{Type: "serie", Title: "Obsession", Year: &year}
	if got := movieYearFallbackQuery(item, "Obsession 2026"); got != "" {
		t.Fatalf("fallback = %q, want empty", got)
	}
}

func TestMovieYearFallbackQuerySkipsQueryWithoutYear(t *testing.T) {
	year := int32(2026)
	item := storage.MediaItem{Type: "movie", Title: "Obsession", Year: &year}
	if got := movieYearFallbackQuery(item, "Obsession"); got != "" {
		t.Fatalf("fallback = %q, want empty", got)
	}
}
