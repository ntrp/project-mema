package subsource

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

type stub struct{ reqs []*http.Request }

func (s *stub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.reqs = append(s.reqs, req)
	switch req.URL.Path {
	case "/api/v1/movies/search":
		return jsonResp(`{"data":[{"movieId":7,"title":"Fixture","releaseYear":2024}]}`), nil
	case "/api/v1/subtitles":
		return jsonResp(`{"success":true,"data":[{"subtitleId":42,"language":"english","releaseInfo":["Fixture.2024.WEB"],"link":"/subtitle/42"}]}`), nil
	case "/api/v1/subtitles/42/download":
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipBody()))}, nil
	default:
		return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("missing"))}, nil
	}
}

func TestSearchUsesSubsourceV1MultiStepFlow(t *testing.T) {
	s := &stub{}
	key := "token"
	year := int32(2024)
	got, err := Adapter.Search(context.Background(), s, providercore.Config{BaseURL: "https://api.subsource.net", APIKey: &key}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture", LanguageID: "english", Year: &year, MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt123"}}})
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if len(got) != 1 || got[0].FileID != 42 || got[0].ReleaseName != "Fixture.2024.WEB" {
		t.Fatalf("unexpected candidates: %#v", got)
	}
	if len(s.reqs) != 2 || s.reqs[0].URL.Path != "/api/v1/movies/search" || s.reqs[0].URL.Query().Get("searchType") != "imdb" || s.reqs[0].URL.Query().Get("api_key") != "token" || s.reqs[1].URL.Query().Get("movieId") != "7" {
		t.Fatalf("unexpected requests: %#v", s.reqs)
	}
}

func TestDownloadExtractsArchive(t *testing.T) {
	s := &stub{}
	key := "token"
	dl, err := Adapter.Download(context.Background(), s, providercore.Config{BaseURL: "https://api.subsource.net", APIKey: &key}, providercore.Candidate{FileID: 42})
	if err != nil {
		t.Fatalf("Download error: %v", err)
	}
	if !strings.Contains(string(dl.Content), "subsource subtitle") || s.reqs[0].URL.Query().Get("api_key") != "token" {
		t.Fatalf("unexpected download: %q request=%s", dl.Content, s.reqs[0].URL)
	}
}

func TestRequiresAPIKey(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &stub{}, providercore.Config{}, providercore.SearchRequest{Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "apiKey") {
		t.Fatalf("expected apiKey prerequisite, got %v", err)
	}
}
func jsonResp(body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body))}
}
func zipBody() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("fixture.srt")
	_, _ = w.Write([]byte("subsource subtitle"))
	_ = zw.Close()
	return buf.Bytes()
}
