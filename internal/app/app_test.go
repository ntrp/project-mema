package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"media-manager/internal/config"
)

func TestDevelopmentRouterDoesNotServeFrontendFallback(t *testing.T) {
	api := chi.NewRouter()
	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := appRouter(config.Config{AppEnv: "development"}, api)

	apiResponse := httptest.NewRecorder()
	handler.ServeHTTP(apiResponse, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if apiResponse.Code != http.StatusNoContent {
		t.Fatalf("api status = %d", apiResponse.Code)
	}

	pageResponse := httptest.NewRecorder()
	handler.ServeHTTP(pageResponse, httptest.NewRequest(http.MethodGet, "/movies/example", nil))
	if pageResponse.Code != http.StatusNotFound {
		t.Fatalf("frontend fallback status = %d", pageResponse.Code)
	}
}
