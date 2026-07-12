package subsro

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
	if strings.Contains(req.URL.Path, "download") {
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Disposition": []string{"attachment; filename=subsro.zip"}}, Body: io.NopCloser(bytes.NewReader(zipBody()))}, nil
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"subtitles":[{"id":9,"language":"ron","release":"Fixture.2024.RO","download_url":"/download/9.zip"}]}`))}, nil
}

func TestSearchRequiresIMDbAndAddsAPIKey(t *testing.T) {
	s := &stub{}
	key := "secret"
	got, err := Adapter.Search(context.Background(), s, providercore.Config{BaseURL: "https://api.subs.ro", APIKey: &key}, providercore.SearchRequest{Title: "Fixture", LanguageID: "ron", MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt1234567"}}})
	if err != nil { t.Fatalf("Search error: %v", err) }
	if len(got) != 1 || got[0].ProviderName != "subsro" || got[0].FileID != 9 || got[0].LanguageID != "ron" { t.Fatalf("unexpected candidates: %#v", got) }
	q := s.reqs[0].URL.Query()
	if q.Get("key") != "secret" || q.Get("imdb") != "tt1234567" { t.Fatalf("unexpected query: %s", s.reqs[0].URL.RawQuery) }
}

func TestDownloadAppendsKeyAndExtractsArchive(t *testing.T) {
	s := &stub{}
	key := "secret"
	dl, err := Adapter.Download(context.Background(), s, providercore.Config{BaseURL: "https://api.subs.ro", APIKey: &key}, providercore.Candidate{SourceURL: "/download/9.zip"})
	if err != nil { t.Fatalf("Download error: %v", err) }
	if !strings.Contains(string(dl.Content), "subsro subtitle") { t.Fatalf("unexpected content: %q", dl.Content) }
	if s.reqs[0].URL.Query().Get("key") != "secret" { t.Fatalf("missing key in download URL: %s", s.reqs[0].URL.String()) }
}

func TestSearchPrerequisites(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &stub{}, providercore.Config{}, providercore.SearchRequest{})
	if err == nil || !strings.Contains(err.Error(), "apiKey") { t.Fatalf("expected apiKey prerequisite, got %v", err) }
	key := "secret"
	_, err = Adapter.Search(context.Background(), &stub{}, providercore.Config{APIKey: &key}, providercore.SearchRequest{})
	if err == nil || !strings.Contains(err.Error(), "imdb") { t.Fatalf("expected imdb prerequisite, got %v", err) }
}

func zipBody() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("fixture.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nsubsro subtitle\n"))
	_ = zw.Close()
	return buf.Bytes()
}
