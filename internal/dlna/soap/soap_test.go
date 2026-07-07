package soap

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseRequestReadsSOAPActionArguments(t *testing.T) {
	body := `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><u:Browse xmlns:u="urn:test"><ObjectID>0</ObjectID><BrowseFlag>BrowseDirectChildren</BrowseFlag></u:Browse></s:Body></s:Envelope>`
	request := httptest.NewRequest("POST", "/control", strings.NewReader(body))
	request.Header.Set("SOAPACTION", `"urn:test#Browse"`)

	action, err := ParseRequest(request, "urn:test")
	if err != nil {
		t.Fatalf("ParseRequest returned error: %v", err)
	}
	if action.Name != "Browse" || action.Args["ObjectID"] != "0" || action.Args["BrowseFlag"] != "BrowseDirectChildren" {
		t.Fatalf("action = %#v", action)
	}
}

func TestDispatcherReturnsFaultForUnknownAction(t *testing.T) {
	dispatcher := NewDispatcher()
	dispatcher.Register("/control", "urn:test", map[string]HandlerFunc{})
	request := httptest.NewRequest("POST", "/control", strings.NewReader(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><u:Missing xmlns:u="urn:test"/></s:Body></s:Envelope>`))
	request.Header.Set("SOAPACTION", `"urn:test#Missing"`)
	response := httptest.NewRecorder()

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>401</errorCode>") {
		t.Fatalf("fault response = %d %s", response.Code, response.Body.String())
	}
}

func TestDispatcherReturnsFaultForInvalidArguments(t *testing.T) {
	dispatcher := NewDispatcher()
	dispatcher.Register("/control", "urn:test", map[string]HandlerFunc{
		"NeedsID": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			if _, err := RequiredArg(args, "ObjectID"); err != nil {
				return nil, err
			}
			return map[string]string{"Result": "ok"}, nil
		},
	})
	request := httptest.NewRequest("POST", "/control", strings.NewReader(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><u:NeedsID xmlns:u="urn:test"/></s:Body></s:Envelope>`))
	request.Header.Set("SOAPACTION", `"urn:test#NeedsID"`)
	response := httptest.NewRecorder()

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>402</errorCode>") {
		t.Fatalf("fault response = %d %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "<errorDescription>Missing argument: ObjectID</errorDescription>") {
		t.Fatalf("fault response missing argument description: %s", response.Body.String())
	}
}

func TestWriteFaultMarshalsUPnPErrorXML(t *testing.T) {
	response := httptest.NewRecorder()

	WriteFault(response, Error{Code: 401, Description: "Invalid Action"})

	body := response.Body.String()
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d body=%s", response.Code, body)
	}
	for _, want := range []string{
		"<faultstring>UPnPError</faultstring>",
		`<UPnPError xmlns="urn:schemas-upnp-org:control-1-0">`,
		"<errorCode>401</errorCode>",
		"<errorDescription>Invalid Action</errorDescription>",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("fault XML missing %q:\n%s", want, body)
		}
	}
}

func TestDispatcherWritesActionResponse(t *testing.T) {
	dispatcher := NewDispatcher()
	dispatcher.Register("/control", "urn:test", map[string]HandlerFunc{
		"Ping": func(ctx context.Context, args map[string]string) (map[string]string, error) {
			return map[string]string{"Result": "ok"}, nil
		},
	})
	request := httptest.NewRequest("POST", "/control", strings.NewReader(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><u:Ping xmlns:u="urn:test"/></s:Body></s:Envelope>`))
	request.Header.Set("SOAPACTION", `"urn:test#Ping"`)
	response := httptest.NewRecorder()

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "<Result>ok</Result>") {
		t.Fatalf("response = %d %s", response.Code, response.Body.String())
	}
}
