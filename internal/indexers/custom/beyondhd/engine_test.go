package beyondhd

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesBeyondHDResponse(t *testing.T) {
	key := "bhd-api-key"
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %s", req.Method)
		}
		if req.URL.Path != "/api/torrents/"+key {
			t.Fatalf("path = %s", req.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload["rsskey"] != "rss-key" || payload["search"] != "Example Movie" {
			t.Fatalf("payload = %#v", payload)
		}
		body := `{"status_code":1,"status_message":"ok","results":[{"name":"Example.Movie.2026.1080p","info_hash":"abc","category":"Movies","size":4096,"times_completed":8,"seeders":15,"leechers":5,"created_at":"2026-07-04T12:00:00Z","download_url":"https://beyond-hd.test/download/1","url":"https://beyond-hd.test/torrents/1","imdb_id":"tt1234567","freeleech":true,"limited":false}]}`
		return testutil.Response(http.StatusOK, body), nil
	})
	fields, _ := json.Marshal([]struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}{{Name: "rssKey", Value: "rss-key"}})
	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "bhd",
		Name:     "BeyondHD",
		Protocol: "torrent",
		BaseURL:  "https://beyond-hd.test/",
		APIKey:   &key,
		Fields:   fields,
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].DownloadURL != "https://beyond-hd.test/download/1" {
		t.Fatalf("download = %q", releases[0].DownloadURL)
	}
	if releases[0].Peers == nil || *releases[0].Peers != 20 {
		t.Fatalf("peers = %#v", releases[0].Peers)
	}
}
