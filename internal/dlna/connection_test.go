package dlna

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSourceProtocolInfoAdvertisesServedFormats(t *testing.T) {
	source := SourceProtocolInfo()
	for _, want := range []string{
		"video/mp4",
		"video/x-matroska",
		"video/mp2t",
		"application/vnd.apple.mpegurl",
	} {
		if !strings.Contains(source, want) {
			t.Fatalf("SourceProtocolInfo missing %q: %s", want, source)
		}
	}
	if strings.Contains(source, "image/") || strings.Contains(source, "audio/") {
		t.Fatalf("SourceProtocolInfo advertises unsupported class: %s", source)
	}
}

func TestConnectionManagerSOAPActions(t *testing.T) {
	dispatcher := NewManager(nil, "http://127.0.0.1:18080").SOAPDispatcher()
	response := httptest.NewRecorder()
	request := soapRequest(
		"/dlna/control/connection-manager",
		"urn:schemas-upnp-org:service:ConnectionManager:1#GetProtocolInfo",
		`<u:GetProtocolInfo xmlns:u="urn:schemas-upnp-org:service:ConnectionManager:1"/>`,
	)

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "<Source>http-get:*:video/mp4") ||
		!strings.Contains(response.Body.String(), "<Sink></Sink>") {
		t.Fatalf("GetProtocolInfo response = %s", response.Body.String())
	}

	response = httptest.NewRecorder()
	request = soapRequest(
		"/dlna/control/connection-manager",
		"urn:schemas-upnp-org:service:ConnectionManager:1#GetCurrentConnectionInfo",
		`<u:GetCurrentConnectionInfo xmlns:u="urn:schemas-upnp-org:service:ConnectionManager:1"><ConnectionID>0</ConnectionID></u:GetCurrentConnectionInfo>`,
	)
	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "<Direction>Output</Direction>") ||
		!strings.Contains(response.Body.String(), "<Status>OK</Status>") {
		t.Fatalf("GetCurrentConnectionInfo response = %d %s", response.Code, response.Body.String())
	}
}

func TestConnectionManagerInvalidConnectionReturns706(t *testing.T) {
	dispatcher := NewManager(nil, "http://127.0.0.1:18080").SOAPDispatcher()
	response := httptest.NewRecorder()
	request := soapRequest(
		"/dlna/control/connection-manager",
		"urn:schemas-upnp-org:service:ConnectionManager:1#GetCurrentConnectionInfo",
		`<u:GetCurrentConnectionInfo xmlns:u="urn:schemas-upnp-org:service:ConnectionManager:1"><ConnectionID>99</ConnectionID></u:GetCurrentConnectionInfo>`,
	)

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>706</errorCode>") {
		t.Fatalf("fault = %d %s", response.Code, response.Body.String())
	}
}

func soapRequest(path string, action string, body string) *http.Request {
	envelope := `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body>` +
		body + `</s:Body></s:Envelope>`
	request := httptest.NewRequest("POST", path, strings.NewReader(envelope))
	request.Header.Set("SOAPACTION", `"`+action+`"`)
	return request
}
