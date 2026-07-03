package indexers

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTorznabCaps(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if got := r.URL.Query().Get("t"); got != "caps" {
			t.Fatalf("t = %q", got)
		}
		return response(http.StatusOK, `<caps><server title="Indexer" version="1.0"/><limits max="100" default="50"/><categories><category id="5000" name="TV"><subcat id="5070" name="Anime"/></category></categories></caps>`), nil
	})

	result := NewService(client).Test(context.Background(), Config{
		Protocol: "torrent",
		BaseURL:  "http://indexer.local/api",
	})

	if !result.Success {
		t.Fatalf("expected success, got %#v", result)
	}
	if got := result.Details["categoryCount"]; got != 2 {
		t.Fatalf("categoryCount = %v", got)
	}
}

func TestCapsDetectsWebUIHTML(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		resp := response(http.StatusOK, `<!doctype html><html><body>Prowlarr</body></html>`)
		resp.Header.Set("Content-Type", "text/html; charset=utf-8")
		return resp, nil
	})

	result := NewService(client).Test(context.Background(), Config{
		Protocol: "torrent",
		BaseURL:  "https://prowlarr.local/",
	})

	if result.Success {
		t.Fatalf("expected failure, got %#v", result)
	}
	if !strings.Contains(result.Message, "not the Prowlarr web UI root") {
		t.Fatalf("unexpected message %q", result.Message)
	}
}

func TestRSSProtocolIsNotAStandaloneIndexer(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		return response(http.StatusOK, `<rss><channel><title>Feed</title><item><title>One</title></item></channel></rss>`), nil
	})

	result := NewService(client).Test(context.Background(), Config{
		Protocol: "rss",
		BaseURL:  "http://rss.local/feed",
	})

	if result.Success {
		t.Fatalf("expected failure, got %#v", result)
	}
}

type fakeHTTPDoer func(req *http.Request) (*http.Response, error)

func (f fakeHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func response(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
