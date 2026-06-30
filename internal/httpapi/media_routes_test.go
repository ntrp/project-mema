package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"media-manager/internal/config"
)

func TestSessionRouteAllowsAnonymousProbe(t *testing.T) {
	router := chi.NewRouter()
	HandlerFromMux(NewServer(config.Config{}, nil, nil, nil, nil, nil, nil), router)

	request := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected anonymous session probe to return 200, got status %d and body %q", response.Code, response.Body.String())
	}
	var body SessionResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode session response: %v", err)
	}
	if body.Authenticated {
		t.Fatal("expected anonymous session probe to be unauthenticated")
	}
}

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
