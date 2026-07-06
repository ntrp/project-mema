package decisions

import "media-manager/internal/storage"

func animeAbsoluteSearch(item storage.MediaItem, criteria ReleaseSearchCriteria) bool {
	return item.ContentKind == "anime" &&
		item.NumberingStrategy != nil &&
		*item.NumberingStrategy == "anidb_absolute" &&
		criteria.Kind == "episode" &&
		criteria.SeasonNumber == nil
}

func resourceTitleMatches(criteria ReleaseSearchCriteria, parsedTitle string, releaseTitle string) bool {
	candidateTitle := normalizedResourceTitle(parsedTitle)
	if candidateTitle == "" {
		candidateTitle = normalizedResourceTitle(releaseTitle)
	}
	for _, expected := range append([]string{criteria.Title}, criteria.Aliases...) {
		expectedTitle := normalizedResourceTitle(expected)
		if expectedTitle != "" && expectedTitle == candidateTitle {
			return true
		}
	}
	return false
}

func normalizedResourceTitle(title string) string {
	return normalizedToken(cleanReleaseResourceTitle(title))
}

func releaseAliasTexts(values []storage.MediaItemAlias) []string {
	aliases := []string{}
	for _, value := range values {
		if releaseAliasKind(value.Kind) {
			aliases = append(aliases, value.Alias)
		}
	}
	return aliases
}

func releaseAliasKind(kind string) bool {
	switch kind {
	case "release_title", "canonical", "english", "romaji", "synonym":
		return true
	default:
		return false
	}
}
