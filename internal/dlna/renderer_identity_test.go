package dlna

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"media-manager/internal/dlna/controlpoint"
)

func TestRendererOverrideDeniedUsesVerifiedUUIDAndDenyDominates(t *testing.T) {
	overrides := []rendererDeviceOverride{
		{Allowed: false, IPAddress: "192.0.2.4"},
		{Allowed: true, RendererUUID: "uuid:allowed"},
		{Allowed: false, RendererUUID: "uuid:blocked"},
	}
	tests := map[string]struct {
		request RendererRequest
		denied  bool
	}{
		"ip denied":                       {RendererRequest{ClientIP: "192.0.2.4"}, true},
		"unverified allow cannot beat ip": {RendererRequest{RendererUUID: "uuid:allowed", ClientIP: "192.0.2.4"}, true},
		"verified uuid denied":            {RendererRequest{RendererUUID: "uuid:blocked", ClientIP: "192.0.2.5", rendererUUIDVerified: true}, true},
		"ip deny beats verified allow":    {RendererRequest{RendererUUID: "uuid:allowed", ClientIP: "192.0.2.4", rendererUUIDVerified: true}, true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := rendererOverrideDenied(overrides, test.request); got != test.denied {
				t.Fatalf("denied = %v, want %v", got, test.denied)
			}
		})
	}
}

func TestRendererRequestDeniedDoesNotTrustUnboundRememberedUUID(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.profileCache = rendererProfileCacheState{loaded: true, overrides: []rendererDeviceOverride{
		{Allowed: false, IPAddress: "192.0.2.8"},
		{Allowed: true, RendererUUID: "uuid:allowed"},
	}}
	manager.recentClients["192.0.2.8"] = ClientStatus{RendererUUID: "uuid:allowed"}
	request := httptest.NewRequest("GET", "/dlna/resource/item", nil)
	request.RemoteAddr = "192.0.2.8:1234"
	if !manager.rendererRequestDenied(request) {
		t.Fatal("unbound remembered UUID bypassed IP denial")
	}
}

func TestRendererRequestDeniedIPDenyBeatsPersistedBoundRememberedUUID(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.profileCache = rendererProfileCacheState{loaded: true, overrides: []rendererDeviceOverride{
		{Allowed: false, IPAddress: "192.0.2.8"},
		{Allowed: true, RendererUUID: "uuid:allowed", IPAddress: "192.0.2.8"},
	}}
	manager.recentClients["192.0.2.8"] = ClientStatus{RendererUUID: "uuid:allowed"}
	request := httptest.NewRequest("GET", "/dlna/resource/item", nil)
	request.RemoteAddr = "192.0.2.8:1234"
	if !manager.rendererRequestDenied(request) {
		t.Fatal("persisted UUID/IP allow bypassed IP denial")
	}
}

func TestSpoofedDiscoveryAndUUIDHeaderCannotBypassIPDeny(t *testing.T) {
	manager, calls := deniedTestManager(t)
	manager.profileCache.overrides = append(manager.profileCache.overrides, rendererDeviceOverride{Allowed: true, RendererUUID: "uuid:allowed"})
	manager.controlPoint = controlpoint.NewService(nil, nil, nil)
	manager.controlPoint.Registry().Upsert(controlpoint.Renderer{UDN: "uuid:allowed"}, net.ParseIP("127.0.0.1"), time.Minute)
	request := soapRequest(
		"/dlna/control/content-directory",
		"urn:schemas-upnp-org:service:ContentDirectory:1#Browse",
		"<u:Browse xmlns:u=\"urn:schemas-upnp-org:service:ContentDirectory:1\"/>",
	)
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("X-Mema-DLNA-Renderer-UUID", "uuid:allowed")
	response := httptest.NewRecorder()
	manager.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>606</errorCode>") {
		t.Fatalf("response = %d %s", response.Code, response.Body.String())
	}
	if calls.Load() != 0 {
		t.Fatalf("catalogue calls = %d", calls.Load())
	}
}

func TestSpoofedUUIDCannotBeatIPDenyForSOAP(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.profileCache = rendererProfileCacheState{loaded: true, overrides: []rendererDeviceOverride{
		{Allowed: false, IPAddress: "127.0.0.1"},
		{Allowed: true, RendererUUID: "uuid:allowed"},
	}}
	request := soapRequest(
		"/dlna/control/content-directory",
		"urn:schemas-upnp-org:service:ContentDirectory:1#GetSystemUpdateID",
		"<u:GetSystemUpdateID xmlns:u=\"urn:schemas-upnp-org:service:ContentDirectory:1\"/>",
	)
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("X-Mema-DLNA-Renderer-UUID", "uuid:allowed")
	response := httptest.NewRecorder()
	manager.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>606</errorCode>") {
		t.Fatalf("response = %d %s", response.Code, response.Body.String())
	}
}

func TestPersistedBoundUUIDAllowCannotBeatIPDenyForSOAP(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.profileCache = rendererProfileCacheState{loaded: true, overrides: []rendererDeviceOverride{
		{Allowed: false, IPAddress: "127.0.0.1"},
		{Allowed: true, RendererUUID: "uuid:allowed", IPAddress: "127.0.0.1"},
	}}
	request := soapRequest(
		"/dlna/control/content-directory",
		"urn:schemas-upnp-org:service:ContentDirectory:1#GetSystemUpdateID",
		"<u:GetSystemUpdateID xmlns:u=\"urn:schemas-upnp-org:service:ContentDirectory:1\"/>",
	)
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("X-Mema-DLNA-Renderer-UUID", "uuid:allowed")
	response := httptest.NewRecorder()
	manager.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>606</errorCode>") {
		t.Fatalf("response = %d %s", response.Code, response.Body.String())
	}
}

func TestSpoofedRendererUUIDHeadersCannotBypassIPDeny(t *testing.T) {
	headers := map[string]string{
		"X-Mema-DLNA-Renderer-UUID": "uuid:allowed",
		"X-Renderer-UUID":           "uuid:allowed",
		"USN":                       "uuid:allowed::urn:schemas-upnp-org:device:MediaRenderer:1",
	}
	for name, value := range headers {
		t.Run(name, func(t *testing.T) {
			manager := NewManager(nil, "http://127.0.0.1:18080")
			manager.profileCache = rendererProfileCacheState{loaded: true, overrides: []rendererDeviceOverride{
				{Allowed: false, IPAddress: "127.0.0.1"},
				{Allowed: true, RendererUUID: "uuid:allowed"},
			}}
			request := httptest.NewRequest(http.MethodGet, "/dlna/rootDesc.xml", nil)
			request.RemoteAddr = "127.0.0.1:1234"
			request.Header.Set(name, value)
			response := httptest.NewRecorder()
			manager.Handler().ServeHTTP(response, request)
			if response.Code != http.StatusForbidden {
				t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
			}
			clients := manager.Status().RecentClients
			if len(clients) != 1 || clients[0].RendererUUID != "" {
				t.Fatalf("spoofed UUID was remembered: %#v", clients)
			}
		})
	}
}

func TestManualUUIDProfileRequiresPersistedSourceIPBinding(t *testing.T) {
	profiles := []RendererProfile{{ID: "special", Name: "Special"}, {ID: "generic", Name: "Generic"}}
	request := httptest.NewRequest(http.MethodGet, "/dlna/rootDesc.xml", nil)
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("X-Renderer-UUID", "uuid:renderer")

	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.profileCache = rendererProfileCacheState{loaded: true, profiles: profiles, overrides: []rendererDeviceOverride{
		{Allowed: true, RendererUUID: "uuid:renderer", ProfileID: "special"},
	}}
	if got := manager.RendererProfileFromRequest(request).ID; got != "generic" {
		t.Fatalf("unbound header selected profile %q", got)
	}

	manager.profileCache.overrides[0].IPAddress = "127.0.0.1"
	if got := manager.RendererProfileFromRequest(request).ID; got != "special" {
		t.Fatalf("bound UUID selected profile %q", got)
	}
}

func TestDiagnosticsDropStaleUnverifiedUUID(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.profileCache = rendererProfileCacheState{loaded: true}
	manager.recentClients["127.0.0.1"] = ClientStatus{IP: "127.0.0.1", RendererUUID: "uuid:stale"}
	request := httptest.NewRequest(http.MethodGet, "/dlna/rootDesc.xml", nil)
	request.RemoteAddr = "127.0.0.1:1234"
	manager.Handler().ServeHTTP(httptest.NewRecorder(), request)
	clients := manager.Status().RecentClients
	if len(clients) != 1 || clients[0].RendererUUID != "" {
		t.Fatalf("stale UUID was restored: %#v", clients)
	}
}
