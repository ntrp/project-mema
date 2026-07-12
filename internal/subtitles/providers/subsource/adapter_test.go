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
	case "/api/search":
		if req.Header.Get("Authorization") != "token" { return &http.Response{StatusCode: 401, Body: io.NopCloser(strings.NewReader("no"))}, nil }
		return jsonResp(`{"results":[{"linkName":"fixture-2024","title":"Fixture"}]}`), nil
	case "/api/movie":
		return jsonResp(`{"subtitles":[{"id":42,"language":"eng","release":"Fixture.2024.WEB","download":"/api/downloadSubtitle?subtitleId=42"}]}`), nil
	case "/api/downloadSubtitle":
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Disposition": []string{"attachment; filename=subsource.zip"}}, Body: io.NopCloser(bytes.NewReader(zipBody()))}, nil
	default:
		return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("missing"))}, nil
	}
}

func TestSearchUsesSubsourceMultiStepFlow(t *testing.T) {
	s := &stub{}
	key := "token"
	got, err := Adapter.Search(context.Background(), s, providercore.Config{BaseURL: "https://subsource.net", APIKey: &key}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture", LanguageID: "eng"})
	if err != nil { t.Fatalf("Search error: %v", err) }
	if len(got) != 1 || got[0].ProviderName != "subsource" || got[0].FileID != 42 || got[0].ReleaseName != "Fixture.2024.WEB" { t.Fatalf("unexpected candidates: %#v", got) }
	if len(s.reqs) != 2 || s.reqs[0].URL.Path != "/api/search" || s.reqs[1].URL.Query().Get("movieName") != "fixture-2024" { t.Fatalf("unexpected requests: %#v", s.reqs) }
}

func TestDownloadExtractsArchive(t *testing.T) {
	s := &stub{}
	key := "token"
	dl, err := Adapter.Download(context.Background(), s, providercore.Config{BaseURL: "https://subsource.net", APIKey: &key}, providercore.Candidate{SourceURL: "/api/downloadSubtitle?subtitleId=42"})
	if err != nil { t.Fatalf("Download error: %v", err) }
	if !strings.Contains(string(dl.Content), "subsource subtitle") { t.Fatalf("unexpected content: %q", dl.Content) }
	if s.reqs[0].URL.Query().Get("subtitleId") != "42" { t.Fatalf("missing subtitle id: %s", s.reqs[0].URL.String()) }
}

func TestRequiresAPIKey(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &stub{}, providercore.Config{}, providercore.SearchRequest{Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "apiKey") { t.Fatalf("expected apiKey prerequisite, got %v", err) }
}

func jsonResp(body string) *http.Response { return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body))} }

func zipBody() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("fixture.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nsubsource subtitle\n"))
	_ = zw.Close()
	return buf.Bytes()
}
