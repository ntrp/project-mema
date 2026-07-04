package speedappapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesSpeedAppResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/api/torrent" {
			t.Fatalf("path = %q", req.URL.Path)
		}
		if req.Header.Get("Authorization") != "Bearer key-1" {
			t.Fatalf("auth = %q", req.Header.Get("Authorization"))
		}
		if req.URL.Query().Get("search") != "Example Movie" {
			t.Fatalf("search = %q", req.URL.Query().Get("search"))
		}
		return testutil.Response(http.StatusOK, `[{
			"id":42,
			"url":"https://speed.test/torrent/42",
			"name":"[REQUEST] Example.Movie.2026.1080p",
			"size":2048,
			"created_at":"2026-07-03T00:00:00Z",
			"times_completed":3,
			"seeders":5,
			"leechers":2,
			"category":{"id":8}
		}]`), nil
	})

	releases, err := New(Options{Name: "SpeedApp", DefaultBaseURL: "https://speed.test/"}, client).
		Search(context.Background(), config(t, field("apiKey", "key-1")), "Example Movie", "movie")
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
	if release.DownloadURL != "https://speed.test/api/torrent/42/download" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.Peers == nil || *release.Peers != 7 {
		t.Fatalf("peers = %#v", release.Peers)
	}
}

func TestSearchLogsInWhenAPIKeyIsMissing(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/api/login":
			var payload map[string]string
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			if payload["username"] != "user-1" || payload["password"] != "pass-1" {
				t.Fatalf("payload = %#v", payload)
			}
			return testutil.Response(http.StatusOK, `{"token":"token-1"}`), nil
		case "/api/torrent":
			if req.Header.Get("Authorization") != "Bearer token-1" {
				t.Fatalf("auth = %q", req.Header.Get("Authorization"))
			}
			return testutil.Response(http.StatusOK, `[]`), nil
		default:
			t.Fatalf("unexpected path %s", req.URL.Path)
			return testutil.Response(http.StatusNotFound, ""), nil
		}
	})

	_, err := New(Options{Name: "SpeedApp", DefaultBaseURL: "https://speed.test/"}, client).
		Search(context.Background(), config(t, field("email", "user-1"), field("password", "pass-1")), "", "movie")
	if err != nil {
		t.Fatal(err)
	}
}

func config(t *testing.T, fields ...map[string]any) engine.Config {
	t.Helper()
	raw, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	return engine.Config{ID: "idx", Name: "SpeedApp", Protocol: "torrent", BaseURL: "https://speed.test/", Fields: raw}
}

func field(name string, value any) map[string]any {
	return map[string]any{"name": name, "value": value}
}
