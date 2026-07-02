package httpapi

import (
	"strings"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func parsedReleaseMetadataResponse(
	parsed decisions.ParsedRelease,
	languages []storage.Language,
) ParsedReleaseMetadata {
	return ParsedReleaseMetadata{
		Release:   parsedReleaseInfoResponse(parsed),
		Quality:   parsedQualityInfoResponse(parsed),
		Languages: parsedLanguageDisplayNames(parsed.Languages, languages),
		Details: ParsedReleaseDetails{
			ReleaseType:       parsed.ReleaseType,
			CustomFormatNames: []string{},
			MatchedSpecCount:  0,
		},
	}
}

func parsedLanguageDisplayNames(values []string, languages []storage.Language) []string {
	displayNames := make([]string, 0, len(values))
	for _, value := range values {
		displayNames = appendStringOnce(displayNames, languageDisplayName(value, languages))
	}
	return displayNames
}

func languageDisplayName(value string, languages []storage.Language) string {
	normalized := normalizedLanguageValue(value)
	for _, language := range languages {
		if normalized == normalizedLanguageValue(language.Code) ||
			normalized == normalizedLanguageValue(language.DisplayName) {
			return language.DisplayName
		}
		for _, alias := range language.Aliases {
			if normalized == normalizedLanguageValue(alias) {
				return language.DisplayName
			}
		}
	}
	return strings.TrimSpace(value)
}

func normalizedLanguageValue(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), "-"))
}

func appendStringOnce(values []string, value string) []string {
	if value == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func parsedReleaseInfoResponse(parsed decisions.ParsedRelease) ParsedReleaseInfo {
	return ParsedReleaseInfo{
		ReleaseTitle:  parsed.ReleaseTitle,
		MovieTitle:    parsed.MovieTitle,
		SeriesTitle:   parsed.SeriesTitle,
		Year:          parsed.Year,
		SeasonNumber:  parsed.SeasonNumber,
		EpisodeNumber: parsed.EpisodeNumber,
		SeasonPack:    parsed.SeasonPack,
		Edition:       parsed.Edition,
		ReleaseGroup:  parsed.ReleaseGroup,
		ReleaseHash:   parsed.ReleaseHash,
	}
}

func parsedQualityInfoResponse(parsed decisions.ParsedRelease) ParsedQualityInfo {
	return ParsedQualityInfo{
		QualityId:     parsed.QualityID,
		Quality:       parsed.Quality,
		Source:        parsed.Source,
		Resolution:    parsed.Resolution,
		VideoCodec:    parsed.VideoCodec,
		AudioCodec:    parsed.AudioCodec,
		AudioChannels: parsed.AudioChannels,
		Version:       parsed.Version,
		Proper:        parsed.Proper,
		Repack:        parsed.Repack,
		Real:          parsed.Real,
	}
}
