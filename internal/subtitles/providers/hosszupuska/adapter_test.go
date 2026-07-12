package hosszupuska

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type providerStub struct{ rawQuery string }

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.rawQuery = req.URL.RawQuery
	body := `<table><tr><td><img src="css/infooldal.png"></td></tr><tr onmouseover="this.style.backgroundImage='url(css/over2.jpg)"><td></td><td><b>Fixture s01e02</b> Fixture (WEB-DL, HDTV)</td><td><img src="flags/2.gif"></td><td></td><td></td><td></td><td><a href="download.php?file=0124336.zip">download</a></td></tr></table>`
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestSearchUsesHosszupuskaEpisodeRoute(t *testing.T) {
	svc := &providerStub{}
	s, e := int32(1), int32(2)
	candidates, err := Adapter.Search(context.Background(), svc, providercore.Config{}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture", SeasonNumber: &s, EpisodeNumber: &e})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if !strings.Contains(svc.rawQuery, "cim=Fixture") || !strings.Contains(svc.rawQuery, "evad=01") || !strings.Contains(svc.rawQuery, "resz=02") {
		t.Fatalf("unexpected query: %s", svc.rawQuery)
	}
	if len(candidates) != 1 || candidates[0].LanguageID != "eng" || candidates[0].FileID != 124336 {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
}

func TestRejectsUnsupportedMediaType(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected unsupported media type error, got %v", err)
	}
}
