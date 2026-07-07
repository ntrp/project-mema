package storage

import (
	"strconv"
	"strings"
)

func renderSeriesTemplate(template string, item MediaItem, season int32, episode int32) string {
	return renderSeriesTemplateWithQuality(template, item, season, episode, "")
}

func renderSeriesTemplateWithQuality(
	template string,
	item MediaItem,
	season int32,
	episode int32,
	qualityFull string,
) string {
	values := map[string]string{
		"series_title":  item.Title,
		"season":        strconv.Itoa(int(season)),
		"episode":       strconv.Itoa(int(episode)),
		"episode_title": episodeTitle(item, season, episode),
		"quality_full":  strings.TrimSpace(qualityFull),
		"air_date":      episodeAirDate(item, season, episode),
	}
	rendered := fileNamingTokenPattern.ReplaceAllStringFunc(template, func(token string) string {
		key, format := splitTemplateToken(strings.Trim(token, "{}"))
		value, ok := values[key]
		if !ok {
			return token
		}
		return formatTemplateValue(value, format)
	})
	return strings.Join(strings.Fields(rendered), " ")
}

func splitTemplateToken(raw string) (string, string) {
	parts := strings.SplitN(raw, ":", 2)
	key := normalizeTemplateTokenName(parts[0])
	if len(parts) == 2 {
		return key, parts[1]
	}
	return key, ""
}

func formatTemplateValue(value string, format string) string {
	if format == "" || !strings.HasPrefix(format, "0") {
		return value
	}
	if len(value) >= len(format) {
		return value
	}
	return strings.Repeat("0", len(format)-len(value)) + value
}

func episodeTitle(item MediaItem, season int32, episode int32) string {
	if matched := findEpisode(item, season, episode); matched != nil {
		return matched.Name
	}
	return ""
}

func episodeAirDate(item MediaItem, season int32, episode int32) string {
	if matched := findEpisode(item, season, episode); matched != nil && matched.AirDate != nil {
		return *matched.AirDate
	}
	return ""
}

func findEpisode(item MediaItem, season int32, episode int32) *MediaEpisode {
	for _, itemSeason := range item.Seasons {
		if itemSeason.SeasonNumber != season {
			continue
		}
		for index := range itemSeason.Episodes {
			if itemSeason.Episodes[index].EpisodeNumber == episode {
				return &itemSeason.Episodes[index]
			}
		}
	}
	return nil
}
