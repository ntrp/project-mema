package decisions

import "media-manager/internal/storage"

func applyLanguageCatalog(parsed ParsedRelease, catalog []storage.Language) ParsedRelease {
	if len(catalog) == 0 {
		return parsed
	}
	for _, language := range catalog {
		if languageMatchesTitle(parsed.ReleaseTitle, language) {
			parsed.Languages = appendLanguage(parsed.Languages, language.DisplayName)
			parsed.Languages = appendLanguage(parsed.Languages, language.Code)
		}
	}
	return parsed
}

func languageMatchesTitle(title string, language storage.Language) bool {
	tokens := tokenSet(title)
	for _, alias := range language.Aliases {
		token := normalizedToken(alias)
		if token == "" {
			continue
		}
		if len(token) <= 3 {
			if _, ok := tokens[token]; ok {
				return true
			}
			continue
		}
		if containsAnyNormalized(title, alias) {
			return true
		}
	}
	return false
}

func appendLanguage(values []string, value string) []string {
	if value == "" || containsString(values, value) {
		return values
	}
	return append(values, value)
}
