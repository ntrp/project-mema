package httpapi

import (
	"os"
	"path/filepath"
	"testing"

	"media-manager/internal/storage"
)

func TestMediaFileSubtitleSatisfactionUsesEmbeddedTracks(t *testing.T) {
	result := mediaFileSubtitleSatisfaction([]MediaFileTrack{
		{Type: Subtitle, Language: stringPtr("eng")},
	}, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
	}, "mixed", nil)

	if result.State != MediaFileSubtitleSatisfactionStateSatisfied {
		t.Fatalf("expected satisfied subtitle state, got %#v", result)
	}
	if len(result.MatchedLanguages) != 1 || result.MatchedLanguages[0] != "english" {
		t.Fatalf("expected english match, got %#v", result.MatchedLanguages)
	}
}

func TestExternalSubtitleLanguagesForPathUsesStoredSidecarLanguage(t *testing.T) {
	root := t.TempDir()
	mediaPath := filepath.Join(root, "Movie.mkv")
	subtitlePath := filepath.Join(root, "English.srt")
	if err := os.WriteFile(subtitlePath, []byte("subtitle"), 0o600); err != nil {
		t.Fatal(err)
	}
	language := "english"

	languages := externalSubtitleLanguagesForPath(nil, []storage.MediaItemSidecar{{
		MediaFilePath: mediaPath,
		FilePath:      subtitlePath,
		SidecarType:   storage.MediaSidecarSubtitle,
		LanguageID:    &language,
	}}, mediaPath)

	if len(languages) != 1 || languages[0] != "english" {
		t.Fatalf("expected stored sidecar language, got %#v", languages)
	}
	result := mediaFileSubtitleSatisfaction(nil, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
	}, "external", languages)
	if result.State != MediaFileSubtitleSatisfactionStateSatisfied {
		t.Fatalf("expected sidecar subtitle to satisfy external target, got %#v", result)
	}
}

func TestExternalSubtitleLanguagesForPathUsesAvailableSidecarLanguage(t *testing.T) {
	root := t.TempDir()
	mediaPath := filepath.Join(root, "Movie.mkv")
	subtitlePath := filepath.Join(root, "English.srt")
	if err := os.WriteFile(subtitlePath, []byte("subtitle"), 0o600); err != nil {
		t.Fatal(err)
	}

	languages := externalSubtitleLanguagesForPath(nil, nil, mediaPath)

	if len(languages) != 1 || languages[0] != "english" {
		t.Fatalf("expected available sidecar language, got %#v", languages)
	}
	result := mediaFileSubtitleSatisfaction(nil, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
	}, "external", languages)
	if result.State != MediaFileSubtitleSatisfactionStateSatisfied {
		t.Fatalf("expected available sidecar subtitle to satisfy external target, got %#v", result)
	}
}

func TestMediaFileSubtitleSatisfactionRequiresImportForEmbeddedModeExternalSidecar(t *testing.T) {
	result := mediaFileSubtitleSatisfaction(nil, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
	}, "embedded", []string{"english"})

	if result.State != MediaFileSubtitleSatisfactionStateMissing {
		t.Fatalf("expected external sidecar source to require embedded import, got %#v", result)
	}
	if len(result.MissingLanguages) != 1 || result.MissingLanguages[0] != "english" {
		t.Fatalf("expected english missing until import, got %#v", result.MissingLanguages)
	}
}

func TestMediaFileSubtitleSatisfactionReportsMissingTargets(t *testing.T) {
	result := mediaFileSubtitleSatisfaction([]MediaFileTrack{
		{Type: Subtitle, Language: stringPtr("jpn")},
	}, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
	}, "embedded", nil)

	if result.State != MediaFileSubtitleSatisfactionStateMissing {
		t.Fatalf("expected missing subtitle state, got %#v", result)
	}
	if len(result.MissingLanguages) != 1 || result.MissingLanguages[0] != "english" {
		t.Fatalf("expected english missing, got %#v", result.MissingLanguages)
	}
}

func TestMediaFileSubtitleSatisfactionDeduplicatesTargets(t *testing.T) {
	result := mediaFileSubtitleSatisfaction(nil, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
		{LanguageID: "eng"},
	}, "external", nil)

	if len(result.MissingLanguages) != 1 || result.MissingLanguages[0] != "english" {
		t.Fatalf("expected one missing english target, got %#v", result.MissingLanguages)
	}
	if len(result.WantedLanguages) != 1 || result.WantedLanguages[0] != "english" {
		t.Fatalf("expected one wanted english target, got %#v", result.WantedLanguages)
	}
}

func TestMediaFileSubtitleSatisfactionUsesExternalRecords(t *testing.T) {
	result := mediaFileSubtitleSatisfaction(nil, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
	}, "external", []string{"eng"})

	if result.State != MediaFileSubtitleSatisfactionStateSatisfied {
		t.Fatalf("expected satisfied subtitle state, got %#v", result)
	}
}

func TestMediaFileSubtitleSatisfactionReportsEveryTargetMissing(t *testing.T) {
	result := mediaFileSubtitleSatisfaction(nil, []storage.MediaProfileSubtitleTarget{
		{LanguageID: "english"},
	}, "external", nil)

	if result.State != MediaFileSubtitleSatisfactionStateMissing {
		t.Fatalf("expected missing subtitle state, got %#v", result)
	}
	if len(result.MissingLanguages) != 1 || result.MissingLanguages[0] != "english" {
		t.Fatalf("expected missing target subtitles, got %#v", result.MissingLanguages)
	}
}
