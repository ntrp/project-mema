package yifysubtitles

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type providerStub struct{ paths []string }

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.paths = append(s.paths, req.URL.Path)
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("zip"))}, nil
	}
	switch req.URL.Path {
	case "/movie-imdb/tt1234567":
		return html(`<a href="/subtitle/fixture-english">English</a>`), nil
	case "/subtitle/fixture-english":
		return html(`<main data-language="eng" data-release="Fixture.YIFY"><a href="/subtitle/fixture-english.zip">Download</a></main>`), nil
	default:
		return html(`<html>ok</html>`), nil
	}
}

func html(body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"text/html"}}, Body: io.NopCloser(strings.NewReader(body))}
}

func TestSearchUsesIMDbDetailAndZipDownloadURL(t *testing.T) {
	stub := &providerStub{}
	candidates, err := Adapter.Search(context.Background(), stub, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", LanguageID: "eng", MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt1234567"}}})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].SourceURL != "https://example.test/subtitle/fixture-english.zip" || candidates[0].SourceRef != "https://example.test/subtitle/fixture-english" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	if strings.Join(stub.paths, ",") != "/movie-imdb/tt1234567,/subtitle/fixture-english" {
		t.Fatalf("unexpected requests: %#v", stub.paths)
	}
}

func TestRejectsUnsupportedMediaType(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected unsupported media type error, got %v", err)
	}
}
