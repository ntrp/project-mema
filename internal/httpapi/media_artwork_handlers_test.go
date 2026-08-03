package httpapi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAniDBArtworkIsServedFromTheLocalCache(t *testing.T) {
	mediaDir := t.TempDir()
	artworkDir := filepath.Join(mediaDir, ".mema", "anidb", "artwork")
	if err := os.MkdirAll(artworkDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(artworkDir, "456.jpg"), []byte("jpeg"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := testConfig()
	cfg.MediaDataDir = mediaDir
	router := chi.NewRouter()
	HandlerFromMux(NewServer(cfg, testSettingsStore(t), nil, nil, nil, nil, nil), router)
	anonymous := httptest.NewRecorder()
	router.ServeHTTP(anonymous, httptest.NewRequest(http.MethodGet, "/media/artwork/anidb/456.jpg", nil))
	if anonymous.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status = %d", anonymous.Code)
	}

	request := httptest.NewRequest(http.MethodGet, "/media/artwork/anidb/456.jpg", nil)
	request.AddCookie(authenticatedCookie(t, router))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Body.String() != "jpeg" {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
	if response.Header().Get("Cache-Control") != "private, max-age=31536000, immutable" {
		t.Fatalf("cache control = %q", response.Header().Get("Cache-Control"))
	}
}
