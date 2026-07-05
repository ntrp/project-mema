package httpapi

import (
	"net/http"
	"strings"
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
		Type:  MediaTypeMovie,
	}, http.StatusOK, &search)
	assertHasProviderResult(t, search.Results, "Example Movie")

	var autocomplete MediaGroupedSearchResponse
	client.doJSON(t, http.MethodGet, "/media/autocomplete?query=Example&includeLibrary=false", nil, http.StatusOK, &autocomplete)
	assertHasProviderGroup(t, autocomplete.Groups)

	query := "Example"
	mediaType := MediaTypeMovie
	limit := int32(5)
	var advanced MediaGroupedSearchResponse
	client.doJSON(t, http.MethodPost, "/media/advanced-search", MediaAdvancedSearchRequest{
		Query: &query,
		Type:  &mediaType,
		Limit: &limit,
	}, http.StatusOK, &advanced)
	assertHasProviderGroup(t, advanced.Groups)

	includeMedia := false
	includePeople := true
	var peopleSearch MediaGroupedSearchResponse
	client.doJSON(t, http.MethodPost, "/media/advanced-search", MediaAdvancedSearchRequest{
		Query:         &query,
		IncludeMedia:  &includeMedia,
		IncludePeople: &includePeople,
		Limit:         &limit,
	}, http.StatusOK, &peopleSearch)
	assertHasProviderPerson(t, peopleSearch.Groups, "Example Actor")

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
	if details.Title != "Example Movie" || details.ExternalId != "936075" || details.Type != MediaTypeMovie {
		t.Fatalf("metadata details = %#v", details)
	}
	if details.Crew == nil || len(*details.Crew) == 0 || (*details.Crew)[0].ExternalId == nil || *(*details.Crew)[0].ExternalId != "2001" {
		t.Fatalf("metadata crew links not mapped: %#v", details.Crew)
	}

	var collection MediaCollection
	client.doJSON(t, http.MethodGet, "/media/collections/tmdb/123", nil, http.StatusOK, &collection)
	if collection.Name != "Example Collection" || collection.Provider != Tmdb {
		t.Fatalf("metadata collection = %#v", collection)
	}
	assertHasProviderResult(t, collection.Results, "Example Movie")
}

func TestScenarioSCNMedia008FilteredMovieDiscoverUsesMetadataCache(t *testing.T) {
	provider := testmocks.NewProviderServer()
	t.Cleanup(provider.Close)
	client := newAcceptanceClientWithProviders(t, "SCN-MEDIA-008", provider)
	createScenarioMetadataProvider(t, client, provider.URL+"/tmdb/3")

	path := "/media/discover/movies/search?page=1&sort=popularity.desc&genres=Drama&runtimeMin=0&runtimeMax=400&scoreMin=0&scoreMax=10&minVoteCount=0"
	var first DiscoverMovieSearchResponse
	client.doJSON(t, http.MethodGet, path, nil, http.StatusOK, &first)
	assertHasProviderResult(t, first.Results, "Example Movie")

	var second DiscoverMovieSearchResponse
	client.doJSON(t, http.MethodGet, path, nil, http.StatusOK, &second)
	assertHasProviderResult(t, second.Results, "Example Movie")

	var cache MetadataCacheResponse
	client.doJSON(t, http.MethodGet, "/settings/metadata-cache?cacheLimit=20&historyLimit=20", nil, http.StatusOK, &cache)
	assertHasFilteredMovieCacheEntry(t, cache)
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

func assertHasProviderPerson(t *testing.T, groups []MediaSearchGroup, name string) {
	t.Helper()
	for _, group := range groups {
		if group.People == nil {
			continue
		}
		for _, person := range *group.People {
			if person.Name == name && person.ExternalProvider == "tmdb" {
				return
			}
		}
	}
	t.Fatalf("provider person %q not found: %#v", name, groups)
}

func assertHasFilteredMovieCacheEntry(t *testing.T, cache MetadataCacheResponse) {
	t.Helper()
	foundEntry := false
	for _, entry := range cache.Entries {
		if entry.CacheKind == MetadataCacheEntryCacheKindDiscover &&
			entry.MediaType == MediaTypeMovie &&
			strings.HasPrefix(entry.Query, "discover:movies:") &&
			entry.ItemCount > 0 {
			foundEntry = true
			break
		}
	}
	if !foundEntry {
		t.Fatalf("filtered movie discover cache entry not found: %#v", cache.Entries)
	}

	missFound := false
	hitFound := false
	for _, entry := range cache.HistoryEntries {
		if entry.CacheKind != MetadataSearchHistoryEntryCacheKindDiscover ||
			entry.MediaType != MediaTypeMovie ||
			!strings.HasPrefix(entry.Query, "discover:movies:") {
			continue
		}
		if entry.CacheHit {
			hitFound = true
		} else {
			missFound = true
		}
	}
	if !missFound || !hitFound {
		t.Fatalf("filtered movie discover history miss/hit not found: %#v", cache.HistoryEntries)
	}
}
