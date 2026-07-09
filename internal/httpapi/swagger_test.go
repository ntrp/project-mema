package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwaggerUIHandlerServesOpenAPISpecViewer(t *testing.T) {
	response := httptest.NewRecorder()

	SwaggerUIHandler(response, httptest.NewRequest(http.MethodGet, "/api/docs", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	body := response.Body.String()
	if !strings.Contains(body, "SwaggerUIBundle") || !strings.Contains(body, "/api/openapi.yaml") {
		t.Fatalf("swagger ui body missing spec viewer: %q", body)
	}
}

func TestOpenAPISpecHandlerServesContractSource(t *testing.T) {
	response := httptest.NewRecorder()

	OpenAPISpecHandler(response, httptest.NewRequest(http.MethodGet, "/api/openapi.yaml", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body = %q", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "openapi: 3.") {
		t.Fatalf("openapi body missing version header: %q", response.Body.String())
	}
}
