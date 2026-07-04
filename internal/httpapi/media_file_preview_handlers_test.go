package httpapi

import (
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestSCNMedia012PreviewSegmentArgsSelectRequestedAudioStreamAndCopiesVideo(t *testing.T) {
	track := int32(2)
	args := mediaPreviewHLSSegmentArgs(
		"/media/movie.mkv",
		&track,
		120,
		6,
		mediaPreviewDecision{plan: mediaPreviewTranscodePlan{videoCodec: "copy", audioCodec: "aac"}},
	)

	if !slices.Contains(args, "0:2") {
		t.Fatalf("expected ffmpeg args to map requested audio stream, got %#v", args)
	}
	if !hasArgPair(args, "-c:v", "copy") || !hasArgPair(args, "-bsf:v", "h264_mp4toannexb") {
		t.Fatalf("expected ffmpeg args to copy browser-compatible video into TS, got %#v", args)
	}
	if !hasArgPair(args, "-c:a", "aac") || !slices.Contains(args, "-ac") {
		t.Fatalf("expected ffmpeg args to encode browser-compatible audio, got %#v", args)
	}
	if !hasArgPair(args, "-f", "mpegts") {
		t.Fatalf("expected MPEG-TS HLS segment output, got %#v", args)
	}
}

func TestSCNMedia012PreviewSegmentArgsFallsBackToFirstAudioStream(t *testing.T) {
	args := mediaPreviewHLSSegmentArgs(
		"/media/movie.mkv",
		nil,
		0,
		6,
		mediaPreviewDecision{plan: mediaPreviewTranscodePlan{videoCodec: "libx264", audioCodec: "aac"}},
	)

	if !slices.Contains(args, "0:a:0?") {
		t.Fatalf("expected optional first audio stream mapping, got %#v", args)
	}
	if !hasArgPair(args, "-c:v", "libx264") {
		t.Fatalf("expected ffmpeg args to transcode incompatible video, got %#v", args)
	}
	if !slices.Contains(args, "-force_key_frames") {
		t.Fatalf("expected ffmpeg args to force segment keyframes, got %#v", args)
	}
}

func TestSCNMedia012PreviewHLSPlaylistBuildsSegmentUrls(t *testing.T) {
	request := httptest.NewRequest("GET", "http://internal/api/media/items/abc/files/preview?path=%2Fmedia%2Fmovie.mkv", nil)
	track := int32(2)
	playlist := mediaPreviewHLSPlaylistText(request, "/media/movie.mkv", &track, Webkit, []mediaPreviewHLSSegment{
		{start: 0, duration: 6},
		{start: 6, duration: 6},
		{start: 12, duration: 1},
	})

	if !strings.Contains(playlist, "#EXT-X-PLAYLIST-TYPE:VOD") {
		t.Fatalf("expected VOD HLS playlist, got %s", playlist)
	}
	if strings.Count(playlist, "preview-segment?") != 3 {
		t.Fatalf("expected 3 segment URLs, got %s", playlist)
	}
	if !strings.Contains(playlist, "audioTrackIndex=2") ||
		!strings.Contains(playlist, "clientProfile=webkit") ||
		!strings.Contains(playlist, "segmentStartSeconds=12.000") {
		t.Fatalf("expected selected audio, client profile, and final segment start in playlist, got %s", playlist)
	}
}

func TestSCNMedia012PreviewHLSSegmentsUseKeyframeBoundaries(t *testing.T) {
	segments := mediaPreviewHLSSegments(20, []float64{0, 5, 7, 14, 19})
	want := []mediaPreviewHLSSegment{
		{start: 0, duration: 7},
		{start: 7, duration: 7},
		{start: 14, duration: 6},
	}

	if !equalPreviewSegments(segments, want) {
		t.Fatalf("segments = %#v, want %#v", segments, want)
	}
	if target := mediaPreviewHLSTargetDuration(segments); target != 7 {
		t.Fatalf("target duration = %d, want 7", target)
	}
}

func TestSCNMedia012PreviewKeyframesNormalizeProbeFrames(t *testing.T) {
	keyframes := normalizedPreviewKeyframes([]ffprobeKeyframe{
		{BestEffortTimestampTime: "7.000"},
		{BestEffortTimestampTime: "0.000"},
		{BestEffortTimestampTime: "7.020"},
		{PktPtsTime: "14.500"},
		{BestEffortTimestampTime: "invalid"},
	})
	want := []float64{0, 7, 14.5}

	if !slices.Equal(keyframes, want) {
		t.Fatalf("keyframes = %#v, want %#v", keyframes, want)
	}
}

func TestSCNMedia012PreviewDecisionCopiesCompatibleStreams(t *testing.T) {
	track := int32(1)
	decision := mediaPreviewDecisionFromTracks("/media/movie.mkv", []MediaFileTrack{
		{Type: Video, Codec: previewString("h264"), PixelFormat: previewString("yuv420p")},
		{Type: Audio, Index: &track, Codec: previewString("aac")},
	}, &track, Browser)

	if decision.mode != Remux || decision.deliveryProtocol != mediaPreviewDeliveryHLS {
		t.Fatalf("decision = %#v, want HLS remux", decision)
	}
	if decision.plan.videoCodec != "copy" || decision.plan.audioCodec != "copy" {
		t.Fatalf("plan = %#v, want stream copy", decision.plan)
	}
}

func TestSCNMedia012PreviewDecisionTranscodesIncompatibleStreams(t *testing.T) {
	decision := mediaPreviewDecisionFromTracks("/media/movie.mkv", []MediaFileTrack{
		{Type: Video, Codec: previewString("hevc"), PixelFormat: previewString("yuv420p10le")},
		{Type: Audio, Codec: previewString("dts")},
	}, nil, Browser)

	if decision.mode != Transcode {
		t.Fatalf("streaming mode = %q, want transcode", decision.mode)
	}
	if decision.plan.videoCodec != "libx264" || decision.plan.audioCodec != "aac" {
		t.Fatalf("plan = %#v, want h264/aac transcode", decision.plan)
	}
	if !slices.Contains(decision.reasons, "video_codec_not_supported") || !slices.Contains(decision.reasons, "audio_codec_not_supported") {
		t.Fatalf("reasons = %#v, want video and audio codec reasons", decision.reasons)
	}
}

func TestSCNMedia012PreviewDecisionTranscodesVideoForWebKitHLS(t *testing.T) {
	track := int32(1)
	decision := mediaPreviewDecisionFromTracks("/media/movie.mkv", []MediaFileTrack{
		{Type: Video, Codec: previewString("h264"), PixelFormat: previewString("yuv420p")},
		{Type: Audio, Index: &track, Codec: previewString("aac")},
	}, &track, Webkit)

	if decision.mode != Transcode || decision.deliveryProtocol != mediaPreviewDeliveryHLS {
		t.Fatalf("decision = %#v, want HLS transcode", decision)
	}
	if decision.plan.videoCodec != "libx264" || decision.plan.audioCodec != "copy" {
		t.Fatalf("plan = %#v, want video transcode and audio copy", decision.plan)
	}
	if !slices.Contains(decision.reasons, mediaPreviewReasonWebKitHLS) {
		t.Fatalf("reasons = %#v, want WebKit HLS reason", decision.reasons)
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
	}, &track, Browser)

	if info.StreamingMode != Direct || info.DeliveryProtocol != mediaPreviewDeliveryFile {
		t.Fatalf("preview info = %#v, want direct file playback", info)
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
	}, &track, Browser)

	if info.StreamingMode != Remux || info.DeliveryProtocol != mediaPreviewDeliveryHLS {
		t.Fatalf("preview info = %#v, want HLS remux", info)
	}
	if info.LiveBitRate == nil || *info.LiveBitRate != "4640000" {
		t.Fatalf("live bit rate = %#v, want 4640000", info.LiveBitRate)
	}
	if info.VideoTrack == nil || info.AudioTrack == nil {
		t.Fatalf("expected selected video and audio tracks, got %#v", info)
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

func TestSCNMedia012ProbeContainerInfoReportsFormatFields(t *testing.T) {
	container := mediaFileContainerInfo(ffprobeFormat{
		BitRate:    "5500000",
		Format:     "matroska,webm",
		FormatName: "Matroska / WebM",
	})

	if *container.bitRate != "5500000" || *container.format != "matroska,webm" || *container.formatName != "Matroska / WebM" {
		t.Fatalf("container = %#v", container)
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

func equalPreviewSegments(left []mediaPreviewHLSSegment, right []mediaPreviewHLSSegment) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
