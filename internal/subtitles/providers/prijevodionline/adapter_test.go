package prijevodionline

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
	switch {
	case strings.Contains(req.URL.Path, "/serije/index/"):
		return resp(`<table><tr id="serija-42"><td class="naziv"><a href="/serije/view/42/show-slug">Fixture Show</a></td></tr></table>`), nil
	case strings.Contains(req.URL.Path, "/serije/view/"):
		return resp(`<script>epizode.key = '0123456789abcdef0123456789abcdef';</script><div id="epizode"><h3 id="sezona-1"></h3><div id="epizoda-7"><li class="broj">2.</li></div></div>`), nil
	default:
		return resp(`<table><tr id="prijevod-99"><td class="naziv"><a href="/download/sub-hr">Sub</a></td><td class="status">provjereno</td></tr><tr id="prijevod-opis-99"><td class="opis">WEB / BluRay</td></tr></table>`), nil
	}
}

func TestSearchAndDownloadFixture(t *testing.T) {
	s, e := int32(1), int32(2)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture Show", LanguageID: "hrv", SeasonNumber: &s, EpisodeNumber: &e})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].ProviderName != "prijevodionline" || candidates[0].LanguageID != "hrv" {
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

func TestMissingEpisodePrerequisite(t *testing.T) {
	_, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "serie"})
	if err == nil {
		t.Fatal("expected prerequisite error")
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
