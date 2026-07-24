package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestAniListSearchUsesRecordedMultilingualFixture(t *testing.T) {
	service := fixtureAniListService(t, "search_multilingual.json", http.StatusOK, nil)
	results, err := service.Search(t.Context(), Config{Type: "anilist", BaseURL: "https://anilist.test/graphql"}, SearchRequest{
		Query: "Spirited Away", MediaType: "movie",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("results = %#v", results)
	}
	result := results[0]
	if result.Title != "Spirited Away" || result.ExternalID != "199" || result.ReleaseDate == nil || *result.ReleaseDate != "2001-07-20" {
		t.Fatalf("result = %#v", result)
	}
	if result.PosterPath == nil || *result.PosterPath != "https://img.test/199-xl.jpg" || len(result.Aliases) != 5 {
		t.Fatalf("multilingual result = %#v", result)
	}
	if result.Aliases[0].Kind != "romaji" || result.Aliases[0].Language == nil || *result.Aliases[0].Language != "ja-Latn" {
		t.Fatalf("romaji alias = %#v", result.Aliases[0])
	}
}

func TestAniListDetailsNormalizesAliasesMappingsAndRawProvenance(t *testing.T) {
	service := fixtureAniListService(t, "details_multilingual.json", http.StatusOK, nil)
	details, err := service.Details(t.Context(), Config{Type: "anilist", BaseURL: "https://anilist.test/graphql"}, DetailsRequest{
		MediaType: "serie", ExternalID: "154587",
	})
	if err != nil {
		t.Fatal(err)
	}
	if details.Type != "serie" || details.Title != "Frieren: Beyond Journey's End" || len(details.Aliases) != 5 {
		t.Fatalf("details = %#v", details)
	}
	if details.FirstAirDate == nil || *details.FirstAirDate != "2023-09-29" || details.RuntimeMinutes == nil || *details.RuntimeMinutes != 24 {
		t.Fatalf("dates/runtime = %#v", details)
	}
	if len(details.ProviderMappings) != 3 || details.ProviderMappings[0].ProviderName != "anilist" || details.ProviderMappings[1].ProviderName != "anidb" || details.ProviderMappings[2].ProviderName != "tmdb" {
		t.Fatalf("mappings = %#v", details.ProviderMappings)
	}
	raw, ok := details.ProviderMappings[0].Source["raw"].(map[string]any)
	if !ok || raw["idMal"] != float64(52991) {
		t.Fatalf("raw provenance = %#v", details.ProviderMappings[0].Source)
	}
}

func TestAniListPartialFixtureOmitsUnavailableValues(t *testing.T) {
	service := fixtureAniListService(t, "details_partial.json", http.StatusOK, nil)
	details, err := service.Details(t.Context(), Config{Type: "anilist", BaseURL: "https://anilist.test/graphql"}, DetailsRequest{
		MediaType: "serie", ExternalID: "170000",
	})
	if err != nil {
		t.Fatal(err)
	}
	if details.Year == nil || *details.Year != 2026 || details.ReleaseDate != nil || details.FirstAirDate != nil || details.PosterPath != nil || details.RuntimeMinutes != nil {
		t.Fatalf("partial details fabricated values: %#v", details)
	}
	if len(details.ProviderMappings) != 1 || len(details.Aliases) != 2 {
		t.Fatalf("partial details mappings/aliases = %#v", details)
	}
}

func TestAniListExternalMappingsRequireProviderOwnedHosts(t *testing.T) {
	tests := []struct {
		name, site, rawURL, provider, entityType, externalID string
	}{
		{"AniDB canonical", "ignored", "https://anidb.net/anime/17617", "anidb", "serie", "17617"},
		{"AniDB www", "ignored", "https://www.anidb.net/anime/17617", "anidb", "serie", "17617"},
		{"TMDB canonical", "ignored", "https://themoviedb.org/movie/129", "tmdb", "movie", "129"},
		{"TMDB www", "ignored", "https://www.themoviedb.org/tv/209867", "tmdb", "serie", "209867"},
		{"TVDB canonical", "ignored", "https://thetvdb.com/series/123", "tvdb", "serie", "123"},
		{"TVDB www", "ignored", "https://www.thetvdb.com/series/123", "tvdb", "serie", "123"},
		{"malicious site label", "AniDB", "https://example.com/anime/17617", "", "", ""},
		{"lookalike host", "AniDB", "https://evil-anidb.example/anime/17617", "", "", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			provider, entityType, externalID := supportedExternalLink(test.site, test.rawURL, "serie")
			if provider != test.provider || entityType != test.entityType || externalID != test.externalID {
				t.Fatalf("mapping = %q/%q/%q", provider, entityType, externalID)
			}
		})
	}
}

func TestAniListDoesNotPersistMappingsFromSpoofedExternalLinks(t *testing.T) {
	var media anilistMedia
	if err := json.Unmarshal([]byte(`{
		"id": 1,
		"externalLinks": [
			{"site":"AniDB","url":"https://example.com/anime/17617"},
			{"site":"Not AniDB","url":"https://evil-anidb.example/anime/17617"}
		]
	}`), &media); err != nil {
		t.Fatal(err)
	}
	mappings := anilistMappings(media, "serie")
	if len(mappings) != 1 || mappings[0].ProviderName != "anilist" {
		t.Fatalf("spoofed mappings persisted: %#v", mappings)
	}
}

func TestAniListFailuresAreTyped(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		status     int
		headers    http.Header
		assertType func(error) bool
	}{
		{"graphql", `{"errors":[{"message":"bad query"}],"data":null}`, http.StatusOK, nil, func(err error) bool { var target ProviderGraphQLError; return errors.As(err, &target) }},
		{"malformed", `{"data":`, http.StatusOK, nil, func(err error) bool { var target ProviderPayloadError; return errors.As(err, &target) }},
		{"empty", `{"data":{"Media":null}}`, http.StatusOK, nil, func(err error) bool { var target ProviderPayloadError; return errors.As(err, &target) }},
		{"http", `{}`, http.StatusBadGateway, nil, func(err error) bool { var target ProviderHTTPError; return errors.As(err, &target) }},
		{"rate limit", `{}`, http.StatusTooManyRequests, http.Header{"Retry-After": []string{"10"}}, IsRateLimited},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := responseAniListService(test.body, test.status, test.headers)
			_, err := service.Details(t.Context(), Config{Type: "anilist", BaseURL: "https://anilist.test/graphql"}, DetailsRequest{MediaType: "serie", ExternalID: "1"})
			if err == nil || !test.assertType(err) {
				t.Fatalf("error = %T %v", err, err)
			}
		})
	}
}

func TestAniListMissingSearchDataIsTypedPayloadFailure(t *testing.T) {
	for name, body := range map[string]string{
		"null data":  `{"data":null}`,
		"empty data": `{"data":{}}`,
		"null page":  `{"data":{"Page":null}}`,
	} {
		t.Run(name, func(t *testing.T) {
			service := responseAniListService(body, http.StatusOK, nil)
			_, err := service.Search(t.Context(), Config{Type: "anilist"}, SearchRequest{Query: "missing", MediaType: "serie"})
			var payloadErr ProviderPayloadError
			if !errors.As(err, &payloadErr) {
				t.Fatalf("error = %T %v", err, err)
			}
		})
	}

	service := responseAniListService(`{"data":{"Page":{"media":[]}}}`, http.StatusOK, nil)
	results, err := service.Search(t.Context(), Config{Type: "anilist"}, SearchRequest{Query: "no match", MediaType: "serie"})
	if err != nil || len(results) != 0 {
		t.Fatalf("valid empty page = %#v, %v", results, err)
	}
}

func TestAniListServiceRegistrationAndTypedErrorMessages(t *testing.T) {
	service := fixtureAniListService(t, "search_multilingual.json", http.StatusOK, nil)
	for _, providerType := range []string{"tmdb", "tvdb", "anilist"} {
		if _, err := service.provider(Config{Type: providerType}); err != nil {
			t.Fatalf("provider %q: %v", providerType, err)
		}
	}
	if _, err := service.provider(Config{Type: "unknown"}); !errors.Is(err, ErrUnsupportedProvider) {
		t.Fatalf("unsupported provider error = %v", err)
	}
	if _, err := service.Discover(t.Context(), Config{Type: "anilist"}, DiscoverRequest{}); !errors.Is(err, ErrUnsupportedProvider) {
		t.Fatalf("AniList discover error = %v", err)
	}
	if results, err := (tmdbProvider{service: service}).Search(t.Context(), SearchRequest{MediaType: "person"}); err != nil || results != nil {
		t.Fatalf("unsupported TMDB search = %#v, %v", results, err)
	}
	if (ProviderHTTPError{StatusCode: 502}).Error() == "" || (ProviderGraphQLError{Message: "bad"}).Error() == "" || (ProviderPayloadError{Message: "bad"}).Error() == "" {
		t.Fatal("typed provider errors must have messages")
	}
	if result := service.Test(t.Context(), Config{Type: "anilist", BaseURL: "https://anilist.test/graphql"}); !result.Success {
		t.Fatalf("service test = %#v", result)
	}
	failed := responseAniListService(`{"errors":[{"message":"bad"}]}`, http.StatusOK, nil).Test(t.Context(), Config{Type: "anilist"})
	if failed.Success || failed.Message == "" {
		t.Fatalf("failed service test = %#v", failed)
	}
}

func TestAniListRejectsInvalidIDBeforeTransport(t *testing.T) {
	service := NewService(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("transport called")
		return nil, nil
	})}, nil)
	_, err := service.Details(context.Background(), Config{Type: "anilist"}, DetailsRequest{ExternalID: "not-a-number"})
	var payloadErr ProviderPayloadError
	if !errors.As(err, &payloadErr) {
		t.Fatalf("error = %T %v", err, err)
	}
}

func fixtureAniListService(t *testing.T, name string, status int, headers http.Header) *Service {
	t.Helper()
	payload, err := os.ReadFile("../testmocks/testdata/anilist/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return responseAniListService(string(payload), status, headers)
}

func responseAniListService(body string, status int, headers http.Header) *Service {
	return NewService(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		var query anilistQuery
		if err := json.NewDecoder(request.Body).Decode(&query); err != nil {
			return nil, err
		}
		if request.Method != http.MethodPost || !strings.Contains(query.Query, "synonyms") || query.Variables["type"] == "MANGA" {
			return &http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader(`{}`)), Header: http.Header{}}, nil
		}
		if headers == nil {
			headers = http.Header{"Content-Type": []string{"application/json"}}
		}
		return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: headers}, nil
	})}, nil)
}
