package subf2m

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
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipBytes()))}, nil
	}
	if strings.Contains(req.URL.Path, "searchbytitle") {
		return resp(`<li><div class="title"><a href="/subtitles/fixture-season-1">Fixture Show - First Season (2020)</a></div></li>`), nil
	}
	if strings.Contains(req.URL.Path, "download-page") {
		return resp(`<a id="downloadButton" href="/download/archive.zip">Download</a>`), nil
	}
	return resp(`<h2><a href="/title/tt1">tt1</a></h2><li class="item"><ul class="scrolllist"><li>Fixture Show S01E02 WEB</li></ul><a class="download icon-download" href="/download-page">download</a></li>`), nil
}

func TestSearchAndDownloadFixture(t *testing.T) {
	s, e, y := int32(1), int32(2), int32(2020)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture Show", LanguageID: "eng", Year: &y, SeasonNumber: &s, EpisodeNumber: &e})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].ProviderName != "subf2m" || candidates[0].LanguageID != "eng" {
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
