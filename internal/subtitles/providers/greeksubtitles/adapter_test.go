package greeksubtitles

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type providerStub struct{ urls []string }

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.urls = append(s.urls, req.URL.String())
	body := `<table><tr><td class="latest_name"><img src="/images/en.gif"><a href="/subtitle/123/fixture">Fixture.Release</a></td></tr><tr><td><a href="search.php?name=Fixture&page=2">Next</a></td></tr></table>`
	if strings.Contains(req.URL.RawQuery, "page=2") {
		body = `<table><tr><td class="latest_name"><img src="/images/gr.gif"><a href="/subtitle/124/fixture2">Fixture.Greek</a></td></tr></table>`
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestSearchParsesGreekSubtitlesPagination(t *testing.T) {
	svc := &providerStub{}
	candidates, err := Adapter.Search(context.Background(), svc, providercore.Config{BaseURL: "http://gr.greek-subtitles.com"}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture"})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 2 || candidates[0].SourceURL != "http://www.greeksubtitles.info/getp.php?id=123" || candidates[1].LanguageID != "ell" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	if len(svc.urls) != 2 {
		t.Fatalf("expected pagination, got %v", svc.urls)
	}
}
