package torrentscsv

import (
	"context"
	"net/http"
	"testing"
	"time"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchMapsTorrentsCSVResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/service/search" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.URL.Query().Get("size") != "100" || r.URL.Query().Get("q") != "Example Movie" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		body := `{"torrents":[{"infohash":"ABC123","name":"Example Movie 2026","size_bytes":2048,"created_unix":1783036800,"leechers":2,"seeders":11,"completed":5}]}`
		return testutil.Response(http.StatusOK, body), nil
	})

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "idx-tcsv",
		Name:     "TorrentsCSV",
		Protocol: "torrent",
		BaseURL:  "https://torrents-csv.test/",
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example Movie 2026" || release.DownloadURL != "magnet:?xt=urn:btih:ABC123" {
		t.Fatalf("release = %#v", release)
	}
	if release.SizeBytes != 2048 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.Seeders == nil || *release.Seeders != 11 || release.Peers == nil || *release.Peers != 13 {
		t.Fatalf("seeders = %#v peers = %#v", release.Seeders, release.Peers)
	}
	want := time.Date(2026, time.July, 3, 0, 0, 0, 0, time.UTC)
	if release.PublishedAt == nil || !release.PublishedAt.Equal(want) {
		t.Fatalf("published = %#v", release.PublishedAt)
	}
}
