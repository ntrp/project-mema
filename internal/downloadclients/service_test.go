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
