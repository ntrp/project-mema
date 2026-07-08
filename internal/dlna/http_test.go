package dlna

import (
	"bytes"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestBaseURLHonorsForwardedHeaders(t *testing.T) {
	request := httptest.NewRequest("GET", "http://internal/dlna/rootDesc.xml", nil)
	request.RemoteAddr = "127.0.0.1:1234"
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
	request.RemoteAddr = "127.0.0.1:1234"
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
	request.RemoteAddr = "127.0.0.1:1234"
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "<deviceType>urn:schemas-upnp-org:device:MediaServer:1</deviceType>") {
		t.Fatalf("root response missing MediaServer device:\n%s", response.Body.String())
	}
	if response.Header().Get("Content-Length") == "" || response.Header().Get("Connection") != "close" {
		t.Fatalf("root response headers = %#v", response.Header())
	}
}

func TestHandlerServesAdvertisedIcon(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	request := httptest.NewRequest("GET", "http://internal/dlna/icon-48.png", nil)
	request.RemoteAddr = "127.0.0.1:1234"
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "image/png" {
		t.Fatalf("Content-Type = %q", response.Header().Get("Content-Type"))
	}
	image, err := png.DecodeConfig(bytes.NewReader(response.Body.Bytes()))
	if err != nil {
		t.Fatalf("icon did not decode: %v", err)
	}
	if image.Width != 48 || image.Height != 48 {
		t.Fatalf("icon size = %dx%d", image.Width, image.Height)
	}
}

func TestHandlerServesAdvertisedEventURLs(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.events.client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil
	})
	for _, path := range []string{
		"/dlna/events/content-directory",
		"/dlna/events/connection-manager",
		"/dlna/events/media-receiver-registrar",
	} {
		request := httptest.NewRequest("SUBSCRIBE", "http://internal"+path, nil)
		request.RemoteAddr = "127.0.0.1:1234"
		request.Header.Set("CALLBACK", "<http://127.0.0.1/callback>")
		response := httptest.NewRecorder()

		manager.Handler().ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("%s status = %d body=%s", path, response.Code, response.Body.String())
		}
	}
}
