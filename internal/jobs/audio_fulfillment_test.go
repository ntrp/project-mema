package jobs

import (
	"strings"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestAudioFulfillmentNeedsSkipSatisfiedArtifact(t *testing.T) {
	language := "english"
	profile := storage.MediaProfile{
		AudioTargets: []storage.MediaProfileAudioTarget{{LanguageID: language, Score: 25}},
	}
	item := storage.MediaItem{
		Title: "Scenario Movie",
		ComponentSources: []storage.MediaComponentSource{{
			Artifacts: []storage.MediaComponentArtifact{{
				StreamType: "audio",
				Language:   &language,
			}},
		}},
	}

	if needs := AudioFulfillmentNeeds(item, profile); len(needs) != 0 {
		t.Fatalf("needs = %#v", needs)
	}
}

func TestAudioFulfillmentNeedsBuildAlternateReleaseQuery(t *testing.T) {
	codec := "eac3"
	profile := storage.MediaProfile{
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:     "english",
			TargetCodec:    &codec,
			TargetChannels: []string{"5.1"},
			Score:          25,
		}},
	}

	needs := AudioFulfillmentNeeds(storage.MediaItem{Type: "movie", Title: "Scenario Movie"}, profile)

	if len(needs) != 1 {
		t.Fatalf("needs = %#v", needs)
	}
	if !strings.Contains(needs[0].Query, "english eac3 5.1") {
		t.Fatalf("query = %q", needs[0].Query)
	}
}

func TestPlanAudioFulfillmentChoosesScoredCandidate(t *testing.T) {
	codec := "eac3"
	profile := storage.MediaProfile{
		QualityIDs: []string{"webdl-1080p"},
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:     "english",
			TargetCodec:    &codec,
			TargetChannels: []string{"5.1"},
			Score:          25,
		}},
	}

	candidates := PlanAudioFulfillment(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		&profile,
		[]storage.ReleaseCandidateInput{
			{Title: "Scenario.Movie.2026.English.1080p.WEBDL.AAC.2.0"},
			{Title: "Scenario.Movie.2026.English.1080p.WEBDL.EAC3.5.1", IndexerName: "idx", DownloadURL: "https://dl"},
		},
		[]storage.Language{{Code: "EN", DisplayName: "English", Aliases: []string{"English"}}},
	)

	if len(candidates) != 1 || candidates[0].ReleaseTitle == "" {
		t.Fatalf("candidates = %#v", candidates)
	}
	if !strings.Contains(candidates[0].Metadata, "audioFulfillment") {
		t.Fatalf("metadata = %q", candidates[0].Metadata)
	}
}

func TestAudioFulfillmentSourceInputPreservesReleaseMetadata(t *testing.T) {
	candidate := AudioFulfillmentCandidate{
		ReleaseTitle: "Audio.Release",
		DownloadURL:  "https://download",
		Metadata:     `{"kind":"audioFulfillment"}`,
	}

	input := AudioFulfillmentSourceInput(candidate, "/library/audio.mkv", `{"streams":[]}`)

	if input.SourceRole != "audio" || input.SourceFilePath != "/library/audio.mkv" {
		t.Fatalf("input = %#v", input)
	}
	if input.ReleaseTitle == nil || *input.ReleaseTitle != "Audio.Release" {
		t.Fatalf("release title = %#v", input.ReleaseTitle)
	}
	if input.ReleaseID == nil || *input.ReleaseID != "https://download" {
		t.Fatalf("release id = %#v", input.ReleaseID)
	}
}

func TestAudioFulfillmentNeedsUseInventory(t *testing.T) {
	codec := "eac3"
	minimum := int32(640)
	profile := storage.MediaProfile{
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:         "english",
			TargetCodec:        &codec,
			TargetChannels:     []string{"5.1"},
			MinimumBitrateKbps: &minimum,
		}},
	}
	item := storage.MediaItem{
		ID: uuid.New(),
		ComponentSources: []storage.MediaComponentSource{{
			StreamInventory: `{"streams":[{"type":"audio","language":"english","codec":"DD+","channels":"5.1","bitrateKbps":768}]}`,
		}},
	}

	if needs := AudioFulfillmentNeeds(item, profile); len(needs) != 0 {
		t.Fatalf("needs = %#v", needs)
	}
}
