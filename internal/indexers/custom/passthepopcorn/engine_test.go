package passthepopcorn

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesPassThePopcornResponse(t *testing.T) {
	key := "ptp-key"
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/torrents.php" {
			t.Fatalf("path = %s", req.URL.Path)
		}
		if req.Header.Get("ApiUser") != "ptp-user" || req.Header.Get("ApiKey") != key {
			t.Fatalf("headers ApiUser=%q ApiKey=%q", req.Header.Get("ApiUser"), req.Header.Get("ApiKey"))
		}
		if req.URL.Query().Get("searchstr") != "Example Movie" {
			t.Fatalf("searchstr = %q", req.URL.Query().Get("searchstr"))
		}
		body := `{"TotalResults":"1","Movies":[{"GroupId":"10","CategoryId":"1","Title":"Example Movie","Year":"2026","ImdbId":"1234567","Torrents":[{"Id":99,"Size":"2048","UploadTime":"2026-07-04 12:00:00","Snatched":"7","Seeders":"12","Leechers":"3","ReleaseName":"Example.Movie.2026.1080p","FreeleechType":"Freeleech"}]}]}`
		return testutil.Response(http.StatusOK, body), nil
	})
	fields, _ := json.Marshal([]struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}{{Name: "apiUser", Value: "ptp-user"}})
	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "ptp",
		Name:     "PassThePopcorn",
		Protocol: "torrent",
		BaseURL:  "https://passthepopcorn.test/",
		APIKey:   &key,
		Fields:   fields,
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].InfoURL != "https://passthepopcorn.test/torrents.php?id=10&torrentid=99" {
		t.Fatalf("info = %q", releases[0].InfoURL)
	}
	if releases[0].DownloadURL != "https://passthepopcorn.test/torrents.php?action=download&id=99" {
		t.Fatalf("download = %q", releases[0].DownloadURL)
	}
}
