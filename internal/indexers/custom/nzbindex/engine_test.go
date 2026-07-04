package nzbindex

import (
	"context"
	"net/http"
	"testing"
	"time"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchMapsNzbIndexResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v3/search/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("key") != "secret" || r.URL.Query().Get("q") != "Example Movie" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		body := `{"results":[{"id":"abc","name":"\"Example.Movie.2026.mkv\"","posted":1783036800000,"size":4096,"file_count":4,"group_ids":[1]}]}`
		return testutil.Response(http.StatusOK, body), nil
	})
	apiKey := "secret"

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "idx-nzbindex",
		Name:     "NZBIndex",
		Protocol: "usenet",
		BaseURL:  "https://nzbindex.test/",
		APIKey:   &apiKey,
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example.Movie.2026" || release.DownloadURL != "https://nzbindex.test/download/abc" {
		t.Fatalf("release = %#v", release)
	}
	if release.InfoURL != "https://nzbindex.test/collection/abc" || release.SizeBytes != 4096 {
		t.Fatalf("release = %#v", release)
	}
	want := time.Date(2026, time.July, 3, 0, 0, 0, 0, time.UTC)
	if release.PublishedAt == nil || !release.PublishedAt.Equal(want) {
		t.Fatalf("published = %#v", release.PublishedAt)
	}
}
