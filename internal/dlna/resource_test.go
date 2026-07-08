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
	request.RemoteAddr = "127.0.0.1:1234"
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
	request.RemoteAddr = "127.0.0.1:1234"
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

func TestResourceTranscodeRangeUsesSeekableMatroskaCache(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	manager.remuxDir = t.TempDir()
	installFakeFFmpeg(t)
	request := httptest.NewRequest("GET", "/dlna/resource/"+url.PathEscape(resourceID)+"?mode=transcode", nil)
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("Range", "bytes=1-3")
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusPartialContent {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "video/x-matroska" {
		t.Fatalf("Content-Type = %q", response.Header().Get("Content-Type"))
	}
	if response.Header().Get("ContentFeatures.DLNA.ORG") != "DLNA.ORG_OP=01;DLNA.ORG_CI=1" {
		t.Fatalf("ContentFeatures = %q", response.Header().Get("ContentFeatures.DLNA.ORG"))
	}
	if response.Header().Get("Accept-Ranges") != "bytes" {
		t.Fatalf("Accept-Ranges = %q", response.Header().Get("Accept-Ranges"))
	}
	if response.Body.String() != "emu" {
		t.Fatalf("body = %q", response.Body.String())
	}
}

func TestResourceRemuxRangeUsesSeekableMPEGTSCache(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	manager.remuxDir = t.TempDir()
	installFakeFFmpeg(t)
	request := httptest.NewRequest("GET", "/dlna/resource/"+url.PathEscape(resourceID)+"?mode=remux", nil)
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("Range", "bytes=1-3")
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusPartialContent {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "video/mp2t" {
		t.Fatalf("Content-Type = %q", response.Header().Get("Content-Type"))
	}
	if response.Body.String() != "emu" {
		t.Fatalf("body = %q", response.Body.String())
	}
}

func TestInitialRangeDoesNotForceRemuxCache(t *testing.T) {
	for _, value := range []string{"", "bytes=0-", "bytes=0-65535"} {
		if isSeekRange(value) {
			t.Fatalf("range %q should not force cache", value)
		}
	}
	for _, value := range []string{"bytes=1-", "bytes=65536-"} {
		if !isSeekRange(value) {
			t.Fatalf("range %q should force cache", value)
		}
	}
}

func TestProfileSeekModeCanDisableByteSeek(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	manager.profileCache = rendererProfileCacheState{
		loaded: true,
		profiles: []RendererProfile{{
			ID: "generic", Name: "Generic DLNA",
			DeliveryRules: RendererDeliveryRules{DirectPlay: true, SeekMode: seekModeNone},
		}},
	}
	request := httptest.NewRequest("GET", "/dlna/resource/"+url.PathEscape(resourceID), nil)
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("Range", "bytes=1-3")
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusRequestedRangeNotSatisfiable {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
}

func TestResourceSegmentHeadDoesNotStartTranscode(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	target := "/dlna/resource/" + url.PathEscape(resourceID) +
		"/segment?segmentStartSeconds=0&segmentDurationSeconds=6"
	request := httptest.NewRequest("HEAD", target, nil)
	request.RemoteAddr = "127.0.0.1:1234"
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

func installFakeFFmpeg(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "ffmpeg")
	script := "#!/bin/sh\n" +
		"for last do :; done\n" +
		"printf remux > \"$last\"\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}
