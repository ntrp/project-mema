package httpapi

import (
	"net/http"
	"testing"

	"media-manager/internal/testmocks"
)

func TestScenarioSCNMedia008SignedInUsersSearchAndInspectProviderMetadata(t *testing.T) {
	provider := testmocks.NewProviderServer()
	t.Cleanup(provider.Close)
	client := newAcceptanceClientWithProviders(t, "SCN-MEDIA-008", provider)
	createScenarioMetadataProvider(t, client, provider.URL+"/tmdb/3")

	var search MediaSearchResponse
	client.doJSON(t, http.MethodPost, "/media/search", MediaSearchRequest{
		Query: "Example Movie",
		Type:  Movie,
	}, http.StatusOK, &search)
	assertHasProviderResult(t, search.Results, "Example Movie")

	var autocomplete MediaGroupedSearchResponse
	client.doJSON(t, http.MethodGet, "/media/autocomplete?query=Example&includeLibrary=false", nil, http.StatusOK, &autocomplete)
	assertHasProviderGroup(t, autocomplete.Groups)

	query := "Example"
	mediaType := Movie
	limit := int32(5)
	var advanced MediaGroupedSearchResponse
	client.doJSON(t, http.MethodPost, "/media/advanced-search", MediaAdvancedSearchRequest{
		Query: &query,
		Type:  &mediaType,
		Limit: &limit,
	}, http.StatusOK, &advanced)
	assertHasProviderGroup(t, advanced.Groups)

	var discover MediaDiscoverSection
	client.doJSON(t, http.MethodGet, "/media/discover/movie-popular?limit=2&page=1", nil, http.StatusOK, &discover)
	assertHasProviderResult(t, discover.Results, "Example Movie")

	var allDiscover MediaDiscoverResponse
	client.doJSON(t, http.MethodGet, "/media/discover", nil, http.StatusOK, &allDiscover)
	if len(allDiscover.Sections) == 0 {
		t.Fatalf("discover response has no sections: %#v", allDiscover)
	}

	var details MediaMetadataDetails
	client.doJSON(t, http.MethodGet, "/media/metadata/tmdb/movie/936075", nil, http.StatusOK, &details)
	if details.Title != "Example Movie" || details.ExternalId != "936075" || details.Type != Movie {
		t.Fatalf("metadata details = %#v", details)
	}

	var collection MediaCollection
	client.doJSON(t, http.MethodGet, "/media/collections/tmdb/123", nil, http.StatusOK, &collection)
	if collection.Name != "Example Collection" || collection.Provider != Tmdb {
		t.Fatalf("metadata collection = %#v", collection)
	}
	assertHasProviderResult(t, collection.Results, "Example Movie")
}

func createScenarioMetadataProvider(t *testing.T, client acceptanceClient, baseURL string) {
	t.Helper()
	var existing MetadataProviderListResponse
	client.doJSON(t, http.MethodGet, "/settings/metadata-providers", nil, http.StatusOK, &existing)
	for _, provider := range existing.Providers {
		if provider.Name == "Scenario TMDb" || provider.Name == "AAA Scenario TMDb" {
			client.doJSON(t, http.MethodDelete, "/settings/metadata-providers/"+provider.Id.String(), nil, http.StatusNoContent, nil)
		}
	}

	apiKey := "scenario-metadata-key"
	var created MetadataProvider
	client.doJSON(t, http.MethodPost, "/settings/metadata-providers", MetadataProviderRequest{
		Name:     "AAA Scenario TMDb",
		Type:     Tmdb,
		BaseUrl:  baseURL,
		ApiKey:   &apiKey,
		Enabled:  true,
		Priority: 0,
	}, http.StatusCreated, &created)
}

func assertHasProviderGroup(t *testing.T, groups []MediaSearchGroup) {
	t.Helper()
	for _, group := range groups {
		if group.SourceType == Provider {
			assertHasProviderResult(t, group.Results, "Example Movie")
			return
		}
	}
	t.Fatalf("provider group not found: %#v", groups)
}

func assertHasProviderResult(t *testing.T, results []MediaSearchResult, title string) {
	t.Helper()
	for _, result := range results {
		if result.Title == title && result.ExternalProvider != nil && *result.ExternalProvider == "tmdb" {
			return
		}
	}
	t.Fatalf("provider result %q not found: %#v", title, results)
}
