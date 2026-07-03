package storage

import "testing"

func TestNormalizeTagNames(t *testing.T) {
	tags := normalizeTagNames([]string{
		"  Anime  ",
		"anime",
		"",
		"  4K   Preferred ",
		"4k preferred",
		"Documentary",
	})

	expectStrings(t, tags, []string{"Anime", "4K Preferred", "Documentary"})
}

func TestNormalizeTagName(t *testing.T) {
	if got := normalizeTagName("  Family   Movies  "); got != "Family Movies" {
		t.Fatalf("expected compacted tag name, got %q", got)
	}
}
