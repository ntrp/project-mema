package dlna

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"time"

	"media-manager/internal/dlna/soap"
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
	}
	m.recentClients[status.IP] = status
}

func (m *Manager) beginStream(r *http.Request, path string, transcode bool) func() {
	request := RendererRequestFromHTTP(r)
	profile := m.RendererProfile(request)

	m.mu.Lock()
	m.nextStreamID++
	stream := StreamStatus{
		ID:        strconv.Itoa(m.nextStreamID),
		ClientIP:  request.ClientIP,
		Path:      path,
		ProfileID: profile.ID,
		StartedAt: time.Now().UTC(),
	}
	m.activeStreams[stream.ID] = stream
	if transcode {
		m.activeTranscodes[stream.ID] = stream
	}
	m.mu.Unlock()

	return func() {
		m.mu.Lock()
		delete(m.activeStreams, stream.ID)
		delete(m.activeTranscodes, stream.ID)
		m.mu.Unlock()
	}
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
