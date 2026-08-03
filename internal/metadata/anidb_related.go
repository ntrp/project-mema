package metadata

import (
	"strconv"
	"strings"
)

type anidbRelated struct {
	ID    int    `xml:"id,attr"`
	Type  string `xml:"type,attr"`
	Title string `xml:",chardata"`
}

type anidbSimilar struct {
	ID       int    `xml:"id,attr"`
	Approval int32  `xml:"approval,attr"`
	Total    int32  `xml:"total,attr"`
	Title    string `xml:",chardata"`
}

type anidbRating struct {
	Count int32   `xml:"count,attr"`
	Value float64 `xml:",chardata"`
}

type anidbRatings struct {
	Permanent anidbRating `xml:"permanent"`
	Temporary anidbRating `xml:"temporary"`
	Review    anidbRating `xml:"review"`
}

func anidbMediaType(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "movie") {
		return "movie"
	}
	return "serie"
}

func anidbRelatedResults(values []anidbRelated, mediaType string) []SearchResult {
	results := make([]SearchResult, 0, len(values))
	for _, value := range values {
		if result, ok := anidbRelatedResult(value.ID, value.Title, mediaType); ok {
			results = append(results, result)
		}
	}
	return results
}

func anidbSimilarResults(values []anidbSimilar, mediaType string) []SearchResult {
	results := make([]SearchResult, 0, len(values))
	for _, value := range values {
		if result, ok := anidbRelatedResult(value.ID, value.Title, mediaType); ok {
			results = append(results, result)
		}
	}
	return results
}

func anidbRelatedResult(id int, title, mediaType string) (SearchResult, bool) {
	title = strings.TrimSpace(title)
	if id <= 0 || title == "" {
		return SearchResult{}, false
	}
	externalID := strconv.Itoa(id)
	return SearchResult{
		Title:            title,
		Type:             mediaType,
		ExternalProvider: "anidb",
		ExternalID:       externalID,
		ExternalURL:      optionalString("https://anidb.net/a" + externalID),
	}, true
}

func anidbVoteAverage(ratings anidbRatings) *float64 {
	for _, rating := range []anidbRating{ratings.Permanent, ratings.Temporary, ratings.Review} {
		if rating.Value > 0 {
			value := rating.Value
			return &value
		}
	}
	return nil
}
