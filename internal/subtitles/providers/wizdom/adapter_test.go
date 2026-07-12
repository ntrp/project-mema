package wizdom

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

type providerStub struct{ path, referer string }

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.path, s.referer = req.URL.Path, req.Header.Get("Referer")
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipBody()))}, nil
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"subs":[{"id":321,"version":"Fixture.WEB"}]}`))}, nil
}

func TestSearchUsesWizdomReleaseAPI(t *testing.T) {
	stub := &providerStub{}
	candidates, err := Adapter.Search(context.Background(), stub, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt7654321"}}})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].FileID != 321 || candidates[0].SourceURL != "https://example.test/api/files/sub/321" || candidates[0].LanguageID != "heb" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	if stub.path != "/api/releases/tt7654321" {
		t.Fatalf("unexpected path: %s", stub.path)
	}
}

func TestRequiresIMDbID(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected imdb prerequisite error, got %v", err)
	}
}

func TestDownloadExtractsArchiveAndSetsReferer(t *testing.T) {
	stub := &providerStub{}
	dl, err := Adapter.Download(context.Background(), stub, providercore.Config{}, providercore.Candidate{SourceURL: "https://example.test/api/files/sub/321", SourceRef: "https://example.test/movies/tt1", ReleaseName: "Fixture"})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if !strings.Contains(string(dl.Content), "fixture") || stub.referer == "" {
		t.Fatalf("unexpected download=%q referer=%q", dl.Content, stub.referer)
	}
}

func zipBody() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("fixture.srt")
	_, _ = w.Write([]byte("fixture"))
	_ = zw.Close()
	return buf.Bytes()
}
