package httpapi

import (
	"slices"
	"testing"
)

func TestSCNMedia012PreviewArgsSelectRequestedAudioStream(t *testing.T) {
	track := int32(2)
	args := mediaPreviewArgs("/media/movie.mkv", &track)

	if !slices.Contains(args, "0:2") {
		t.Fatalf("expected ffmpeg args to map requested audio stream, got %#v", args)
	}
	if !slices.Contains(args, "aac") {
		t.Fatalf("expected ffmpeg args to encode browser-compatible audio, got %#v", args)
	}
	if !slices.Contains(args, "frag_keyframe+empty_moov+default_base_moof") {
		t.Fatalf("expected fragmented MP4 flags, got %#v", args)
	}
}

func TestSCNMedia012PreviewArgsFallsBackToFirstAudioStream(t *testing.T) {
	args := mediaPreviewArgs("/media/movie.mkv", nil)

	if !slices.Contains(args, "0:a:0?") {
		t.Fatalf("expected optional first audio stream mapping, got %#v", args)
	}
}
