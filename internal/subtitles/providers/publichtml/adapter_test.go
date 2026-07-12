package publichtml

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type roundTripService func(*http.Request, string, bool) (*http.Response, error)

func (f roundTripService) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	return f(req, providerType, isDownload)
}

func TestAdapterSearchParsesFixtureAndDownload(t *testing.T) {
	adapter := New(Spec{Key: "fixture", BaseURL: "https://fixture.example", MediaTypes: []string{"movie"}})
	var sawSearch bool
	svc := roundTripService(func(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
		if providerType != "fixture" {
			t.Fatalf("providerType = %s", providerType)
		}
		if isDownload {
			return textResponse(200, "1\n00:00:01,000 --> 00:00:02,000\nhello\n"), nil
		}
		if req.URL.Path == "/search" && req.URL.Query().Get("q") == "Fixture Movie" && req.URL.Query().Get("imdb") == "tt123" {
			sawSearch = true
		}
		return textResponse(200, `<ul><li class="subtitle" data-lang="ell" data-format="srt" data-release="Fixture.2024"><a href="/download/fixture.srt">download</a></li></ul>`), nil
	})
	candidates, err := adapter.Search(context.Background(), svc, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture Movie", LanguageID: "ell", MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt123"}}})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if !sawSearch || len(candidates) != 1 || candidates[0].SourceURL != "https://fixture.example/download/fixture.srt" || candidates[0].ReleaseName != "Fixture.2024" {
		t.Fatalf("unexpected candidates: %#v sawSearch=%v", candidates, sawSearch)
	}
	download, err := adapter.Download(context.Background(), svc, providercore.Config{}, candidates[0])
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if !strings.Contains(string(download.Content), "hello") {
		t.Fatalf("unexpected content: %q", string(download.Content))
	}
}

func TestAdapterExtractsArchiveDownloads(t *testing.T) {
	adapter := New(Spec{Key: "fixture", BaseURL: "https://fixture.example", MediaTypes: []string{"movie"}, Archive: true})
	archive := zipFixture(t, "nested/movie.srt", "subtitle body")
	svc := roundTripService(func(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
		return binaryResponse(200, "application/zip", archive), nil
	})
	download, err := adapter.Download(context.Background(), svc, providercore.Config{}, providercore.Candidate{SourceURL: "https://fixture.example/download/movie.zip"})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if string(download.Content) != "subtitle body" {
		t.Fatalf("unexpected archive content: %q", download.Content)
	}
}

func TestAdapterRejectsUnsupportedMediaType(t *testing.T) {
	adapter := New(Spec{Key: "fixture", BaseURL: "https://fixture.example", MediaTypes: []string{"serie"}})
	_, err := adapter.Search(context.Background(), nil, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected prerequisite error, got %v", err)
	}
}

func TestAdapterClassifiesBrokenLayout(t *testing.T) {
	adapter := New(Spec{Key: "fixture", BaseURL: "https://fixture.example", MediaTypes: []string{"movie"}})
	svc := roundTripService(func(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
		return textResponse(200, `<html><body>no links</body></html>`), nil
	})
	_, err := adapter.Search(context.Background(), svc, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "provider_broken_upstream") {
		t.Fatalf("expected broken upstream error, got %v", err)
	}
}

func textResponse(status int, body string) *http.Response {
	return binaryResponse(status, "text/html", []byte(body))
}

func binaryResponse(status int, contentType string, body []byte) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{contentType}}, Body: io.NopCloser(bytes.NewReader(body))}
}

func zipFixture(t *testing.T, name string, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	writer, err := zw.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
