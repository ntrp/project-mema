package nebulance

import (
	"context"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesNebulanceResponse(t *testing.T) {
	key := "nbl-key"
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/api.php" {
			t.Fatalf("path = %s", req.URL.Path)
		}
		if req.URL.Query().Get("api_key") != key {
			t.Fatalf("api_key = %q", req.URL.Query().Get("api_key"))
		}
		if req.URL.Query().Get("release") != "Example Show" {
			t.Fatalf("release = %q", req.URL.Query().Get("release"))
		}
		body := `{"total_results":1,"items":[{"rls_name":"Example.Show.S01E01.1080p","group_name":"Example Show","group_id":42,"size":4096,"seed":8,"leech":2,"snatch":12,"download":"https://nebulance.test/download/42","file_list":["episode.mkv"],"rls_utc":"2026-07-04T12:00:00Z"}]}`
		return testutil.Response(http.StatusOK, body), nil
	})
	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "nbl",
		Name:     "Nebulance",
		Protocol: "torrent",
		BaseURL:  "https://nebulance.test/",
		APIKey:   &key,
	}, "Example Show", "tv")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].InfoURL != "https://nebulance.test/torrents.php?id=42" {
		t.Fatalf("info = %q", releases[0].InfoURL)
	}
	if releases[0].PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}
