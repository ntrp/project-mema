package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/acceptance"
	"media-manager/internal/config"
	"media-manager/internal/storage"
	"media-manager/internal/testdb"
)

func TestScenarioSCNAuth001AnonymousSessionIsUnauthenticated(t *testing.T) {
	requireAcceptanceScenario(t, "SCN-AUTH-001", "api")
	router := chi.NewRouter()
	HandlerFromMux(NewServer(config.Config{}, nil, nil, nil, nil, nil, nil), router)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/auth/session", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", response.Code, response.Body.String())
	}
	var body SessionResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Authenticated {
		t.Fatal("session should be unauthenticated")
	}
}

func TestScenarioSCNAuth002AdminCanSignIn(t *testing.T) {
	requireAcceptanceScenario(t, "SCN-AUTH-002", "api")
	router := authRouter(t)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, loginRequest("admin", "admin"))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", response.Code, response.Body.String())
	}
	if response.Result().Cookies()[0].Name != sessionCookieName {
		t.Fatalf("cookies = %#v", response.Result().Cookies())
	}
	var body SessionResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !body.Authenticated || body.User == nil || body.User.Username != "admin" {
		t.Fatalf("session response = %#v", body)
	}
}

func TestScenarioSCNAuth002AdminCanSignOut(t *testing.T) {
	requireAcceptanceScenario(t, "SCN-AUTH-002", "api")
	router := authRouter(t)

	login := httptest.NewRecorder()
	router.ServeHTTP(login, loginRequest("admin", "admin"))
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %q", login.Code, login.Body.String())
	}
	cookies := login.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("login response did not include a session cookie")
	}

	logoutRequest := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	logoutRequest.AddCookie(cookies[0])
	logout := httptest.NewRecorder()
	router.ServeHTTP(logout, logoutRequest)
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout status = %d, body = %q", logout.Code, logout.Body.String())
	}
	if !hasExpiredSessionCookie(logout.Result().Cookies()) {
		t.Fatalf("logout cookies = %#v", logout.Result().Cookies())
	}

	sessionRequest := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	sessionRequest.AddCookie(cookies[0])
	sessionResponse := httptest.NewRecorder()
	router.ServeHTTP(sessionResponse, sessionRequest)
	if sessionResponse.Code != http.StatusOK {
		t.Fatalf("session status = %d, body = %q", sessionResponse.Code, sessionResponse.Body.String())
	}
	var body SessionResponse
	if err := json.Unmarshal(sessionResponse.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Authenticated {
		t.Fatalf("session after logout = %#v", body)
	}
}

func TestScenarioSCNAuth003InvalidCredentialsAreRejected(t *testing.T) {
	requireAcceptanceScenario(t, "SCN-AUTH-003", "api")
	router := authRouter(t)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, loginRequest("admin", "wrong-password"))

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, body = %q", response.Code, response.Body.String())
	}
	if len(response.Result().Cookies()) != 0 {
		t.Fatalf("unexpected cookies = %#v", response.Result().Cookies())
	}
}

func TestScenarioSCNAuth001ExpiredSessionIsCleared(t *testing.T) {
	requireAcceptanceScenario(t, "SCN-AUTH-001", "api")
	router := chi.NewRouter()
	server := NewServer(config.Config{AppEnv: "development", SessionTTL: time.Hour}, nil, nil, nil, nil, nil, nil)
	expiredID := "expired-session"
	server.sessions.put(expiredID, session{
		user:      UserSummary{Username: "admin", Role: Admin},
		expiresAt: time.Now().Add(-time.Hour),
	})
	HandlerFromMux(server, router)

	request := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: expiredID})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", response.Code, response.Body.String())
	}
	var body SessionResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Authenticated {
		t.Fatalf("expired session response = %#v", body)
	}
	if !hasExpiredSessionCookie(response.Result().Cookies()) {
		t.Fatalf("expired-session cookies = %#v", response.Result().Cookies())
	}
}

func requireAcceptanceScenario(t *testing.T, id string, tag string) {
	t.Helper()
	scenario, err := acceptance.RequireScenario("features/behavior", id)
	if err != nil {
		t.Fatal(err)
	}
	if !scenario.HasTag(tag) {
		t.Fatalf("%s missing @%s tag", id, tag)
	}
}

func authRouter(t *testing.T) http.Handler {
	t.Helper()
	router := chi.NewRouter()
	HandlerFromMux(NewServer(testConfig(), testSettingsStore(t), nil, nil, nil, nil, nil), router)
	return router
}

func testConfig() config.Config {
	return config.Config{AppEnv: "development", SessionTTL: 24 * time.Hour}
}

func testSettingsStore(t *testing.T) *storage.SettingsStore {
	t.Helper()
	databaseURL := testdb.Create(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := storage.EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	store := storage.NewSettingsStore(pool)
	if err := store.EnsureDefaultAdminUser(ctx, "admin", "admin"); err != nil {
		t.Fatal(err)
	}
	return store
}

func loginRequest(username string, password string) *http.Request {
	body, _ := json.Marshal(LoginRequest{Username: username, Password: password})
	return httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
}

func hasExpiredSessionCookie(cookies []*http.Cookie) bool {
	for _, cookie := range cookies {
		if cookie.Name == sessionCookieName && cookie.MaxAge < 0 {
			return true
		}
	}
	return false
}
