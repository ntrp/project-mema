package scenehd

import (
	"context"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesSceneHDResponse(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/browse.php" {
			t.Fatalf("path = %s", req.URL.Path)
		}
		if req.URL.Query().Get("passkey") != key {
			t.Fatalf("passkey = %q", req.URL.Query().Get("passkey"))
		}
		if req.URL.Query().Get("search") != "Example Movie" {
			t.Fatalf("search = %q", req.URL.Query().Get("search"))
		}
		body := `[{"id":123,"name":"Example.Movie.2026.1080p","added":"2026-07-04 12:00:00","size":2048,"times_completed":7,"numfiles":2,"seeders":11,"leechers":3,"category":"1","is_freeleech":1}]`
		return testutil.Response(http.StatusOK, body), nil
	})
	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "scenehd",
		Name:     "SceneHD",
		Protocol: "torrent",
		BaseURL:  "https://scenehd.test/",
		APIKey:   &key,
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].Title != "Example.Movie.2026.1080p" {
		t.Fatalf("title = %q", releases[0].Title)
	}
	if releases[0].DownloadURL != "https://scenehd.test/download.php?id=123&passkey="+key {
		t.Fatalf("download = %q", releases[0].DownloadURL)
	}
	if releases[0].Peers == nil || *releases[0].Peers != 14 {
		t.Fatalf("peers = %#v", releases[0].Peers)
	}
}
