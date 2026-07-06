package jobs

import (
	"strings"
	"testing"

	"media-manager/internal/storage"
)

func TestSubtitleFulfillmentNeedsUseBothProvidersForMixedMode(t *testing.T) {
	item := storage.MediaItem{
		Type:                         "movie",
		Title:                        "Scenario Movie",
		FilePaths:                    []string{"/library/Scenario.Movie.mkv"},
		SubtitlePreferredMode:        "mixed",
		AllowSubtitleReleaseFallback: true,
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{
			LanguageID: "english",
			Formats:    []string{"vtt"},
		}},
	}

	needs := SubtitleFulfillmentNeeds(item)

	if len(needs) != 2 {
		t.Fatalf("needs = %#v", needs)
	}
	if needs[0].Mode != "provider" || needs[1].Mode != "alternateRelease" {
		t.Fatalf("needs = %#v", needs)
	}
}

func TestSubtitleFulfillmentNeedsSkipSatisfiedExternalTarget(t *testing.T) {
	item := storage.MediaItem{
		Type:                  "movie",
		Title:                 "Scenario Movie",
		FilePaths:             []string{"/library/Scenario.Movie.mkv"},
		SubtitlePreferredMode: "external",
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{
			LanguageID: "english",
		}},
		ExternalSubtitles: []storage.MediaItemSubtitle{{
			LanguageID: "english",
			FilePath:   "/library/Scenario.Movie.english.srt",
		}},
	}

	if needs := SubtitleFulfillmentNeeds(item); len(needs) != 0 {
		t.Fatalf("needs = %#v", needs)
	}
}

func TestSubtitleFulfillmentNeedsUseEmbeddedInventory(t *testing.T) {
	item := storage.MediaItem{
		Type:                  "movie",
		Title:                 "Scenario Movie",
		SubtitlePreferredMode: "embedded",
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{
			LanguageID: "english",
			Formats:    []string{"ass"},
		}},
		ComponentSources: []storage.MediaComponentSource{{
			StreamInventory: `{"streams":[{"type":"subtitle","language":"english","format":"ass"}]}`,
		}},
	}

	if needs := SubtitleFulfillmentNeeds(item); len(needs) != 0 {
		t.Fatalf("needs = %#v", needs)
	}
}

func TestPlanSubtitleFulfillmentChoosesAlternateRelease(t *testing.T) {
	item := storage.MediaItem{
		Type:                         "movie",
		Title:                        "Scenario Movie",
		SubtitlePreferredMode:        "embedded",
		AllowSubtitleReleaseFallback: true,
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{
			LanguageID: "english",
			Formats:    []string{"ass"},
		}},
	}

	candidates := PlanSubtitleFulfillment(item, []storage.ReleaseCandidateInput{
		{Title: "Scenario.Movie.2026.English.ASS.1080p", IndexerName: "idx", DownloadURL: "https://dl"},
	})

	if len(candidates) != 1 || !strings.Contains(candidates[0].Metadata, "subtitleFulfillment") {
		t.Fatalf("candidates = %#v", candidates)
	}
}

func TestSubtitleFulfillmentSourceInputPreservesReleaseMetadata(t *testing.T) {
	candidate := SubtitleFulfillmentCandidate{
		ReleaseTitle: "Subtitle.Release",
		DownloadURL:  "https://download",
		Metadata:     `{"kind":"subtitleFulfillment"}`,
	}

	input := SubtitleFulfillmentSourceInput(candidate, "/library/subs.mkv", `{"streams":[]}`)

	if input.SourceRole != "subtitle" || input.SourceFilePath != "/library/subs.mkv" {
		t.Fatalf("input = %#v", input)
	}
	if input.ReleaseTitle == nil || *input.ReleaseTitle != "Subtitle.Release" {
		t.Fatalf("release title = %#v", input.ReleaseTitle)
	}
	if input.ReleaseID == nil || *input.ReleaseID != "https://download" {
		t.Fatalf("release id = %#v", input.ReleaseID)
	}
}
