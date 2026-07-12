package subx

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type httpService struct{ client *http.Client }
func (s httpService) DoProviderRequest(req *http.Request, _ string, _ bool) (*http.Response, error) { return s.client.Do(req) }

func TestSearchUsesBearerAndFiltersEpisodeWithSeasonPackFallback(t *testing.T) {
	var auth, videoType, imdb string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		videoType = r.URL.Query().Get("video_type")
		imdb = r.URL.Query().Get("imdb_id")
		_, _ = w.Write([]byte(`{"items":[
			{"id":1,"title":"wrong","season":1,"episode":2},
			{"id":2,"title":"pack","description":"latam","season":1,"episode":null},
			{"id":3,"title":"other season","season":2,"episode":1}]}`))
	}))
	defer server.Close()
	key := "secret"
	s, e := int32(1), int32(1)
	got, err := adapter{}.Search(context.Background(), httpService{server.Client()}, providercore.Config{BaseURL: server.URL, APIKey: &key}, providercore.SearchRequest{MediaType: "serie", Title: "Show", SeasonNumber: &s, EpisodeNumber: &e, MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt123"}}})
	if err != nil { t.Fatalf("Search error: %v", err) }
	if auth != "Bearer secret" || videoType != "episode" || imdb != "tt123" { t.Fatalf("request auth/type/imdb = %q/%q/%q", auth, videoType, imdb) }
	if len(got) != 1 || got[0].FileID != 2 || !strings.Contains(got[0].SourceURL, "/api/subtitles/2/download") { t.Fatalf("candidates = %#v", got) }
}

func TestSearchReturnsExactEpisodeAndSpainSpanish(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"items":[{"id":5,"title":"exact","description":"Castellano","season":1,"episode":1}]}`))
	}))
	defer server.Close()
	key := "secret"; s, e := int32(1), int32(1)
	got, err := adapter{}.Search(context.Background(), httpService{server.Client()}, providercore.Config{BaseURL: server.URL, APIKey: &key}, providercore.SearchRequest{MediaType: "serie", SeasonNumber: &s, EpisodeNumber: &e})
	if err != nil { t.Fatal(err) }
	if len(got) != 1 || got[0].LanguageID != "es" || got[0].FileID != 5 { t.Fatalf("candidates = %#v", got) }
}

func TestDownloadExtractsArchive(t *testing.T) {
	archive := zipBytes(t, map[string]string{"sub.srt": "hello"})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(archive) }))
	defer server.Close()
	key := "secret"
	got, err := adapter{}.Download(context.Background(), httpService{server.Client()}, providercore.Config{APIKey: &key}, providercore.Candidate{SourceURL: server.URL + "/file.zip"})
	if err != nil { t.Fatal(err) }
	if string(got.Content) != "hello" { t.Fatalf("content = %q", got.Content) }
}

func zipBytes(t *testing.T, files map[string]string) []byte {
	t.Helper(); var buf bytes.Buffer; zw := zip.NewWriter(&buf)
	for name, content := range files { w, err := zw.Create(name); if err != nil { t.Fatal(err) }; _, _ = w.Write([]byte(content)) }
	if err := zw.Close(); err != nil { t.Fatal(err) }
	return buf.Bytes()
}
