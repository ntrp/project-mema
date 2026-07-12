package providercore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/security"
)

func TestHTTPClientValidatesAllowlistBeforeRequest(t *testing.T) {
	transport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("blocked request reached transport")
		return nil, nil
	})
	client := NewHTTPClient("opensubtitlescom", &http.Client{Transport: transport})
	_, _, err := client.Get(context.Background(), "https://evil.example/subs", 1024)
	if !errors.Is(err, security.ErrOutboundURLBlocked) {
		t.Fatalf("expected outbound block, got %v", err)
	}
}

func TestHTTPClientAllowsConfiguredHostAndHeaders(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("User-Agent") != "mema-test" {
			t.Fatalf("headers = %#v", request.Header)
		}
		return response(http.StatusOK, "ok"), nil
	})
	client := NewHTTPClient("opensubtitlescom", &http.Client{Transport: transport})
	client.SetHeader("User-Agent", "mema-test")
	body, res, err := client.Get(context.Background(), "https://api.opensubtitles.com/search", 1024)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(body) != "ok" || res.StatusCode != http.StatusOK {
		t.Fatalf("body=%q response=%#v", body, res)
	}
}

func TestHTTPClientValidatesRedirects(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusFound,
			Header:     http.Header{"Location": []string{"http://localhost/blocked"}},
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    request,
		}, nil
	})
	client := NewHTTPClient("opensubtitlescom", &http.Client{Transport: transport})
	_, _, err := client.Get(context.Background(), "https://api.opensubtitles.com/redirect", 1024)
	if !errors.Is(err, security.ErrOutboundURLBlocked) {
		t.Fatalf("expected redirect block, got %v", err)
	}
}

func TestHTTPClientResponseFailures(t *testing.T) {
	for name, res := range map[string]*http.Response{
		"status": response(http.StatusInternalServerError, "bad"),
		"large":  response(http.StatusOK, "too large"),
	} {
		transport := roundTripFunc(func(*http.Request) (*http.Response, error) { return res, nil })
		client := NewHTTPClient("opensubtitlescom", &http.Client{Transport: transport})
		_, _, err := client.Get(context.Background(), "https://api.opensubtitles.com/test", 3)
		if err == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
}

func TestDownloadHTTPClientUsesDownloadHosts(t *testing.T) {
	transport := roundTripFunc(func(*http.Request) (*http.Response, error) { return response(http.StatusOK, "download"), nil })
	client := NewDownloadHTTPClient("opensubtitlescom", &http.Client{Transport: transport})
	body, _, err := client.Get(context.Background(), "https://dl.opensubtitles.com/file", 1024)
	if err != nil {
		t.Fatalf("download Get failed: %v", err)
	}
	if string(body) != "download" {
		t.Fatalf("body = %q", body)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }

func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: http.Header{}}
}
