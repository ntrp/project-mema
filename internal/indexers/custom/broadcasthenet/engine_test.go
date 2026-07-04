package broadcasthenet

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"media-manager/internal/indexers/custom/testutil"
	"media-manager/internal/indexers/engine"
)

func TestSearchMapsBroadcastheNetResponse(t *testing.T) {
	client := testutil.FakeHTTPDoer(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost || req.URL.String() != "https://api.broadcasthe.test" {
			t.Fatalf("request = %s %s", req.Method, req.URL.String())
		}
		if req.Header.Get("Accept") != "application/json-rpc, application/json" {
			t.Fatalf("accept = %q", req.Header.Get("Accept"))
		}
		var body map[string]any
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		assertRPCPayload(t, body)
		return testutil.Response(http.StatusOK, `{
			"result": {
				"results": 1,
				"torrents": {
					"1": {
						"GroupID": 10,
						"TorrentID": 20,
						"Category": "Episode",
						"Snatched": 3,
						"Seeders": 11,
						"Leechers": 2,
						"Source": "WEB",
						"Container": "MKV",
						"Codec": "H.264",
						"Resolution": "1080p",
						"Origin": "Internal",
						"ReleaseName": "Example.Show.S01E01.1080p",
						"Size": 2048,
						"Time": 1783166400,
						"InfoHash": "abc",
						"DownloadURL": "http://broadcasthe.net/torrents.php?action=download&id=20"
					}
				}
			},
			"error": null
		}`), nil
	})

	releases, err := New(client).Search(context.Background(), engine.Config{
		ID:       "idx-btn",
		Name:     "BroadcastheNet",
		Protocol: "torrent",
		BaseURL:  "https://api.broadcasthe.test/",
		APIKey:   stringPtr("key-1"),
	}, "Example Show", "series")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example.Show.S01E01.1080p" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "https://broadcasthe.net/torrents.php?action=download&id=20" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.InfoURL != "https://broadcasthe.net/torrents.php?id=10&torrentid=20" {
		t.Fatalf("info = %q", release.InfoURL)
	}
	if release.Seeders == nil || *release.Seeders != 11 || release.Peers == nil || *release.Peers != 13 {
		t.Fatalf("seeders = %#v peers = %#v", release.Seeders, release.Peers)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}

func assertRPCPayload(t *testing.T, body map[string]any) {
	t.Helper()
	if body["method"] != "getTorrents" || body["jsonrpc"] != "2.0" {
		t.Fatalf("body = %#v", body)
	}
	params, ok := body["params"].([]any)
	if !ok || len(params) != 4 || params[0] != "key-1" {
		t.Fatalf("params = %#v", body["params"])
	}
	filters, ok := params[1].(map[string]any)
	if !ok || filters["Search"] != "Example%Show" {
		t.Fatalf("filters = %#v", params[1])
	}
	if params[2].(float64) != 100 || params[3].(float64) != 0 {
		t.Fatalf("pagination params = %#v", params)
	}
}

func stringPtr(value string) *string {
	return &value
}
