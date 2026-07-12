package subtis

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type httpService struct{ client *http.Client }
func (s httpService) DoProviderRequest(req *http.Request, _ string, _ bool) (*http.Response, error) { return s.client.Do(req) }

func TestSearchCascadesHashBytesNameAlternative(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		if strings.Contains(r.URL.Path, "/file/name/") {
			_, _ = w.Write([]byte(`{"subtitle":{"subtitle_link":"https://cdn/sub.srt"},"title":{"title_name":"Movie Release"}}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()
	got, err := adapter{}.Search(context.Background(), httpService{server.Client()}, providercore.Config{BaseURL: server.URL}, providercore.SearchRequest{MediaType: "movie", LanguageID: "es", FilePath: "/tmp/Movie.Name.mkv", MediaContext: providercore.MediaContext{File: providercore.FileContext{SizeBytes: 123, Hashes: map[string]string{"opensubtitles": "abc"}}}})
	if err != nil { t.Fatal(err) }
	if len(got) != 1 || got[0].SourceURL != "https://cdn/sub.srt" || got[0].ReleaseName != "Movie Release" { t.Fatalf("candidates = %#v", got) }
	if len(paths) != 3 || !strings.Contains(paths[0], "/hash/abc") || !strings.Contains(paths[1], "/bytes/123") || !strings.Contains(paths[2], "/name/Movie.Name.mkv") { t.Fatalf("cascade paths = %#v", paths) }
}

func TestSearchUsesAlternativeAsFuzzyMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/alternative/") { http.NotFound(w, r); return }
		_, _ = w.Write([]byte(`{"subtitle":{"subtitle_link":"https://cdn/alt.srt"},"title":{"title_name":"Alt"}}`))
	}))
	defer server.Close()
	got, err := adapter{}.Search(context.Background(), httpService{server.Client()}, providercore.Config{BaseURL: server.URL}, providercore.SearchRequest{MediaType: "movie", FilePath: "/tmp/Movie.mkv"})
	if err != nil { t.Fatal(err) }
	if len(got) != 1 || !strings.Contains(got[0].ReleaseName, "[fuzzy match]") { t.Fatalf("candidates = %#v", got) }
}

func TestDownloadReturnsSubtitleContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:02,000\nHi")) }))
	defer server.Close()
	got, err := adapter{}.Download(context.Background(), httpService{server.Client()}, providercore.Config{}, providercore.Candidate{SourceURL: server.URL + "/sub.srt"})
	if err != nil { t.Fatal(err) }
	if !strings.Contains(string(got.Content), "Hi") { t.Fatalf("content = %q", got.Content) }
}

func TestSearchComputesHashFromFileWhenMissing(t *testing.T) {
	file := t.TempDir() + "/movie.bin"
	if err := os.WriteFile(file, []byte("0123456789abcdef"), 0o600); err != nil { t.Fatal(err) }
	var firstPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if firstPath == "" { firstPath = r.URL.Path }
		http.NotFound(w, r)
	}))
	defer server.Close()
	_, err := adapter{}.Search(context.Background(), httpService{server.Client()}, providercore.Config{BaseURL: server.URL}, providercore.SearchRequest{MediaType: "movie", FilePath: file})
	if err != nil { t.Fatal(err) }
	if !strings.Contains(firstPath, "/subtitle/find/file/hash/") { t.Fatalf("first path = %s", firstPath) }
}
