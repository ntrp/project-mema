package animesubinfo

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
	body := `<table class="Napisy" style="text-align:center"><tr class="KNap"><td>Kimetsu ep02</td><td>date</td><td></td><td>Advanced SSA</td></tr><tr class="KNap"><td>Demon Slayer ep02</td><td><a>Uploader</a></td><td></td><td></td><td>12 KB</td></tr><tr class="KNap"><td>Alt Title</td><td></td><td></td><td>308 razy</td></tr><tr class="KKom"><td class="KNap" align="left">Synchro: [SubsPlease]</td><td><form method="POST"><input name="id" value="42"><input name="sh" value="abc123"></form></td></tr></table>`
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestSearchUsesPolishAnimeSubInfoFormRoute(t *testing.T) {
	svc := &providerStub{}
	e := int32(2)
	candidates, err := Adapter.Search(context.Background(), svc, providercore.Config{}, providercore.SearchRequest{MediaType: "serie", Title: "Kimetsu", LanguageID: "pol", EpisodeNumber: &e})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if !strings.Contains(svc.rawQuery, "pTitle=org") || !strings.Contains(svc.rawQuery, "pSortuj=pobrn") {
		t.Fatalf("unexpected query: %s", svc.rawQuery)
	}
	if len(candidates) != 1 || candidates[0].FileID != 42 || candidates[0].SourceRef != "abc123" || candidates[0].DownloadCount != 308 {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
}
