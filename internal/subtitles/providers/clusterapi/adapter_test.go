package clusterapi

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type fixtureService struct{ client *http.Client }

func (s fixtureService) DoProviderRequest(req *http.Request, _ string, _ bool) (*http.Response, error) {
	return s.client.Do(req)
}

func TestSearchParsesStructuredJSONAndSendsAuth(t *testing.T) {
	var sawToken bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") == "secret" { sawToken = true }
		if got := r.URL.Query().Get("title"); got != "Example" { t.Fatalf("title query = %q", got) }
		_, _ = w.Write([]byte(`{"results":[{"id":12,"language":"en","release_name":"Example.WEB","download_url":"/download/12.zip","download_count":7}]}`))
	}))
	defer server.Close()
	secret := "secret"
	adapter := Adapter{Spec: Spec{Key: "subx", DefaultBaseURL: server.URL, SearchPath: "/search", RequiredSecret: "apiKey", SecretHeader: "X-API-Key"}}
	candidates, err := adapter.Search(context.Background(), fixtureService{server.Client()}, providercore.Config{APIKey: &secret}, providercore.SearchRequest{MediaType: "movie", Title: "Example", LanguageID: "en"})
	if err != nil { t.Fatalf("Search returned error: %v", err) }
	if !sawToken { t.Fatal("expected API key header") }
	if len(candidates) != 1 || candidates[0].FileID != 12 || candidates[0].ReleaseName != "Example.WEB" || candidates[0].SourceURL != "/download/12.zip" { t.Fatalf("unexpected candidates: %#v", candidates) }
}

func TestDownloadExtractsBestSubtitleFromZip(t *testing.T) {
	archive := zipBytes(t, map[string]string{"readme.txt": "ignore", "movie.srt": "subtitle"})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(archive) }))
	defer server.Close()
	adapter := Adapter{Spec: Spec{Key: "subdl", DefaultBaseURL: server.URL}}
	download, err := adapter.Download(context.Background(), fixtureService{server.Client()}, providercore.Config{}, providercore.Candidate{SourceURL: server.URL + "/file.zip"})
	if err != nil { t.Fatalf("Download returned error: %v", err) }
	if string(download.Content) != "subtitle" { t.Fatalf("content = %q", download.Content) }
}

func TestValidationHonorsMediaAndIMDbPrerequisites(t *testing.T) {
	adapter := Adapter{Spec: Spec{Key: "subsro", SeriesOnly: true, RequireIMDb: true}}
	_, err := adapter.Search(context.Background(), fixtureService{http.DefaultClient}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Bad"})
	if err == nil || !strings.Contains(err.Error(), "unsupported media type") { t.Fatalf("expected unsupported media error, got %v", err) }
	_, err = adapter.Search(context.Background(), fixtureService{http.DefaultClient}, providercore.Config{}, providercore.SearchRequest{MediaType: "serie", Title: "No ID"})
	if err == nil || !strings.Contains(err.Error(), "imdb id is required") { t.Fatalf("expected imdb prerequisite, got %v", err) }
}

func TestWhisperSearchChecksCommandAndReturnsSyntheticCandidate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`ok`)) }))
	defer server.Close()
	adapter := Adapter{Spec: Spec{Key: "whisperai", DefaultBaseURL: server.URL, SearchPath: "/transcribe", TestPath: "/health", Local: true, CommandName: "ffmpeg", CommandArgs: []string{"-version"}}}
	config := providercore.Config{BaseURL: server.URL, CommandRunner: func(context.Context, string, ...string) ([]byte, error) { return []byte("ok"), nil }}
	candidates, err := adapter.Search(context.Background(), fixtureService{server.Client()}, config, providercore.SearchRequest{MediaType: "movie", Title: "Local", LanguageID: "en"})
	if err != nil { t.Fatalf("Search returned error: %v", err) }
	if len(candidates) != 1 || candidates[0].ProviderName != "whisperai" || !strings.Contains(candidates[0].SourceURL, "/transcribe") { t.Fatalf("unexpected candidates: %#v", candidates) }
}

func TestClassifiesAuthAndRateLimitHTTPFailures(t *testing.T) {
	for status, want := range map[int]string{403: "provider_prerequisite_missing", 429: "provider_broken_upstream"} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(status) }))
		adapter := Adapter{Spec: Spec{Key: "subsource", DefaultBaseURL: server.URL, SearchPath: "/search"}}
		_, err := adapter.Search(context.Background(), fixtureService{server.Client()}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Example"})
		server.Close()
		if err == nil || !strings.Contains(err.Error(), want) { t.Fatalf("status %d: expected %s, got %v", status, want, err) }
	}
}

func zipBytes(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil { t.Fatal(err) }
		if _, err := w.Write([]byte(content)); err != nil { t.Fatal(err) }
	}
	if err := zw.Close(); err != nil { t.Fatal(err) }
	return buf.Bytes()
}
