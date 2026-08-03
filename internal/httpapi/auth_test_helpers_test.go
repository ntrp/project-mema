package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func authenticatedCookie(t *testing.T, router http.Handler) *http.Cookie {
	t.Helper()
	login := httptest.NewRecorder()
	router.ServeHTTP(login, loginRequest("admin", "admin"))
	cookies := login.Result().Cookies()
	if login.Code != http.StatusOK || len(cookies) == 0 {
		t.Fatalf("login status=%d cookies=%d", login.Code, len(cookies))
	}
	return cookies[0]
}
