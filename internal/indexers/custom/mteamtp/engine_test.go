package mteamtp

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchParsesMTeamResponse(t *testing.T) {
	key := "mteam-key"
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.m-team.test" {
			t.Fatalf("host = %s", req.URL.Host)
		}
		if req.URL.Path != "/api/torrent/search" {
			t.Fatalf("path = %s", req.URL.Path)
		}
		if req.Header.Get("x-api-key") != key {
			t.Fatalf("x-api-key = %q", req.Header.Get("x-api-key"))
		}
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		if payload["keyword"] != "Example" {
			t.Fatalf("keyword = %#v", payload["keyword"])
		}
		body := `{"message":"SUCCESS","data":{"data":[{"id":"77","name":"Example Movie 2026 1080p","smallDescr":"demo","category":"401","numfiles":"3","size":"8192","status":{"createdDate":"2026-07-04 12:00:00","discount":"FREE","timesCompleted":"5","seeders":"14","leechers":"4"}}]}}`
		return testutil.Response(http.StatusOK, body), nil
	})
	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "mteam",
		Name:     "M-Team - TP",
		Protocol: "torrent",
		BaseURL:  "https://kp.m-team.test/",
		APIKey:   &key,
	}, "Example", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if releases[0].DownloadURL != "https://api.m-team.test/api/torrent/genDlToken?id=77" {
		t.Fatalf("download = %q", releases[0].DownloadURL)
	}
	if releases[0].Peers == nil || *releases[0].Peers != 18 {
		t.Fatalf("peers = %#v", releases[0].Peers)
	}
}
