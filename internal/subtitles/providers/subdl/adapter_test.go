package subdl

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

type testService struct{ client *http.Client }
func (s testService) DoProviderRequest(req *http.Request, _ string, _ bool) (*http.Response, error) { return s.client.Do(req) }

func TestSearchUsesBazarrParamsMergesFallbacksAndUnpackFiles(t *testing.T) {
	key := "key"
	paths := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.RawQuery)
		if r.URL.Query().Get("api_key") != key || r.URL.Query().Get("bazarr") != "1" || r.URL.Query().Get("unpack") != "1" { t.Fatalf("query = %s", r.URL.RawQuery) }
		switch len(paths) {
		case 1:
			if r.URL.Query().Get("episode_number") != "2" || r.URL.Query().Get("season_number") != "1" { t.Fatalf("query = %s", r.URL.RawQuery) }
			_, _ = w.Write([]byte(`{"success":true,"subtitles":[{"name":"pack","url":"/packs.zip","language":"EN","episode_from":1,"episode_end":3,"releases":["Show.S01.Pack"],"unpack_files":[{"episode":2,"url":"/direct.srt","file_n_id":"2"}]}]}`))
		case 2:
			if r.URL.Query().Get("episode_number") != "" || r.URL.Query().Get("season_number") != "1" { t.Fatalf("season query = %s", r.URL.RawQuery) }
			_, _ = w.Write([]byte(`{"success":true,"subtitles":[{"name":"season","url":"/season.zip","language":"EN","releases":["Show.Season"]}]}`))
		case 3:
			if r.URL.Query().Get("episode_number") != "12" || r.URL.Query().Get("season_number") != "" { t.Fatalf("absolute query = %s", r.URL.RawQuery) }
			_, _ = w.Write([]byte(`{"success":true,"subtitles":[{"name":"pack","url":"/duplicate.zip","language":"EN"}]}`))
		}
	}))
	defer server.Close()
	season, episode, absolute := int32(1), int32(2), int32(12)
	request := providercore.SearchRequest{MediaType: "serie", Title: "Show", LanguageID: "EN", SeasonNumber: &season, EpisodeNumber: &episode, MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt1"}, EpisodeNumbering: []providercore.EpisodeNumbering{{AbsoluteNumber: &absolute}}}}
	got, err := adapter.Search(context.Background(), testService{server.Client()}, providercore.Config{BaseURL: server.URL, APIKey: &key}, request)
	if err != nil { t.Fatal(err) }
	if len(got) != 2 { t.Fatalf("unexpected candidates: %#v", got) }
	if got[0].SourceURL != "/direct.srt" || !strings.Contains(got[0].SourceRef, "false") { t.Fatalf("unexpected unpack candidate: %#v", got[0]) }
}

func TestMovieSearchIncludesTMDBAndAuthRequirement(t *testing.T) {
	_, err := adapter.Search(context.Background(), testService{http.DefaultClient}, providercore.Config{}, providercore.SearchRequest{MediaType: "movie"})
	if err == nil { t.Fatal("expected api key error") }
	key := "key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("type") != "movie" || r.URL.Query().Get("imdb_id") != "tt1" || r.URL.Query().Get("tmdb_id") != "55" { t.Fatalf("query = %s", r.URL.RawQuery) }
		_, _ = w.Write([]byte(`{"status":true,"subtitles":[]}`))
	}))
	defer server.Close()
	_, err = adapter.Search(context.Background(), testService{server.Client()}, providercore.Config{BaseURL: server.URL, APIKey: &key}, providercore.SearchRequest{MediaType: "movie", LanguageID: "EN", MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"imdb": "tt1", "tmdb": "55"}}})
	if err != nil { t.Fatal(err) }
}

func TestDownloadExtractsZipAndAcceptsUnpackedRawSubtitle(t *testing.T) {
	archive := zipData(t, map[string]string{"sub.srt": "subtitle"})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".srt") { _, _ = w.Write([]byte("raw")); return }
		_, _ = w.Write(archive)
	}))
	defer server.Close()
	got, err := adapter.Download(context.Background(), testService{server.Client()}, providercore.Config{}, providercore.Candidate{SourceURL: server.URL + "/file.zip", ReleaseName: "release"})
	if err != nil { t.Fatal(err) }
	if string(got.Content) != "subtitle" { t.Fatalf("content = %q", got.Content) }
	got, err = adapter.Download(context.Background(), testService{server.Client()}, providercore.Config{}, providercore.Candidate{SourceURL: server.URL + "/direct.srt", ReleaseName: "release"})
	if err != nil { t.Fatal(err) }
	if string(got.Content) != "raw" { t.Fatalf("raw content = %q", got.Content) }
}

func TestClassifiesForbiddenAndRateLimit(t *testing.T) {
	key := "key"
	for status, body := range map[int]string{403: ``, 429: `{"error":"rate_limit"}`} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(status); _, _ = w.Write([]byte(body)) }))
		_, err := adapter.Search(context.Background(), testService{server.Client()}, providercore.Config{BaseURL: server.URL, APIKey: &key}, providercore.SearchRequest{MediaType: "movie", LanguageID: "EN"})
		server.Close()
		if err == nil { t.Fatalf("status %d: expected error", status) }
	}
}

func zipData(t *testing.T, files map[string]string) []byte {
	t.Helper(); var buf bytes.Buffer; zw := zip.NewWriter(&buf)
	for name, content := range files { w, err := zw.Create(name); if err != nil { t.Fatal(err) }; _, _ = w.Write([]byte(content)) }
	if err := zw.Close(); err != nil { t.Fatal(err) }
	return buf.Bytes()
}
