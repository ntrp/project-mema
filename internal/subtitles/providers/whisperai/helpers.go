package whisperai

import (
	"context"
	"strings"
	"time"

	"media-manager/internal/subtitles/providercore"
)

var alpha3ToAlpha2 = map[string]string{"eng": "en", "spa": "es", "deu": "de", "fra": "fr", "ita": "it", "por": "pt", "rus": "ru", "jpn": "ja", "kor": "ko", "zho": "zh", "und": ""}
var alpha2ToAlpha3 = map[string]string{"en": "eng", "es": "spa", "de": "deu", "fr": "fra", "it": "ita", "pt": "por", "ru": "rus", "ja": "jpn", "ko": "kor", "zh": "zho"}
var languageMapping = map[string]string{"gsw": "deu", "und": "eng"}
var ambiguous = map[string]bool{"alg": true, "art": true, "ath": true, "aus": true, "mis": true, "mul": true, "sgn": true, "und": true, "zxx": true}

func normalizeLang(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if len(s) == 2 {
		return alpha2ToAlpha3[s]
	}
	if mapped := languageMapping[s]; mapped != "" {
		return mapped
	}
	return s
}

func firstAudioLanguage(config providercore.Config) string {
	for _, language := range providercore.NewConfig(config).StringsSetting("audioLanguages") {
		if normalized := normalizeLang(language); normalized != "" {
			return normalized
		}
	}
	return normalizeLang(providercore.NewConfig(config).StringSetting("audioLanguage"))
}

func firstPath(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func iso6392(code string) string {
	if code == "deu" {
		return "ger"
	}
	if code == "fra" {
		return "fre"
	}
	return code
}

func timeoutContext(ctx context.Context, config providercore.Config, key string, fallback int) (context.Context, context.CancelFunc) {
	seconds := providercore.NewConfig(config).IntSetting(key)
	if seconds <= 0 {
		seconds = fallback
	}
	return context.WithTimeout(ctx, time.Duration(seconds)*time.Second)
}
