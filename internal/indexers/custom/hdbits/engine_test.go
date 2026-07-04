package hdbits

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesHDBitsResponse(t *testing.T) {
	key := "hdbits-key"
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("method = %s", req.Method)
		}
		if req.URL.Path != "/api/torrents" {
			t.Fatalf("path = %s", req.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload["username"] != "alice" || payload["passkey"] != key {
			t.Fatalf("auth payload = %#v", payload)
		}
		if payload["search"] != "Example Movie" {
			t.Fatalf("search = %#v", payload["search"])
		}
		body := `{"status":0,"data":[{"id":"123","hash":"abc","name":"Example.Movie.2026.1080p","filename":"Example.Movie.2026.1080p.mkv.torrent","size":2048,"seeders":12,"leechers":3,"times_completed":5,"numfiles":1,"type_category":1,"type_medium":3,"added":"2026-07-04T12:00:00Z"}]}`
		return testutil.Response(http.StatusOK, body), nil
	})
	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "hdb",
		Name:     "HDBits",
		Protocol: "torrent",
		BaseURL:  "https://hdbits.test/",
		APIKey:   &key,
		Fields:   json.RawMessage(`[{"name":"username","value":"alice"},{"name":"useFilenames","value":true}]`),
	}, "Example.Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].Title != "Example.Movie.2026.1080p.mkv" {
		t.Fatalf("title = %q", releases[0].Title)
	}
	if releases[0].DownloadURL != "https://hdbits.test/download.php?id=123&passkey=hdbits-key" {
		t.Fatalf("download = %q", releases[0].DownloadURL)
	}
}
