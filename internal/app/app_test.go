package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRouteDLNAAllowsCustomUPnPMethods(t *testing.T) {
	dlnaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "SUBSCRIBE" || r.URL.Path != "/events/content-directory" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	appHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach app handler")
	})
	handler := routeDLNA(appHandler, dlnaHandler)
	request := httptest.NewRequest("SUBSCRIBE", "/dlna/events/content-directory", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestRouteDLNARootPathsBeforeFrontend(t *testing.T) {
	dlnaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/control/content-directory" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
	})
	appHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach app handler")
	})
	handler := routeDLNA(appHandler, dlnaHandler)
	request := httptest.NewRequest("POST", "/control/content-directory", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "application/xml" {
		t.Fatalf("response = %d %q", response.Code, response.Header().Get("Content-Type"))
	}
}

func TestSCNSystem006EnsureMediaDataDirCreatesNestedPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "media", "movies")

	if err := ensureMediaDataDir(path); err != nil {
		t.Fatalf("ensureMediaDataDir returned error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected media data directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("%s is not a directory", path)
	}
}

func TestSCNSystem006EnsureMediaDataDirReportsFilesystemFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "media")
	if err := os.WriteFile(path, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write fixture file: %v", err)
	}

	err := ensureMediaDataDir(filepath.Join(path, "movies"))
	if err == nil {
		t.Fatal("expected setup error")
	}
	if !strings.Contains(err.Error(), "media data directory setup failed") {
		t.Fatalf("error = %q", err.Error())
	}
}
