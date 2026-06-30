package jobs

import (
	"context"
	"log/slog"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

const (
	jobEventInfo    = "info"
	jobEventWarning = "warning"
	jobEventError   = "error"
)

func publishSystemEvent(
	ctx context.Context,
	store *storage.SettingsStore,
	broker *events.Broker,
	severity string,
	category string,
	message string,
	data map[string]any,
) {
	if store == nil {
		return
	}
	event, err := store.CreateSystemEvent(ctx, storage.SystemEventInput{
		Severity: severity,
		Category: category,
		Message:  message,
		Data:     data,
	})
	if err != nil {
		slog.Error("job system event record failed", "severity", severity, "category", category, "message", message, "error", err)
		return
	}
	if broker != nil {
		broker.Publish("system.event.created", map[string]any{
			"id":        event.ID,
			"severity":  event.Severity,
			"category":  event.Category,
			"message":   event.Message,
			"data":      event.Data,
			"createdAt": event.CreatedAt,
		})
	}
}
