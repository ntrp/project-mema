package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"media-manager/internal/config"
)

func TestMetadataDetailsRouteIsMounted(t *testing.T) {
	router := chi.NewRouter()
	HandlerFromMux(NewServer(config.Config{}, nil, nil, nil, nil, nil, nil), router)

	request := httptest.NewRequest(http.MethodGet, "/media/metadata/tmdb/movie/936075", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected metadata details route to reach auth middleware, got status %d and body %q", response.Code, response.Body.String())
	}
}

func TestSettingsUsersRouteIsMounted(t *testing.T) {
	router := chi.NewRouter()
	HandlerFromMux(NewServer(config.Config{}, nil, nil, nil, nil, nil, nil), router)

	request := httptest.NewRequest(http.MethodGet, "/settings/users", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected settings users route to reach auth middleware, got status %d and body %q", response.Code, response.Body.String())
	}
}

func TestSettingsProfilesRouteIsMounted(t *testing.T) {
	router := chi.NewRouter()
	HandlerFromMux(NewServer(config.Config{}, nil, nil, nil, nil, nil, nil), router)

	request := httptest.NewRequest(http.MethodGet, "/settings/profiles", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected settings profiles route to reach auth middleware, got status %d and body %q", response.Code, response.Body.String())
	}
}

func TestMediaRequestsRouteIsMounted(t *testing.T) {
	router := chi.NewRouter()
	HandlerFromMux(NewServer(config.Config{}, nil, nil, nil, nil, nil, nil), router)

	request := httptest.NewRequest(http.MethodGet, "/media/requests", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected media requests route to reach auth middleware, got status %d and body %q", response.Code, response.Body.String())
	}
}
