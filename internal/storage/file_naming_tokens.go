package storage

import "strings"

var legacyFileNamingTokens = map[string]string{
	"Air-Date":      "air_date",
	"Episode Title": "episode_title",
	"Movie Title":   "movie_title",
	"Quality Full":  "quality_full",
	"Release Year":  "release_year",
	"Series Title":  "series_title",
	"Year":          "year",
}

func normalizeTemplateTokens(template string) string {
	return fileNamingTokenPattern.ReplaceAllStringFunc(template, func(token string) string {
		key := strings.Trim(token, "{}")
		return "{" + normalizeTemplateTokenName(key) + "}"
	})
}

func normalizeTemplateTokenName(key string) string {
	if normalized, ok := legacyFileNamingTokens[key]; ok {
		return normalized
	}
	return strings.ToLower(strings.NewReplacer(" ", "_", "-", "_").Replace(key))
}
