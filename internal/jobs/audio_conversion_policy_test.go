package jobs

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestAudioConversionPolicyBlocksDisabledConversion(t *testing.T) {
	codec := "aac"
	decision := DecideAudioConversion(AudioConversionInput{
		Policy:      "disabled",
		SourceCodec: "flac",
		TargetCodec: &codec,
	})

	if decision.Status != "blocked" || decision.Allowed {
		t.Fatalf("decision = %#v", decision)
	}
}

func TestAudioConversionPolicyAllowsLosslessSource(t *testing.T) {
	codec := "aac"
	decision := DecideAudioConversion(AudioConversionInput{
		Policy:      "losslessToLossy",
		SourceCodec: "flac",
		TargetCodec: &codec,
	})

	if decision.Status != "allowed" || !decision.Allowed || !decision.Needed {
		t.Fatalf("decision = %#v", decision)
	}
}

func TestAudioConversionPolicyBlocksLossySourceWhenLosslessRequired(t *testing.T) {
	codec := "eac3"
	decision := DecideAudioConversion(AudioConversionInput{
		Policy:      "losslessToLossy",
		SourceCodec: "aac",
		TargetCodec: &codec,
	})

	if decision.Status != "blocked" || decision.Allowed {
		t.Fatalf("decision = %#v", decision)
	}
}

func TestAudioConversionPolicyAllowsLossySourceWhenConfigured(t *testing.T) {
	codec := "eac3"
	decision := DecideAudioConversion(AudioConversionInput{
		Policy:      "lossyToLossy",
		SourceCodec: "aac",
		TargetCodec: &codec,
	})

	if decision.Status != "allowed" || !decision.Allowed {
		t.Fatalf("decision = %#v", decision)
	}
}

func TestFfmpegAudioConversionArgs(t *testing.T) {
	decision := AudioConversionDecision{
		Allowed:           true,
		TargetCodec:       "eac3",
		TargetChannels:    "5.1",
		TargetBitrateKbps: 640,
	}

	args, err := FfmpegAudioConversionArgs("/library/in.mka", "/library/out.mka", decision)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"-y", "-i", "/library/in.mka", "-map", "0:a:0",
		"-c:a", "eac3", "-ac", "6", "-b:a", "640k", "/library/out.mka",
	}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %#v, want %#v", args, want)
	}
}

func TestAudioConversionProvenanceKeepsSourceLineage(t *testing.T) {
	sourceID := uuid.New()
	artifactID := uuid.New()
	provenance := AudioConversionProvenance(
		storage.MediaComponentArtifact{
			ID:         artifactID,
			SourceID:   sourceID,
			StreamID:   2,
			StreamType: "audio",
		},
		AudioConversionDecision{
			Status:      "allowed",
			Policy:      "losslessToLossy",
			SourceCodec: "flac",
			TargetCodec: "aac",
		},
	)

	if provenance["sourceId"] != sourceID.String() || provenance["artifactId"] != artifactID.String() {
		t.Fatalf("provenance = %#v", provenance)
	}
}
