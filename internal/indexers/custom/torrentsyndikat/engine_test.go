package torrentsyndikat

import (
	"context"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesTorrentSyndikatResponse(t *testing.T) {
	key := "ts-key"
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/api_9djWe8Tb2NE3p6opyqnh/v1/browse.php" {
			t.Fatalf("path = %s", req.URL.Path)
		}
		if req.URL.Query().Get("apikey") != key {
			t.Fatalf("apikey = %q", req.URL.Query().Get("apikey"))
		}
		if req.URL.Query().Get("searchstring") != "Example S01*" {
			t.Fatalf("searchstring = %q", req.URL.Query().Get("searchstring"))
		}
		body := `{"rows":[{"id":"55","name":"Example.S01.1080p","category":53,"added":1783166400,"size":8192,"numfiles":4,"seeders":9,"leechers":1,"snatched":3}]}`
		return testutil.Response(http.StatusOK, body), nil
	})
	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "ts",
		Name:     "TorrentSyndikat",
		Protocol: "torrent",
		BaseURL:  "https://torrent-syndikat.test/",
		APIKey:   &key,
	}, "Example S01", "tv")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].DownloadURL != "https://torrent-syndikat.test/download.php?apikey=ts-key&id=55" {
		t.Fatalf("download = %q", releases[0].DownloadURL)
	}
	if releases[0].Seeders == nil || *releases[0].Seeders != 9 {
		t.Fatalf("seeders = %#v", releases[0].Seeders)
	}
}
