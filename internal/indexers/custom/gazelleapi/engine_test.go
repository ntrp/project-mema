package gazelleapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesNestedGazelleTorrentsWithCookie(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/ajax.php" {
			t.Fatalf("path = %q", req.URL.Path)
		}
		if req.Header.Get("Cookie") != "session=abc" {
			t.Fatalf("cookie = %q", req.Header.Get("Cookie"))
		}
		if req.URL.Query().Get("action") != "browse" {
			t.Fatalf("action = %q", req.URL.Query().Get("action"))
		}
		if req.URL.Query().Get("searchstr") != "Example Movie" {
			t.Fatalf("searchstr = %q", req.URL.Query().Get("searchstr"))
		}
		return testutil.Response(http.StatusOK, `{"status":"success","response":{"results":[{
			"groupId":"10",
			"groupName":"Example Album",
			"artist":"Example Artist",
			"groupYear":"2026",
			"torrents":[{
				"torrentId":42,
				"media":"WEB",
				"encoding":"FLAC",
				"format":"Lossless",
				"hasCue":true,
				"fileCount":12,
				"time":"2026-07-03T00:00:00Z",
				"size":"2048",
				"seeders":"7",
				"leechers":"2"
			}]
		}]}}`), nil
	})

	releases, err := New(Options{Name: "Gazelle", DefaultBaseURL: "https://gazelle.test/"}, client).
		Search(context.Background(), config(t, field("cookie", "session=abc")), "Example.Movie", "music")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example Artist - Example Album (2026) [Lossless FLAC] [WEB] [Cue]" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://gazelle.test/torrents.php?action=download&id=42" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.InfoURL != "https://gazelle.test/torrents.php?id=10&torrentid=42" {
		t.Fatalf("info = %q", release.InfoURL)
	}
	if release.SizeBytes != 2048 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.Peers == nil || *release.Peers != 9 {
		t.Fatalf("peers = %#v", release.Peers)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}

func TestSearchLogsInWhenCookieIsMissing(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/login.php":
			if req.Method != http.MethodPost {
				t.Fatalf("login method = %s", req.Method)
			}
			if err := req.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if req.Form.Get("username") != "user-1" || req.Form.Get("password") != "pass-1" {
				t.Fatalf("login form = %#v", req.Form)
			}
			resp := testutil.Response(http.StatusOK, `{"status":"success"}`)
			resp.Header.Add("Set-Cookie", "session=logged-in; Path=/")
			return resp, nil
		case "/ajax.php":
			if !strings.Contains(req.Header.Get("Cookie"), "session=logged-in") {
				t.Fatalf("cookie = %q", req.Header.Get("Cookie"))
			}
			return testutil.Response(http.StatusOK, `{"status":"success","response":{"results":[]}}`), nil
		default:
			t.Fatalf("unexpected path %s", req.URL.Path)
			return testutil.Response(http.StatusNotFound, ""), nil
		}
	})

	_, err := New(Options{Name: "Gazelle", DefaultBaseURL: "https://gazelle.test/"}, client).
		Search(context.Background(), config(t, field("username", "user-1"), field("password", "pass-1")), "", "movie")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchParsesGreatPosterWallStyleNumbers(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Query().Get("freetorrent") != "1" {
			t.Fatalf("freetorrent = %q", req.URL.Query().Get("freetorrent"))
		}
		return testutil.Response(http.StatusOK, `{"status":"success","response":{"results":[{
			"groupId":10,
			"groupName":"Example Movie",
			"torrents":[{
				"torrentId":42,
				"fileName":"Example.Movie.2026.2160p",
				"time":"2026-07-03T00:00:00Z",
				"size":4096,
				"seeders":3,
				"leechers":1
			}]
		}]}}`), nil
	})

	releases, err := New(Options{Name: "GreatPosterWall", DefaultBaseURL: "https://gpw.test/", FreeleechParam: "freetorrent"}, client).
		Search(context.Background(), config(t, field("cookie", "session=abc"), field("freeleechOnly", true)), "Example", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].Title != "Example.Movie.2026.2160p" {
		t.Fatalf("title = %q", releases[0].Title)
	}
	if releases[0].InfoURL != "https://gazelle.test/torrents.php?id=10&torrentid=42" {
		t.Fatalf("info = %q", releases[0].InfoURL)
	}
}

func TestSearchUsesAuthorizationHeaderAndTokenDownloads(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("Authorization") != "token key-1" {
			t.Fatalf("auth = %q", req.Header.Get("Authorization"))
		}
		return testutil.Response(http.StatusOK, `{"status":"success","response":{"results":[{
			"groupId":"10",
			"groupName":"Example Album",
			"artist":"Example Artist",
			"groupYear":"2026",
			"releaseType":"Album",
			"torrents":[{
				"torrentId":42,
				"media":"WEB",
				"encoding":"FLAC",
				"format":"Lossless",
				"hasLog":true,
				"logScore":100,
				"hasCue":true,
				"time":"2026-07-03T00:00:00Z",
				"size":"2048",
				"seeders":"7",
				"leechers":"2",
				"canUseToken":true
			}]
		}]}}`), nil
	})

	releases, err := New(Options{
		Name:                    "Orpheus",
		DefaultBaseURL:          "https://orpheus.test/",
		AuthHeader:              "Authorization",
		AuthTokenPrefix:         "token ",
		DownloadPath:            "/ajax.php",
		PreferMusicTitle:        true,
		SupportsFreeleechTokens: true,
	}, client).Search(context.Background(), config(t,
		field("apiKey", "key-1"),
		field("useFreeleechToken", 1),
	), "Example", "music")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example Artist - Example Album (2026) [Album] [Lossless FLAC / WEB / Log (100%) / Cue]" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://gazelle.test/ajax.php?action=download&id=42&usetoken=1" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
}

func TestSearchUsesFreeloadQueryAndFilter(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("Authorization") != "key-1" {
			t.Fatalf("auth = %q", req.Header.Get("Authorization"))
		}
		if req.URL.Query().Get("freetorrent") != "4" {
			t.Fatalf("freetorrent = %q", req.URL.Query().Get("freetorrent"))
		}
		return testutil.Response(http.StatusOK, `{"status":"success","response":{"results":[{
			"groupId":"10",
			"groupName":"Example Album",
			"artist":"Example Artist",
			"groupYear":"2026",
			"torrents":[{
				"torrentId":41,
				"media":"WEB",
				"encoding":"MP3",
				"format":"320",
				"time":"2026-07-03T00:00:00Z",
				"size":"1024",
				"seeders":"1",
				"leechers":"1",
				"isFreeload":false
			},{
				"torrentId":42,
				"media":"WEB",
				"encoding":"FLAC",
				"format":"Lossless",
				"time":"2026-07-03T00:00:00Z",
				"size":"2048",
				"seeders":"7",
				"leechers":"2",
				"isFreeload":true
			}]
		}]}}`), nil
	})

	releases, err := New(Options{
		Name:           "Redacted",
		DefaultBaseURL: "https://redacted.test/",
		AuthHeader:     "Authorization",
		FreeleechParam: "freetorrent",
		FreeleechValue: "4",
		DownloadPath:   "/ajax.php",
	}, client).Search(context.Background(), config(t,
		field("apiKey", "key-1"),
		field("freeloadOnly", true),
	), "", "music")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].DownloadURL != "https://gazelle.test/ajax.php?action=download&id=42" {
		t.Fatalf("download = %q", releases[0].DownloadURL)
	}
}

func config(t *testing.T, fields ...map[string]any) engine.Config {
	t.Helper()
	raw, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	return engine.Config{
		ID:       "idx",
		Name:     "Gazelle",
		Protocol: "torrent",
		BaseURL:  "https://gazelle.test/",
		Fields:   raw,
	}
}

func field(name string, value any) map[string]any {
	return map[string]any{"name": name, "value": value}
}
