package knaben

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchMapsKnabenResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.String() != apiEndpoint {
			t.Fatalf("request = %s %s", r.Method, r.URL.String())
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["query"] != "Example Movie" || body["search_field"] != "title" {
			t.Fatalf("body = %#v", body)
		}
		payload := `{"hits":[{"title":"Example Movie 2026","categoryId":[3001000],"hash":"ABC123","details":"https://knaben.test/details/1","link":"","magnetUrl":"magnet:?xt=urn:btih:abc","bytes":8192,"seeders":9,"peers":4,"date":"2026-07-02T18:40:00"}]}`
		return testutil.Response(http.StatusOK, payload), nil
	})

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "idx-knaben",
		Name:     "Knaben",
		Protocol: "torrent",
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example Movie 2026" || release.DownloadURL != "magnet:?xt=urn:btih:abc" {
		t.Fatalf("release = %#v", release)
	}
	if release.InfoURL != "https://knaben.test/details/1" || release.SizeBytes != 8192 {
		t.Fatalf("release = %#v", release)
	}
	if release.Seeders == nil || *release.Seeders != 9 || release.Peers == nil || *release.Peers != 13 {
		t.Fatalf("seeders = %#v peers = %#v", release.Seeders, release.Peers)
	}
	want := time.Date(2026, time.July, 2, 17, 40, 0, 0, time.UTC)
	if release.PublishedAt == nil || !release.PublishedAt.Equal(want) {
		t.Fatalf("published = %#v", release.PublishedAt)
	}
}
