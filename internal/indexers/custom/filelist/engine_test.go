package filelist

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesFileListResponse(t *testing.T) {
	key := "filelist-key"
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		username, password, ok := req.BasicAuth()
		if !ok || username != "alice" || password != key {
			t.Fatalf("basic auth = %q %q %v", username, password, ok)
		}
		if req.URL.Query().Get("action") != "search-torrents" {
			t.Fatalf("action = %q", req.URL.Query().Get("action"))
		}
		if req.URL.Query().Get("query") != "Example Movie" {
			t.Fatalf("query = %q", req.URL.Query().Get("query"))
		}
		body := `[{"id":44,"name":"Example.Movie.2026.1080p","size":4096,"seeders":10,"leechers":2,"times_completed":6,"files":2,"upload_date":"2026-07-04 12:00:00","category":"Filme HD","freeleech":true,"doubleup":false}]`
		return testutil.Response(http.StatusOK, body), nil
	})
	fields, _ := json.Marshal([]struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}{{Name: "username", Value: "alice"}})
	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "fl",
		Name:     "FileList",
		Protocol: "torrent",
		BaseURL:  "https://filelist.test/",
		APIKey:   &key,
		Fields:   fields,
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].InfoURL != "https://filelist.test/details.php?id=44" {
		t.Fatalf("info = %q", releases[0].InfoURL)
	}
	if releases[0].Peers == nil || *releases[0].Peers != 12 {
		t.Fatalf("peers = %#v", releases[0].Peers)
	}
}
