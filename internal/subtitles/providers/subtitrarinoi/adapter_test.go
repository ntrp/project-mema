package subtitrarinoi

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
	if req.Method != http.MethodPost || req.FormValue("cautare") != "1234567" || req.Header.Get("X-Requested-With") != "XMLHttpRequest" {
		return &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(""))}, nil
	}
	html := `<div id="round"><div id="content-main"><a href="/movie">Film Title (2020)</a><p></p><p></p><p></p><p></p><p>Uploader: ion</p></div><div class="buton"><a href="download/film.zip">download</a></div><div id="content-right"><p>Downloads: 42</p><a><img src="imdb.png"/></a></div></div><div>WEB-DL;Group</div>`
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(html))}, nil
}

func TestAjaxSearchByImdbAndDownloadArchive(t *testing.T) {
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", Title: "Ignored", MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt1234567"}}})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search = %#v, %v", candidates, err)
	}
	if candidates[0].ProviderName != "subtitrarinoi" || candidates[0].LanguageID != "ron" || candidates[0].DownloadCount != 42 || !strings.Contains(candidates[0].ReleaseName, "WEB-DL") {
		t.Fatalf("unexpected candidate: %#v", candidates[0])
	}
	dl, err := Adapter.Download(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, candidates[0])
	if err != nil || !bytes.Contains(dl.Content, []byte("salut")) {
		t.Fatalf("Download = %q, %v", dl.Content, err)
	}
}

func zipFixture() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("film.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nsalut\n"))
	_ = zw.Close()
	return buf.Bytes()
}
