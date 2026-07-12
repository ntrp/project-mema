package yifysubtitles

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

type providerStub struct{ paths []string }

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.paths = append(s.paths, req.URL.Path)
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipBody()))}, nil
	}
	switch req.URL.Path {
	case "/movie-imdb/tt1234567":
		return html(`<table class="other-subs"><tbody><tr><td>8</td><td>English</td><td><a href="/subtitle/fixture-english">subtitle Fixture.YIFY</a></td><td><span class="hi-subtitle"></span></td><td>uploader</td></tr></tbody></table>`), nil
	case "/subtitle/fixture-english":
		return html(`<a class="download-subtitle" href="/subtitle/fixture-english.zip">Download</a>`), nil
	default:
		return html(`<html>ok</html>`), nil
	}
}
func html(body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"text/html"}}, Body: io.NopCloser(strings.NewReader(body))}
}
func TestSearchAndDownloadUseYIFYFlow(t *testing.T) {
	stub := &providerStub{}
	candidates, err := Adapter.Search(context.Background(), stub, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", LanguageID: "eng", MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt1234567"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 1 || candidates[0].SourceURL != "https://example.test/subtitle/fixture-english" || candidates[0].DownloadCount != 8 {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	dl, err := Adapter.Download(context.Background(), stub, providercore.Config{}, candidates[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(dl.Content), "fixture subtitle") {
		t.Fatalf("unexpected content: %q", dl.Content)
	}
	if strings.Join(stub.paths, ",") != "/movie-imdb/tt1234567,/subtitle/fixture-english,/subtitle/fixture-english.zip" {
		t.Fatalf("unexpected requests: %#v", stub.paths)
	}
}
func TestRejectsUnsupportedMediaType(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "serie"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected unsupported media type error, got %v", err)
	}
}
func zipBody() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("fixture.srt")
	_, _ = w.Write([]byte("fixture subtitle"))
	_ = zw.Close()
	return buf.Bytes()
}
