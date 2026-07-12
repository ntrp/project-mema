package subtitriid

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

type providerStub struct{}

func (providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipFixture()))}, nil
	}
	switch req.URL.Path {
	case "/search/":
		if req.URL.Query().Get("q") != "Straume" {
			return &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(""))}, nil
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`<div class="eBlock"><div class="eTitle"><a href="/load/straume/1">Straume</a></div></div>`))}, nil
	case "/load/straume/1":
		html := `<h1 class="main-header">Flow / Straume</h1><span id="film-page-year">2024</span><div id="actors-page"><a href="https://www.imdb.com/title/tt4772188/">imdb</a></div><a class="hvr" href="/download/straume.zip">download</a>`
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(html))}, nil
	}
	return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader(""))}, nil
}

func TestSearchFetchesDetailPageAndDownloadsArchive(t *testing.T) {
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", Title: "Straume"})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search = %#v, %v", candidates, err)
	}
	if candidates[0].ProviderName != "subtitriid" || candidates[0].LanguageID != "lav" || !strings.Contains(candidates[0].ReleaseName, "Straume") {
		t.Fatalf("unexpected candidate: %#v", candidates[0])
	}
	dl, err := Adapter.Download(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, candidates[0])
	if err != nil || !bytes.Contains(dl.Content, []byte("sveiki")) {
		t.Fatalf("Download = %q, %v", dl.Content, err)
	}
}

func zipFixture() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("straume.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nsveiki\n"))
	_ = zw.Close()
	return buf.Bytes()
}
