package nekur

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
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/zip"}}, Body: io.NopCloser(bytes.NewReader(zipBytes()))}, nil
	}
	html := `<table><tbody><tr><td class="title"><a href="/lv.zip">My Movie <span>x</span></a></td><td class="year">(2020)</td><td></td><td><a href="/title/tt1/">imdb</a></td><td class="fps">23.976</td><td class="notes">BluRay</td></tr></tbody></table>`
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(html))}, nil
}

func TestSearchAndDownloadFixture(t *testing.T) {
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "http://example.test"}, providercore.SearchRequest{MediaType: "movie", Title: "My Movie", LanguageID: "lav"})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].ProviderName != "nekur" || candidates[0].LanguageID != "lav" || !strings.Contains(candidates[0].ReleaseName, "BluRay") {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	dl, err := Adapter.Download(context.Background(), providerStub{}, providercore.Config{}, candidates[0])
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if !bytes.Contains(dl.Content, []byte("fixture")) {
		t.Fatalf("unexpected content: %q", dl.Content)
	}
}

func TestUnsupportedMedia(t *testing.T) {
	_, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "serie"})
	if err == nil {
		t.Fatal("expected unsupported media error")
	}
}

func zipBytes() []byte {
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	w, _ := zw.Create("sub.srt")
	_, _ = w.Write([]byte("1\nfixture\n"))
	_ = zw.Close()
	return b.Bytes()
}
