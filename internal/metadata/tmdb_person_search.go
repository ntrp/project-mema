package metadata

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (s *Service) SearchPeople(ctx context.Context, config Config, query string) ([]PersonSearchResult, error) {
	if config.Type != "tmdb" {
		return []PersonSearchResult{}, nil
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return []PersonSearchResult{}, nil
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "search", "person")
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("query", query)
	endpoint = endpoint + "?" + values.Encode()

	var payload tmdbPersonSearchResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	return tmdbPersonSearchResults(payload.Results, 20), nil
}

func tmdbPersonSearchResults(items []tmdbPersonSearchResult, limit int) []PersonSearchResult {
	results := make([]PersonSearchResult, 0, min(len(items), limit))
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		results = append(results, PersonSearchResult{
			Name:             name,
			ExternalProvider: "tmdb",
			ExternalID:       strconv.FormatInt(item.ID, 10),
			ProfilePath:      optionalString(item.ProfilePath),
			Popularity:       optionalFloat64(item.Popularity),
			KnownFor:         tmdbPersonKnownFor(item.KnownFor),
		})
		if len(results) >= limit {
			break
		}
	}
	return results
}

func tmdbPersonKnownFor(items []tmdbMedia) []string {
	values := []string{}
	for _, item := range items {
		title := strings.TrimSpace(item.Title)
		if title == "" {
			title = strings.TrimSpace(item.Name)
		}
		if title != "" {
			values = append(values, title)
		}
	}
	return values
}
