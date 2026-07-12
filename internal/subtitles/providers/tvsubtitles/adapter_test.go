package tvsubtitles

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
	s.paths = append(s.paths, req.Method+" "+req.URL.Path)
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("1\n00:00:01,000 --> 00:00:02,000\nfixture\n"))}, nil
	}
	switch req.URL.Path {
	case "/search.php":
		return html(`<a href="/tvshow-42.html">Fixture</a>`), nil
	case "/season-1.html":
		return html(`<a href="/episode-1x02.html">Fixture 1x02</a>`), nil
	case "/episode-1x02.html":
		return html(`<a data-lang="eng" href="/download-99.html">Fixture.Release</a>`), nil
	default:
		return html(`<html>ok</html>`), nil
	}
}

func html(body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"text/html"}}, Body: io.NopCloser(strings.NewReader(body))}
}

func TestSearchUsesPostTraversal(t *testing.T) {
	season, episode := int32(1), int32(2)
	stub := &providerStub{}
	candidates, err := Adapter.Search(context.Background(), stub, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture", LanguageID: "eng", SeasonNumber: &season, EpisodeNumber: &episode})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].SourceURL != "https://example.test/download-99.html" || candidates[0].LanguageID != "eng" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	want := []string{"POST /search.php", "GET /season-1.html", "GET /episode-1x02.html"}
	if strings.Join(stub.paths, ",") != strings.Join(want, ",") {
		t.Fatalf("unexpected traversal: %#v", stub.paths)
	}
}

func TestRejectsUnsupportedMediaType(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected unsupported media type error, got %v", err)
	}
}
