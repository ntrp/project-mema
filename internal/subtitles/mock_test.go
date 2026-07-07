package subtitles

import (
	"context"
	"strings"
	"testing"
)

func TestMockSubtitleProviderMatchesTitleAndLanguage(t *testing.T) {
	service := NewService(nil)
	candidates, err := service.Search(context.Background(), Config{
		Name: "Mock Subtitles",
		Type: "mock",
		MockSubtitles: []MockSubtitle{{
			Title:      "Scenario Movie",
			LanguageID: "english",
			Format:     "vtt",
		}},
	}, SearchRequest{Title: " scenario   movie ", LanguageID: "eng"})
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 1 || candidates[0].Format != "vtt" {
		t.Fatalf("candidates = %#v", candidates)
	}
}

func TestMockSubtitleProviderDownloadRepeatsMockCue(t *testing.T) {
	service := NewService(nil)
	download, err := service.Download(context.Background(), Config{Type: "mock"}, Candidate{Format: "srt"})
	if err != nil {
		t.Fatal(err)
	}
	content := string(download.Content)
	if !strings.Contains(content, "00:00:00,000 --> 00:00:01,000\nmock") {
		t.Fatalf("missing first mock cue: %q", content)
	}
	if !strings.Contains(content, "00:00:03,000 --> 00:00:04,000\nmock") {
		t.Fatalf("missing repeated mock cue: %q", content)
	}
}
