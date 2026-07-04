package torrentpotato

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchMapsTorrentPotatoResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Get("user") != "demo" || r.URL.Query().Get("passkey") != "secret" || r.URL.Query().Get("search") != "Example Movie" {
			t.Fatalf("query = %s", r.URL.RawQuery)
		}
		body := `{"results":[{"release_name":"Example &amp; Movie 2026","details_url":"https://potato.test/details/1","download_url":"magnet:?xt=urn:btih:abcdef","size":700,"leechers":3,"seeders":8,"publish_date":"2026-07-02T18:40:00Z"}]}`
		return testutil.Response(http.StatusOK, body), nil
	})
	fields, _ := json.Marshal([]struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}{{Name: "user", Value: "demo"}, {Name: "passkey", Value: "secret"}})

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "idx-potato",
		Name:     "TorrentPotato",
		Protocol: "torrent",
		BaseURL:  "https://potato.test/api",
		Fields:   fields,
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example & Movie 2026" || release.GUID != "potato-abcdef" {
		t.Fatalf("release = %#v", release)
	}
	if release.SizeBytes != 700000000 || release.Seeders == nil || *release.Seeders != 8 || release.Peers == nil || *release.Peers != 11 {
		t.Fatalf("release = %#v", release)
	}
	want := time.Date(2026, time.July, 2, 18, 40, 0, 0, time.UTC)
	if release.PublishedAt == nil || !release.PublishedAt.Equal(want) {
		t.Fatalf("published = %#v", release.PublishedAt)
	}
}
