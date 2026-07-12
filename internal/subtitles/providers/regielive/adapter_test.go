package regielive

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type testService struct{ client *http.Client }
func (s testService) DoProviderRequest(req *http.Request, _ string, _ bool) (*http.Response, error) { return s.client.Do(req) }

func TestSearchSendsRegieLivePayloadAndParsesNestedResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("RL-API") != apiHeader { t.Fatalf("missing api header") }
		if r.URL.Query().Get("nume") != "Show" || r.URL.Query().Get("sezon") != "1" || r.URL.Query().Get("episod") != "2" || r.URL.Query().Get("an") != "2024" { t.Fatalf("query = %s", r.URL.RawQuery) }
		_, _ = w.Write([]byte(`{"rezultate":{"film":{"subtitrari":{"1":{"titlu":"Show.S01E02","url":"https://subtitrari.regielive.ro/dl.zip","rating":{"nota":9}}}}}}`))
	}))
	defer server.Close()
	year, season, episode := int32(2024), int32(1), int32(2)
	got, err := adapter.Search(context.Background(), testService{server.Client()}, providercore.Config{BaseURL: server.URL}, providercore.SearchRequest{MediaType: "serie", Title: "Show", Year: &year, SeasonNumber: &season, EpisodeNumber: &episode})
	if err != nil { t.Fatal(err) }
	if len(got) != 1 || got[0].LanguageID != "ro" || got[0].DownloadCount != 9 { t.Fatalf("unexpected candidates: %#v", got) }
}

func TestDownloadUsesBrowserHeadersAndSkipsTxtMembers(t *testing.T) {
	archive := zipData(t, map[string]string{"note.txt": "ignore", ".hidden.srt": "ignore", "movie.ass": "subtitle"})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Referer") != "https://subtitrari.regielive.ro" { t.Fatalf("missing referer") }
		_, _ = w.Write(archive)
	}))
	defer server.Close()
	got, err := adapter.Download(context.Background(), testService{server.Client()}, providercore.Config{}, providercore.Candidate{SourceURL: server.URL + "/download"})
	if err != nil { t.Fatal(err) }
	if string(got.Content) != "subtitle" { t.Fatalf("content = %q", got.Content) }
}

func TestDownloadRejectsProvider500Text(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("500")) }))
	defer server.Close()
	_, err := adapter.Download(context.Background(), testService{server.Client()}, providercore.Config{}, providercore.Candidate{SourceURL: server.URL})
	if err == nil { t.Fatal("expected error") }
}

func zipData(t *testing.T, files map[string]string) []byte {
	t.Helper(); var buf bytes.Buffer; zw := zip.NewWriter(&buf)
	for name, content := range files { w, err := zw.Create(name); if err != nil { t.Fatal(err) }; _, _ = w.Write([]byte(content)) }
	if err := zw.Close(); err != nil { t.Fatal(err) }
	return buf.Bytes()
}
