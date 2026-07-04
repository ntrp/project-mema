package playback

import (
	"slices"
	"testing"
)

func TestSCNPlayback001DirectPlaysCompatibleBrowserMp4(t *testing.T) {
	info := BuildVideoStreamInfo(browserMovie("mp4", "h264", "aac", "yuv420p"), MediaOptions{})

	if info.PlayMethod != PlayMethodDirectPlay {
		t.Fatalf("play method = %q, want direct play", info.PlayMethod)
	}
	if info.Protocol != ProtocolHTTP || info.OutputVideoCodec != "copy" || info.OutputAudioCodec != "copy" {
		t.Fatalf("stream info = %#v, want direct HTTP stream copy", info)
	}
}

func TestSCNPlayback001DirectStreamsCompatibleCodecsInUnsupportedContainer(t *testing.T) {
	info := BuildVideoStreamInfo(browserMovie("mkv", "h264", "aac", "yuv420p"), MediaOptions{})

	if info.PlayMethod != PlayMethodDirectStream || info.Protocol != ProtocolHLS {
		t.Fatalf("stream info = %#v, want HLS direct stream", info)
	}
	if info.OutputVideoCodec != "copy" || info.OutputAudioCodec != "copy" {
		t.Fatalf("stream info = %#v, want stream copy", info)
	}
	if !slices.Contains(info.TranscodeReasons, ReasonContainerNotSupported) {
		t.Fatalf("reasons = %#v, want container reason", info.TranscodeReasons)
	}
}

func TestSCNPlayback001TranscodesOnlyUnsupportedAudio(t *testing.T) {
	info := BuildVideoStreamInfo(browserMovie("mkv", "h264", "dts", "yuv420p"), MediaOptions{})

	if info.PlayMethod != PlayMethodTranscode {
		t.Fatalf("play method = %q, want transcode", info.PlayMethod)
	}
	if info.OutputVideoCodec != "copy" || info.OutputAudioCodec != "aac" {
		t.Fatalf("stream info = %#v, want video copy and audio transcode", info)
	}
}

func TestSCNPlayback001TranscodesUnsupportedVideoPixelFormat(t *testing.T) {
	info := BuildVideoStreamInfo(browserMovie("mp4", "h264", "aac", "yuv420p10le"), MediaOptions{})

	if info.PlayMethod != PlayMethodTranscode || info.OutputVideoCodec != "h264" {
		t.Fatalf("stream info = %#v, want video transcode", info)
	}
	if !slices.Contains(info.TranscodeReasons, ReasonVideoPixelFormatUnsupported) {
		t.Fatalf("reasons = %#v, want pixel format reason", info.TranscodeReasons)
	}
}

func TestSCNPlayback001DirectStreamsSelectedSecondaryAudio(t *testing.T) {
	index := 2
	source := browserMovie("mp4", "h264", "aac", "yuv420p")
	source.AudioStreams = append(source.AudioStreams, MediaStream{Index: index, Type: StreamAudio, Codec: "aac"})
	info := BuildVideoStreamInfo(source, MediaOptions{AudioStreamIndex: &index})

	if info.PlayMethod != PlayMethodDirectStream || info.OutputAudioCodec != "copy" {
		t.Fatalf("stream info = %#v, want HLS stream copy for selected secondary audio", info)
	}
	if !slices.Contains(info.TranscodeReasons, ReasonSecondaryAudioNotSupported) {
		t.Fatalf("reasons = %#v, want secondary audio reason", info.TranscodeReasons)
	}
}

func browserMovie(container, videoCodec, audioCodec, pixelFormat string) MediaSource {
	return MediaSource{
		Container:            container,
		Video:                &MediaStream{Index: 0, Type: StreamVideo, Codec: videoCodec, PixelFormat: pixelFormat},
		AudioStreams:         []MediaStream{{Index: 1, Type: StreamAudio, Codec: audioCodec}},
		SupportsDirectPlay:   true,
		SupportsDirectStream: true,
		SupportsTranscoding:  true,
	}
}
