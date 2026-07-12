package bayflix

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

type providerStub struct{ seen string }

func (s *providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	s.seen = req.URL.String()
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(zipBytes(timedText())))}, nil
	}
	body := `[{"_id":"dune-2021","title":"Dune","subtitle_link":"https://bayflix.sb/api/subtitles/download/dune-2021","release_name":["Dune.2021.1080p","Dune.S01E02"],"release_date":"2021-10-20","description":"FPS: 23.976\nCDs: 1","downloads":7}]`
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
}

func TestSearchUsesBayflixAPI(t *testing.T) {
	svc := &providerStub{}
	y := int32(2021)
	candidates, err := Adapter.Search(context.Background(), svc, providercore.Config{}, providercore.SearchRequest{MediaType: "movie", Title: "Dune", LanguageID: "bul", Year: &y})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if !strings.Contains(svc.seen, "/api/subtitles/search?title=Dune") {
		t.Fatalf("unexpected route: %s", svc.seen)
	}
	if len(candidates) != 1 || candidates[0].ProviderName != "bayflix" || candidates[0].SourceURL != "https://bayflix.sb/api/subtitles/download/dune-2021" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
}

func TestSearchFiltersEpisodeAndDownloadArchive(t *testing.T) {
	svc := &providerStub{}
	s, e := int32(1), int32(2)
	candidates, err := Adapter.Search(context.Background(), svc, providercore.Config{}, providercore.SearchRequest{MediaType: "serie", Title: "Dune", SeasonNumber: &s, EpisodeNumber: &e})
	if err != nil || len(candidates) != 1 {
		t.Fatalf("Search = %#v, %v", candidates, err)
	}
	dl, err := Adapter.Download(context.Background(), svc, providercore.Config{}, candidates[0])
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if string(dl.Content) != timedText() {
		t.Fatalf("unexpected content: %q", dl.Content)
	}
}

func TestRejectsUnsupportedMediaType(t *testing.T) {
	_, err := Adapter.Search(context.Background(), &providerStub{}, providercore.Config{}, providercore.SearchRequest{MediaType: "audio", Title: "Fixture"})
	if err == nil || !strings.Contains(err.Error(), "provider_prerequisite_missing") {
		t.Fatalf("expected error, got %v", err)
	}
}

func timedText() string { return "1\n00:00:01,000 --> 00:00:02,000\nfixture\n" }
func zipBytes(content string) []byte {
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	w, _ := zw.Create("fixture.srt")
	_, _ = w.Write([]byte(content))
	_ = zw.Close()
	return b.Bytes()
}
