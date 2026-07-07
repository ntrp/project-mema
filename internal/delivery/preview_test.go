package delivery

import (
	"slices"
	"strings"
	"testing"
)

func TestSCNMedia012PreviewSegmentArgsSelectRequestedAudioStreamAndCopiesVideo(t *testing.T) {
	track := int32(2)
	args := SegmentArgs(
		"/media/movie.mkv",
		&track,
		120,
		6,
		Decision{Plan: TranscodePlan{VideoCodec: "copy", AudioCodec: "aac"}},
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
	args := SegmentArgs(
		"/media/movie.mkv",
		nil,
		0,
		6,
		Decision{Plan: TranscodePlan{VideoCodec: "libx264", AudioCodec: "aac"}},
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
	track := int32(2)
	playlist := HLSPlaylistText(PlaylistRequest{
		Path:          "/api/media/items/abc/files/preview",
		FilePath:      "/media/movie.mkv",
		AudioTrack:    &track,
		ClientProfile: ClientWebKit,
		Segments: []HLSSegment{
			{Start: 0, Duration: 6},
			{Start: 6, Duration: 6},
			{Start: 12, Duration: 1},
		},
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
	segments := HLSSegments(20, []float64{0, 5, 7, 14, 19})
	want := []HLSSegment{
		{Start: 0, Duration: 7},
		{Start: 7, Duration: 7},
		{Start: 14, Duration: 6},
	}

	if !slices.Equal(segments, want) {
		t.Fatalf("segments = %#v, want %#v", segments, want)
	}
	if target := HLSTargetDuration(segments); target != 7 {
		t.Fatalf("target duration = %d, want 7", target)
	}
}

func TestSCNMedia012PreviewKeyframesNormalizeProbeFrames(t *testing.T) {
	keyframes := NormalizeKeyframes([]ffprobeKeyframe{
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
	decision := DecisionFromTracks("/media/movie.mkv", []Track{
		{Type: TrackVideo, Codec: previewString("h264"), PixelFormat: previewString("yuv420p")},
		{Type: TrackAudio, Index: &track, Codec: previewString("aac")},
	}, &track, ClientBrowser)

	if decision.Mode != ModeRemux || decision.DeliveryProtocol != ProtocolHLS {
		t.Fatalf("decision = %#v, want HLS remux", decision)
	}
	if decision.Plan.VideoCodec != "copy" || decision.Plan.AudioCodec != "copy" {
		t.Fatalf("plan = %#v, want stream copy", decision.Plan)
	}
}

func TestSCNMedia012PreviewDecisionTranscodesIncompatibleStreams(t *testing.T) {
	decision := DecisionFromTracks("/media/movie.mkv", []Track{
		{Type: TrackVideo, Codec: previewString("hevc"), PixelFormat: previewString("yuv420p10le")},
		{Type: TrackAudio, Codec: previewString("dts")},
	}, nil, ClientBrowser)

	if decision.Mode != ModeTranscode {
		t.Fatalf("streaming mode = %q, want transcode", decision.Mode)
	}
	if decision.Plan.VideoCodec != "libx264" || decision.Plan.AudioCodec != "aac" {
		t.Fatalf("plan = %#v, want h264/aac transcode", decision.Plan)
	}
	if !slices.Contains(decision.Reasons, "video_codec_not_supported") ||
		!slices.Contains(decision.Reasons, "audio_codec_not_supported") {
		t.Fatalf("reasons = %#v, want video and audio codec reasons", decision.Reasons)
	}
}

func TestSCNMedia012PreviewDecisionTranscodesVideoForWebKitHLS(t *testing.T) {
	track := int32(1)
	decision := DecisionFromTracks("/media/movie.mkv", []Track{
		{Type: TrackVideo, Codec: previewString("h264"), PixelFormat: previewString("yuv420p")},
		{Type: TrackAudio, Index: &track, Codec: previewString("aac")},
	}, &track, ClientWebKit)

	if decision.Mode != ModeTranscode || decision.DeliveryProtocol != ProtocolHLS {
		t.Fatalf("decision = %#v, want HLS transcode", decision)
	}
	if decision.Plan.VideoCodec != "libx264" || decision.Plan.AudioCodec != "copy" {
		t.Fatalf("plan = %#v, want video transcode and audio copy", decision.Plan)
	}
	if !slices.Contains(decision.Reasons, ReasonWebKitHLS) {
		t.Fatalf("reasons = %#v, want WebKit HLS reason", decision.Reasons)
	}
}

func TestSCNMedia012ProbeDurationIgnoresInvalidValues(t *testing.T) {
	if duration := OptionalDuration("5400.25"); duration == nil || *duration != 5400.25 {
		t.Fatalf("duration = %#v, want 5400.25", duration)
	}
	if duration := OptionalDuration("unknown"); duration != nil {
		t.Fatalf("duration = %#v, want nil", duration)
	}
	if duration := OptionalDuration("0"); duration != nil {
		t.Fatalf("duration = %#v, want nil", duration)
	}
}

func TestSCNMedia012ProbeContainerInfoReportsFormatFields(t *testing.T) {
	container := containerInfo(ffprobeFormat{
		BitRate:    "5500000",
		Format:     "matroska,webm",
		FormatName: "Matroska / WebM",
	})

	if *container.BitRate != "5500000" || *container.Format != "matroska,webm" ||
		*container.FormatName != "Matroska / WebM" {
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
