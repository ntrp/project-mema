package indexers

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestSCNIntegrations001TorznabSearchMapsReleaseFeed(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query().Get("t"); got != "movie" {
			t.Fatalf("search type = %q, want movie", got)
		}
		if got := r.URL.Query().Get("cat"); got != "2000,2010" {
			t.Fatalf("categories = %q", got)
		}
		body := `<rss><channel>
			<item>
				<title> Scenario.Movie.2026.1080p.WEBDL </title>
				<link>https://indexer.local/download/1</link>
				<guid>https://indexer.local/details/1</guid>
				<pubDate>Fri, 03 Jul 2026 01:02:03 +0000</pubDate>
				<size>12345</size>
				<attr name="seeders" value="42"/>
				<attr name="peers" value="7"/>
			</item>
			<item>
				<title>Missing download</title>
				<guid>opaque-guid</guid>
			</item>
		</channel></rss>`
		return response(http.StatusOK, body), nil
	})

	releases, err := NewService(client).Search(context.Background(), Config{
		ID:         "idx-1",
		Name:       "Local Torznab",
		Protocol:   "torrent",
		BaseURL:    "http://indexer.local/api",
		APIKey:     stringPtr("secret"),
		Categories: []int32{2000, 2010},
	}, "Scenario Movie", "movie")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(releases) != 1 {
		t.Fatalf("releases = %#v, want one valid release", releases)
	}
	release := releases[0]
	if release.Title != "Scenario.Movie.2026.1080p.WEBDL" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://indexer.local/download/1" {
		t.Fatalf("download url = %q", release.DownloadURL)
	}
	if release.InfoURL != "https://indexer.local/download/1" {
		t.Fatalf("info url = %q", release.InfoURL)
	}
	if release.SizeBytes != 12345 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.Seeders == nil || *release.Seeders != 42 {
		t.Fatalf("seeders = %#v", release.Seeders)
	}
	if release.Peers == nil || *release.Peers != 7 {
		t.Fatalf("peers = %#v", release.Peers)
	}
	wantPublished := time.Date(2026, time.July, 3, 1, 2, 3, 0, time.UTC)
	if release.PublishedAt == nil || !release.PublishedAt.Equal(wantPublished) {
		t.Fatalf("published = %#v, want %s", release.PublishedAt, wantPublished)
	}
}

func TestSCNIntegrations001TorznabRecentUsesRSSRequestWithoutQuery(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query().Get("t"); got != "search" {
			t.Fatalf("recent type = %q, want search", got)
		}
		if got := r.URL.Query().Get("q"); got != "" {
			t.Fatalf("recent query = %q, want empty", got)
		}
		if got := r.URL.Query().Get("cat"); got != "2000,2010" {
			t.Fatalf("categories = %q", got)
		}
		body := `<rss><channel><item>
			<title>Scenario.Movie.2026.1080p.WEBDL</title>
			<link>https://indexer.local/download/1</link>
			<guid>rss-guid-1</guid>
		</item></channel></rss>`
		return response(http.StatusOK, body), nil
	})

	releases, err := NewService(client).Recent(context.Background(), Config{
		ID:         "idx-1",
		Name:       "Local Torznab",
		Protocol:   "torrent",
		BaseURL:    "http://indexer.local/api",
		APIKey:     stringPtr("secret"),
		Categories: []int32{2000, 2010},
	})
	if err != nil {
		t.Fatalf("recent failed: %v", err)
	}
	if len(releases) != 1 || releases[0].GUID != "rss-guid-1" {
		t.Fatalf("releases = %#v", releases)
	}
}

func TestSCNIntegrations001TorznabSearchReturnsStatusErrorWithRetryAfter(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		resp := response(http.StatusTooManyRequests, "rate limited")
		resp.Header.Set("Retry-After", "5")
		return resp, nil
	})

	_, err := NewService(client).Search(context.Background(), Config{
		Protocol: "torrent",
		BaseURL:  "http://indexer.local/api",
	}, "Scenario Movie", "movie")
	if err == nil {
		t.Fatal("expected status error")
	}
	code := StatusCode(err)
	if code == nil || *code != http.StatusTooManyRequests {
		t.Fatalf("status code = %#v, want 429", code)
	}
	if RetryAfter(err) != 5*time.Second {
		t.Fatalf("retry after = %s, want 5s", RetryAfter(err))
	}
}

func stringPtr(value string) *string {
	return &value
}
