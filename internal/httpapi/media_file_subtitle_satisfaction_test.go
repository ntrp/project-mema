package httpapi

import (
	"testing"

	"media-manager/internal/storage"
)

func TestMediaFileSubtitleSatisfactionUsesEmbeddedTracks(t *testing.T) {
	result := mediaFileSubtitleSatisfaction([]MediaFileTrack{
		{Type: Subtitle, Language: stringPtr("eng")},
	}, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english", Source: "any"},
	}, nil)

	if result.State != MediaFileSubtitleSatisfactionStateSatisfied {
		t.Fatalf("expected satisfied subtitle state, got %#v", result)
	}
	if len(result.MatchedLanguages) != 1 || result.MatchedLanguages[0] != "english" {
		t.Fatalf("expected english match, got %#v", result.MatchedLanguages)
	}
}

func TestMediaFileSubtitleSatisfactionReportsMissingTargets(t *testing.T) {
	result := mediaFileSubtitleSatisfaction([]MediaFileTrack{
		{Type: Subtitle, Language: stringPtr("jpn")},
	}, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english", Source: "embedded"},
	}, nil)

	if result.State != MediaFileSubtitleSatisfactionStateMissing {
		t.Fatalf("expected missing subtitle state, got %#v", result)
	}
	if len(result.MissingLanguages) != 1 || result.MissingLanguages[0] != "english" {
		t.Fatalf("expected english missing, got %#v", result.MissingLanguages)
	}
}

func TestMediaFileSubtitleSatisfactionUsesExternalRecords(t *testing.T) {
	result := mediaFileSubtitleSatisfaction(nil, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english", Source: "external"},
	}, []string{"eng"})

	if result.State != MediaFileSubtitleSatisfactionStateSatisfied {
		t.Fatalf("expected satisfied subtitle state, got %#v", result)
	}
}

func TestMediaFileSubtitleSatisfactionReportsEveryTargetMissing(t *testing.T) {
	result := mediaFileSubtitleSatisfaction(nil, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english", Source: "external"},
	}, nil)

	if result.State != MediaFileSubtitleSatisfactionStateMissing {
		t.Fatalf("expected missing subtitle state, got %#v", result)
	}
	if len(result.MissingLanguages) != 1 || result.MissingLanguages[0] != "english" {
		t.Fatalf("expected missing target subtitles, got %#v", result.MissingLanguages)
	}
}
