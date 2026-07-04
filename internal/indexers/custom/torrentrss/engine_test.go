package torrentrss

import (
	"context"
	"net/http"
	"testing"
	"time"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchFiltersAndMapsTorrentRssFeed(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://rss.test/feed" {
			t.Fatalf("url = %s", r.URL.String())
		}
		body := `<rss><channel>
			<item>
				<title>Example Movie 2026</title>
				<link>https://rss.test/details/1</link>
				<guid>https://rss.test/details/1</guid>
				<pubDate>Fri, 03 Jul 2026 00:00:00 +0000</pubDate>
				<enclosure url="magnet:?xt=urn:btih:abc" length="2048"/>
			</item>
			<item><title>Other Movie</title><link>magnet:?xt=urn:btih:other</link></item>
		</channel></rss>`
		resp := testutil.Response(http.StatusOK, body)
		resp.Header.Set("Content-Type", "application/rss+xml")
		return resp, nil
	})

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "idx-rss",
		Name:     "Torrent RSS Feed",
		Protocol: "torrent",
		BaseURL:  "https://rss.test/feed",
	}, "Example", "movie")
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
	if release.InfoURL != "https://rss.test/details/1" || release.SizeBytes != 2048 {
		t.Fatalf("release = %#v", release)
	}
	want := time.Date(2026, time.July, 3, 0, 0, 0, 0, time.UTC)
	if release.PublishedAt == nil || !release.PublishedAt.Equal(want) {
		t.Fatalf("published = %#v", release.PublishedAt)
	}
}
