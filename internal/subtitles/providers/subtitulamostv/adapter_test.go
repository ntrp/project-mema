package subtitulamostv

import (
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
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("1\n00:00:01,000 --> 00:00:02,000\nhola\n"))}, nil
	}
	switch req.URL.Path {
	case "/search/query":
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`[{"show_id":7,"show_name":"Las Encinas"}]`))}, nil
	case "/shows/7":
		html := `<div id="season-choices"><a class="selected" href="/shows/7/season/2">2</a></div><div id="episode-choices"><a class="selected" href="/shows/7/season/2/episode/3">3</a></div><div class="language-container"><div class="language-name">Español</div><div class="version-container"><p>ignored</p><p>WEB-DL</p><a href="/download/77">download</a></div></div>`
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(html))}, nil
	}
	return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader(""))}, nil
}

func TestSearchTraversesShowEpisodeAndDownloadsText(t *testing.T) {
	season, episode := int32(2), int32(3)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Las Encinas", SeasonNumber: &season, EpisodeNumber: &episode})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search = %#v, %v", candidates, err)
	}
	if candidates[0].ProviderName != "subtitulamostv" || candidates[0].LanguageID != "spa" || candidates[0].SourceURL != "/download/77" {
		t.Fatalf("unexpected candidate: %#v", candidates[0])
	}
	dl, err := Adapter.Download(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, candidates[0])
	if err != nil || !bytes.Contains(dl.Content, []byte("hola")) {
		t.Fatalf("Download = %q, %v", dl.Content, err)
	}
}
