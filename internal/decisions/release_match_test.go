package decisions

import (
	"testing"

	"github.com/google/uuid"

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
	seasonID := uuid.New()
	episodeID := uuid.New()
	match := EvaluateReleaseMatch(
		storage.MediaItem{Type: "serie", Title: "The Show"},
		storage.ReleaseCandidate{
			Title:            "The.Show.S01E02.1080p.WEBDL",
			SeasonID:         &seasonID,
			EpisodeID:        &episodeID,
			SearchKind:       "episode",
			RequestedSeason:  &season,
			RequestedEpisode: &episode,
		},
	)
	if match.Severity != "info" {
		t.Fatalf("expected info, got %q: %v", match.Severity, match.Details)
	}
	if match.MatchedSeasonID == nil || *match.MatchedSeasonID != seasonID || match.MatchedEpisodeID == nil || *match.MatchedEpisodeID != episodeID {
		t.Fatalf("expected matched persisted ids, got %#v", match)
	}
}

func TestEvaluateReleaseMatchAcceptsExactSeriesTitle(t *testing.T) {
	season := int32(1)
	episode := int32(1)
	match := EvaluateReleaseMatch(
		storage.MediaItem{Type: "serie", Title: "Friends"},
		storage.ReleaseCandidate{
			Title:            "Friends.S01E01.1080p.WEBDL",
			SearchKind:       "episode",
			RequestedSeason:  &season,
			RequestedEpisode: &episode,
		},
	)
	if match.Severity != "info" {
		t.Fatalf("expected info, got %q: %v", match.Severity, match.Details)
	}
}

func TestEvaluateReleaseMatchRejectsSeriesTitleContainingExpectedTitle(t *testing.T) {
	season := int32(1)
	episode := int32(1)
	item := storage.MediaItem{Type: "serie", Title: "Friends"}
	releases := []string{
		"Friends.Like.These.The.Murder.of.Skylar.Neese.S01E01.The.Disappearance.2160p.DSNP.WEB-DL.DD+5.1.DoVi.H.265-playWEB",
		"Graceful.Friends.S01E01.1080p.LINETV.WEB-DL.AAC2.0.H.264-MWeb",
	}

	for _, title := range releases {
		match := EvaluateReleaseMatch(
			item,
			storage.ReleaseCandidate{
				Title:            title,
				SearchKind:       "episode",
				RequestedSeason:  &season,
				RequestedEpisode: &episode,
			},
		)
		if match.Severity != "error" {
			t.Fatalf("expected error for %q, got %q: %v", title, match.Severity, match.Details)
		}
	}
}

func TestEvaluateReleaseMatchWarnsSeasonPackForEpisodeSearch(t *testing.T) {
	season := int32(1)
	episode := int32(2)
	match := EvaluateReleaseMatch(
		storage.MediaItem{Type: "serie", Title: "The Show"},
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

func TestEvaluateReleaseMatchRejectsDisabledQuality(t *testing.T) {
	profile := storage.MediaProfile{QualityIDs: []string{"webdl-1080p"}}
	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "The Movie"},
		storage.ReleaseCandidate{Title: "The.Movie.2026.Remux.2160p"},
		&profile,
		nil,
	)
	if match.Severity != "error" {
		t.Fatalf("expected error, got %q: %v", match.Severity, match.Details)
	}
}

func TestEvaluateReleaseMatchRejectsBelowCustomFormatMinimum(t *testing.T) {
	formatID := uuid.MustParse("00000000-0000-4000-8000-000000000203")
	profile := storage.MediaProfile{
		QualityIDs:               []string{"webdl-1080p"},
		MinimumCustomFormatScore: 50,
		CustomFormatScores: []storage.MediaProfileCustomFormatScore{
			{CustomFormatID: formatID, Score: -100},
		},
	}
	formats := []storage.CustomFormat{{
		ID:   formatID,
		Name: "Bad group",
		IncludeSpecs: []storage.CustomFormatSpec{{
			ID: "bad", Name: "Bad", Type: "releaseTitle", Value: "BadGroup", Required: true,
		}},
	}}
	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "movie", Title: "The Movie"},
		storage.ReleaseCandidate{Title: "The.Movie.2026.WEB-DL.1080p.BadGroup"},
		&profile,
		formats,
	)
	if match.Severity != "error" {
		t.Fatalf("expected error, got %q: %v", match.Severity, match.Details)
	}
}

func TestEvaluateReleaseMatchScoresProfileCustomFormats(t *testing.T) {
	formatID := uuid.MustParse("493b6d1d-bec3-c336-4c59-d7607f7e3405")
	season := int32(1)
	episode := int32(1)
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-2160p"},
		CustomFormatScores: []storage.MediaProfileCustomFormatScore{
			{CustomFormatID: formatID, Score: 1000},
		},
	}
	formats := []storage.CustomFormat{{
		ID:   formatID,
		Name: "HDR",
		IncludeSpecs: []storage.CustomFormatSpec{{
			ID: "hdr", Name: "HDR", Type: "releaseTitle", Value: `\b(HDR)\b`,
		}},
	}}

	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "serie", Title: "Friends"},
		storage.ReleaseCandidate{
			Title:            "Friends.S01E01.NORDiC.2160p.MAX.WEB-DL.DV.HDR.H.265-NORViNE",
			SearchKind:       "episode",
			RequestedSeason:  &season,
			RequestedEpisode: &episode,
		},
		&profile,
		formats,
	)
	if match.Severity != "info" {
		t.Fatalf("expected info, got %q: %v", match.Severity, match.Details)
	}
	if match.CustomFormatScore != 1000 {
		t.Fatalf("custom format score = %d, want 1000", match.CustomFormatScore)
	}
	if match.Score != 1000 {
		t.Fatalf("score = %d, want 1000", match.Score)
	}
	if len(match.CustomFormatContributors) != 1 || match.CustomFormatContributors[0].Label != "HDR" {
		t.Fatalf("custom format contributors = %#v, want HDR", match.CustomFormatContributors)
	}
}

func TestEvaluateReleaseMatchScoresSeededHDRCustomFormat(t *testing.T) {
	formatID := uuid.MustParse("493b6d1d-bec3-c336-4c59-d7607f7e3405")
	season := int32(1)
	episode := int32(1)
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-2160p"},
		CustomFormatScores: []storage.MediaProfileCustomFormatScore{
			{CustomFormatID: formatID, Score: 1000},
		},
	}
	formats := []storage.CustomFormat{{
		ID:   formatID,
		Name: "HDR",
		IncludeSpecs: []storage.CustomFormatSpec{
			{ID: "dv-with-hdr10-fallback-0", Name: "DV With HDR10 fallback", Type: "releaseTitle", Value: `^(?=.*\b(dv|dovi|dolby[ .]?v(ision)?)\b)(?!(?=.*\b(WEB[ ._-]?(DL|Rip)?)\b)(?!.*\b(hulu)\b))`},
			{ID: "hdr-1", Name: "HDR", Type: "releaseTitle", Value: `\b(HDR)\b`},
			{ID: "hdr10-2", Name: "HDR10", Type: "releaseTitle", Value: `\b(HDR10(?![+]|P(lus)?))`},
			{ID: "hdr10-3", Name: "HDR10+", Type: "releaseTitle", Value: `\b(HDR10(?=[+]|P(lus)?))`},
			{ID: "hlg-4", Name: "HLG", Type: "releaseTitle", Value: `\b(HLG)\b`},
			{ID: "pq-5", Name: "PQ", Type: "releaseTitle", Value: `\b(PQ)\b`},
			{ID: "rlsgrp-missing-hdr-6", Name: "RlsGrp (Missing HDR)", Type: "releaseTitle", Value: `^(?=.*\b(FraMeSToR|HQMUX|SiCFoI)\b)(?=.*\b(2160p)\b)(?!.*\b(HDR10([+]|P(lus)?)))(?!.*\b(SDR)\b).*`},
		},
	}}

	match := EvaluateReleaseMatchWithContext(
		storage.MediaItem{Type: "serie", Title: "Friends"},
		storage.ReleaseCandidate{
			Title:            "Friends.S01E01.NORDiC.2160p.MAX.WEB-DL.DV.HDR.H.265-NORViNE",
			SearchKind:       "episode",
			RequestedSeason:  &season,
			RequestedEpisode: &episode,
		},
		&profile,
		formats,
	)
	if match.CustomFormatScore != 1000 {
		t.Fatalf("custom format score = %d, want 1000; contributors = %#v", match.CustomFormatScore, match.CustomFormatContributors)
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
