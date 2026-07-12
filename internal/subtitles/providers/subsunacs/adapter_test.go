package subsunacs

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
	if req.Method != http.MethodPost || req.FormValue("m") != "Show 01 02" || req.FormValue("l") != "1" {
		return &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(""))}, nil
	}
	html := `<table><tr onmouseover="x"><td class="tdMovie"><a class="tooltip" href="/subs/99.zip" title="notes">Show Release<span class="smGray">&nbsp;(2024)</span></a></td><td>1</td><td>25</td><td><img title="9.5"></td><td>12</td><td>uploader</td></tr></table>`
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(html))}, nil
}

func TestEpisodeFormRowsAndArchive(t *testing.T) {
	season, episode := int32(1), int32(2)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Show", LanguageID: "eng", SeasonNumber: &season, EpisodeNumber: &episode})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search = %#v, %v", candidates, err)
	}
	if candidates[0].ProviderName != "subsunacs" || candidates[0].DownloadCount != 12 || candidates[0].LanguageID != "eng" {
		t.Fatalf("unexpected candidate: %#v", candidates[0])
	}
	dl, err := Adapter.Download(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, candidates[0])
	if err != nil || !bytes.Contains(dl.Content, []byte("subs")) {
		t.Fatalf("Download = %q, %v", dl.Content, err)
	}
}

func zipFixture() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("show.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nsubs\n"))
	_ = zw.Close()
	return buf.Bytes()
}
