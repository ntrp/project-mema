package unit3d

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesUNIT3DResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/api/torrents/filter" {
			t.Fatalf("path = %q", req.URL.Path)
		}
		if req.URL.Query().Get("api_token") != "key-1" {
			t.Fatalf("api_token = %q", req.URL.Query().Get("api_token"))
		}
		if req.URL.Query().Get("name") != "Example Movie" {
			t.Fatalf("name = %q", req.URL.Query().Get("name"))
		}
		if got := req.URL.Query()["categories[]"]; len(got) != 1 || got[0] != "2000" {
			t.Fatalf("categories = %#v", got)
		}
		return testutil.Response(http.StatusOK, `{"data":[{
			"id":"hash-1",
			"attributes":{
				"name":"Example.Movie.2026.1080p",
				"size":2048,
				"num_file":1,
				"times_completed":5,
				"seeders":7,
				"leechers":2,
				"created_at":"2026-07-03T00:00:00Z",
				"download_link":"https://unit3d.test/download/1",
				"details_link":"https://unit3d.test/torrents/1"
			}
		}]}`), nil
	})

	releases, err := New(client).Search(context.Background(), config(t), "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example.Movie.2026.1080p" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://unit3d.test/download/1" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.InfoURL != "https://unit3d.test/torrents/1" {
		t.Fatalf("info = %q", release.InfoURL)
	}
	if release.SizeBytes != 2048 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.Seeders == nil || *release.Seeders != 7 {
		t.Fatalf("seeders = %#v", release.Seeders)
	}
	if release.Peers == nil || *release.Peers != 9 {
		t.Fatalf("peers = %#v", release.Peers)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}

func config(t *testing.T) engine.Config {
	t.Helper()
	raw, err := json.Marshal([]map[string]any{{"name": "apiKey", "value": "key-1"}})
	if err != nil {
		t.Fatal(err)
	}
	return engine.Config{
		ID:         "idx",
		Name:       "UNIT3D",
		Protocol:   "torrent",
		BaseURL:    "https://unit3d.test/",
		Categories: []int32{2000},
		Fields:     raw,
	}
}
