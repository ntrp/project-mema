package decisions

import (
	"testing"

	"media-manager/internal/storage"
)

func TestEvaluateReleaseMatchRejectsWrongMovieYear(t *testing.T) {
	year := int32(2026)
	match := EvaluateReleaseMatch(
		storage.MediaItem{Type: "movie", Title: "The Movie", Year: &year},
		storage.ReleaseCandidate{Title: "The.Movie.2025.1080p.WEBDL"},
	)
	if match.Severity != "error" {
		t.Fatalf("expected error, got %q", match.Severity)
	}
}

func TestEvaluateReleaseMatchAcceptsRequestedEpisode(t *testing.T) {
	season := int32(1)
	episode := int32(2)
	match := EvaluateReleaseMatch(
		storage.MediaItem{Type: "series", Title: "The Show"},
		storage.ReleaseCandidate{
			Title:            "The.Show.S01E02.1080p.WEBDL",
			SearchKind:       "episode",
			RequestedSeason:  &season,
			RequestedEpisode: &episode,
		},
	)
	if match.Severity != "info" {
		t.Fatalf("expected info, got %q: %v", match.Severity, match.Details)
	}
}

func TestEvaluateReleaseMatchWarnsSeasonPackForEpisodeSearch(t *testing.T) {
	season := int32(1)
	episode := int32(2)
	match := EvaluateReleaseMatch(
		storage.MediaItem{Type: "series", Title: "The Show"},
		storage.ReleaseCandidate{
			Title:            "The.Show.S01.1080p.WEBDL",
			SearchKind:       "episode",
			RequestedSeason:  &season,
			RequestedEpisode: &episode,
		},
	)
	if match.Severity != "warning" {
		t.Fatalf("expected warning, got %q: %v", match.Severity, match.Details)
	}
}

func TestSearchQueriesForSeasonFallback(t *testing.T) {
	season := int32(1)
	queries := SearchQueriesForCriteria(ReleaseSearchCriteria{
		Kind:         "season",
		Title:        "The Show",
		SeasonNumber: &season,
	}, "")

	want := []string{"The Show s1", "The Show S01"}
	if len(queries) != len(want) {
		t.Fatalf("queries = %#v, want %#v", queries, want)
	}
	for index, expected := range want {
		if queries[index] != expected {
			t.Fatalf("queries[%d] = %q, want %q", index, queries[index], expected)
		}
	}
}
