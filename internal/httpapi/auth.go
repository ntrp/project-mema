package httpapi

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

const sessionCookieName = "session"

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16*1024)).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid login request")
		return
	}

	userRecord, err := s.settings.GetUserByUsername(r.Context(), body.Username)
	if err != nil || !storage.VerifyPassword(body.Password, userRecord.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Invalid username or password")
		return
	}

	expiresAt := s.now().Add(s.cfg.SessionTTL).UTC()
	sessionID, err := newSessionID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "session_create_failed", "Could not create session")
		return
	}

	user := UserSummary{
		Id:          openapi_types.UUID(userRecord.ID),
		Username:    userRecord.Username,
		DisplayName: optionalString(userRecord.DisplayName),
		PictureUrl:  optionalString(userRecord.PictureURL),
		Role:        UserRole(userRecord.Role),
	}
	s.sessions.put(sessionID, session{user: user, expiresAt: expiresAt})
	http.SetCookie(w, s.sessionCookie(sessionID, expiresAt))

	writeJSON(w, http.StatusOK, SessionResponse{
		Authenticated: true,
		ExpiresAt:     &expiresAt,
		User:          &user,
	})
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		s.sessions.delete(cookie.Value)
	}

	http.SetCookie(w, s.expiredSessionCookie())
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeJSON(w, http.StatusOK, SessionResponse{Authenticated: false})
		return
	}

	session, ok := s.sessions.get(cookie.Value, s.now())
	if !ok {
		http.SetCookie(w, s.expiredSessionCookie())
		writeJSON(w, http.StatusOK, SessionResponse{Authenticated: false})
		return
	}
	writeJSON(w, http.StatusOK, SessionResponse{
		Authenticated: true,
		ExpiresAt:     &session.expiresAt,
		User:          &session.user,
	})
}

func (s *Server) requireSession(w http.ResponseWriter, r *http.Request) (session, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return session{}, false
	}

	currentSession, ok := s.sessions.get(cookie.Value, s.now())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Session expired or invalid")
		return session{}, false
	}

	return currentSession, true
}

func (s *Server) requireAdmin(w http.ResponseWriter, r *http.Request) (session, bool) {
	currentSession, ok := s.requireSession(w, r)
	if !ok {
		return session{}, false
	}
	if currentSession.user.Role != Admin {
		writeError(w, http.StatusForbidden, "forbidden", "Admin role required")
		return session{}, false
	}
	return currentSession, true
}

func (s *Server) sessionCookie(value string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   !s.cfg.IsDevelopment(),
	}
}

func (s *Server) expiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   !s.cfg.IsDevelopment(),
	}
}

type session struct {
	user      UserSummary
	expiresAt time.Time
}

func (s session) userID() openapi_types.UUID {
	return s.user.Id
}

type sessionStore struct {
	mu       sync.Mutex
	sessions map[string]session
}

func newSessionStore() *sessionStore {
	return &sessionStore{sessions: map[string]session{}}
}

func (s *sessionStore) put(id string, session session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = session
}

func (s *sessionStore) get(id string, now time.Time) (session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return session, false
	}
	if !session.expiresAt.After(now) {
		delete(s.sessions, id)
		return session, false
	}
	return session, true
}

func (s *sessionStore) delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
}

func newSessionID() (string, error) {
	var bytes [32]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes[:]), nil
}
