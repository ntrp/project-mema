package dlna

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"media-manager/internal/dlna/content"
	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func TestSubtitleRouteServesExternalSRT(t *testing.T) {
	manager, resourceID := subtitleTestManager(t)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/dlna/subtitle/"+url.PathEscape(resourceID)+"/0", nil)
	request.RemoteAddr = "127.0.0.1:1234"

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "application/x-subrip; charset=utf-8" {
		t.Fatalf("content type = %q", response.Header().Get("Content-Type"))
	}
	if response.Body.String() != "1\n00:00:00,000 --> 00:00:01,000\nHi\n" {
		t.Fatalf("body = %q", response.Body.String())
	}
}

func TestSubtitleConvertHeadDoesNotStartConversion(t *testing.T) {
	manager, resourceID := subtitleTestManagerWithFormat(t, "ass")
	response := httptest.NewRecorder()
	request := httptest.NewRequest("HEAD", "/dlna/subtitle/"+url.PathEscape(resourceID)+"/0", nil)
	request.RemoteAddr = "127.0.0.1:1234"

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "text/vtt; charset=utf-8" {
		t.Fatalf("HEAD response = %d %s", response.Code, response.Header().Get("Content-Type"))
	}
	if response.Body.Len() != 0 {
		t.Fatalf("HEAD body = %q", response.Body.String())
	}
}

func subtitleTestManager(t *testing.T) (*Manager, string) {
	return subtitleTestManagerWithFormat(t, "srt")
}

func subtitleTestManagerWithFormat(t *testing.T, format string) (*Manager, string) {
	t.Helper()
	dir := t.TempDir()
	mediaPath := filepath.Join(dir, "Movie.mkv")
	subtitlePath := filepath.Join(dir, "Movie.eng."+format)
	if err := os.WriteFile(mediaPath, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(subtitlePath, []byte("1\n00:00:00,000 --> 00:00:01,000\nHi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lang := "eng"
	item := storage.MediaItem{
		ID:        uuid.New(),
		Type:      "movie",
		Title:     "Movie",
		FilePaths: []string{mediaPath},
		Sidecars: []storage.MediaItemSidecar{{
			MediaFilePath: mediaPath,
			FilePath:      subtitlePath,
			SidecarType:   storage.MediaSidecarSubtitle,
			LanguageID:    &lang,
			Format:        &format,
		}},
	}
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.source = contentFakeSource{items: []storage.MediaItem{item}}
	movie, err := manager.contentTree().BrowseChildren(t.Context(), content.EncodeID(content.RootContainerRef("movies")))
	if err != nil || len(movie) != 1 {
		t.Fatalf("movie = %#v err=%v", movie, err)
	}
	files, err := manager.contentTree().BrowseChildren(t.Context(), movie[0].ID)
	if err != nil || len(files) != 1 {
		t.Fatalf("files = %#v err=%v", files, err)
	}
	return manager, files[0].ID
}
