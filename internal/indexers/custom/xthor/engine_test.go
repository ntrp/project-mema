package xthor

import (
	"context"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchMapsXthorResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()
		if req.Method != http.MethodGet || req.URL.String() == "" {
			t.Fatalf("request = %s %s", req.Method, req.URL.String())
		}
		if query.Get("passkey") != "key-1" || query.Get("search") != "Example Movie" {
			t.Fatalf("query = %s", req.URL.RawQuery)
		}
		if query.Get("category") != "1+2" || query.Get("freeleech") != "1" {
			t.Fatalf("query = %s", req.URL.RawQuery)
		}
		return testutil.Response(http.StatusOK, `{
			"Error": {"Code": 0, "Descr": "ok"},
			"Torrents": [
				{
					"Id": 1,
					"Category": 106,
					"Seeders": 1,
					"Leechers": 1,
					"Name": "Hidden",
					"Size": 1,
					"Added": 1783166400,
					"Download_link": "https://api.xthor.test/download/hidden"
				},
				{
					"Id": 55,
					"Category": 1,
					"Seeders": 9,
					"Leechers": 1,
					"Name": "Example.Movie.2026.1080p",
					"Times_completed": 4,
					"Size": 4096,
					"Added": 1783166400,
					"Freeleech": 1,
					"Numfiles": 2,
					"Download_link": "https://api.xthor.test/download/55",
					"Tmdb_id": 123
				}
			]
		}`), nil
	})

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:         "idx-xthor",
		Name:       "Xthor",
		Protocol:   "torrent",
		BaseURL:    "https://api.xthor.test/",
		APIKey:     stringPtr("key-1"),
		Categories: []int32{1, 2},
		Fields:     []byte(`[{"name":"freeleechOnly","value":true}]`),
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example.Movie.2026.1080p" || release.DownloadURL != "https://api.xthor.test/download/55" {
		t.Fatalf("release = %#v", release)
	}
	if release.InfoURL != "https://xthor.test/details.php?id=55" || release.GUID != release.InfoURL {
		t.Fatalf("release = %#v", release)
	}
	if release.Seeders == nil || *release.Seeders != 9 || release.Peers == nil || *release.Peers != 10 {
		t.Fatalf("seeders = %#v peers = %#v", release.Seeders, release.Peers)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}

func stringPtr(value string) *string {
	return &value
}
