package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestScenarioSCNSystem004StaticHandlerServesShellAndAssets(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "200.html"), []byte("app shell"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "asset.txt"), []byte("asset"), 0o644); err != nil {
		t.Fatal(err)
	}
	handler := StaticHandler(dir)

	for _, path := range []string{"/", "/missing/route"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK || response.Body.String() != "app shell" {
			t.Fatalf("%s returned %d %q", path, response.Code, response.Body.String())
		}
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/asset.txt", nil))
	if response.Code != http.StatusOK || response.Body.String() != "asset" {
		t.Fatalf("asset returned %d %q", response.Code, response.Body.String())
	}
}
