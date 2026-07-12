package subssabbz

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
	if req.Method != http.MethodPost || req.FormValue("act") != "search" || req.FormValue("select-language") != "1" {
		return &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(""))}, nil
	}
	html := `<table><tr class="subs-row"><td></td><td class="c2field"><a href="/download/42.zip" onmouseover="ddrivetip('note','#1')">Movie Release</a> (2024)</td><td></td><td></td><td></td><td></td><td>1</td><td>23.976</td><td>uploader</td><td><a href="https://imdb.com/title/tt1234567/">imdb</a></td></tr></table>`
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(html))}, nil
}

func TestPostSearchRowsAndDownloadArchive(t *testing.T) {
	year := int32(2024)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", Title: "Movie", LanguageID: "eng", Year: &year})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search = %#v, %v", candidates, err)
	}
	if candidates[0].ProviderName != "subssabbz" || candidates[0].ReleaseName != "Movie Release" || !strings.HasPrefix(candidates[0].SourceURL, "https://example.test/download/") {
		t.Fatalf("unexpected candidate: %#v", candidates[0])
	}
	dl, err := Adapter.Download(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, candidates[0])
	if err != nil || !bytes.Contains(dl.Content, []byte("fixture")) {
		t.Fatalf("Download = %q, %v", dl.Content, err)
	}
}

func zipFixture() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("movie.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nfixture\n"))
	_ = zw.Close()
	return buf.Bytes()
}
