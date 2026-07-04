package subsplease

import (
	"context"
	"net/http"
	"testing"
	"time"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchMapsSubsPleaseResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("f") != "search" || r.URL.Query().Get("s") != "Example Show" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		body := `{"example":{"release_date":"2026-07-02T18:40:00Z","show":"Example Show","episode":"03","page":"example-show","downloads":[{"res":"1080","magnet":"magnet:?xt=urn:btih:abc&xl=4096"}]}}`
		return testutil.Response(http.StatusOK, body), nil
	})

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "idx-subs",
		Name:     "SubsPlease",
		Protocol: "torrent",
		BaseURL:  "https://subsplease.test/",
	}, "[SubsPlease] Example Show 1080p", "series")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "[SubsPlease] Example Show - 03 (1080p)" || release.DownloadURL != "magnet:?xt=urn:btih:abc&xl=4096" {
		t.Fatalf("release = %#v", release)
	}
	if release.InfoURL != "https://subsplease.test/shows/example-show/" {
		t.Fatalf("info url = %q", release.InfoURL)
	}
	if release.SizeBytes != 4096 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	want := time.Date(2026, time.July, 2, 18, 40, 0, 0, time.UTC)
	if release.PublishedAt == nil || !release.PublishedAt.Equal(want) {
		t.Fatalf("published = %#v", release.PublishedAt)
	}
}
