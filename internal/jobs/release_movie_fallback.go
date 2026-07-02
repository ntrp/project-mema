package jobs

import (
	"regexp"
	"strconv"
	"strings"

	"media-manager/internal/storage"
)

func movieYearFallbackQuery(item storage.MediaItem, query string) string {
	if item.Type != "movie" || item.Year == nil {
		return ""
	}
	year := strconv.Itoa(int(*item.Year))
	withoutYear := yearTokenPattern(year).ReplaceAllString(query, " ")
	withoutYear = strings.Join(strings.Fields(withoutYear), " ")
	if withoutYear == strings.TrimSpace(query) {
		return ""
	}
	return withoutYear
}

func yearTokenPattern(year string) *regexp.Regexp {
	return regexp.MustCompile(`(^|[\s._-])` + regexp.QuoteMeta(year) + `($|[\s._-])`)
}
