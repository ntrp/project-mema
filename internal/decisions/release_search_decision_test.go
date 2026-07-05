package decisions

import (
	"testing"
	"time"

	"media-manager/internal/storage"
)

func TestSCNMedia002SearchCriteriaForMoviesAndSeries(t *testing.T) {
	year := int32(2026)
	movie := storage.MediaItem{Type: "movie", Title: "Scenario Movie", Year: &year}
	series := storage.MediaItem{Type: "serie", Title: "Scenario Series", Year: &year}

	if query := SearchQueryForMediaItem(movie); query != "Scenario Movie 2026" {
		t.Fatalf("movie query = %q", query)
	}
	if query := SearchQueryForMediaItem(storage.MediaItem{Title: "No Year"}); query != "No Year" {
		t.Fatalf("no-year query = %q", query)
	}

	episode := SearchCriteriaForQuery(series, "Scenario.Series.S02E03")
	if episode.Kind != "episode" || *episode.SeasonNumber != 2 || *episode.EpisodeNumber != 3 {
		t.Fatalf("episode criteria = %#v", episode)
	}
	season := SearchCriteriaForQuery(series, "Scenario.Series.S02")
	if season.Kind != "season" || *season.SeasonNumber != 2 || season.EpisodeNumber != nil {
		t.Fatalf("season criteria = %#v", season)
	}
	if all := SearchCriteriaForQuery(series, "Scenario Series"); all.Kind != "serie" {
		t.Fatalf("series criteria = %#v", all)
	}
}

func TestSCNMedia002SearchQueryVariantsDeduplicateOriginal(t *testing.T) {
	season := int32(2)
	episode := int32(3)
	queries := SearchQueriesForCriteria(ReleaseSearchCriteria{
		Kind:          "episode",
		Title:         "Scenario Series",
		SeasonNumber:  &season,
		EpisodeNumber: &episode,
	}, "Scenario Series S02E03")

	want := []string{"Scenario Series S02E03", "Scenario Series s2e3"}
	if len(queries) != len(want) {
		t.Fatalf("queries = %#v, want %#v", queries, want)
	}
	for index, expected := range want {
		if queries[index] != expected {
			t.Fatalf("queries[%d] = %q, want %q", index, queries[index], expected)
		}
	}
}

func TestSCNMedia002CandidateInputMatchPreservesMetadata(t *testing.T) {
	seeders := int32(12)
	peers := int32(20)
	publishedAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	match := EvaluateReleaseCandidateInputMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidateInput{
			Title:           "Scenario.Movie.2026.1080p.WEBDL",
			IndexerProtocol: "torrent",
			SizeBytes:       5,
			Seeders:         &seeders,
			Peers:           &peers,
			PublishedAt:     &publishedAt,
		},
		&storage.MediaProfile{QualityIDs: []string{"webdl-1080p"}},
		nil,
	)

	if match.Severity != "info" {
		t.Fatalf("expected info, got %q: %v", match.Severity, match.Details)
	}
	if len(match.RankContributors) == 0 {
		t.Fatalf("expected rank contributors, got %#v", match.RankContributors)
	}
}

func TestSCNMedia002CandidateInputMatchDefaultsToNoProfileContext(t *testing.T) {
	match := EvaluateReleaseCandidateInputMatch(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidateInput{Title: "Scenario.Movie.2026.1080p.WEBDL"},
	)

	if match.Severity != "info" {
		t.Fatalf("expected info, got %q: %v", match.Severity, match.Details)
	}
	if match.QualityID != "webdl-1080p" {
		t.Fatalf("quality id = %q, want webdl-1080p", match.QualityID)
	}
}

func TestSCNMedia002CurrentFileDowngradeReturnsWarning(t *testing.T) {
	match := EvaluateReleaseMatch(
		storage.MediaItem{
			Type:      "movie",
			Title:     "Scenario Movie",
			FilePaths: []string{"Scenario.Movie.2026.2160p.Remux.mkv"},
		},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.1080p.WEBDL"},
	)

	if match.Severity != "warning" {
		t.Fatalf("expected warning, got %q: %v", match.Severity, match.Details)
	}
}

func TestSCNMedia002SeasonPackPreferenceBreaksReleaseTie(t *testing.T) {
	season := int32(1)
	profile := storage.MediaProfile{
		QualityIDs:           []string{"webdl-1080p"},
		SeriesPackPreference: "preferPacks",
	}
	decision, ok := NewEngine().ChooseReleaseWithProfile(
		storage.MediaItem{Type: "serie", Title: "Scenario Series"},
		&profile,
		nil,
		[]storage.ReleaseCandidateInput{
			{
				Title:           "Scenario.Series.S01E01.1080p.WEBDL",
				SearchKind:      "season",
				RequestedSeason: &season,
				SizeBytes:       10,
			},
			{
				Title:           "Scenario.Series.S01.1080p.WEBDL",
				SearchKind:      "season",
				RequestedSeason: &season,
				SizeBytes:       10,
			},
		},
	)

	if !ok {
		t.Fatal("expected release decision")
	}
	if decision.Release.Title != "Scenario.Series.S01.1080p.WEBDL" {
		t.Fatalf("expected season pack, got %q", decision.Release.Title)
	}
}

func TestSCNMedia002EngineDetectsHighestMatchingQuality(t *testing.T) {
	engine := NewEngine()

	if quality := engine.detectQuality("Scenario.Movie.2026.WEBDL-1080p-GRP"); quality.id != "webdl-1080p" {
		t.Fatalf("webdl quality = %#v", quality)
	}
	if quality := engine.detectQuality("Scenario.Movie.2026.Remux-2160p-GRP"); quality.id != "remux-2160p" {
		t.Fatalf("remux quality = %#v", quality)
	}
	if quality := engine.detectQuality("Scenario.Movie.2026.SourceFree-GRP"); quality.id != "" {
		t.Fatalf("unknown quality = %#v", quality)
	}
}
