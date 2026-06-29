package metadata

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTMDBSeriesDetailsLoadsSeasonEpisodes(t *testing.T) {
	seasonRequested := false
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := ""
		status := http.StatusOK
		switch r.URL.Path {
		case "/tv/123":
			body = `{
				"id": 123,
				"name": "Example Series",
				"first_air_date": "2026-01-01",
				"number_of_seasons": 1,
				"number_of_episodes": 1,
				"seasons": [
					{"name": "Season 1", "season_number": 1, "episode_count": 1}
				],
				"credits": {"cast": []}
			}`
		case "/tv/123/season/1":
			seasonRequested = true
			body = `{
				"episodes": [
					{
						"name": "Pilot",
						"episode_number": 1,
						"overview": "The story starts.",
						"air_date": "2026-01-01",
						"still_path": "/pilot.jpg"
					}
				]
			}`
		default:
			status = http.StatusNotFound
		}
		return &http.Response{
			StatusCode: status,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}

	apiKey := "test-key"
	service := NewService(client, nil)
	details, err := service.Details(context.Background(), Config{
		Type:    "tmdb",
		BaseURL: "https://metadata.test",
		APIKey:  &apiKey,
	}, DetailsRequest{
		MediaType:  "series",
		ExternalID: "123",
	})
	if err != nil {
		t.Fatalf("Details() error = %v", err)
	}
	if !seasonRequested {
		t.Fatal("Details() did not request season episode data")
	}
	if len(details.Seasons) != 1 {
		t.Fatalf("len(details.Seasons) = %d, want 1", len(details.Seasons))
	}
	if len(details.Seasons[0].Episodes) != 1 {
		t.Fatalf("len(details.Seasons[0].Episodes) = %d, want 1", len(details.Seasons[0].Episodes))
	}
	episode := details.Seasons[0].Episodes[0]
	if episode.Name != "Pilot" || episode.EpisodeNumber != 1 {
		t.Fatalf("episode = %#v, want Pilot episode 1", episode)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
