package storage

import "testing"

func TestNormalizeTemplateConvertsLegacyTokensToSnakeCase(t *testing.T) {
	got := normalizeTemplate("{Movie Title} ({Release Year}) {Quality Full}")
	want := "{movie_title} ({release_year}) {quality_full}"
	if got != want {
		t.Fatalf("normalizeTemplate() = %q, want %q", got, want)
	}
}

func TestRenderMediaTemplateSupportsSnakeCaseAndLegacyTokens(t *testing.T) {
	year := int32(2024)
	input := MediaItemInput{Title: "Example Movie", Year: &year}

	for _, template := range []string{
		"{movie_title} ({release_year})",
		"{Movie Title} ({Release Year})",
	} {
		got := renderMediaTemplate(template, input)
		want := "Example Movie (2024)"
		if got != want {
			t.Fatalf("renderMediaTemplate(%q) = %q, want %q", template, got, want)
		}
	}
}
