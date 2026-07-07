package dlna

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestBaseURLHonorsForwardedHeaders(t *testing.T) {
	request := httptest.NewRequest("GET", "http://internal/dlna/rootDesc.xml", nil)
	request.Host = "127.0.0.1:18080"
	request.Header.Set("X-Forwarded-Proto", "https")
	request.Header.Set("X-Forwarded-Host", "mema.local")

	if got := requestBaseURL(request); got != "https://mema.local" {
		t.Fatalf("requestBaseURL = %q", got)
	}
}

func TestHandlerServesSOAPControlAction(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	body := `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><u:GetSystemUpdateID xmlns:u="urn:schemas-upnp-org:service:ContentDirectory:1"/></s:Body></s:Envelope>`
	request := httptest.NewRequest("POST", "http://internal/dlna/control/content-directory", strings.NewReader(body))
	request.Header.Set("SOAPACTION", `"urn:schemas-upnp-org:service:ContentDirectory:1#GetSystemUpdateID"`)
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "<Id>0</Id>") {
		t.Fatalf("SOAP response missing update id:\n%s", response.Body.String())
	}
}

func TestHandlerServesMountedDLNARoutes(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	request := httptest.NewRequest("GET", "http://internal/dlna/rootDesc.xml", nil)
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "<deviceType>urn:schemas-upnp-org:device:MediaServer:1</deviceType>") {
		t.Fatalf("root response missing MediaServer device:\n%s", response.Body.String())
	}
}
