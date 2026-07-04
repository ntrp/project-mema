package binsearch

import (
	"context"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesBinSearchRows(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("adv_col") != "on" {
			t.Fatalf("adv_col = %q", req.URL.Query().Get("adv_col"))
		}
		if req.URL.Query().Get("q") != "Example Movie" {
			t.Fatalf("q = %q", req.URL.Query().Get("q"))
		}
		return testutil.Response(http.StatusOK, `<table class="xMenuT"><tbody>
			<tr><th>header</th></tr>
			<tr>
				<td><input type="checkbox" name="abc123"><a href="/details/abc123">info</a><span class="s">"Example.Movie.2026.mkv"</span><span class="d">size: 1.5 GB</span></td>
				<td></td>
				<td></td>
				<td></td>
				<td></td>
				<td>Fri, 03 Jul 2026 00:00:00 +0000</td>
			</tr>
		</tbody></table>`), nil
	})

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "idx",
		Name:     "BinSearch",
		Protocol: "usenet",
		BaseURL:  "https://binsearch.test/",
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example.Movie.2026" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://binsearch.test/?action=nzb&abc123=1" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.InfoURL != "https://binsearch.test/details/abc123" {
		t.Fatalf("info = %q", release.InfoURL)
	}
	if release.SizeBytes != 1610612736 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}
