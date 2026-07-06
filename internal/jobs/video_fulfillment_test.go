package jobs

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestPlanVideoFulfillmentDetectsTargetMismatch(t *testing.T) {
	sourceID := uuid.New()
	profile := storage.MediaProfile{
		VideoTarget: storage.MediaProfileVideoTarget{
			Codecs:       []string{"hevc"},
			HDRFormats:   []string{"hdr10"},
			PixelFormats: []string{"yuv420p10le"},
		},
	}
	item := storage.MediaItem{
		ComponentSources: []storage.MediaComponentSource{{
			ID:              sourceID,
			SourceRole:      "baseVideo",
			RetentionState:  "retained",
			RetainedPath:    "/library/base.mkv",
			StreamInventory: `{"streams":[{"type":"video","codec":"h264","hdrFormat":"sdr","pixelFormat":"yuv420p"}]}`,
		}},
	}

	plan := PlanVideoFulfillment(item, &profile)

	if plan.Status != "transcodeRequired" || plan.SourcePath != "/library/base.mkv" {
		t.Fatalf("plan = %#v", plan)
	}
	if plan.Provenance["sourceId"] != sourceID.String() {
		t.Fatalf("provenance = %#v", plan.Provenance)
	}
}

func TestPlanVideoFulfillmentSatisfied(t *testing.T) {
	profile := storage.MediaProfile{
		VideoTarget: storage.MediaProfileVideoTarget{
			Codecs:       []string{"hevc"},
			HDRFormats:   []string{"hdr10"},
			PixelFormats: []string{"yuv420p10le"},
		},
	}
	item := storage.MediaItem{
		ComponentSources: []storage.MediaComponentSource{{
			SourceRole:      "baseVideo",
			RetentionState:  "retained",
			RetainedPath:    "/library/base.mkv",
			StreamInventory: `{"streams":[{"type":"video","codec":"hevc","hdrFormat":"hdr10","pixelFormat":"yuv420p10le"}]}`,
		}},
	}

	plan := PlanVideoFulfillment(item, &profile)

	if plan.Status != "satisfied" {
		t.Fatalf("plan = %#v", plan)
	}
}

func TestFfprobeVideoArgs(t *testing.T) {
	args, err := FfprobeVideoArgs("/library/base.mkv")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"-v", "error", "-select_streams", "v:0", "-show_streams", "-of", "json", "/library/base.mkv"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %#v, want %#v", args, want)
	}
}

func TestFfmpegVideoTranscodeArgs(t *testing.T) {
	plan := VideoFulfillmentPlan{
		Status:      "transcodeRequired",
		TargetCodec: "hevc",
		TargetPixel: "yuv420p10le",
	}

	args, err := FfmpegVideoTranscodeArgs("/library/in.mkv", "/library/out.mkv", plan)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"-y", "-i", "/library/in.mkv", "-map", "0", "-c", "copy",
		"-c:v", "libx265", "-pix_fmt", "yuv420p10le", "/library/out.mkv",
	}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %#v, want %#v", args, want)
	}
}
