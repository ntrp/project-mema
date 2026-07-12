package subs4free

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
	if isDownload && req.Method == http.MethodPost {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipBytes()))}, nil
	}
	if isDownload {
		return resp(`<input name="id" value="123"><input type="image" width="10px" height="10px">`), nil
	}
	if strings.Contains(req.URL.Path, "search_report") {
		return resp(`<select name="Mov_sel"><option value="?p=movie/fixture">Fixture Movie 2020</option></select>`), nil
	}
	if strings.Contains(req.URL.Path, "anti-block") || strings.Contains(req.URL.Path, "favicon") {
		return resp(`ok`), nil
	}
	return resp(`<td id="dates_header"><table><u>Fixture Movie</u><div>x (2020)</div></table></td><div class="movie-details"><span>WEB</span><a href="/download-page">dl</a><div class="movie-info"><p><a>uploader</a></p></div><i class="sprite engif"></i></div>`), nil
}

func TestSearchAndDownloadFixture(t *testing.T) {
	y := int32(2020)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture Movie", LanguageID: "eng", Year: &y})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].ProviderName != "subs4free" || candidates[0].LanguageID != "eng" {
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

func resp(s string) *http.Response {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(s))}
}
func zipBytes() []byte {
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	w, _ := zw.Create("sub.srt")
	_, _ = w.Write([]byte("1\nfixture\n"))
	_ = zw.Close()
	return b.Bytes()
}
