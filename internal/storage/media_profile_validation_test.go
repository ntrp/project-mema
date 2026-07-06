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
			{LanguageID: " English ", Score: 25, Required: true, SubtitleType: "embedded"},
			{LanguageID: "english", Score: 10, Required: false, SubtitleType: "external"},
			{LanguageID: "German", Score: 5, Required: false, SubtitleType: "not-real"},
		},
		ComponentTargets: []MediaProfileComponentTarget{
			{ComponentType: "video", Required: true, LanguageID: stringPtr("English"), Channels: stringPtr("5.1"), Source: "release"},
			{ComponentType: "audio", Required: true, LanguageID: stringPtr(" English "), Codec: stringPtr(" AAC "), Channels: stringPtr("5.1"), Source: "release", FallbackBehavior: "preferExisting"},
			{ComponentType: "subtitle", Required: false, LanguageID: stringPtr("German"), Source: "release", FallbackBehavior: "allowMissing"},
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
	if input.SubtitleLanguages[0].LanguageID != "english" || input.SubtitleLanguages[0].SubtitleType != "embedded" || input.SubtitleLanguages[0].Score != 25 {
		t.Fatalf("unexpected first subtitle language: %#v", input.SubtitleLanguages[0])
	}
	if input.SubtitleLanguages[1].LanguageID != "german" || input.SubtitleLanguages[1].SubtitleType != "any" || input.SubtitleLanguages[1].Score != 5 {
		t.Fatalf("unexpected second subtitle language: %#v", input.SubtitleLanguages[1])
	}
	if len(input.ComponentTargets) != 3 || input.ComponentTargets[0].LanguageID != nil {
		t.Fatalf("unexpected component targets: %#v", input.ComponentTargets)
	}
	if *input.ComponentTargets[1].LanguageID != "english" || *input.ComponentTargets[1].Codec != "aac" || input.ComponentTargets[1].FallbackBehavior != "preferExisting" {
		t.Fatalf("unexpected audio component target: %#v", input.ComponentTargets[1])
	}
	if input.ComponentTargets[2].Source != "subtitleProvider" {
		t.Fatalf("expected subtitle provider source, got %#v", input.ComponentTargets[2])
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

func TestNormalizeMediaProfileInputRejectsInvalidComponentTarget(t *testing.T) {
	_, err := normalizeMediaProfileInput(MediaProfileInput{
		ComponentTargets: []MediaProfileComponentTarget{{ComponentType: "not-real"}},
	}, []string{"webdl-1080p"})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid component target error, got %v", err)
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
