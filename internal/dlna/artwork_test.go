package dlna

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"media-manager/internal/dlna/content"
	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func TestArtworkServesFallbackIcon(t *testing.T) {
	manager, objectID := artworkTestManager(t, "")
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/dlna/artwork/"+url.PathEscape(objectID), nil)

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("fallback response = %d %s", response.Code, response.Header().Get("Content-Type"))
	}
	if len(response.Body.Bytes()) < 8 || string(response.Body.Bytes()[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Fatalf("fallback body is not png: %x", response.Body.Bytes())
	}
}

func TestArtworkServesLocalMetadataArtwork(t *testing.T) {
	dir := t.TempDir()
	poster := filepath.Join(dir, "poster.png")
	if err := os.WriteFile(poster, fallbackPNG, 0o644); err != nil {
		t.Fatal(err)
	}
	manager, objectID := artworkTestManager(t, poster)
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/dlna/artwork/"+url.PathEscape(objectID), nil)

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Body.Len() != len(fallbackPNG) {
		t.Fatalf("local artwork response = %d len=%d", response.Code, response.Body.Len())
	}
}

func TestThumbnailHeadDoesNotGenerate(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	response := httptest.NewRecorder()
	request := httptest.NewRequest("HEAD", "/dlna/artwork/"+url.PathEscape(resourceID)+"?kind=thumbnail", nil)

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "image/jpeg" {
		t.Fatalf("thumbnail HEAD = %d %s", response.Code, response.Header().Get("Content-Type"))
	}
	if response.Body.Len() != 0 {
		t.Fatalf("HEAD body = %q", response.Body.String())
	}
}

func TestThumbnailCachePathIncludesFileIdentityAndModTime(t *testing.T) {
	now := time.Unix(100, 0)
	left := thumbnailCachePath("cache", "/media/a.mkv", now, 10)
	right := thumbnailCachePath("cache", "/media/a.mkv", now.Add(time.Second), 10)
	other := thumbnailCachePath("cache", "/media/a.mkv", now, 11)

	if left == right || left == other {
		t.Fatalf("cache paths not unique: %q %q %q", left, right, other)
	}
}

func artworkTestManager(t *testing.T, poster string) (*Manager, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "Scenario.Movie.mp4")
	if err := os.WriteFile(path, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}
	item := storage.MediaItem{
		ID:         uuid.New(),
		Type:       "movie",
		Title:      "Scenario Movie",
		FilePaths:  []string{path},
		PosterPath: optionalPoster(poster),
	}
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.source = contentFakeSource{items: []storage.MediaItem{item}}
	children, err := manager.contentTree().BrowseChildren(t.Context(), content.EncodeID(content.RootContainerRef("movies")))
	if err != nil || len(children) != 1 {
		t.Fatalf("children = %#v err=%v", children, err)
	}
	return manager, children[0].ID
}

func optionalPoster(path string) *string {
	if path == "" {
		return nil
	}
	return &path
}
