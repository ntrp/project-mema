package jobs

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestSCNMedia002DedupeReleaseCandidatesKeepsBestDuplicate(t *testing.T) {
	guid := "same-release"
	season := int32(1)
	episode := int32(2)
	item := storage.MediaItem{Type: "serie", Title: "Scenario Series"}
	releases := []storage.ReleaseCandidateInput{
		{
			Title:            "Scenario.Series.S01.1080p.WEBDL",
			GUID:             &guid,
			SearchKind:       "episode",
			RequestedSeason:  &season,
			RequestedEpisode: &episode,
		},
		{
			Title:            "Scenario.Series.S01E02.1080p.WEBDL",
			GUID:             &guid,
			SearchKind:       "episode",
			RequestedSeason:  &season,
			RequestedEpisode: &episode,
		},
	}

	deduped := dedupeReleaseCandidates(item, nil, nil, nil, releases)

	if len(deduped) != 1 {
		t.Fatalf("deduped len = %d, want 1: %#v", len(deduped), deduped)
	}
	if deduped[0].Title != "Scenario.Series.S01E02.1080p.WEBDL" {
		t.Fatalf("deduped title = %q", deduped[0].Title)
	}
}

func TestSCNMedia002SortReleaseCandidatesUsesQualityThenCustomFormatThenSeeders(t *testing.T) {
	formatID := uuid.MustParse("00000000-0000-4000-8000-000000000301")
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p", "remux-2160p"},
		CustomFormatScores: []storage.MediaProfileCustomFormatScore{
			{CustomFormatID: formatID, Score: 100},
		},
	}
	formats := []storage.CustomFormat{{
		ID:   formatID,
		Name: "Preferred group",
		IncludeSpecs: []storage.CustomFormatSpec{{
			ID: "preferred", Name: "Preferred", Type: "releaseTitle", Value: "Preferred", Required: true,
		}},
	}}
	lowSeeders := int32(2)
	highSeeders := int32(12)
	releases := []storage.ReleaseCandidateInput{
		{Title: "Scenario.Movie.2026.1080p.WEBDL.LowSeed", SizeBytes: 30, Seeders: &lowSeeders},
		{Title: "Scenario.Movie.2026.1080p.WEBDL.Preferred", SizeBytes: 20, Seeders: &lowSeeders},
		{Title: "Scenario.Movie.2026.Remux.2160p", SizeBytes: 10, Seeders: &lowSeeders},
		{Title: "Scenario.Movie.2026.1080p.WEBDL.HighSeed", SizeBytes: 40, Seeders: &highSeeders},
	}

	sortReleaseCandidates(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		&profile,
		formats,
		nil,
		releases,
	)

	want := []string{
		"Scenario.Movie.2026.Remux.2160p",
		"Scenario.Movie.2026.1080p.WEBDL.Preferred",
		"Scenario.Movie.2026.1080p.WEBDL.HighSeed",
		"Scenario.Movie.2026.1080p.WEBDL.LowSeed",
	}
	for index, title := range want {
		if releases[index].Title != title {
			t.Fatalf("release[%d] = %q, want %q; all = %#v", index, releases[index].Title, title, releases)
		}
	}
}

func TestSCNMedia002ReleaseCandidateRankingFallbackHelpers(t *testing.T) {
	infoURL := " https://indexer.local/info/1 "
	downloadURL := " HTTPS://INDEXER.LOCAL/DOWNLOAD/1 "

	if got := releaseDedupeKey(storage.ReleaseCandidateInput{InfoURL: &infoURL}); got != "https://indexer.local/info/1" {
		t.Fatalf("info url dedupe key = %q", got)
	}
	if got := releaseDedupeKey(storage.ReleaseCandidateInput{DownloadURL: downloadURL}); got != "https://indexer.local/download/1" {
		t.Fatalf("download url dedupe key = %q", got)
	}
	if got := releaseDedupeKey(storage.ReleaseCandidateInput{}); got != "" {
		t.Fatalf("empty dedupe key = %q", got)
	}

	profile := storage.MediaProfile{QualityIDs: []string{"webdl-1080p", "remux-2160p"}}
	if qualityRank("remux-2160p", &profile) <= qualityRank("webdl-1080p", &profile) {
		t.Fatal("profile quality rank should follow profile ordering")
	}
	if qualityRank("missing-quality", &profile) != 0 {
		t.Fatal("unknown profile quality should rank as zero")
	}
	if severityRank("info") <= severityRank("warning") || severityRank("warning") <= severityRank("error") {
		t.Fatal("severity ranks should prefer info over warning over errors")
	}
}
