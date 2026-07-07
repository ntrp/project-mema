package dlna

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMatchRendererProfileKnownClients(t *testing.T) {
	tests := []struct {
		name    string
		request RendererRequest
		want    string
	}{
		{name: "generic", request: RendererRequest{UserAgent: "unknown"}, want: "generic"},
		{name: "vlc", request: RendererRequest{UserAgent: "VLC/3.0"}, want: "vlc"},
		{name: "kodi", request: RendererRequest{FriendlyName: "Kodi Media Center"}, want: "kodi"},
		{name: "samsung", request: RendererRequest{UserAgent: "Samsung Tizen DLNADOC"}, want: "samsung"},
		{name: "lg", request: RendererRequest{UserAgent: "LG webOS TV"}, want: "lg"},
		{name: "sony", request: RendererRequest{UserAgent: "Sony BRAVIA"}, want: "sony"},
		{name: "bubbleupnp", request: RendererRequest{UserAgent: "BubbleUPnP/4.3"}, want: "bubbleupnp"},
		{name: "chromecast", request: RendererRequest{Headers: http.Header{"X-Device": []string{"Google Cast"}}}, want: "chromecast"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := MatchRendererProfile(test.request, nil).ID; got != test.want {
				t.Fatalf("profile = %q, want %q", got, test.want)
			}
		})
	}
}

func TestRendererProfileOverrideBeatsUserAgent(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.SetRendererProfileOverride("192.0.2.44", "chromecast")

	profile := manager.RendererProfile(RendererRequest{ClientIP: "192.0.2.44", UserAgent: "VLC/3.0"})

	if profile.ID != "chromecast" {
		t.Fatalf("profile = %q, want chromecast", profile.ID)
	}
}

func TestProfileDrivesProtocolInfoAndDeliveryPlan(t *testing.T) {
	generic := MatchRendererProfile(RendererRequest{}, nil)
	chromecast := MatchRendererProfile(RendererRequest{UserAgent: "Chromecast"}, nil)

	if !strings.HasPrefix(SourceProtocolInfoForProfile(generic), "http-get:*:video/mp4") {
		t.Fatalf("generic protocol info = %s", SourceProtocolInfoForProfile(generic))
	}
	if !strings.HasPrefix(SourceProtocolInfoForProfile(chromecast), "http-get:*:application/vnd.apple.mpegurl") {
		t.Fatalf("chromecast protocol info = %s", SourceProtocolInfoForProfile(chromecast))
	}
	if DeliveryClientProfile(generic) != "browser" || DeliveryClientProfile(chromecast) != "webkit" {
		t.Fatalf("delivery profiles = %s/%s", DeliveryClientProfile(generic), DeliveryClientProfile(chromecast))
	}
}

func TestConnectionManagerUsesRequestProfile(t *testing.T) {
	dispatcher := NewManager(nil, "http://127.0.0.1:18080").SOAPDispatcher()
	response := httptest.NewRecorder()
	request := soapRequest(
		"/dlna/control/connection-manager",
		"urn:schemas-upnp-org:service:ConnectionManager:1#GetProtocolInfo",
		`<u:GetProtocolInfo xmlns:u="urn:schemas-upnp-org:service:ConnectionManager:1"/>`,
	)
	request.Header.Set("User-Agent", "Chromecast")

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "<Source>http-get:*:application/vnd.apple.mpegurl") {
		t.Fatalf("GetProtocolInfo response = %d %s", response.Code, response.Body.String())
	}
}

func TestHandlerAppliesProfileHeadersAndDisablesEventing(t *testing.T) {
	handler := NewManager(nil, "http://127.0.0.1:18080").Handler()
	request := httptest.NewRequest("SUBSCRIBE", "/dlna/events/content-directory", nil)
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("User-Agent", "Chromecast")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("TransferMode.DLNA.ORG") != "Streaming" {
		t.Fatalf("headers = %#v", response.Header())
	}
}
