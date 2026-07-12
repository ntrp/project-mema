package soustitreseu

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
	if strings.Contains(req.URL.Path, "search.html") {
		return resp(`<div class="serie"><h3><a href="series/fixture.html">Fixture Show</a></h3></div>`), nil
	}
	return resp(`<a class="subList" href="Fixture.Show.S01E02.en.zip"><span class="episodenum">1×2</span></a>`), nil
}

func TestSearchAndDownloadFixture(t *testing.T) {
	s, e := int32(1), int32(2)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture Show", LanguageID: "eng", SeasonNumber: &s, EpisodeNumber: &e})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].ProviderName != "soustitreseu" || candidates[0].LanguageID != "eng" {
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
	w, _ := zw.Create("Fixture.Show.S01E02.en.srt")
	_, _ = w.Write([]byte("1\nfixture\n"))
	_ = zw.Close()
	return b.Bytes()
}
