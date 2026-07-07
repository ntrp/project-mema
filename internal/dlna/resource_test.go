package dlna

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"media-manager/internal/dlna/content"
	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func TestResourceServesDirectFileRanges(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	request := httptest.NewRequest("GET", "/dlna/resource/"+url.PathEscape(resourceID), nil)
	request.Header.Set("Range", "bytes=0-3")
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusPartialContent {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if response.Body.String() != "0123" {
		t.Fatalf("body = %q", response.Body.String())
	}
	if response.Header().Get("Accept-Ranges") != "bytes" {
		t.Fatalf("Accept-Ranges = %q", response.Header().Get("Accept-Ranges"))
	}
}

func TestResourceHLSHeadDoesNotStartTranscode(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	request := httptest.NewRequest("HEAD", "/dlna/resource/"+url.PathEscape(resourceID)+"?mode=hls", nil)
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "application/vnd.apple.mpegurl" {
		t.Fatalf("Content-Type = %q", response.Header().Get("Content-Type"))
	}
	if response.Body.Len() != 0 {
		t.Fatalf("HEAD body = %q", response.Body.String())
	}
}

func TestResourceSegmentHeadDoesNotStartTranscode(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	target := "/dlna/resource/" + url.PathEscape(resourceID) +
		"/segment?segmentStartSeconds=0&segmentDurationSeconds=6"
	request := httptest.NewRequest("HEAD", target, nil)
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "video/mp2t" {
		t.Fatalf("Content-Type = %q", response.Header().Get("Content-Type"))
	}
	if response.Body.Len() != 0 {
		t.Fatalf("HEAD body = %q", response.Body.String())
	}
}

func resourceTestManager(t *testing.T, payload string) (*Manager, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "Scenario.Movie.mp4")
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	item := storage.MediaItem{
		ID:        uuid.New(),
		Type:      "movie",
		Title:     "Scenario Movie",
		FilePaths: []string{path},
	}
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.source = contentFakeSource{items: []storage.MediaItem{item}}
	children, err := manager.contentTree().BrowseChildren(context.Background(), content.EncodeID(content.RootContainerRef("movies")))
	if err != nil || len(children) != 1 {
		t.Fatalf("children = %#v err=%v", children, err)
	}
	files, err := manager.contentTree().BrowseChildren(context.Background(), children[0].ID)
	if err != nil || len(files) != 1 {
		t.Fatalf("files = %#v err=%v", files, err)
	}
	if strings.Contains(files[0].ID, path) {
		t.Fatalf("resource id exposes path: %q", files[0].ID)
	}
	return manager, files[0].ID
}
