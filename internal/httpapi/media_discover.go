package httpapi

import (
	"context"

	"media-manager/internal/storage"
)

func providerHasCredentials(provider storage.MetadataProvider) bool {
	return optionalTrimmedString(provider.APIKey) != nil || optionalTrimmedString(provider.AccessToken) != nil
}

func discoverProvider(providers []storage.MetadataProvider) (storage.MetadataProvider, bool) {
	for _, provider := range providers {
		if provider.Enabled && provider.Type == "tmdb" && providerHasCredentials(provider) {
			return provider, true
		}
	}
	return storage.MetadataProvider{}, false
}

func metadataProviderByType(providers []storage.MetadataProvider, providerType string) (storage.MetadataProvider, bool) {
	for _, provider := range providers {
		if provider.Enabled && provider.Type == providerType && providerHasCredentials(provider) {
			return provider, true
		}
	}
	return storage.MetadataProvider{}, false
}

func (s *Server) discoverSectionResponse(
	ctx context.Context,
	providers []storage.MetadataProvider,
	section discoverSection,
	limit int,
	page int,
	blacklist []storage.DiscoverBlacklistItem,
) MediaDiscoverSection {
	providerName := "TMDB"
	results := []MediaSearchResult{}
	if provider, ok := discoverProvider(providers); ok {
		providerName = provider.Name
		for _, request := range section.requests {
			providerResults, err := s.discoverMetadataProvider(ctx, provider, request.mediaType, request.id, limit, page)
			if err != nil {
				continue
			}
			for _, result := range providerResults {
				results = append(results, metadataSearchResultResponse(result))
			}
		}
	}
	return MediaDiscoverSection{
		Id:           section.responseID,
		Title:        section.title,
		ProviderName: providerName,
		MediaType:    MediaDiscoverMediaType(section.mediaType),
		Results:      filterDiscoverBlacklist(dedupeMediaSearchResults(results), blacklist),
	}
}

func dedupeMediaSearchResults(results []MediaSearchResult) []MediaSearchResult {
	seen := map[string]struct{}{}
	deduped := make([]MediaSearchResult, 0, len(results))
	for _, result := range results {
		key := string(result.Type) + ":" + valueOrEmpty(result.ExternalProvider) + ":" + valueOrEmpty(result.ExternalId)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		deduped = append(deduped, result)
	}
	return deduped
}

type discoverSection struct {
	responseID string
	title      string
	mediaType  string
	requests   []discoverSectionRequest
}

type discoverSectionRequest struct {
	mediaType string
	id        string
}

var discoverSections = []discoverSection{
	{responseID: "trending", title: "Trending", mediaType: "mixed", requests: []discoverSectionRequest{
		{mediaType: "mixed", id: "trending"},
	}},
	{responseID: "movie-popular", title: "Popular Movies", mediaType: "movie", requests: []discoverSectionRequest{{mediaType: "movie", id: "popular"}}},
	{responseID: "movie-upcoming", title: "Upcoming Movies", mediaType: "movie", requests: []discoverSectionRequest{{mediaType: "movie", id: "upcoming"}}},
	{responseID: "movie-top-rated", title: "Top Rated Movies", mediaType: "movie", requests: []discoverSectionRequest{{mediaType: "movie", id: "top_rated"}}},
	{responseID: "series-popular", title: "Popular Series", mediaType: "series", requests: []discoverSectionRequest{{mediaType: "series", id: "popular"}}},
	{responseID: "series-on-the-air", title: "Airing Series", mediaType: "series", requests: []discoverSectionRequest{{mediaType: "series", id: "on_the_air"}}},
	{responseID: "series-top-rated", title: "Top Rated Series", mediaType: "series", requests: []discoverSectionRequest{{mediaType: "series", id: "top_rated"}}},
}

func discoverSectionByID(id string) (discoverSection, bool) {
	for _, section := range discoverSections {
		if section.responseID == id {
			return section, true
		}
	}
	return discoverSection{}, false
}
