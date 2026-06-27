package httpapi

import (
	"net/http"

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
