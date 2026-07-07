package dlna

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"time"

	"media-manager/internal/dlna/soap"
	"media-manager/internal/storage"
)

type ClientStatus struct {
	IP             string
	UserAgent      string
	ProfileID      string
	LastSOAPAction string
	LastError      *string
	LastSeen       time.Time
}

type StreamStatus struct {
	ID        string
	ClientIP  string
	Path      string
	ProfileID string
	StartedAt time.Time
}

func (m *Manager) recordHTTPClient(r *http.Request) {
	m.recordClient(r, "", nil)
}

func (m *Manager) recordSOAPAction(ctx context.Context, action string, err error) {
	r, ok := soap.RequestFromContext(ctx)
	if !ok {
		return
	}
	m.recordClient(r, action, err)
}

func (m *Manager) recordClient(r *http.Request, action string, err error) {
	request := RendererRequestFromHTTP(r)
	profile := m.RendererProfile(request)
	status := ClientStatus{
		IP:        request.ClientIP,
		UserAgent: request.UserAgent,
		ProfileID: profile.ID,
		LastSeen:  time.Now().UTC(),
	}
	if action != "" {
		status.LastSOAPAction = action
	}
	if err != nil {
		message := err.Error()
		status.LastError = &message
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if existing, ok := m.recentClients[status.IP]; ok {
		if status.UserAgent == "" {
			status.UserAgent = existing.UserAgent
		}
		if status.LastSOAPAction == "" {
			status.LastSOAPAction = existing.LastSOAPAction
		}
		if status.LastError == nil {
			status.LastError = existing.LastError
		}
	}
	if action != "" {
		m.status.LastSOAPAction = &action
		m.audit(r.Context(), "DLNA SOAP action", map[string]any{
			"clientIP": status.IP,
			"profile":  status.ProfileID,
			"action":   action,
			"result":   auditResult(err),
		})
	}
	m.recentClients[status.IP] = status
}

func (m *Manager) beginStream(r *http.Request, objectID string, delivery string, transcode bool) (func(), bool) {
	request := RendererRequestFromHTTP(r)
	profile := m.RendererProfile(request)

	m.mu.Lock()
	if len(m.activeStreams) >= maxActiveStreams {
		m.mu.Unlock()
		m.audit(r.Context(), "DLNA stream rejected", map[string]any{
			"clientIP": request.ClientIP,
			"profile":  profile.ID,
			"objectID": objectID,
			"delivery": delivery,
			"result":   "limit",
		})
		return nil, false
	}
	m.nextStreamID++
	stream := StreamStatus{
		ID:        strconv.Itoa(m.nextStreamID),
		ClientIP:  request.ClientIP,
		Path:      objectID,
		ProfileID: profile.ID,
		StartedAt: time.Now().UTC(),
	}
	m.activeStreams[stream.ID] = stream
	if transcode {
		m.activeTranscodes[stream.ID] = stream
	}
	m.mu.Unlock()
	m.audit(r.Context(), "DLNA stream started", map[string]any{
		"clientIP": request.ClientIP,
		"profile":  profile.ID,
		"objectID": objectID,
		"delivery": delivery,
		"result":   "started",
	})

	return func() {
		m.mu.Lock()
		delete(m.activeStreams, stream.ID)
		delete(m.activeTranscodes, stream.ID)
		m.mu.Unlock()
		m.audit(r.Context(), "DLNA stream finished", map[string]any{
			"clientIP": request.ClientIP,
			"profile":  profile.ID,
			"objectID": objectID,
			"delivery": delivery,
			"result":   "finished",
		})
	}, true
}

func (m *Manager) audit(ctx context.Context, message string, data map[string]any) {
	if m.store == nil {
		return
	}
	_, _ = m.store.CreateSystemEvent(ctx, storage.SystemEventInput{
		Severity: "info",
		Category: "dlna",
		Message:  message,
		Data:     data,
	})
}

func auditResult(err error) string {
	if err != nil {
		return "error"
	}
	return "ok"
}

func sortedClients(values map[string]ClientStatus) []ClientStatus {
	clients := make([]ClientStatus, 0, len(values))
	for _, value := range values {
		clients = append(clients, value)
	}
	sort.Slice(clients, func(i, j int) bool {
		return clients[i].LastSeen.After(clients[j].LastSeen)
	})
	if len(clients) > 12 {
		return clients[:12]
	}
	return clients
}

func sortedStreams(values map[string]StreamStatus) []StreamStatus {
	streams := make([]StreamStatus, 0, len(values))
	for _, value := range values {
		streams = append(streams, value)
	}
	sort.Slice(streams, func(i, j int) bool {
		return streams[i].StartedAt.Before(streams[j].StartedAt)
	})
	return streams
}
