package greeksubs

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type providerStub struct {
	urls    []string
	methods []string
}

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.urls = append(s.urls, req.URL.String())
	s.methods = append(s.methods, req.Method)
	if isDownload && req.Method == http.MethodPost {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("subtitle"))}, nil
	}
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`<input name="langcode" value="en"><input name="uid" value="u"><input name="output" value="o"><input name="dll" value="d">`))}, nil
	}
	body := `<input id="secCode" value="sec"><table id="elSub"><tbody><tr><td><img alt="en"></td><td></td><td><a onclick="downloadMe('abc')">dl</a><span>Fixture.Release</span></td></tr></tbody></table>`
	if strings.Contains(req.URL.Path, "/en/view/") {
		body = `<div class="col-lg-offset-2 col-md-8 text-center top30 bottom10"><a href="https://greeksubs.net/episode">Season 1 Episode 2</a></div>`
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestSearchUsesImdbAndEpisodePage(t *testing.T) {
	svc := &providerStub{}
	s, e := int32(1), int32(2)
	candidates, err := Adapter.Search(context.Background(), svc, providercore.Config{}, providercore.SearchRequest{MediaType: "serie", SeasonNumber: &s, EpisodeNumber: &e, MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt1234567"}}})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].LanguageID != "eng" || !strings.Contains(candidates[0].SourceURL, "/dll/abc/0/sec") {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	if !strings.Contains(strings.Join(svc.urls, ","), "/en/view/1234567") {
		t.Fatalf("missing imdb route: %v", svc.urls)
	}
}

func TestDownloadPostsSingleUseForm(t *testing.T) {
	svc := &providerStub{}
	dl, err := Adapter.Download(context.Background(), svc, providercore.Config{}, providercore.Candidate{SourceURL: "https://greeksubs.net/dll/abc/0/sec", SourceRef: "https://greeksubs.net/en/view/123"})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if string(dl.Content) != "subtitle" || strings.Join(svc.methods, ",") != "GET,POST" {
		t.Fatalf("unexpected download: %q methods=%v", dl.Content, svc.methods)
	}
}

func TestRequiresImdbID(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected error, got %v", err)
	}
}
