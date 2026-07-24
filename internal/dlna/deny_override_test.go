package dlna

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"media-manager/internal/dlna/content"
	"media-manager/internal/dlna/controlpoint"
	"media-manager/internal/storage"
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

func TestDeniedBrowseAndSearchReturnSOAP606BeforeCatalogueWork(t *testing.T) {
	manager, calls := deniedTestManager(t)
	requests := map[string]*http.Request{
		"browse": soapRequest(
			"/dlna/control/content-directory",
			"urn:schemas-upnp-org:service:ContentDirectory:1#Browse",
			strings.TrimSuffix(strings.TrimPrefix(browseEnvelope(content.RootID, string(content.BrowseDirectChildren), "0", "1"), `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body>`), `</s:Body></s:Envelope>`),
		),
		"search": soapRequest(
			"/dlna/control/content-directory",
			"urn:schemas-upnp-org:service:ContentDirectory:1#Search",
			searchEnvelope(content.RootID, `dc:title contains "Scenario"`, "0", "10"),
		),
	}
	for name, request := range requests {
		t.Run(name, func(t *testing.T) {
			request.RemoteAddr = "127.0.0.1:1234"
			response := httptest.NewRecorder()
			manager.Handler().ServeHTTP(response, request)
			if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>606</errorCode>") {
				t.Fatalf("fault = %d %s", response.Code, response.Body.String())
			}
			if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/xml") {
				t.Fatalf("Content-Type = %q", got)
			}
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("catalogue calls = %d", calls.Load())
	}
}

func TestDeniedResourcesReturn403WithoutStreamOrProcessWork(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "media")
	original := manager.source.(contentFakeSource)
	calls := &atomic.Int32{}
	manager.source = countingContentSource{items: original.items, calls: calls}
	manager.profileCache = rendererProfileCacheState{loaded: true, overrides: []rendererDeviceOverride{{Allowed: false, IPAddress: "127.0.0.1"}}}
	manager.rateLimiter = nil
	marker := installDenyMarkerTools(t)
	paths := []string{
		"/dlna/resource/" + url.PathEscape(resourceID),
		"/dlna/resource/" + url.PathEscape(resourceID) + "?mode=remux",
		"/dlna/resource/" + url.PathEscape(resourceID) + "?mode=transcode",
		"/dlna/artwork/" + url.PathEscape(resourceID) + "?kind=thumbnail",
		"/dlna/subtitle/" + url.PathEscape(resourceID) + "/0",
	}
	for _, path := range paths {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		request.RemoteAddr = "127.0.0.1:1234"
		response := httptest.NewRecorder()
		manager.Handler().ServeHTTP(response, request)
		if response.Code != http.StatusForbidden {
			t.Fatalf("%s status = %d body=%s", path, response.Code, response.Body.String())
		}
	}
	if calls.Load() != 0 {
		t.Fatalf("catalogue calls = %d", calls.Load())
	}
	status := manager.Status()
	if len(status.ActiveStreams) != 0 || len(status.ActiveTranscodes) != 0 {
		t.Fatalf("streams = %#v transcodes = %#v", status.ActiveStreams, status.ActiveTranscodes)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("media process was started: %v", err)
	}
}

func TestDefaultAndExplicitlyAllowedRequestsStillPass(t *testing.T) {
	tests := map[string][]rendererDeviceOverride{
		"default":    nil,
		"allowed ip": {{Allowed: true, IPAddress: "127.0.0.1"}},
	}
	for name, overrides := range tests {
		t.Run(name, func(t *testing.T) {
			manager := NewManager(nil, "http://127.0.0.1:18080")
			manager.profileCache = rendererProfileCacheState{loaded: true, overrides: overrides}
			request := httptest.NewRequest(http.MethodGet, "/dlna/rootDesc.xml", nil)
			request.RemoteAddr = "127.0.0.1:1234"
			response := httptest.NewRecorder()
			manager.Handler().ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
			}
		})
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

func TestDeniedRequestIsExplicitInDiagnostics(t *testing.T) {
	manager, _ := deniedTestManager(t)
	request := soapRequest(
		"/dlna/control/content-directory",
		"urn:schemas-upnp-org:service:ContentDirectory:1#Browse",
		"<u:Browse xmlns:u=\"urn:schemas-upnp-org:service:ContentDirectory:1\"/>",
	)
	request.RemoteAddr = "127.0.0.1:1234"
	manager.Handler().ServeHTTP(httptest.NewRecorder(), request)
	clients := manager.Status().RecentClients
	if len(clients) != 1 || clients[0].LastError == nil || *clients[0].LastError != errRendererOverrideDenied.Error() {
		t.Fatalf("clients = %#v", clients)
	}
	if clients[0].LastSOAPAction != "Browse" {
		t.Fatalf("LastSOAPAction = %q", clients[0].LastSOAPAction)
	}
}

func TestRendererOverridesFromStorageMapsAllowedAndIdentities(t *testing.T) {
	uuidValue, ip := "uuid:renderer", "192.0.2.10"
	got := rendererOverridesFromStorage([]storage.DLNARendererDeviceOverride{{
		RendererUUID: &uuidValue,
		IPAddress:    &ip,
		ProfileID:    "generic",
		Allowed:      false,
	}})
	if len(got) != 1 || got[0].RendererUUID != uuidValue || got[0].IPAddress != ip || got[0].ProfileID != "generic" || got[0].Allowed {
		t.Fatalf("overrides = %#v", got)
	}
}

type countingContentSource struct {
	items []storage.MediaItem
	calls *atomic.Int32
}

func (s countingContentSource) ListMediaItems(context.Context) ([]storage.MediaItem, error) {
	s.calls.Add(1)
	return s.items, nil
}

func deniedTestManager(t *testing.T) (*Manager, *atomic.Int32) {
	t.Helper()
	calls := &atomic.Int32{}
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.source = countingContentSource{calls: calls}
	manager.profileCache = rendererProfileCacheState{loaded: true, overrides: []rendererDeviceOverride{{Allowed: false, IPAddress: "127.0.0.1"}}}
	manager.rateLimiter = nil
	return manager, calls
}

func installDenyMarkerTools(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	marker := filepath.Join(dir, "started")
	payload := []byte("#!/bin/sh\ntouch \"" + marker + "\"\nexit 1\n")
	for _, name := range []string{"ffmpeg", "ffprobe"} {
		if err := os.WriteFile(filepath.Join(dir, name), payload, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return marker
}
