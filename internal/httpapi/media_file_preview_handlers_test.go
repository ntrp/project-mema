package httpapi

import (
	"slices"
	"testing"
)

func TestSCNMedia012PreviewArgsSelectRequestedAudioStreamAndCopiesVideo(t *testing.T) {
	track := int32(2)
	args := mediaPreviewArgsWithPlan(
		"/media/movie.mkv",
		&track,
		mediaPreviewTranscodePlan{videoCodec: "copy", audioCodec: "aac"},
	)

	if !slices.Contains(args, "0:2") {
		t.Fatalf("expected ffmpeg args to map requested audio stream, got %#v", args)
	}
	if !hasArgPair(args, "-c:v", "copy") {
		t.Fatalf("expected ffmpeg args to copy compatible video, got %#v", args)
	}
	if !slices.Contains(args, "aac") {
		t.Fatalf("expected ffmpeg args to encode browser-compatible audio, got %#v", args)
	}
	if !slices.Contains(args, "frag_keyframe+empty_moov+default_base_moof") {
		t.Fatalf("expected fragmented MP4 flags, got %#v", args)
	}
}

func TestSCNMedia012PreviewArgsFallsBackToFirstAudioStream(t *testing.T) {
	args := mediaPreviewArgsWithPlan(
		"/media/movie.mkv",
		nil,
		mediaPreviewTranscodePlan{videoCodec: "libx264", audioCodec: "aac"},
	)

	if !slices.Contains(args, "0:a:0?") {
		t.Fatalf("expected optional first audio stream mapping, got %#v", args)
	}
	if !hasArgPair(args, "-c:v", "libx264") {
		t.Fatalf("expected ffmpeg args to transcode incompatible video, got %#v", args)
	}
	if !slices.Contains(args, "-preset") {
		t.Fatalf("expected ffmpeg args to include encoder preset, got %#v", args)
	}
}

func TestSCNMedia012PreviewPlanCopiesCompatibleStreams(t *testing.T) {
	track := int32(1)
	plan := mediaPreviewPlanFromTracks([]MediaFileTrack{
		{Type: Video, Codec: previewString("h264"), PixelFormat: previewString("yuv420p")},
		{Type: Audio, Index: &track, Codec: previewString("aac")},
	}, &track)

	if plan.videoCodec != "copy" {
		t.Fatalf("video codec = %q, want copy", plan.videoCodec)
	}
	if plan.audioCodec != "copy" {
		t.Fatalf("audio codec = %q, want copy", plan.audioCodec)
	}
}

func TestSCNMedia012PreviewPlanTranscodesIncompatibleStreams(t *testing.T) {
	plan := mediaPreviewPlanFromTracks([]MediaFileTrack{
		{Type: Video, Codec: previewString("hevc"), PixelFormat: previewString("yuv420p10le")},
		{Type: Audio, Codec: previewString("dts")},
	}, nil)

	if plan.videoCodec != "libx264" {
		t.Fatalf("video codec = %q, want libx264", plan.videoCodec)
	}
	if plan.audioCodec != "aac" {
		t.Fatalf("audio codec = %q, want aac", plan.audioCodec)
	}
}

func hasArgPair(args []string, key string, value string) bool {
	for index := 0; index < len(args)-1; index += 1 {
		if args[index] == key && args[index+1] == value {
			return true
		}
	}
	return false
}

func previewString(value string) *string {
	return &value
}
