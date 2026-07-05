package storage

import (
	"errors"
	"strings"
	"testing"
)

func TestNormalizeProfileQualityIDs(t *testing.T) {
	qualityIDs, err := normalizeProfileQualityIDs([]string{" webdl-1080p ", "webdl-1080p", "", "bluray-2160p"})
	if err != nil {
		t.Fatalf("normalize quality ids: %v", err)
	}

	expectStrings(t, qualityIDs, []string{"webdl-1080p", "bluray-2160p"})
}

func TestNormalizeProfileQualityIDsRejectsUnknownQuality(t *testing.T) {
	if _, err := normalizeProfileQualityIDs([]string{"not-real"}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input for unknown quality, got %v", err)
	}
}

func TestNormalizeMediaProfileInput(t *testing.T) {
	upgradeUntil := " webdl-1080p "
	input, err := normalizeMediaProfileInput(MediaProfileInput{
		UpgradeUntilQualityID:             &upgradeUntil,
		PreferredProtocol:                 " torrent ",
		SeriesPackPreference:              "preferPacks",
		MinimumCustomFormatScoreIncrement: 5,
		TargetLanguageScores: []MediaProfileLanguageScore{
			{LanguageID: " English ", Score: 100, Required: true},
			{LanguageID: "english", Score: 80},
			{LanguageID: "Brazilian Portuguese", Score: 50},
		},
		SubtitleLanguages: []MediaProfileSubtitleLanguage{
			{LanguageID: " English ", Required: true, SubtitleType: "embedded"},
			{LanguageID: "english", Required: false, SubtitleType: "external"},
			{LanguageID: "German", Required: false, SubtitleType: "not-real"},
		},
	}, []string{"webdl-1080p", "bluray-2160p"})
	if err != nil {
		t.Fatalf("normalize media profile input: %v", err)
	}

	if input.UpgradeUntilQualityID == nil || *input.UpgradeUntilQualityID != "webdl-1080p" {
		t.Fatalf("expected trimmed upgrade quality, got %#v", input.UpgradeUntilQualityID)
	}
	if input.PreferredProtocol != "torrent" || input.SeriesPackPreference != "preferPacks" {
		t.Fatalf("unexpected protocol or pack preference: %#v", input)
	}
	expectStrings(t, input.TargetLanguages, []string{"english", "brazilian-portuguese"})
	if !input.TargetLanguageScores[0].Required || input.TargetLanguageScores[0].Score != 100 {
		t.Fatalf("expected first language score to retain score and required flag: %#v", input.TargetLanguageScores[0])
	}
	if len(input.SubtitleLanguages) != 2 {
		t.Fatalf("expected deduped subtitle languages, got %#v", input.SubtitleLanguages)
	}
	if input.SubtitleLanguages[0].LanguageID != "english" || input.SubtitleLanguages[0].SubtitleType != "embedded" {
		t.Fatalf("unexpected first subtitle language: %#v", input.SubtitleLanguages[0])
	}
	if input.SubtitleLanguages[1].LanguageID != "german" || input.SubtitleLanguages[1].SubtitleType != "any" {
		t.Fatalf("unexpected second subtitle language: %#v", input.SubtitleLanguages[1])
	}
}

func TestNormalizeMediaProfileInputRejectsInvalidUpgradeQuality(t *testing.T) {
	upgradeUntil := "raw-hd"
	_, err := normalizeMediaProfileInput(MediaProfileInput{
		UpgradeUntilQualityID: &upgradeUntil,
	}, []string{"webdl-1080p"})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid upgrade quality error, got %v", err)
	}
}

func TestMediaProfileSlug(t *testing.T) {
	if got := normalizeMediaProfileName("  Main   Profile  "); got != "Main Profile" {
		t.Fatalf("expected compacted profile name, got %q", got)
	}
	if got := mediaProfileSlug("Main Profile: 4K + Anime"); got != "main-profile-4k-anime" {
		t.Fatalf("expected normalized slug, got %q", got)
	}
	long := mediaProfileSlug(strings.Repeat("a", 120))
	if len(long) != 80 {
		t.Fatalf("expected slug trimmed to 80 chars, got %d", len(long))
	}
}
