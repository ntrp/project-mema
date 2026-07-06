package decisions

import (
	"regexp"
	"strings"
)

var (
	seasonEpisodePattern   = regexp.MustCompile(`(?i)\bs(\d{1,2})(?:e(\d{1,3}))?\b`)
	absoluteEpisodePattern = regexp.MustCompile(`(?i)(?:^|[\s._-])(?:e|ep|episode)?[\s._-]?(\d{1,3})(?:v\d+)?(?:[\s._-]|$)`)
)

func releaseSeriesTitle(title string) string {
	index := len(title)
	if match := seasonEpisodePattern.FindStringIndex(title); match != nil && match[0] > 0 {
		index = match[0]
	} else if match := absoluteEpisodePattern.FindStringIndex(title); match != nil && match[0] > 0 {
		index = match[0]
	} else if year := yearPattern.FindStringIndex(title); year != nil && year[0] > 0 {
		index = year[0]
	}
	return cleanReleaseResourceTitle(title[:index])
}

func detectSeasonEpisode(title string) (*int32, *int32) {
	match := seasonEpisodePattern.FindStringSubmatch(title)
	if len(match) == 0 {
		return nil, nil
	}
	season := parsePositiveInt32(match[1])
	if season == nil {
		return nil, nil
	}
	if len(match) < 3 || match[2] == "" {
		return season, nil
	}
	return season, parsePositiveInt32(match[2])
}

func detectAbsoluteEpisode(title string) *int32 {
	match := absoluteEpisodePattern.FindStringSubmatch(title)
	if len(match) < 2 {
		return nil
	}
	return parsePositiveInt32(match[1])
}

func cleanReleaseResourceTitle(title string) string {
	return strings.Join(strings.Fields(releaseSeparator.ReplaceAllString(title, " ")), " ")
}

func parsePositiveInt32(value string) *int32 {
	var parsed int32
	for _, char := range value {
		if char < '0' || char > '9' {
			return nil
		}
		parsed = parsed*10 + int32(char-'0')
	}
	if parsed <= 0 {
		return nil
	}
	return &parsed
}
