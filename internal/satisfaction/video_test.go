package satisfaction

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func TestVideoTargetMissingWithoutPersistedVideoTrack(t *testing.T) {
	result := EvaluateVideoTarget(mediaItem(), mediaProfile(), storage.MediaFileFact{ID: uuid.New(), FilePath: "/media/movie.mkv"})

	if result.Target.State != targets.StateMissing {
		t.Fatalf("state = %s", result.Target.State)
	}
	if len(result.Candidates) != 0 {
		t.Fatalf("candidates = %#v", result.Candidates)
	}
}

func TestVideoTargetSatisfiedFromPersistedFacts(t *testing.T) {
	quality := "webdl-1080p"
	codec := "h264"
	pixel := "yuv420p"
	result := EvaluateVideoTarget(mediaItem(), mediaProfile(), videoFact(quality, codec, "", pixel, "mkv"))

	if result.Target.State != targets.StateSatisfied {
		t.Fatalf("state = %s reasons=%v", result.Target.State, result.Target.Reasons)
	}
	if len(result.Candidates) != 1 || result.Candidates[0].VisualState != targets.VisualMatching {
		t.Fatalf("candidates = %#v", result.Candidates)
	}
}

func TestVideoTargetPartialWithPreciseFailedFields(t *testing.T) {
	quality := "webdl-1080p"
	codec := "hevc"
	pixel := "yuv420p10le"
	result := EvaluateVideoTarget(mediaItem(), mediaProfile(), videoFact(quality, codec, "", pixel, "mkv"))

	if result.Target.State != targets.StatePartial {
		t.Fatalf("state = %s", result.Target.State)
	}
	if len(result.FailedRequirements) != 2 {
		t.Fatalf("failed = %#v", result.FailedRequirements)
	}
}

func TestVideoTargetPartialWhenTrackResolutionBelowSelectedQuality(t *testing.T) {
	quality := "webdl-1080p"
	codec := "h264"
	pixel := "yuv420p"
	fact := videoFact(quality, codec, "", pixel, "mkv")
	fact.Tracks[0].Width = int32Ptr(1280)
	fact.Tracks[0].Height = int32Ptr(536)

	result := EvaluateVideoTarget(mediaItem(), mediaProfile(), fact)

	if result.Target.State != targets.StatePartial {
		t.Fatalf("target = %#v", result.Target)
	}
	if len(result.FailedRequirements) != 1 || result.FailedRequirements[0] != "video resolution is below selected quality 1080p" {
		t.Fatalf("failed = %#v", result.FailedRequirements)
	}
}

func TestVideoTargetPendingForKnownContainerRemux(t *testing.T) {
	profile := mediaProfile()
	profile.FinalContainer = "mp4"
	quality := "webdl-1080p"
	codec := "h264"
	pixel := "yuv420p"
	result := EvaluateVideoTarget(mediaItem(), profile, videoFact(quality, codec, "", pixel, "mkv"))

	if result.Target.State != targets.StatePending || result.Target.RequiredOperation == nil {
		t.Fatalf("target = %#v", result.Target)
	}
	if result.Target.RequiredOperation.Type != targets.OperationContainerRemux {
		t.Fatalf("operation = %#v", result.Target.RequiredOperation)
	}
}

func TestVideoTargetUpgradeableForQualityTarget(t *testing.T) {
	profile := mediaProfile()
	target := "remux-2160p"
	profile.UpgradeUntilQualityID = &target
	quality := "webdl-1080p"
	codec := "h264"
	pixel := "yuv420p"
	result := EvaluateVideoTarget(mediaItem(), profile, videoFact(quality, codec, "", pixel, "mkv"))

	if result.Target.State != targets.StateUpgradeable {
		t.Fatalf("target = %#v", result.Target)
	}
}

func mediaProfile() *storage.MediaProfile {
	return &storage.MediaProfile{
		FinalContainer: "mkv",
		QualityIDs:     []string{"webdl-1080p", "remux-2160p"},
		VideoTarget: storage.MediaProfileVideoTarget{
			Codecs:       []string{"h264"},
			PixelFormats: []string{"yuv420p"},
		},
	}
}

func mediaItem() storage.MediaItem {
	return storage.MediaItem{ID: uuid.New(), Type: "movie", Title: "Scenario Movie"}
}

func videoFact(quality string, codec string, hdr string, pixel string, container string) storage.MediaFileFact {
	factID := uuid.New()
	return storage.MediaFileFact{
		ID:              factID,
		MediaItemID:     uuid.New(),
		FilePath:        "/media/movie." + container,
		QualityID:       &quality,
		ContainerFormat: &container,
		Tracks: []storage.MediaFileTrackFact{{
			ID:              uuid.New(),
			MediaFileFactID: factID,
			StreamIndex:     0,
			TrackType:       "video",
			Codec:           &codec,
			HDRFormat:       &hdr,
			PixelFormat:     &pixel,
			Width:           int32Ptr(1920),
			Height:          int32Ptr(800),
		}},
	}
}

func int32Ptr(value int32) *int32 {
	return &value
}
