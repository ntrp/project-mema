package httpapi

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"media-manager/internal/logging"
	"media-manager/internal/storage"
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

func (s *Server) GetSystemStatus(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	database, err := s.settings.GetDatabaseStatus(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_status_failed", "Could not load system status")
		return
	}
	writeJSON(w, http.StatusOK, SystemStatusResponse{
		Version:         s.cfg.Version,
		Commit:          s.cfg.Commit,
		DatabaseType:    database.Type,
		DatabaseVersion: database.Version,
		License:         "AGPL-3.0-or-later",
		SourceLocation:  s.cfg.SourceURL,
	})
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

func (s *Server) GetSystemLogFileSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	settings, err := s.settings.GetLogFileSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "log_file_settings_failed", "Could not load log file settings")
		return
	}
	writeJSON(w, http.StatusOK, systemLogFileSettingsResponse(settings))
}

func (s *Server) UpdateSystemLogFileSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body SystemLogFileSettingsRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	settings, err := s.settings.UpdateLogFileSettings(r.Context(), storage.LogFileSettingsInput{
		Enabled:       body.Enabled,
		Directory:     body.Directory,
		RetentionDays: body.RetentionDays,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_log_file_settings", "Log file settings are invalid")
		return
	}
	if err := logging.Default.ConfigureFile(logging.FileSettings{
		Enabled:       settings.Enabled,
		Directory:     settings.Directory,
		RetentionDays: int(settings.RetentionDays),
	}); err != nil {
		writeError(w, http.StatusBadRequest, "log_file_setup_failed", "Log file directory could not be used")
		return
	}
	writeJSON(w, http.StatusOK, systemLogFileSettingsResponse(settings))
}

func (s *Server) ListSystemLogFiles(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	settings, err := s.settings.GetLogFileSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "log_file_settings_failed", "Could not load log file settings")
		return
	}
	files, err := logging.ListFiles(settings.Directory)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "log_file_list_failed", "Could not list log files")
		return
	}
	response := SystemLogFileListResponse{Files: make([]SystemLogFile, 0, len(files))}
	for _, file := range files {
		response.Files = append(response.Files, SystemLogFile{
			Name:       file.Name,
			SizeBytes:  file.SizeBytes,
			ModifiedAt: file.ModifiedAt,
		})
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) DownloadSystemLogFile(w http.ResponseWriter, r *http.Request, name string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	settings, err := s.settings.GetLogFileSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "log_file_settings_failed", "Could not load log file settings")
		return
	}
	path, ok := logging.FilePath(settings.Directory, name)
	if !ok {
		writeError(w, http.StatusNotFound, "log_file_not_found", "Log file was not found")
		return
	}
	if _, err := os.Stat(path); err != nil {
		writeError(w, http.StatusNotFound, "log_file_not_found", "Log file was not found")
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(name)))
	http.ServeFile(w, r, path)
}

func systemLogLevel(level logging.Level) SystemLogLevel {
	return SystemLogLevel(level)
}

func systemLogFileSettingsResponse(settings storage.LogFileSettings) SystemLogFileSettings {
	effectiveDirectory, err := filepath.Abs(settings.Directory)
	if err != nil {
		effectiveDirectory = settings.Directory
	}
	return SystemLogFileSettings{
		Enabled:            settings.Enabled,
		Directory:          settings.Directory,
		EffectiveDirectory: effectiveDirectory,
		RetentionDays:      settings.RetentionDays,
	}
}
