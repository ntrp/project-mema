package metadata

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTVDBMovieDetailsMapsProviderAgnosticDetails(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := ""
		status := http.StatusOK
		switch r.URL.Path {
		case "/login":
			body = `{"data":{"token":"tvdb-token"}}`
		case "/movies/900/extended":
			if r.URL.Query().Get("meta") != "translations" {
				t.Fatalf("meta query = %q, want translations", r.URL.Query().Get("meta"))
			}
			body = tvdbRichMoviePayload
		default:
			status = http.StatusNotFound
		}
		return tvdbTestResponse(status, body), nil
	})}

	apiKey := "test-key"
	service := NewService(client, nil)
	details, err := service.Details(context.Background(), Config{
		Type:    "tvdb",
		BaseURL: "https://metadata.test",
		APIKey:  &apiKey,
	}, DetailsRequest{MediaType: "movie", ExternalID: "900"})
	if err != nil {
		t.Fatalf("Details() error = %v", err)
	}
	if details.ExternalProvider != "tvdb" || details.ExternalID != "900" {
		t.Fatalf("details provider link = %s:%s, want tvdb:900", details.ExternalProvider, details.ExternalID)
	}
	if details.ExternalURL == nil || *details.ExternalURL != "https://thetvdb.com/movies/example-tvdb-movie" {
		t.Fatalf("details provider url = %#v", details.ExternalURL)
	}
	if details.Title != "Example TVDB Movie" || details.ReleaseDate == nil || *details.ReleaseDate != "2026-04-12" {
		t.Fatalf("details = %#v", details)
	}
	if details.Overview == nil || *details.Overview != "Translated TVDB overview." {
		t.Fatalf("overview = %#v", details.Overview)
	}
	if details.PosterPath == nil || *details.PosterPath != "/poster.jpg" {
		t.Fatalf("poster = %#v", details.PosterPath)
	}
	if details.BackdropPath == nil || *details.BackdropPath != "/backdrop.jpg" {
		t.Fatalf("backdrop = %#v", details.BackdropPath)
	}
	if details.TrailerURL == nil || *details.TrailerURL != "https://video.test/trailer" {
		t.Fatalf("trailer = %#v", details.TrailerURL)
	}
	if details.VoteAverage != nil {
		t.Fatalf("TVDB score mapped as vote average: %#v", details.VoteAverage)
	}
	if factValue(details.Facts, "Certification") != "PG" || factValue(details.Facts, "IMDb ID") != "tt1234567" {
		t.Fatalf("facts = %#v", details.Facts)
	}
	if factValue(details.Facts, "Studios") != "Scenario Studio\nScenario Productions" {
		t.Fatalf("studio facts = %#v", details.Facts)
	}
	if factValue(details.Facts, "Revenue") != "$533,300,000.00" || factValue(details.Facts, "Budget") != "$120,000,000.00" {
		t.Fatalf("money facts = %#v", details.Facts)
	}
	if factValue(details.Facts, "Production Countries") != "🇺🇸 United States" {
		t.Fatalf("country facts = %#v", details.Facts)
	}
	if len(details.Keywords) != 2 || details.Keywords[0] != "Space" || details.Keywords[1] != "Based on a Novel" {
		t.Fatalf("keywords = %#v", details.Keywords)
	}
	if len(details.Cast) != 1 || details.Cast[0].Name != "Example Actor" {
		t.Fatalf("details cast = %#v", details.Cast)
	}
}

func TestTVDBMovieDetailsLoadsTranslationFallback(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := ""
		status := http.StatusOK
		switch r.URL.Path {
		case "/login":
			body = `{"data":{"token":"tvdb-token"}}`
		case "/movies/902/extended":
			body = `{"data":{"id":902,"name":"Fallback Movie","year":"2026"}}`
		case "/movies/902/translations/eng":
			body = `{"data":{"language":"eng","overview":"Fallback overview."}}`
		default:
			status = http.StatusNotFound
		}
		return tvdbTestResponse(status, body), nil
	})}

	apiKey := "test-key"
	service := NewService(client, nil)
	details, err := service.Details(context.Background(), Config{
		Type:    "tvdb",
		BaseURL: "https://metadata.test",
		APIKey:  &apiKey,
	}, DetailsRequest{MediaType: "movie", ExternalID: "902"})
	if err != nil {
		t.Fatalf("Details() error = %v", err)
	}
	if details.Overview == nil || *details.Overview != "Fallback overview." {
		t.Fatalf("overview = %#v", details.Overview)
	}
	if details.ExternalURL == nil || *details.ExternalURL != "https://thetvdb.com/dereferrer/movie/902" {
		t.Fatalf("fallback external url = %#v", details.ExternalURL)
	}
}

func TestTVDBMovieDetailsAcceptsObjectFirstRelease(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := ""
		status := http.StatusOK
		switch r.URL.Path {
		case "/login":
			body = `{"data":{"token":"tvdb-token"}}`
		case "/movies/901/extended":
			body = `{"data":{"id":901,"name":"Object Date Movie","year":"2026","first_release":{"country":"usa","date":"2026-05-20"},"runtime":99}}`
		default:
			status = http.StatusNotFound
		}
		return tvdbTestResponse(status, body), nil
	})}

	apiKey := "test-key"
	service := NewService(client, nil)
	details, err := service.Details(context.Background(), Config{
		Type:    "tvdb",
		BaseURL: "https://metadata.test",
		APIKey:  &apiKey,
	}, DetailsRequest{MediaType: "movie", ExternalID: "901"})
	if err != nil {
		t.Fatalf("Details() error = %v", err)
	}
	if details.ReleaseDate == nil || *details.ReleaseDate != "2026-05-20" {
		t.Fatalf("ReleaseDate = %#v, want 2026-05-20", details.ReleaseDate)
	}
}

func factValue(facts []Fact, label string) string {
	for _, fact := range facts {
		if fact.Label == label {
			return fact.Value
		}
	}
	return ""
}

func tvdbTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

const tvdbRichMoviePayload = `{
	"data": {
		"id": 900,
		"slug": "example-tvdb-movie",
		"name": "Example TVDB Movie",
		"year": "2026",
		"first_release": {"country": "usa", "date": "2026-04-12"},
		"runtime": 101,
		"status": {"name": "Released"},
		"originalCountry": "usa",
		"originalLanguage": "eng",
		"score": 9876,
		"boxOffice": "533300000.00",
		"budget": "120000000.00",
		"genres": [{"name": "Adventure"}],
		"artworks": [
			{"type": 14, "image": "/poster.jpg", "width": 680, "height": 1000, "score": 50},
			{"type": 15, "image": "/backdrop.jpg", "width": 1920, "height": 1080, "score": 50}
		],
		"translations": {
			"overviewTranslations": [
				{"language": "eng", "overview": "Translated TVDB overview."}
			]
		},
		"contentRatings": [{"country": "usa", "name": "PG"}],
		"companies": {
			"studio": [{"name": "Scenario Studio"}],
			"network": [{"name": "Scenario Network"}],
			"production": [{"name": "Scenario Productions"}]
		},
		"tagOptions": [{"name": "Space"}],
		"inspirations": [{"type_name": "Based on a Novel"}],
		"production_countries": [{"country": "usa", "name": "United States"}],
		"remoteIds": [{"sourceName": "IMDB", "id": "tt1234567"}],
		"trailers": [{"url": "https://video.test/trailer"}],
		"characters": [{
			"peopleId": 4001, "personName": "Example Actor",
			"name": "Lead",
			"peopleType": "Actor",
			"image": "/actor.jpg"
		}]
	}
}`
