package satisfaction

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func TestAudioTargetMissingWithoutLanguageTrack(t *testing.T) {
	evaluation := EvaluateAudioTargets(mediaItem(), audioProfile(), audioFact())

	if len(evaluation.Results) != 1 || evaluation.Results[0].Target.State != targets.StateMissing {
		t.Fatalf("evaluation = %#v", evaluation)
	}
}

func TestAudioTargetSatisfiedFromOneOfMultipleTracks(t *testing.T) {
	evaluation := EvaluateAudioTargets(mediaItem(), audioProfile(), audioFact(
		audioTrack(0, "english", "aac", "2.0", 192),
		audioTrack(1, "english", "eac3", "5.1", 768),
	))

	if evaluation.Results[0].Target.State != targets.StateSatisfied {
		t.Fatalf("result = %#v", evaluation.Results[0])
	}
	if len(evaluation.Candidates) != 2 {
		t.Fatalf("candidates = %#v", evaluation.Candidates)
	}
}

func TestAudioTargetPartialIncludesFailedRequirements(t *testing.T) {
	evaluation := EvaluateAudioTargets(mediaItem(), audioProfile(), audioFact(
		audioTrack(0, "english", "aac", "2.0", 128),
	))

	result := evaluation.Results[0]
	if result.Target.State != targets.StatePartial || len(result.FailedRequirements) != 3 {
		t.Fatalf("result = %#v", result)
	}
}

func TestAudioTargetPendingWhenTranscodeAllowed(t *testing.T) {
	profile := audioProfile()
	profile.AudioLossyTranscodePolicy = "lossyToLossy"
	evaluation := EvaluateAudioTargets(mediaItem(), profile, audioFact(
		audioTrack(0, "english", "aac", "2.0", 128),
	))

	result := evaluation.Results[0]
	if result.Target.State != targets.StatePending || result.Target.RequiredOperation == nil {
		t.Fatalf("result = %#v", result)
	}
	if result.Target.RequiredOperation.Type != targets.OperationAudioTranscode {
		t.Fatalf("operation = %#v", result.Target.RequiredOperation)
	}
}

func TestAudioUnwantedCandidatesDoNotChangeTargetState(t *testing.T) {
	profile := audioProfile()
	profile.RemoveUnwantedAudio = true
	evaluation := EvaluateAudioTargets(mediaItem(), profile, audioFact(
		audioTrack(0, "english", "eac3", "5.1", 768),
		audioTrack(1, "japanese", "aac", "2.0", 192),
	))

	if evaluation.Results[0].Target.State != targets.StateSatisfied {
		t.Fatalf("result = %#v", evaluation.Results[0])
	}
	unwanted := false
	for _, candidate := range evaluation.Candidates {
		if candidate.VisualState == targets.VisualUnwanted {
			unwanted = true
		}
	}
	if !unwanted {
		t.Fatalf("candidates = %#v", evaluation.Candidates)
	}
}

func audioProfile() *storage.MediaProfile {
	minimum := int32(640)
	return &storage.MediaProfile{
		AudioLossyTranscodePolicy: "disabled",
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:         "english",
			TargetCodec:        stringPtr("eac3"),
			TargetChannels:     []string{"5.1"},
			MinimumBitrateKbps: &minimum,
		}},
	}
}

func audioFact(tracks ...storage.MediaFileTrackFact) storage.MediaFileFact {
	factID := uuid.New()
	for index := range tracks {
		tracks[index].MediaFileFactID = factID
	}
	return storage.MediaFileFact{
		ID:       factID,
		FilePath: "/media/movie.mkv",
		Tracks:   tracks,
	}
}

func audioTrack(index int32, language string, codec string, channels string, bitrate int32) storage.MediaFileTrackFact {
	return storage.MediaFileTrackFact{
		ID:          uuid.New(),
		StreamIndex: index,
		TrackType:   "audio",
		LanguageID:  &language,
		Codec:       &codec,
		Channels:    &channels,
		BitrateKbps: &bitrate,
	}
}

func stringPtr(value string) *string {
	return &value
}
