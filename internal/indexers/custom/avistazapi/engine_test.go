package avistazapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesAvistazAPIResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/api/v1/jackett/torrents" {
			t.Fatalf("path = %q", req.URL.Path)
		}
		if req.Header.Get("Authorization") != "Bearer token-1" {
			t.Fatalf("auth = %q", req.Header.Get("Authorization"))
		}
		if req.URL.Query().Get("search") != "Example Movie" {
			t.Fatalf("search = %q", req.URL.Query().Get("search"))
		}
		if req.URL.Query().Get("in") != "1" {
			t.Fatalf("in = %q", req.URL.Query().Get("in"))
		}
		if req.URL.Query().Get("type") != "1" {
			t.Fatalf("type = %q", req.URL.Query().Get("type"))
		}
		return testutil.Response(http.StatusOK, `{"data":[{
			"url":"https://avistaz.test/torrent/1",
			"download":"https://avistaz.test/download/1",
			"file_name":"Example.Movie.2026.1080p",
			"release_title":"Display Title",
			"info_hash":"abc123",
			"seed":7,
			"leech":2,
			"file_size":2048,
			"created_at_iso":"2026-07-03T00:00:00Z",
			"type":"MOVIE",
			"video_quality":"1080p"
		}]}`), nil
	})

	releases, err := New(Options{Name: "AvistaZ", DefaultBaseURL: "https://avistaz.test/"}, client).
		Search(context.Background(), config(t, []int32{2000}, field("token", "token-1")), "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example.Movie.2026.1080p" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://avistaz.test/download/1" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.InfoURL != "https://avistaz.test/torrent/1" {
		t.Fatalf("info = %q", release.InfoURL)
	}
	if release.SizeBytes != 2048 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.Seeders == nil || *release.Seeders != 7 {
		t.Fatalf("seeders = %#v", release.Seeders)
	}
	if release.Peers == nil || *release.Peers != 9 {
		t.Fatalf("peers = %#v", release.Peers)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}

func TestSearchUsesLoginWhenTokenIsMissing(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/api/v1/jackett/auth":
			if req.Method != http.MethodPost {
				t.Fatalf("auth method = %s", req.Method)
			}
			if err := req.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if req.Form.Get("username") != "user-1" || req.Form.Get("password") != "pass-1" || req.Form.Get("pid") != "pid-1" {
				t.Fatalf("auth form = %#v", req.Form)
			}
			return testutil.Response(http.StatusOK, `{"token":"new-token"}`), nil
		case "/api/v1/jackett/torrents":
			if req.Header.Get("Authorization") != "Bearer new-token" {
				t.Fatalf("auth = %q", req.Header.Get("Authorization"))
			}
			return testutil.Response(http.StatusOK, `{"data":[]}`), nil
		default:
			t.Fatalf("unexpected path %s", req.URL.Path)
			return testutil.Response(http.StatusNotFound, ""), nil
		}
	})

	_, err := New(Options{Name: "AvistaZ", DefaultBaseURL: "https://avistaz.test/"}, client).
		Search(context.Background(), config(t, nil,
			field("username", "user-1"),
			field("password", "pass-1"),
			field("pid", "pid-1"),
		), "", "movie")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAnimeZUsesReleaseTitleAndFormatCategories(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("in") != "" || req.URL.Query().Get("type") != "" {
			t.Fatalf("unexpected base category params %s", req.URL.RawQuery)
		}
		if got := req.URL.Query()["format[]"]; len(got) == 0 {
			t.Fatalf("missing anime format categories in %s", req.URL.RawQuery)
		}
		return testutil.Response(http.StatusOK, `{"data":[{
			"url":"https://animez.test/torrent/1",
			"download":"https://animez.test/download/1",
			"file_name":"fallback-title",
			"release_title":"Anime Release Title",
			"seed":1,
			"leech":1,
			"file_size":1024,
			"created_at_iso":"2026-07-03T00:00:00Z",
			"format":"TV"
		}]}`), nil
	})

	releases, err := New(Options{
		Name:            "AnimeZ",
		DefaultBaseURL:  "https://animez.test/",
		AnimeCategories: true,
		PreferRelease:   true,
	}, client).Search(context.Background(), config(t, []int32{5000}, field("token", "token-1")), "Anime", "series")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].Title != "Anime Release Title" {
		t.Fatalf("title = %q", releases[0].Title)
	}
}

func config(t *testing.T, categories []int32, fields ...map[string]any) engine.Config {
	t.Helper()
	raw, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	return engine.Config{
		ID:         "idx",
		Name:       "AvistaZ",
		Protocol:   "torrent",
		BaseURL:    "https://avistaz.test/",
		Categories: categories,
		Fields:     raw,
	}
}

func field(name string, value any) map[string]any {
	return map[string]any{"name": name, "value": value}
}
