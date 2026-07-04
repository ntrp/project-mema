package headphones

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesHeadphonesFeed(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/api" {
			t.Fatalf("path = %q", req.URL.Path)
		}
		if req.URL.Query().Get("t") != "search" {
			t.Fatalf("t = %q", req.URL.Query().Get("t"))
		}
		if req.URL.Query().Get("extended") != "1" {
			t.Fatalf("extended = %q", req.URL.Query().Get("extended"))
		}
		if req.URL.Query().Get("apikey") != "key-1" {
			t.Fatalf("apikey = %q", req.URL.Query().Get("apikey"))
		}
		if req.URL.Query().Get("q") != "Example Album" {
			t.Fatalf("q = %q", req.URL.Query().Get("q"))
		}
		wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("user-1:pass-1"))
		if req.Header.Get("Authorization") != wantAuth {
			t.Fatalf("auth = %q", req.Header.Get("Authorization"))
		}
		return testutil.Response(http.StatusOK, `<rss xmlns:newznab="http://www.newznab.com/DTD/2010/feeds/attributes/"><channel><item>
			<title>Example Album</title>
			<link>https://hp.test/details/1</link>
			<guid>guid-1</guid>
			<comments>https://hp.test/details/1#comments</comments>
			<pubDate>Fri, 03 Jul 2026 00:00:00 +0000</pubDate>
			<enclosure url="https://hp.test/download/1" length="2048"/>
			<newznab:attr name="size" value="4096"/>
			<newznab:attr name="usenetdate" value="2026-07-03T00:00:00Z"/>
		</item></channel></rss>`), nil
	})

	releases, err := New(client).Search(context.Background(), config(t), "Example Album", "music")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example Album" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://hp.test/download/1" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.InfoURL != "https://hp.test/details/1" {
		t.Fatalf("info = %q", release.InfoURL)
	}
	if release.GUID != "guid-1" {
		t.Fatalf("guid = %q", release.GUID)
	}
	if release.SizeBytes != 4096 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}

func config(t *testing.T) engine.Config {
	t.Helper()
	fields, err := json.Marshal([]map[string]any{
		{"name": "apiKey", "value": "key-1"},
		{"name": "username", "value": "user-1"},
		{"name": "password", "value": "pass-1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return engine.Config{
		ID:       "idx",
		Name:     "Headphones",
		Protocol: "usenet",
		BaseURL:  "https://hp.test/",
		Fields:   fields,
	}
}
