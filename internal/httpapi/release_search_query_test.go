package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSCNMedia002ReleaseSearchQueryHandlesEmptyAndTrimmedInput(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/movies/id/releases/search", nil)
	query, ok := releaseSearchQuery(httptest.NewRecorder(), request)
	if !ok || query != "" {
		t.Fatalf("empty body query = %q, ok = %v", query, ok)
	}

	request = httptest.NewRequest(
		http.MethodPost,
		"/movies/id/releases/search",
		strings.NewReader(`{"query":"  Scenario.Movie.2026  "}`),
	)
	query, ok = releaseSearchQuery(httptest.NewRecorder(), request)
	if !ok || query != "Scenario.Movie.2026" {
		t.Fatalf("trimmed query = %q, ok = %v", query, ok)
	}

	request = httptest.NewRequest(
		http.MethodPost,
		"/movies/id/releases/search",
		strings.NewReader(`{"query":null}`),
	)
	query, ok = releaseSearchQuery(httptest.NewRecorder(), request)
	if !ok || query != "" {
		t.Fatalf("null query = %q, ok = %v", query, ok)
	}
}

func TestSCNMedia002ReleaseSearchQueryRejectsMalformedJSON(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/movies/id/releases/search",
		strings.NewReader(`{"query":`),
	)

	query, ok := releaseSearchQuery(response, request)

	if ok || query != "" {
		t.Fatalf("malformed query = %q, ok = %v", query, ok)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}
