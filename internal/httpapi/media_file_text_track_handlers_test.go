package httpapi

import (
	"slices"
	"testing"
)

func TestSCNMedia013SubtitlePreviewArgsSelectRequestedStream(t *testing.T) {
	args := mediaSubtitlePreviewArgs("/media/movie.mkv", 4)

	if !slices.Contains(args, "0:4") {
		t.Fatalf("expected ffmpeg args to map requested subtitle stream, got %#v", args)
	}
	if !hasArgPair(args, "-f", "webvtt") {
		t.Fatalf("expected ffmpeg args to output WebVTT, got %#v", args)
	}
	if args[len(args)-1] != "pipe:1" {
		t.Fatalf("expected ffmpeg args to stream to stdout, got %#v", args)
	}
}
