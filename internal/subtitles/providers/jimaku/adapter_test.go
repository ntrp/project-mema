package jimaku

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

func TestSearchUsesAniListAndFiltersBazarrStyle(t *testing.T) {
	secret := "token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != secret { t.Fatalf("missing auth header") }
		switch r.URL.Path {
		case "/entries/search":
			if r.URL.Query().Get("anilist_id") != "123" { t.Fatalf("query = %s", r.URL.RawQuery) }
			_, _ = w.Write([]byte(`[{"id":7,"name":"Show","flags":{"movie":false}}]`))
		case "/entries/7/files":
			if r.URL.Query().Get("episode") != "2" { t.Fatalf("files query = %s", r.URL.RawQuery) }
			_, _ = w.Write([]byte(`[
				{"name":"Show.E02.srt","url":"https://jimaku.cc/file.srt","size":900},
				{"name":"Show.E02.zip","url":"https://jimaku.cc/file.zip","size":900},
				{"name":"Show.whisper.srt","url":"https://jimaku.cc/ai.srt","size":900},
				{"name":"Show.tiny.srt","url":"https://jimaku.cc/tiny.srt","size":10}]`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()
	season, episode := int32(1), int32(2)
	got, err := adapter.Search(context.Background(), testService{server.Client()}, providercore.Config{BaseURL: server.URL, APIKey: &secret}, providercore.SearchRequest{MediaType: "serie", Title: "Show", LanguageID: "ja", SeasonNumber: &season, EpisodeNumber: &episode, MediaContext: providercore.MediaContext{ExternalIDs: map[string]string{"anilist": "123"}}})
	if err != nil { t.Fatal(err) }
	if len(got) != 1 || got[0].ReleaseName != "Show.E02.srt" { t.Fatalf("unexpected candidates: %#v", got) }
}

func TestSearchNameFallbackRetriesLiveActionAndArchives(t *testing.T) {
	secret := "token"
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/entries/search":
			calls++
			if calls == 1 { _, _ = w.Write([]byte(`[]`)); return }
			if r.URL.Query().Get("anime") != "false" { t.Fatalf("expected live-action retry") }
			_, _ = w.Write([]byte(`[{"id":8,"flags":{"movie":true}}]`))
		case "/entries/8/files":
			_, _ = w.Write([]byte(`[{"name":"Movie.zip","url":"https://jimaku.cc/movie.zip","size":2000}]`))
		}
	}))
	defer server.Close()
	enable := true
	settings := map[string]providercore.SettingValue{"enableArchivesDownload": {BooleanValue: &enable}}
	got, err := adapter.Search(context.Background(), testService{server.Client()}, providercore.Config{BaseURL: server.URL, APIKey: &secret, Settings: settings}, providercore.SearchRequest{MediaType: "movie", Title: "Movie", LanguageID: "ja"})
	if err != nil { t.Fatal(err) }
	if len(got) != 1 || got[0].Format != "zip" { t.Fatalf("unexpected candidates: %#v", got) }
}

func TestDownloadExtractsArchive(t *testing.T) {
	secret := "token"
	archive := zipData(t, map[string]string{"sub.srt": "hello"})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(archive) }))
	defer server.Close()
	got, err := adapter.Download(context.Background(), testService{server.Client()}, providercore.Config{APIKey: &secret}, providercore.Candidate{ReleaseName: "sub.zip", SourceURL: server.URL + "/sub.zip"})
	if err != nil { t.Fatal(err) }
	if strings.TrimSpace(string(got.Content)) != "hello" { t.Fatalf("content = %q", got.Content) }
}

func zipData(t *testing.T, files map[string]string) []byte {
	t.Helper(); var buf bytes.Buffer; zw := zip.NewWriter(&buf)
	for name, content := range files { w, err := zw.Create(name); if err != nil { t.Fatal(err) }; _, _ = w.Write([]byte(content)) }
	if err := zw.Close(); err != nil { t.Fatal(err) }
	return buf.Bytes()
}
