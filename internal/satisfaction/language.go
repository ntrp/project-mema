package satisfaction

import "strings"

var languageAliases = map[string]string{
	"ar":         "ar",
	"ara":        "ar",
	"arabic":     "ar",
	"chi":        "zh",
	"chinese":    "zh",
	"da":         "da",
	"dan":        "da",
	"danish":     "da",
	"de":         "de",
	"deu":        "de",
	"dut":        "nl",
	"dutch":      "nl",
	"en":         "en",
	"eng":        "en",
	"english":    "en",
	"es":         "es",
	"fi":         "fi",
	"fin":        "fi",
	"finnish":    "fi",
	"fr":         "fr",
	"fra":        "fr",
	"fre":        "fr",
	"french":     "fr",
	"ger":        "de",
	"german":     "de",
	"hi":         "hi",
	"hin":        "hi",
	"hindi":      "hi",
	"it":         "it",
	"ita":        "it",
	"italian":    "it",
	"ja":         "ja",
	"japanese":   "ja",
	"jpn":        "ja",
	"ko":         "ko",
	"kor":        "ko",
	"korean":     "ko",
	"nl":         "nl",
	"nld":        "nl",
	"no":         "no",
	"nor":        "no",
	"norwegian":  "no",
	"pl":         "pl",
	"pol":        "pl",
	"polish":     "pl",
	"por":        "pt",
	"portuguese": "pt",
	"pt":         "pt",
	"ru":         "ru",
	"rus":        "ru",
	"russian":    "ru",
	"spa":        "es",
	"spanish":    "es",
	"sv":         "sv",
	"swe":        "sv",
	"swedish":    "sv",
	"zh":         "zh",
	"zho":        "zh",
}

func LanguageMatches(left string, right string) bool {
	return languageMatchKey(left) == languageMatchKey(right)
}

func languageMatchKey(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.TrimSuffix(normalized, " language")
	if normalized == "" || normalized == "-" {
		return ""
	}
	if alias, ok := languageAliases[normalized]; ok {
		return alias
	}
	return strings.ReplaceAll(normalized, " ", "-")
}
