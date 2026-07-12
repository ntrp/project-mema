package tvsubtitles

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
	s.paths = append(s.paths, req.Method+" "+req.URL.Path)
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipBody()))}, nil
	}
	switch req.URL.Path {
	case "/search1.php":
		return html(`<div class="left"><li><div><a href="/tvshow-42.html">Fixture (2024)</a></div></li></div>`), nil
	case "/tvshow-42-1.html":
		return html(`<table id="table5"><tr><td>1x02</td><td><a href="episode-77.html">Episode</a></td></tr></table>`), nil
	case "/episode-77.html":
		return html(`<a href="subtitle-99.html"><div class="subtitlen"><h5><img src="/images/flags/eng.gif">Fixture.Release</h5><p title="rip">WEB</p></div></a>`), nil
	case "/download-99.html":
		return html("<script>s1 = 'files/';\ns2 = 'fixture.zip';\n</script>"), nil
	default:
		return html(`<html>ok</html>`), nil
	}
}

func html(body string) *http.Response {
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"text/html"}}, Body: io.NopCloser(strings.NewReader(body))}
}

func TestSearchAndDownloadUseProviderTraversal(t *testing.T) {
	season, episode := int32(1), int32(2)
	stub := &providerStub{}
	candidates, err := Adapter.Search(context.Background(), stub, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture", LanguageID: "eng", SeasonNumber: &season, EpisodeNumber: &episode})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].SourceURL != "https://example.test/download-99.html" || candidates[0].LanguageID != "eng" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	download, err := Adapter.Download(context.Background(), stub, providercore.Config{}, candidates[0])
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if !strings.Contains(string(download.Content), "fixture subtitle") {
		t.Fatalf("unexpected download: %q", download.Content)
	}
	want := []string{"POST /search1.php", "GET /tvshow-42-1.html", "GET /episode-77.html", "GET /download-99.html", "GET /files/fixture.zip"}
	if strings.Join(stub.paths, ",") != strings.Join(want, ",") {
		t.Fatalf("unexpected traversal: %#v", stub.paths)
	}
}

func TestRejectsUnsupportedMediaType(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected unsupported media type error, got %v", err)
	}
}

func zipBody() []byte {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("fixture.srt")
	_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nfixture subtitle\n"))
	_ = zw.Close()
	return buf.Bytes()
}
