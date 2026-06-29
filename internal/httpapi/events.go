package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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

func writeSSE(w http.ResponseWriter, flusher http.Flusher, eventType string, data interface{}) {
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
