package animekalesi

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type providerStub struct{ paths []string }

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.paths = append(s.paths, req.URL.Path)
	body := ""
	switch req.URL.Path {
	case "/tum-anime-serileri.html":
		body = `<table><tr><td id="bolumler"><a href="bolumler-fixture.html">Fixture Anime</a></td></tr></table>`
	case "/altyazib-fixture.html":
		body = `<table><tr><td id="ayazi_indir"><a href="indir_bolum-1.html" title="1. Sezon 2. Bölüm Türkçe Altyazısı">indir</a></td></tr></table>`
	case "/indir_bolum-1.html":
		body = `<strong>Altyazı/Çeviri:</strong> FixtureUploader <div id="altyazi_indir"><a href="download/fixture.zip">download</a></div>`
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestSearchTraversesAnimeKalesiPages(t *testing.T) {
	svc := &providerStub{}
	s, e := int32(1), int32(2)
	candidates, err := Adapter.Search(context.Background(), svc, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture Anime", SeasonNumber: &s, EpisodeNumber: &e})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if strings.Join(svc.paths, ",") != "/tum-anime-serileri.html,/altyazib-fixture.html,/indir_bolum-1.html" {
		t.Fatalf("unexpected traversal: %#v", svc.paths)
	}
	if len(candidates) != 1 || candidates[0].LanguageID != "tur" || !strings.Contains(candidates[0].SourceURL, "/download/fixture.zip") {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
}

func TestRejectsUnsupportedMediaType(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected unsupported media type error, got %v", err)
	}
}
