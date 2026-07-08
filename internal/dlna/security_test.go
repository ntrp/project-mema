package dlna

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/storage"
	"media-manager/internal/testdb"
)

func TestHandlerRejectsRequestsOutsideAllowedCIDRs(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	request := httptest.NewRequest(http.MethodGet, "/dlna/rootDesc.xml", nil)
	request.RemoteAddr = "203.0.113.10:1234"
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
}

func TestActiveStreamLimitRejectsBeforeServing(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	manager.mu.Lock()
	for i := 0; i < maxActiveStreams; i++ {
		id := string(rune('a' + i))
		manager.activeStreams[id] = StreamStatus{ID: id}
	}
	manager.mu.Unlock()
	request := httptest.NewRequest(http.MethodGet, "/dlna/resource/"+url.PathEscape(resourceID), nil)
	request.RemoteAddr = "127.0.0.1:1234"
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
}

func TestHandlerRateLimitRejectsNoisyClient(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.rateLimiter = newDLNARateLimiter(1, time.Minute)

	first := httptest.NewRequest(http.MethodGet, "/dlna/rootDesc.xml", nil)
	first.RemoteAddr = "127.0.0.1:1234"
	manager.Handler().ServeHTTP(httptest.NewRecorder(), first)

	second := httptest.NewRequest(http.MethodGet, "/dlna/rootDesc.xml", nil)
	second.RemoteAddr = "127.0.0.1:1234"
	response := httptest.NewRecorder()
	manager.Handler().ServeHTTP(response, second)

	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
}

func TestStopCancelsAndClearsStreamDiagnostics(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	done, ok := manager.beginStream(allowedTestRequest(), resourceID, "direct", true)
	if !ok {
		t.Fatal("beginStream rejected test stream")
	}
	defer done()

	if err := manager.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	status := manager.Status()
	if len(status.ActiveStreams) != 0 || len(status.ActiveTranscodes) != 0 {
		t.Fatalf("status = %#v", status)
	}
}

func TestStreamDiagnosticsDoNotExposeAbsolutePaths(t *testing.T) {
	manager, resourceID := resourceTestManager(t, "0123456789")
	done, ok := manager.beginStream(allowedTestRequest(), resourceID, "direct", false)
	if !ok {
		t.Fatal("beginStream rejected test stream")
	}
	defer done()

	status := manager.Status()

	if len(status.ActiveStreams) != 1 || status.ActiveStreams[0].Path != resourceID {
		t.Fatalf("active streams = %#v", status.ActiveStreams)
	}
	if strings.Contains(status.ActiveStreams[0].Path, "/") {
		t.Fatalf("stream path leaks filesystem path: %#v", status.ActiveStreams[0])
	}
}

func TestSOAPAuditEventRecordsActionWithoutPathLeak(t *testing.T) {
	ctx, store := dlnaTestStore(t)
	manager := NewManager(store, "http://127.0.0.1:18080")
	body := `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><u:GetSystemUpdateID xmlns:u="urn:schemas-upnp-org:service:ContentDirectory:1"/></s:Body></s:Envelope>`
	request := httptest.NewRequest(http.MethodPost, "/dlna/control/content-directory", strings.NewReader(body))
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("SOAPACTION", `"urn:schemas-upnp-org:service:ContentDirectory:1#GetSystemUpdateID"`)
	response := httptest.NewRecorder()

	manager.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	events, err := store.ListSystemEvents(ctx, 10, nil)
	if err != nil {
		t.Fatalf("list events: %v", err)
	}
	if len(events) == 0 || events[0].Category != "dlna" || events[0].Data["action"] != "GetSystemUpdateID" {
		t.Fatalf("events = %#v", events)
	}
	if _, ok := events[0].Data["path"]; ok {
		t.Fatalf("audit event exposes path: %#v", events[0].Data)
	}
}

func allowedTestRequest() *http.Request {
	request := httptest.NewRequest(http.MethodGet, "/dlna/resource/test", nil)
	request.RemoteAddr = "127.0.0.1:1234"
	return request
}

func dlnaTestStore(t *testing.T) (context.Context, *storage.SettingsStore) {
	t.Helper()
	databaseURL := testdb.Create(t)
	ctx := context.Background()
	if err := storage.EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return ctx, storage.NewSettingsStore(pool)
}
