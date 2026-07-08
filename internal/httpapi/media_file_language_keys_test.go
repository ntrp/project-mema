package httpapi

import "testing"

func TestMediaFileLanguageMatchKeyUsesIsoAliases(t *testing.T) {
	tests := map[string]string{
		" English Language ": "en",
		"eng":                "en",
		"German":             "de",
		"ita":                "it",
		"italian":            "it",
		"jpn":                "ja",
		"Korean":             "ko",
		"spa":                "es",
	}

	for value, want := range tests {
		if got := languageMatchKey(value); got != want {
			t.Fatalf("languageMatchKey(%q) = %q, want %q", value, got, want)
		}
	}
}
