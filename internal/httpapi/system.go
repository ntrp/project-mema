package httpapi

import (
	"net/http"

	"media-manager/internal/logging"
	"media-manager/internal/tools"
)

func (s *Server) GetHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{
		Status:  Ok,
		Version: s.cfg.Version,
		Commit:  s.cfg.Commit,
		Time:    s.now().UTC(),
	})
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

func (s *Server) StreamSystemLogs(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
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

	entries, unsubscribe := logging.Default.Subscribe()
	defer unsubscribe()

	writeSSE(w, flusher, "system.log.level", map[string]interface{}{
		"level": logging.Default.Level(),
	})
	for {
		select {
		case <-r.Context().Done():
			return
		case entry, ok := <-entries:
			if !ok {
				return
			}
			writeSSE(w, flusher, "system.log", entry)
		}
	}
}

func (s *Server) GetSystemLogLevel(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	writeJSON(w, http.StatusOK, SystemLogLevelResponse{
		Level: systemLogLevel(logging.Default.Level()),
	})
}

func (s *Server) UpdateSystemLogLevel(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body SystemLogLevelRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	level := logging.Level(body.Level)
	if err := logging.Default.SetLevel(level); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_log_level", "Log level is not supported")
		return
	}

	writeJSON(w, http.StatusOK, SystemLogLevelResponse{
		Level: systemLogLevel(logging.Default.Level()),
	})
}

func systemLogLevel(level logging.Level) SystemLogLevel {
	return SystemLogLevel(level)
}
