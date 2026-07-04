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
		nil,
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

func TestSCNMedia012PreviewArgsSeekFromRequestedStartTime(t *testing.T) {
	startTime := 120.5
	args := mediaPreviewArgsWithPlan(
		"/media/movie.mkv",
		nil,
		&startTime,
		mediaPreviewTranscodePlan{videoCodec: "copy", audioCodec: "copy"},
	)

	if !hasArgPair(args, "-ss", "120.500") {
		t.Fatalf("expected ffmpeg args to seek to requested start time, got %#v", args)
	}
	if !argBefore(args, "-ss", "-i") {
		t.Fatalf("expected ffmpeg seek arg before input, got %#v", args)
	}
	if !hasArgPair(args, "-c:v", "libx264") || !hasArgPair(args, "-c:a", "aac") {
		t.Fatalf("expected seeked preview to transcode streams for stable A/V sync, got %#v", args)
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

func TestSCNMedia012PreviewInfoReportsDirectModeForCompatibleMp4(t *testing.T) {
	track := int32(1)
	info := mediaPreviewInfoFromTracks("/media/movie.mp4", []MediaFileTrack{
		{
			Type:        Video,
			Codec:       previewString("h264"),
			PixelFormat: previewString("yuv420p"),
			BitRate:     previewString("4000000"),
		},
		{Type: Audio, Index: &track, Codec: previewString("aac"), BitRate: previewString("640000")},
	}, &track)

	if info.StreamingMode != Direct {
		t.Fatalf("streaming mode = %q, want direct", info.StreamingMode)
	}
}

func TestSCNMedia012PreviewInfoReportsRemuxModeAndSelectedBitrate(t *testing.T) {
	track := int32(2)
	info := mediaPreviewInfoFromTracks("/media/movie.mkv", []MediaFileTrack{
		{
			Type:        Video,
			Codec:       previewString("h264"),
			PixelFormat: previewString("yuv420p"),
			BitRate:     previewString("4000000"),
		},
		{Type: Audio, Index: &track, Codec: previewString("aac"), BitRate: previewString("640000")},
	}, &track)

	if info.StreamingMode != Remux {
		t.Fatalf("streaming mode = %q, want remux", info.StreamingMode)
	}
	if info.LiveBitRate == nil || *info.LiveBitRate != "4640000" {
		t.Fatalf("live bit rate = %#v, want 4640000", info.LiveBitRate)
	}
	if info.VideoTrack == nil || info.AudioTrack == nil {
		t.Fatalf("expected selected video and audio tracks, got %#v", info)
	}
}

func TestSCNMedia012PreviewInfoReportsTranscodeMode(t *testing.T) {
	info := mediaPreviewInfoFromTracks("/media/movie.mkv", []MediaFileTrack{
		{Type: Video, Codec: previewString("hevc"), PixelFormat: previewString("yuv420p10le")},
		{Type: Audio, Codec: previewString("dts")},
	}, nil)

	if info.StreamingMode != Transcode {
		t.Fatalf("streaming mode = %q, want transcode", info.StreamingMode)
	}
}

func TestSCNMedia012ProbeDurationIgnoresInvalidValues(t *testing.T) {
	if duration := optionalProbeDuration("5400.25"); duration == nil || *duration != 5400.25 {
		t.Fatalf("duration = %#v, want 5400.25", duration)
	}
	if duration := optionalProbeDuration("unknown"); duration != nil {
		t.Fatalf("duration = %#v, want nil", duration)
	}
	if duration := optionalProbeDuration("0"); duration != nil {
		t.Fatalf("duration = %#v, want nil", duration)
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

func argBefore(args []string, first string, second string) bool {
	firstIndex := slices.Index(args, first)
	secondIndex := slices.Index(args, second)
	return firstIndex >= 0 && secondIndex >= 0 && firstIndex < secondIndex
}

func previewString(value string) *string {
	return &value
}
