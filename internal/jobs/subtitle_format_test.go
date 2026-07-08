package jobs

import (
	"strings"
	"testing"
)

func TestConvertSubtitleContentConvertsTextFormats(t *testing.T) {
	content := []byte("1\n00:00:01,000 --> 00:00:02,000\nHello\n")

	converted, format, err := convertSubtitleContent(content, "webvtt")

	if err != nil {
		t.Fatalf("convert subtitle: %v", err)
	}
	if format != "vtt" {
		t.Fatalf("format = %q", format)
	}
	if !strings.HasPrefix(string(converted), "WEBVTT") {
		t.Fatalf("converted = %q", converted)
	}
}

func TestConvertSubtitleContentRejectsBitmapTarget(t *testing.T) {
	_, _, err := convertSubtitleContent([]byte("1\n00:00:01,000 --> 00:00:02,000\nHello\n"), "sup")

	if err == nil || !strings.Contains(err.Error(), "pgs requires bitmap subtitle extraction") {
		t.Fatalf("err = %v", err)
	}
}
