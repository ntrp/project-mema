package metadata

import (
	"strconv"
	"strings"
)

func tmdbDiscoverSeriesResults(items []tmdbMedia, genres map[string]string) []SearchResult {
	results := make([]SearchResult, 0, min(len(items), 20))
	for _, item := range items {
		result := tmdbDiscoverSeriesResult(item, genres)
		if result.Title == "" {
			continue
		}
		results = append(results, result)
		if len(results) >= 20 {
			break
		}
	}
	return results
}

func tmdbDiscoverSeriesResult(item tmdbMedia, genres map[string]string) SearchResult {
	return SearchResult{
		Title:            strings.TrimSpace(item.Name),
		Type:             "serie",
		Year:             yearFromDate(item.FirstAirDate),
		ExternalProvider: "tmdb",
		ExternalID:       strconv.FormatInt(item.ID, 10),
		Overview:         optionalString(item.Overview),
		PosterPath:       optionalString(item.PosterPath),
		BackdropPath:     optionalString(item.BackdropPath),
		Popularity:       optionalFloat64(item.Popularity),
		ReleaseDate:      optionalString(item.FirstAirDate),
		VoteAverage:      optionalFloat64(item.VoteAverage),
		VoteCount:        optionalInt32(item.VoteCount),
		OriginalLanguage: optionalString(item.Language),
		Genres:           genreNames(item.GenreIDs, genres),
	}
}

func seriesStatusIDs(input []string) []string {
	values := []string{}
	for _, status := range input {
		switch strings.ToLower(strings.TrimSpace(status)) {
		case "returning", "returning series", "0":
			values = append(values, "0")
		case "planned", "1":
			values = append(values, "1")
		case "in production", "production", "2":
			values = append(values, "2")
		case "ended", "3":
			values = append(values, "3")
		case "canceled", "cancelled", "4":
			values = append(values, "4")
		case "pilot", "5":
			values = append(values, "5")
		}
	}
	return values
}
