package subsarr

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

type stub struct {
	reqs      []*http.Request
	downloads []bool
}

func (s *stub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.reqs = append(s.reqs, req)
	s.downloads = append(s.downloads, isDownload)
	if strings.Contains(req.URL.Path, "download") {
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Disposition": []string{"attachment; filename=fixture.zip"}}, Body: io.NopCloser(bytes.NewReader(zipBody()))}, nil
	}
	if strings.Contains(req.URL.Path, "health") {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"ok":true}`))}, nil
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"subtitles":[{"id":77,"language":"eng","release":"Fixture.2024","download_url":"/api/subtitles/download/77"}]}`))}, nil
}

func TestSearchBuildsLocalAPIRequest(t *testing.T) {
	s := &stub{}
	year := int32(2024)
	got, err := Adapter.Search(context.Background(), s, providercore.Config{BaseURL: "http://127.0.0.1:7878"}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture", LanguageID: "eng", Year: &year})
	if err != nil { t.Fatalf("Search error: %v", err) }
	if len(got) != 1 || got[0].ProviderName != "subsarr" || got[0].FileID != 77 || got[0].ReleaseName != "Fixture.2024" { t.Fatalf("unexpected candidates: %#v", got) }
	if s.reqs[0].URL.Host != "127.0.0.1:7878" || s.reqs[0].URL.Query().Get("title") != "Fixture" { t.Fatalf("unexpected request URL: %s", s.reqs[0].URL.String()) }
}

func TestDownloadUsesLocalAPINotDownloadPolicyAndExtractsArchive(t *testing.T) {
	s := &stub{}
	dl, err := Adapter.Download(context.Background(), s, providercore.Config{BaseURL: "http://127.0.0.1:7878"}, providercore.Candidate{SourceURL: "/api/subtitles/download/77"})
	if err != nil { t.Fatalf("Download error: %v", err) }
	if !strings.Contains(string(dl.Content), "fixture subtitle") { t.Fatalf("unexpected content: %q", dl.Content) }
	if len(s.downloads) != 1 || s.downloads[0] { t.Fatalf("subsarr local endpoint must be requested as API traffic, got downloads=%v", s.downloads) }
}

func TestRequiresBaseURL(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &stub{}, providercore.Config{}, providercore.SearchRequest{Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "baseUrl") { t.Fatalf("expected baseUrl prerequisite, got %v", err) }
}

func zipBody() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("fixture.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nfixture subtitle\n"))
	_ = zw.Close()
	return buf.Bytes()
}
