package anidex

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesAnidexRows(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("page") != "search" {
			t.Fatalf("page = %q", req.URL.Query().Get("page"))
		}
		if req.URL.Query().Get("q") != "Example Anime" {
			t.Fatalf("q = %q", req.URL.Query().Get("q"))
		}
		if req.URL.Query().Get("a") != "1" {
			t.Fatalf("authorised flag = %q", req.URL.Query().Get("a"))
		}
		return testutil.Response(http.StatusOK, `<div id="content"><table><tbody><tr>
			<td><img title="English"></td>
			<td></td>
			<td><a href="/torrent/1"><span title="Example Anime 01"></span></a></td>
			<td><a href="/dl/1">DL</a><a href="magnet:?xt=urn:btih:abc">Magnet</a></td>
			<td></td>
			<td></td>
			<td>1.5 GiB</td>
			<td title="2026-07-03 00:00:00 UTC"></td>
			<td>9</td>
			<td>4</td>
		</tr></tbody></table></div>`), nil
	})

	releases, err := New(client).Search(context.Background(), config(t, "anidex.test", field("authorisedOnly", true)), "Example Anime", "series")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example Anime 01 [English]" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://anidex.test/dl/1" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.InfoURL != "https://anidex.test/torrent/1" {
		t.Fatalf("info = %q", release.InfoURL)
	}
	if release.SizeBytes != 1610612736 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.Seeders == nil || *release.Seeders != 9 {
		t.Fatalf("seeders = %#v", release.Seeders)
	}
	if release.Peers == nil || *release.Peers != 13 {
		t.Fatalf("peers = %#v", release.Peers)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}

func config(t *testing.T, host string, fields ...map[string]any) engine.Config {
	t.Helper()
	raw, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	return engine.Config{
		ID:       "idx",
		Name:     "Test Indexer",
		Protocol: "torrent",
		BaseURL:  "https://" + host + "/",
		Fields:   raw,
	}
}

func field(name string, value any) map[string]any {
	return map[string]any{"name": name, "value": value}
}
