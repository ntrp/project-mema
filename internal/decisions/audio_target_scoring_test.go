package decisions

import (
	"strings"
	"testing"

	"media-manager/internal/storage"
)

func TestAudioTargetScoringAddsDetailContributors(t *testing.T) {
	preferred := int32(640)
	codec := "eac3"
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-2160p"},
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:           "english",
			Score:                25,
			TargetCodec:          &codec,
			TargetChannels:       []string{"5.1"},
			PreferredBitrateKbps: &preferred,
		}},
	}

	match := EvaluateReleaseMatchWithLanguageContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.English.2160p.WEBDL.EAC3.5.1.640kbps"},
		&profile,
		nil,
		[]storage.Language{{Code: "EN", DisplayName: "English", Aliases: []string{"English"}}},
	)

	if match.Severity != "info" {
		t.Fatalf("expected info match, got %q: %v", match.Severity, match.Details)
	}
	if match.LanguageContributors[0].Score != 25 || match.TargetScore != 75 || match.Score != 100 {
		t.Fatalf("scores = language %#v target %d total %d", match.LanguageContributors, match.TargetScore, match.Score)
	}
	if len(match.TargetContributors) != 3 {
		t.Fatalf("contributors = %#v", match.TargetContributors)
	}
}

func TestAudioTargetScoringRejectsKnownMinimumBitrateMiss(t *testing.T) {
	minimum := int32(640)
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:         "english",
			Score:              25,
			MinimumBitrateKbps: &minimum,
		}},
	}

	match := EvaluateReleaseMatchWithLanguageContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.English.1080p.WEBDL.AAC.2.0.128kbps"},
		&profile,
		nil,
		[]storage.Language{{Code: "EN", DisplayName: "English", Aliases: []string{"English"}}},
	)

	if match.Severity != "error" {
		t.Fatalf("expected bitrate rejection, got %q: %v", match.Severity, match.Details)
	}
	if len(match.Details) == 0 || !strings.Contains(match.Details[0], "below the profile minimum") {
		t.Fatalf("details = %#v", match.Details)
	}
}

func TestAudioTargetScoringDoesNotRejectUnknownBitrate(t *testing.T) {
	minimum := int32(640)
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:         "english",
			Score:              25,
			MinimumBitrateKbps: &minimum,
		}},
	}

	match := EvaluateReleaseMatchWithLanguageContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.English.1080p.WEBDL.AAC.2.0"},
		&profile,
		nil,
		[]storage.Language{{Code: "EN", DisplayName: "English", Aliases: []string{"English"}}},
	)

	if match.Severity != "info" {
		t.Fatalf("expected info match, got %q: %v", match.Severity, match.Details)
	}
}
