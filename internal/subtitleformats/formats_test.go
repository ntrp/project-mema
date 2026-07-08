package subtitleformats

import "testing"

func TestNormalizeSubtitleFormatAliases(t *testing.T) {
	cases := map[string]string{
		"subrip": "subrip",
		".srt":   "subrip",
		"webvtt": "vtt",
		"sup":    "pgs",
	}
	for value, want := range cases {
		if got := Normalize(value); got != want {
			t.Fatalf("Normalize(%q) = %q, want %q", value, got, want)
		}
	}
}

func TestAnyMatchUsesAliases(t *testing.T) {
	if !AnyMatch([]string{"srt"}, "subrip") {
		t.Fatalf("expected subrip to match srt")
	}
	if !AnyMatch([]string{"webvtt"}, "vtt") {
		t.Fatalf("expected vtt to match webvtt")
	}
}

func TestTextFormatsExcludeBitmapSubtitles(t *testing.T) {
	if !Text("subrip") || !Text("srt") || !Text("webvtt") || !Text("ass") || !Text("ssa") {
		t.Fatalf("expected text formats")
	}
	if Text("pgs") || Text("sup") {
		t.Fatalf("expected bitmap formats to be non-text")
	}
}
