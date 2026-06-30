package httpapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

const (
	eventSeverityInfo    = "info"
	eventSeverityWarning = "warning"
	eventSeverityError   = "error"
)

func (s *Server) ListSystemEvents(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	events, err := s.settings.ListSystemEvents(r.Context(), 200)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_event_list_failed", "Could not list system events")
		return
	}
	response := SystemEventListResponse{Events: make([]SystemEvent, 0, len(events))}
	for _, event := range events {
		response.Events = append(response.Events, systemEventResponse(event))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) DeleteSystemEvent(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	eventID := uuid.UUID(id)
	if err := s.settings.DeleteSystemEvent(r.Context(), eventID); err != nil {
		writeSettingsError(w, err, "Could not delete system event")
		return
	}
	s.events.Publish("system.event.deleted", map[string]any{"id": eventID.String()})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetSystemEventSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	settings, err := s.settings.GetSystemEventSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_event_settings_load_failed", "Could not load event settings")
		return
	}
	writeJSON(w, http.StatusOK, systemEventSettingsResponse(settings))
}

func (s *Server) UpdateSystemEventSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body SystemEventSettingsRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	settings, err := s.settings.UpdateSystemEventSettings(r.Context(), storage.SystemEventSettingsInput{
		RetentionDays: body.RetentionDays,
	})
	if err != nil {
		writeSettingsError(w, err, "Could not update event settings")
		return
	}
	writeJSON(w, http.StatusOK, systemEventSettingsResponse(settings))
}

func (s *Server) recordEvent(ctx context.Context, severity, category, message string, data map[string]any) {
	event, err := s.settings.CreateSystemEvent(ctx, storage.SystemEventInput{
		Severity: severity,
		Category: category,
		Message:  message,
		Data:     data,
	})
	if err != nil {
		slog.Error("system event record failed", "severity", severity, "category", category, "message", message, "error", err)
		return
	}
	s.events.Publish("system.event.created", systemEventResponse(event))
}

func systemEventResponse(event storage.SystemEvent) SystemEvent {
	return SystemEvent{
		Id:        openapi_types.UUID(event.ID),
		Severity:  SystemEventSeverity(event.Severity),
		Category:  event.Category,
		Message:   event.Message,
		Data:      event.Data,
		CreatedAt: event.CreatedAt,
	}
}

func systemEventSettingsResponse(settings storage.SystemEventSettings) SystemEventSettings {
	return SystemEventSettings{
		RetentionDays: settings.RetentionDays,
	}
}
