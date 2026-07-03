package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScenarioSCNSystem006AdminInspectsRuntimeSettings(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SYSTEM-006")

	var status SystemStatusResponse
	client.doJSON(t, http.MethodGet, "/system/status", nil, http.StatusOK, &status)
	if status.License != "AGPL-3.0-or-later" || status.DatabaseType == "" {
		t.Fatalf("system status = %#v", status)
	}

	var level SystemLogLevelResponse
	client.doJSON(t, http.MethodGet, "/system/log-level", nil, http.StatusOK, &level)
	client.doJSON(t, http.MethodPut, "/system/log-level", SystemLogLevelRequest{
		Level: level.Level,
	}, http.StatusOK, &level)

	logDir := filepath.Join(t.TempDir(), "logs")
	var fileSettings SystemLogFileSettings
	client.doJSON(t, http.MethodPut, "/system/log-file-settings", SystemLogFileSettingsRequest{
		Enabled:       false,
		Directory:     logDir,
		RetentionDays: 3,
	}, http.StatusOK, &fileSettings)
	if fileSettings.Directory != logDir || fileSettings.RetentionDays != 3 {
		t.Fatalf("file settings = %#v", fileSettings)
	}

	var files SystemLogFileListResponse
	client.doJSON(t, http.MethodGet, "/system/log-files", nil, http.StatusOK, &files)
	if files.Files == nil {
		t.Fatalf("log files should default to an empty list: %#v", files)
	}

	var events SystemEventListResponse
	client.doJSON(t, http.MethodGet, "/system/events?limit=2", nil, http.StatusOK, &events)
	if events.Events == nil {
		t.Fatalf("events should default to an empty list: %#v", events)
	}

	var eventSettings SystemEventSettings
	client.doJSON(t, http.MethodGet, "/system/event-settings", nil, http.StatusOK, &eventSettings)
	eventSettings.RetentionDays = 5
	client.doJSON(t, http.MethodPut, "/system/event-settings", SystemEventSettingsRequest{
		RetentionDays: eventSettings.RetentionDays,
	}, http.StatusOK, &eventSettings)
	if eventSettings.RetentionDays != 5 {
		t.Fatalf("event settings = %#v", eventSettings)
	}
	client.doJSON(t, http.MethodDelete, "/system/events", nil, http.StatusNoContent, nil)
}

func TestScenarioSCNSystem009HealthToolsAndLogFileEdges(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SYSTEM-009")

	healthResponse := httptest.NewRecorder()
	client.router.ServeHTTP(healthResponse, httptest.NewRequest(http.MethodGet, "/health", nil))
	if healthResponse.Code != http.StatusOK {
		t.Fatalf("health status = %d, body = %q", healthResponse.Code, healthResponse.Body.String())
	}
	var health HealthResponse
	if err := json.Unmarshal(healthResponse.Body.Bytes(), &health); err != nil {
		t.Fatalf("decode health response: %v", err)
	}
	if health.Status != Ok || health.Time.IsZero() {
		t.Fatalf("health response = %#v", health)
	}

	var tools ToolStatusResponse
	client.doJSON(t, http.MethodGet, "/system/tools", nil, http.StatusOK, &tools)
	if len(tools.Tools) == 0 {
		t.Fatalf("expected tool status entries, got %#v", tools)
	}

	client.doJSON(t, http.MethodPut, "/system/log-level", SystemLogLevelRequest{
		Level: SystemLogLevel("scenario-invalid"),
	}, http.StatusBadRequest, nil)

	logDir := filepath.Join(t.TempDir(), "downloadable-logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(logDir, "app-scenario.log")
	if err := os.WriteFile(logPath, []byte("scenario log content"), 0o644); err != nil {
		t.Fatal(err)
	}
	client.doJSON(t, http.MethodPut, "/system/log-file-settings", SystemLogFileSettingsRequest{
		Enabled:       false,
		Directory:     logDir,
		RetentionDays: 7,
	}, http.StatusOK, nil)

	download := client.doRaw(t, http.MethodGet, "/system/log-files/app-scenario.log/download", true)
	if download.Code != http.StatusOK {
		t.Fatalf("download status = %d, body = %q", download.Code, download.Body.String())
	}
	if !strings.Contains(download.Header().Get("Content-Disposition"), `filename="app-scenario.log"`) {
		t.Fatalf("download disposition = %q", download.Header().Get("Content-Disposition"))
	}
	if download.Body.String() != "scenario log content" {
		t.Fatalf("download body = %q", download.Body.String())
	}

	missing := client.doRaw(t, http.MethodGet, "/system/log-files/missing.log/download", true)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing log status = %d, body = %q", missing.Code, missing.Body.String())
	}
}

func (c acceptanceClient) doRaw(t *testing.T, method string, path string, authenticated bool) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, nil)
	if authenticated {
		request.AddCookie(c.cookie)
	}
	response := httptest.NewRecorder()
	c.router.ServeHTTP(response, request)
	return response
}
