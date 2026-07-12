package subsynchro

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
	if providerType != "subsynchro" {
		return nil, nil
	}
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipFixture()))}, nil
	}
	if !strings.Contains(req.URL.Path, "/include/ajax/subMarin.php") || req.URL.Query().Get("title") != "Amelie" {
		return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader(""))}, nil
	}
	body := `{"status":200,"data":[{"release":"Amelie.2001.FR","filename":"amelie.srt","telechargement":"/download/amelie.zip","fichier":"zip"}]}`
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestSearchAndDownloadSubMarinJSON(t *testing.T) {
	year := int32(2001)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "movie", Title: "Amelie", Year: &year})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search = %#v, %v", candidates, err)
	}
	if candidates[0].ProviderName != "subsynchro" || candidates[0].LanguageID != "fre" || candidates[0].SourceURL != "/download/amelie.zip" {
		t.Fatalf("unexpected candidate: %#v", candidates[0])
	}
	dl, err := Adapter.Download(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, candidates[0])
	if err != nil || !bytes.Contains(dl.Content, []byte("bonjour")) {
		t.Fatalf("Download = %q, %v", dl.Content, err)
	}
}

func zipFixture() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("amelie.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nbonjour\n"))
	_ = zw.Close()
	return buf.Bytes()
}
