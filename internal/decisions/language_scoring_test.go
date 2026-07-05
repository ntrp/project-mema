package decisions

import (
	"strings"
	"testing"

	"media-manager/internal/storage"
)

func TestSCNMedia002LanguageCatalogAddsAliasesToReleaseMatch(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		TargetLanguageScores: []storage.MediaProfileLanguageScore{
			{LanguageID: "japanese", Score: 100, Required: true},
		},
	}
	languages := []storage.Language{{
		Code:        "JA",
		DisplayName: "Japanese",
		Aliases:     []string{"JPN", "Nihongo"},
	}}

	match := EvaluateReleaseMatchWithLanguageContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.JPN.1080p.WEBDL"},
		&profile,
		nil,
		languages,
	)

	if match.Severity != "info" {
		t.Fatalf("expected info match, got %q: %v", match.Severity, match.Details)
	}
	if match.Score != 100 {
		t.Fatalf("score = %d, want language score 100", match.Score)
	}
	if len(match.LanguageContributors) != 1 || match.LanguageContributors[0].Label != "Japanese" {
		t.Fatalf("language contributors = %#v", match.LanguageContributors)
	}
	if !hasLanguage(match.Languages, "Japanese") || !hasLanguage(match.Languages, "JA") {
		t.Fatalf("languages = %#v", match.Languages)
	}
}

func TestSCNMedia002RequiredLanguageRejectsMissingReleaseLanguage(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		TargetLanguageScores: []storage.MediaProfileLanguageScore{
			{LanguageID: "japanese", Score: 100, Required: true},
		},
	}

	match := EvaluateReleaseMatchWithLanguageContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.English.1080p.WEBDL"},
		&profile,
		nil,
		[]storage.Language{{Code: "JA", DisplayName: "Japanese", Aliases: []string{"JPN"}}},
	)

	if match.Severity != "error" {
		t.Fatalf("expected language rejection, got %q: %v", match.Severity, match.Details)
	}
	if len(match.Details) == 0 || !strings.Contains(match.Details[0], "Japanese is missing") {
		t.Fatalf("details = %#v", match.Details)
	}
}

func TestSCNMedia002SubtitleLanguageScoreContributesToReleaseMatch(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		SubtitleLanguages: []storage.MediaProfileSubtitleLanguage{
			{LanguageID: "english", Score: 25, Required: false, SubtitleType: "any"},
		},
	}

	match := EvaluateReleaseMatchWithLanguageContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.English.1080p.WEBDL"},
		&profile,
		nil,
		[]storage.Language{{Code: "EN", DisplayName: "English", Aliases: []string{"ENG"}}},
	)

	if match.Severity != "info" {
		t.Fatalf("expected info match, got %q: %v", match.Severity, match.Details)
	}
	if match.Score != 25 {
		t.Fatalf("score = %d, want subtitle score 25", match.Score)
	}
	if len(match.LanguageContributors) != 1 || match.LanguageContributors[0].Label != "Subtitle: English" {
		t.Fatalf("language contributors = %#v", match.LanguageContributors)
	}
}

func TestSCNMedia002RequiredSubtitleLanguageRejectsMissingReleaseLanguage(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		SubtitleLanguages: []storage.MediaProfileSubtitleLanguage{
			{LanguageID: "japanese", Score: 25, Required: true, SubtitleType: "any"},
		},
	}

	match := EvaluateReleaseMatchWithLanguageContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.English.1080p.WEBDL"},
		&profile,
		nil,
		[]storage.Language{{Code: "JA", DisplayName: "Japanese", Aliases: []string{"JPN"}}},
	)

	if match.Severity != "error" {
		t.Fatalf("expected subtitle language rejection, got %q: %v", match.Severity, match.Details)
	}
	if len(match.Details) == 0 || !strings.Contains(match.Details[0], "Required subtitle language Japanese is missing") {
		t.Fatalf("details = %#v", match.Details)
	}
}

func TestSCNMedia002RemoveNonEnabledLanguagesRejectsExtraLanguage(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs:                    []string{"webdl-1080p"},
		RemoveNonEnabledLanguages:     true,
		TargetLanguageScores:          []storage.MediaProfileLanguageScore{{LanguageID: "english", Score: 10}},
		MinimumCustomFormatScore:      0,
		PreferredProtocol:             "any",
		SeriesPackPreference:          "auto",
		UpgradesAllowed:               true,
		UpgradeUntilCustomFormatScore: 0,
	}

	match := EvaluateReleaseMatchWithLanguageContext(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{Title: "Scenario.Movie.2026.English.Japanese.1080p.WEBDL"},
		&profile,
		nil,
		[]storage.Language{{Code: "JA", DisplayName: "Japanese", Aliases: []string{"Japanese"}}},
	)

	if match.Severity != "error" {
		t.Fatalf("expected non-enabled language rejection, got %q: %v", match.Severity, match.Details)
	}
}

func TestSCNMedia002ProtocolPreferenceBreaksReleaseTie(t *testing.T) {
	profile := storage.MediaProfile{
		QualityIDs:           []string{"webdl-1080p"},
		PreferredProtocol:    "usenet",
		TargetLanguages:      []string{},
		TargetLanguageScores: []storage.MediaProfileLanguageScore{},
	}
	decision, ok := NewEngine().ChooseReleaseWithProfile(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		&profile,
		nil,
		[]storage.ReleaseCandidateInput{
			{Title: "Scenario.Movie.2026.1080p.WEBDL", IndexerProtocol: "torrent", SizeBytes: 10},
			{Title: "Scenario.Movie.2026.1080p.WEBDL", IndexerProtocol: "nzb", SizeBytes: 10},
		},
	)

	if !ok {
		t.Fatal("expected release decision")
	}
	if decision.Release.IndexerProtocol != "nzb" {
		t.Fatalf("expected usenet release, got %#v", decision.Release)
	}
}

func hasLanguage(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
