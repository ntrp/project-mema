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
		FinalContainer:                    "mp4",
		PreferredProtocol:                 " torrent ",
		SeriesPackPreference:              "preferPacks",
		AudioLossyTranscodePolicy:         "lossyToLossy",
		MinimumCustomFormatScoreIncrement: 5,
		AudioTargets: []MediaProfileAudioTarget{
			{LanguageID: " English ", Score: 100, TargetCodec: stringPtr(" AAC "), TargetChannels: []string{"5.1"}},
			{LanguageID: "english", Score: 80},
			{LanguageID: "Brazilian Portuguese", Score: 50},
		},
		SubtitleMode:                 "embedded",
		AllowSubtitleReleaseFallback: true,
		SubtitleTargets: []MediaProfileSubtitleTarget{
			{LanguageID: " English ", Score: 25, Formats: []string{" SRT "}},
			{LanguageID: "english", Score: 10},
			{LanguageID: "German", Score: 5},
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
	if input.FinalContainer != "mp4" {
		t.Fatalf("expected normalized container, got %q", input.FinalContainer)
	}
	if input.AudioLossyTranscodePolicy != "lossyToLossy" {
		t.Fatalf("expected normalized lossy conversion policy, got %q", input.AudioLossyTranscodePolicy)
	}
	if input.SubtitleMode != "embedded" || !input.AllowSubtitleReleaseFallback {
		t.Fatalf("unexpected subtitle settings: %#v", input)
	}
	if len(input.AudioTargets) != 2 {
		t.Fatalf("expected deduped audio targets, got %#v", input.AudioTargets)
	}
	if input.AudioTargets[0].Score != 100 || input.AudioTargets[0].TargetCodec == nil || *input.AudioTargets[0].TargetCodec != "aac" {
		t.Fatalf("expected first audio target to retain fields: %#v", input.AudioTargets[0])
	}
	if len(input.SubtitleTargets) != 2 {
		t.Fatalf("expected deduped subtitle targets, got %#v", input.SubtitleTargets)
	}
	if input.SubtitleTargets[0].LanguageID != "english" || len(input.SubtitleTargets[0].Formats) != 1 || input.SubtitleTargets[0].Formats[0] != "subrip" {
		t.Fatalf("unexpected first subtitle target: %#v", input.SubtitleTargets[0])
	}
	if input.SubtitleTargets[1].LanguageID != "german" || input.SubtitleTargets[1].Score != 5 {
		t.Fatalf("unexpected second subtitle target: %#v", input.SubtitleTargets[1])
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

func TestNormalizeMediaProfileInputRejectsMissingAudioTarget(t *testing.T) {
	_, err := normalizeMediaProfileInput(MediaProfileInput{
		FinalContainer: "mkv",
	}, []string{"webdl-1080p"})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected missing audio target error, got %v", err)
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
