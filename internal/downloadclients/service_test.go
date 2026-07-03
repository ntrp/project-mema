package downloadclients

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTransmissionHandlesSessionChallenge(t *testing.T) {
	calls := 0
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.String() != "http://transmission.local/transmission/rpc" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if calls == 1 {
			return response(http.StatusConflict, "", map[string]string{transmissionSessionHeader: "session-1"}), nil
		}
		if got := r.Header.Get(transmissionSessionHeader); got != "session-1" {
			t.Fatalf("session header = %q", got)
		}
		body, err := json.Marshal(map[string]interface{}{
			"result": "success",
			"arguments": map[string]interface{}{
				"version":     "4.0.0",
				"rpc-version": 17,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return response(http.StatusOK, string(body), nil), nil
	})

	result := NewService(client).Test(context.Background(), Config{
		Type:    "transmission",
		BaseURL: "http://transmission.local",
	})

	if !result.Success {
		t.Fatalf("expected success, got %#v", result)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestSABnzbdVersion(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("mode"); got != "version" {
			t.Fatalf("mode = %q", got)
		}
		body, err := json.Marshal(map[string]string{"version": "4.3.0"})
		if err != nil {
			t.Fatal(err)
		}
		return response(http.StatusOK, string(body), nil), nil
	})

	apiKey := "secret"
	result := NewService(client).Test(context.Background(), Config{
		Type:    "sabnzbd",
		BaseURL: "http://sabnzbd.local",
		APIKey:  &apiKey,
	})

	if !result.Success {
		t.Fatalf("expected success, got %#v", result)
	}
}

func TestSABnzbdStatusFindsQueueSlotByNZOID(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query().Get("search"); got != "" {
			t.Fatalf("search = %q", got)
		}
		body, err := json.Marshal(map[string]interface{}{
			"queue": map[string]interface{}{
				"slots": []map[string]string{
					{
						"nzo_id":     "download-1",
						"status":     "Downloading",
						"percentage": "42",
					},
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		return response(http.StatusOK, string(body), nil), nil
	})

	apiKey := "secret"
	result := NewService(client).Status(context.Background(), Config{
		Type:    "sabnzbd",
		BaseURL: "http://sabnzbd.local",
		APIKey:  &apiKey,
	}, StatusRequest{DownloadID: "download-1"})

	if !result.Success || !result.Found || result.Status != "downloading" {
		t.Fatalf("unexpected status result %#v", result)
	}
	if result.ProgressPercent == nil || *result.ProgressPercent != 42 {
		t.Fatalf("progress = %#v", result.ProgressPercent)
	}
}

func TestSABnzbdStatusFindsCompletedHistoryByNZOID(t *testing.T) {
	call := 0
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		call++
		if got := r.URL.Query().Get("search"); got != "" {
			t.Fatalf("search = %q", got)
		}
		if r.URL.Query().Get("mode") == "queue" {
			return jsonResponse(t, map[string]interface{}{"queue": map[string]interface{}{"slots": []map[string]string{}}}), nil
		}
		return jsonResponse(t, map[string]interface{}{
			"history": map[string]interface{}{
				"slots": []map[string]string{
					{
						"nzo_id":  "download-2",
						"status":  "Completed",
						"storage": "/downloads/Toy.Story.5.2026.1080p/Toy.Story.5.2026.1080p.mkv",
					},
				},
			},
		}), nil
	})

	apiKey := "secret"
	result := NewService(client).Status(context.Background(), Config{
		Type:    "sabnzbd",
		BaseURL: "http://sabnzbd.local",
		APIKey:  &apiKey,
	}, StatusRequest{DownloadID: "download-2"})

	if call != 2 {
		t.Fatalf("calls = %d", call)
	}
	if !result.Success || !result.Found || result.Status != "completed" {
		t.Fatalf("unexpected status result %#v", result)
	}
	if result.ProgressPercent == nil || *result.ProgressPercent != 100 {
		t.Fatalf("progress = %#v", result.ProgressPercent)
	}
	if len(result.Files) != 1 || result.Files[0].Path == "" || !result.Files[0].Complete {
		t.Fatalf("files = %#v", result.Files)
	}
}

func TestSCNIntegrations003SABnzbdAddAndCancel(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Query().Get("mode") {
		case "addurl":
			if got := r.URL.Query().Get("name"); got != "https://indexer.test/release.nzb" {
				t.Fatalf("name = %q", got)
			}
			if got := r.URL.Query().Get("cat"); got != "movies" {
				t.Fatalf("cat = %q", got)
			}
			return jsonResponse(t, map[string]interface{}{"status": true, "nzo_ids": []string{"nzo-1"}}), nil
		case "queue":
			if r.URL.Query().Get("name") != "delete" || r.URL.Query().Get("value") != "nzo-1" {
				t.Fatalf("unexpected cancel query %s", r.URL.RawQuery)
			}
			return jsonResponse(t, map[string]interface{}{"status": true}), nil
		default:
			t.Fatalf("unexpected mode %q", r.URL.Query().Get("mode"))
			return nil, nil
		}
	})

	category := "movies"
	service := NewService(client)
	config := Config{Type: "sabnzbd", BaseURL: "http://sabnzbd.local", Category: &category}
	added := service.Add(context.Background(), config, AddRequest{URL: "https://indexer.test/release.nzb"})
	if !added.Success || added.DownloadID != "nzo-1" {
		t.Fatalf("add result = %#v", added)
	}
	cancelled := service.Cancel(context.Background(), config, CancelRequest{DownloadID: added.DownloadID})
	if !cancelled.Success {
		t.Fatalf("cancel result = %#v", cancelled)
	}
}

func TestSCNIntegrations003TransmissionAddStatusAndCancel(t *testing.T) {
	call := 0
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		call++
		if call == 1 {
			return response(http.StatusConflict, "", map[string]string{transmissionSessionHeader: "session-1"}), nil
		}
		if call == 2 && r.Header.Get(transmissionSessionHeader) != "session-1" {
			t.Fatalf("session header = %q", r.Header.Get(transmissionSessionHeader))
		}
		if call != 2 && r.Header.Get(transmissionSessionHeader) != "" {
			t.Fatalf("session header = %q", r.Header.Get(transmissionSessionHeader))
		}
		var request struct {
			Method    string                 `json:"method"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		switch request.Method {
		case "torrent-add":
			return jsonResponse(t, map[string]interface{}{
				"result": "success",
				"arguments": map[string]interface{}{
					"torrent-added": map[string]int{"id": 42},
				},
			}), nil
		case "torrent-get":
			return jsonResponse(t, map[string]interface{}{
				"result": "success",
				"arguments": map[string]interface{}{
					"torrents": []map[string]interface{}{
						{
							"id":            42,
							"status":        4,
							"percentDone":   0.5,
							"leftUntilDone": 100,
							"downloadDir":   "/downloads",
							"name":          "Scenario.Movie",
							"files": []map[string]interface{}{
								{"name": "Scenario.Movie.mkv", "length": 1000, "bytesCompleted": 500},
							},
						},
					},
				},
			}), nil
		case "torrent-remove":
			return jsonResponse(t, map[string]interface{}{"result": "success"}), nil
		default:
			t.Fatalf("unexpected method %q", request.Method)
			return nil, nil
		}
	})

	service := NewService(client)
	config := Config{Type: "transmission", BaseURL: "http://transmission.local"}
	added := service.Add(context.Background(), config, AddRequest{URL: "magnet:?xt=urn:btih:scenario"})
	if !added.Success || added.DownloadID != "42" {
		t.Fatalf("add result = %#v", added)
	}
	status := service.Status(context.Background(), config, StatusRequest{DownloadID: added.DownloadID})
	if !status.Success || !status.Found || status.Status != "downloading" {
		t.Fatalf("status result = %#v", status)
	}
	if status.ProgressPercent == nil || *status.ProgressPercent != 50 {
		t.Fatalf("progress = %#v", status.ProgressPercent)
	}
	if len(status.Files) != 0 {
		t.Fatalf("active downloads should not report completed files: %#v", status.Files)
	}
	cancelled := service.Cancel(context.Background(), config, CancelRequest{DownloadID: added.DownloadID})
	if !cancelled.Success {
		t.Fatalf("cancel result = %#v", cancelled)
	}
}

type fakeHTTPDoer func(req *http.Request) (*http.Response, error)

func (f fakeHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func response(statusCode int, body string, headers map[string]string) *http.Response {
	header := http.Header{}
	for key, value := range headers {
		header.Set(key, value)
	}
	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func jsonResponse(t *testing.T, payload interface{}) *http.Response {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	return response(http.StatusOK, string(body), nil)
}
