package downloadclients

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSCNIntegrations003TransmissionCompletedStatusIncludesFiles(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		var request struct {
			Method string `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.Method != "torrent-get" {
			t.Fatalf("method = %q", request.Method)
		}
		return jsonResponse(t, map[string]interface{}{
			"result": "success",
			"arguments": map[string]interface{}{
				"torrents": []map[string]interface{}{
					{
						"id":            42,
						"status":        6,
						"isFinished":    true,
						"percentDone":   1,
						"leftUntilDone": 0,
						"downloadDir":   "/downloads",
						"name":          "Scenario.Movie",
						"files": []map[string]interface{}{
							{"name": "Scenario.Movie/feature.mkv", "length": 1000, "bytesCompleted": 1000},
							{"name": "Scenario.Movie/sample.mkv", "length": 100, "bytesCompleted": 50},
						},
					},
				},
			},
		}), nil
	})

	result := NewService(client).Status(context.Background(), Config{
		Type:    "transmission",
		BaseURL: "http://transmission.local",
	}, StatusRequest{DownloadID: "42"})

	if !result.Success || !result.Found || result.Status != "completed" {
		t.Fatalf("status result = %#v", result)
	}
	if result.ProgressPercent == nil || *result.ProgressPercent != 100 {
		t.Fatalf("progress = %#v", result.ProgressPercent)
	}
	if len(result.Files) != 2 {
		t.Fatalf("files = %#v", result.Files)
	}
	if result.Files[0].Path != "/downloads/Scenario.Movie/feature.mkv" || !result.Files[0].Complete {
		t.Fatalf("first file = %#v", result.Files[0])
	}
	if result.Files[1].Complete {
		t.Fatalf("second file should be incomplete: %#v", result.Files[1])
	}
}

func TestSCNIntegrations003SABnzbdHistoryFailureIsReported(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.Query().Get("mode") == "queue" {
			return jsonResponse(t, map[string]interface{}{"queue": map[string]interface{}{"slots": []map[string]string{}}}), nil
		}
		return jsonResponse(t, map[string]interface{}{
			"history": map[string]interface{}{
				"slots": []map[string]string{{
					"nzo_id":       "nzo-failed",
					"status":       "Failed",
					"fail_message": "Unpack failed",
				}},
			},
		}), nil
	})

	result := NewService(client).Status(context.Background(), Config{
		Type:    "sabnzbd",
		BaseURL: "http://sabnzbd.local",
	}, StatusRequest{DownloadID: "nzo-failed"})

	if !result.Success || !result.Found || result.Status != "failed" {
		t.Fatalf("status result = %#v", result)
	}
	if result.Message != "Unpack failed" {
		t.Fatalf("message = %q", result.Message)
	}
}
