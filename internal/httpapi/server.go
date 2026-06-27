package httpapi

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/config"
	"media-manager/internal/storage"
	"media-manager/internal/tools"
)

const sessionCookieName = "session"

type Server struct {
	cfg      config.Config
	settings *storage.SettingsStore
	sessions *sessionStore
	now      func() time.Time
}

func NewServer(cfg config.Config, settings *storage.SettingsStore) *Server {
	return &Server{
		cfg:      cfg,
		settings: settings,
		sessions: newSessionStore(),
		now:      time.Now,
	}
}

func (s *Server) GetHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{
		Status:  Ok,
		Version: s.cfg.Version,
		Commit:  s.cfg.Commit,
		Time:    s.now().UTC(),
	})
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16*1024)).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid login request")
		return
	}

	if !sameString(body.Username, s.cfg.AdminUsername) || !sameString(body.Password, s.cfg.AdminPassword) {
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
		Id:       openapi_types.UUID(uuid.New()),
		Username: body.Username,
		Role:     Admin,
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
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}

	writeJSON(w, http.StatusOK, SessionResponse{
		Authenticated: true,
		ExpiresAt:     &session.expiresAt,
		User:          &session.user,
	})
}

func (s *Server) ListDownloadClients(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	clients, err := s.settings.ListDownloadClients(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list download clients")
		return
	}

	response := DownloadClientListResponse{Clients: make([]DownloadClient, 0, len(clients))}
	for _, client := range clients {
		response.Clients = append(response.Clients, downloadClientResponse(client))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateDownloadClient(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body DownloadClientRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := downloadClientInput(w, body)
	if !ok {
		return
	}

	client, err := s.settings.CreateDownloadClient(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_create_failed", "Could not create download client")
		return
	}
	writeJSON(w, http.StatusCreated, downloadClientResponse(client))
}

func (s *Server) UpdateDownloadClient(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body DownloadClientRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := downloadClientInput(w, body)
	if !ok {
		return
	}

	client, err := s.settings.UpdateDownloadClient(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update download client")
		return
	}
	writeJSON(w, http.StatusOK, downloadClientResponse(client))
}

func (s *Server) DeleteDownloadClient(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	if err := s.settings.DeleteDownloadClient(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete download client")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ListIndexers(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	indexers, err := s.settings.ListIndexers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list indexers")
		return
	}

	response := IndexerListResponse{Indexers: make([]Indexer, 0, len(indexers))}
	for _, indexer := range indexers {
		response.Indexers = append(response.Indexers, indexerResponse(indexer))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateIndexer(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body IndexerRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := indexerInput(w, body)
	if !ok {
		return
	}

	indexer, err := s.settings.CreateIndexer(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_create_failed", "Could not create indexer")
		return
	}
	writeJSON(w, http.StatusCreated, indexerResponse(indexer))
}

func (s *Server) UpdateIndexer(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body IndexerRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := indexerInput(w, body)
	if !ok {
		return
	}

	indexer, err := s.settings.UpdateIndexer(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update indexer")
		return
	}
	writeJSON(w, http.StatusOK, indexerResponse(indexer))
}

func (s *Server) DeleteIndexer(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	if err := s.settings.DeleteIndexer(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete indexer")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetToolStatus(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	detected := tools.Detect(r.Context(), tools.DefaultTools)
	response := ToolStatusResponse{Tools: make([]ToolStatus, 0, len(detected))}
	for _, tool := range detected {
		item := ToolStatus{
			Name:      ToolName(tool.Name),
			Required:  tool.Required,
			Available: tool.Available,
		}
		if tool.Version != "" {
			item.Version = &tool.Version
		}
		if tool.Path != "" {
			item.Path = &tool.Path
		}
		if tool.Error != "" {
			item.Error = &tool.Error
		}
		response.Tools = append(response.Tools, item)
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) StreamEvents(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming_unavailable", "Streaming is unavailable")
		return
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	writeSSE(w, flusher, "system.heartbeat", map[string]interface{}{"status": "ok"})
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			writeSSE(w, flusher, "system.heartbeat", map[string]interface{}{"status": "ok"})
		}
	}
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

func sameString(left, right string) bool {
	if len(left) != len(right) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(left), []byte(right)) == 1
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64*1024)).Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "Invalid request body")
		return false
	}
	return true
}

func downloadClientInput(w http.ResponseWriter, request DownloadClientRequest) (storage.DownloadClientInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.DownloadClientInput{}, false
	}
	if !request.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Download client type is not supported")
		return storage.DownloadClientInput{}, false
	}
	if baseURL == "" {
		writeError(w, http.StatusBadRequest, "invalid_base_url", "Base URL is required")
		return storage.DownloadClientInput{}, false
	}
	if request.Priority < 0 || request.Priority > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_priority", "Priority must be between 0 and 1000")
		return storage.DownloadClientInput{}, false
	}

	return storage.DownloadClientInput{
		Name:     name,
		Type:     string(request.Type),
		BaseURL:  baseURL,
		Username: optionalTrimmedString(request.Username),
		Password: optionalTrimmedString(request.Password),
		APIKey:   optionalTrimmedString(request.ApiKey),
		Category: optionalTrimmedString(request.Category),
		Enabled:  request.Enabled,
		Priority: request.Priority,
	}, true
}

func indexerInput(w http.ResponseWriter, request IndexerRequest) (storage.IndexerInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.IndexerInput{}, false
	}
	if !request.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Indexer type is not supported")
		return storage.IndexerInput{}, false
	}
	if baseURL == "" {
		writeError(w, http.StatusBadRequest, "invalid_base_url", "Base URL is required")
		return storage.IndexerInput{}, false
	}
	if request.Priority < 0 || request.Priority > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_priority", "Priority must be between 0 and 1000")
		return storage.IndexerInput{}, false
	}

	categories := []int32{}
	if request.Categories != nil {
		categories = append(categories, (*request.Categories)...)
	}

	return storage.IndexerInput{
		Name:       name,
		Type:       string(request.Type),
		BaseURL:    baseURL,
		APIKey:     optionalTrimmedString(request.ApiKey),
		Categories: categories,
		Enabled:    request.Enabled,
		Priority:   request.Priority,
	}, true
}

func downloadClientResponse(client storage.DownloadClient) DownloadClient {
	return DownloadClient{
		Id:        openapi_types.UUID(client.ID),
		Name:      client.Name,
		Type:      DownloadClientType(client.Type),
		BaseUrl:   client.BaseURL,
		Username:  client.Username,
		Password:  client.Password,
		ApiKey:    client.APIKey,
		Category:  client.Category,
		Enabled:   client.Enabled,
		Priority:  client.Priority,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}

func indexerResponse(indexer storage.Indexer) Indexer {
	categories := append([]int32(nil), indexer.Categories...)
	if categories == nil {
		categories = []int32{}
	}
	return Indexer{
		Id:         openapi_types.UUID(indexer.ID),
		Name:       indexer.Name,
		Type:       IndexerType(indexer.Type),
		BaseUrl:    indexer.BaseURL,
		ApiKey:     indexer.APIKey,
		Categories: &categories,
		Enabled:    indexer.Enabled,
		Priority:   indexer.Priority,
		CreatedAt:  indexer.CreatedAt,
		UpdatedAt:  indexer.UpdatedAt,
	}
}

func optionalTrimmedString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func writeSettingsError(w http.ResponseWriter, err error, message string) {
	if errors.Is(err, storage.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", message)
		return
	}
	writeError(w, http.StatusInternalServerError, "settings_update_failed", message)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{Code: code, Message: message})
}

func writeSSE(w http.ResponseWriter, flusher http.Flusher, eventType string, data map[string]interface{}) {
	envelope := map[string]interface{}{
		"id":   uuid.NewString(),
		"type": eventType,
		"time": time.Now().UTC(),
		"data": data,
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return
	}
	_, _ = w.Write([]byte("event: " + eventType + "\n"))
	_, _ = w.Write([]byte("data: " + string(payload) + "\n\n"))
	flusher.Flush()
}
