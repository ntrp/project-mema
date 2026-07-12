package providers

import "testing"

func TestAlpha3LanguageMappings(t *testing.T) {
	cases := map[string]string{
		"Arabic":               "ara",
		"bg":                   "bul",
		"czech":                "ces",
		"dan":                  "dan",
		"German":               "deu",
		"greek":                "ell",
		"English":              "eng",
		"Finnish":              "fin",
		"French":               "fra",
		"Hebrew":               "heb",
		"Hungarian":            "hun",
		"Indonesian":           "ind",
		"Italian":              "ita",
		"Japanese":             "jpn",
		"Korean":               "kor",
		"Dutch":                "nld",
		"Polish":               "pol",
		"Brazilian Portuguese": "pob",
		"Portuguese":           "por",
		"Romanian":             "ron",
		"Russian":              "rus",
		"Spanish":              "spa",
		"Swedish":              "swe",
		"Thai":                 "tha",
		"Turkish":              "tur",
		"Ukrainian":            "ukr",
		"Vietnamese":           "vie",
		"Chinese":              "zho",
		"custom":               "custom",
	}
	for input, want := range cases {
		if got := alpha3Language(input); got != want {
			t.Fatalf("alpha3Language(%q) = %q, want %q", input, got, want)
		}
	}
}
