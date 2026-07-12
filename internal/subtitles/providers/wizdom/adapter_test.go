package wizdom

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type providerStub struct{ rawQuery string }

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.rawQuery = req.URL.RawQuery
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("1\n00:00:01,000 --> 00:00:02,000\nfixture\n"))}, nil
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`[{"id":321,"versioname":"Fixture.WEB","language":"heb"}]`))}, nil
}

func TestSearchUsesWizdomIMDbAPI(t *testing.T) {
	stub := &providerStub{}
	candidates, err := Adapter.Search(context.Background(), stub, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt7654321"}}})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].FileID != 321 || candidates[0].SourceURL != "https://example.test/api/files/sub/321" || candidates[0].LanguageID != "heb" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	if !strings.Contains(stub.rawQuery, "action=by_id") || !strings.Contains(stub.rawQuery, "imdb=tt7654321") {
		t.Fatalf("unexpected query: %s", stub.rawQuery)
	}
}

func TestRequiresIMDbID(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected imdb prerequisite error, got %v", err)
	}
}

func TestDownloadUsesCandidateURL(t *testing.T) {
	dl, err := Adapter.Download(context.Background(), &providerStub{}, providercore.Config{}, providercore.Candidate{SourceURL: "https://example.test/api/files/sub/321", ReleaseName: "Fixture"})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if !strings.Contains(string(dl.Content), "fixture") {
		t.Fatalf("unexpected download: %q", dl.Content)
	}
}
