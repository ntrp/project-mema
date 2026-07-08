package dlna

import (
	"context"
	"net/http"
)

func (m *Manager) commandContext(r *http.Request) (context.Context, context.CancelFunc) {
	m.mu.Lock()
	shutdown := m.streamShutdown
	m.mu.Unlock()
	ctx, cancel := context.WithCancel(r.Context())
	if shutdown == nil {
		return ctx, cancel
	}
	go func() {
		select {
		case <-shutdown.Done():
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}
