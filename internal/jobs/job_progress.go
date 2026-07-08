package jobs

import (
	"context"
	"strings"
	"time"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

func recordJobProgress(ctx context.Context, store *storage.SettingsStore, broker *events.Broker, percent *int32, label string) {
	recordJobProgressData(ctx, store, broker, percent, label, nil)
}

func recordJobProgressData(
	ctx context.Context,
	store *storage.SettingsStore,
	broker *events.Broker,
	percent *int32,
	label string,
	data map[string]any,
) {
	riverJobID, ok := jobExecutionID(ctx)
	if !ok || store == nil {
		return
	}
	progress := normalizedProgressData(percent, label, data)
	execution, err := store.UpdateSystemJobExecutionProgressData(ctx, riverJobID, percent, label, progress)
	if err != nil {
		return
	}
	_, _ = store.CreateSystemJobExecutionLog(ctx, riverJobID, "info", label, map[string]any{"progress": progress})
	publishJobExecutionUpdated(broker, execution)
}

func normalizedProgressData(percent *int32, label string, data map[string]any) map[string]any {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	progress := map[string]any{
		"label":         strings.TrimSpace(label),
		"phase":         strings.TrimSpace(label),
		"lastUpdatedAt": now,
	}
	if percent != nil {
		progress["percent"] = *percent
	}
	for key, value := range data {
		progress[key] = value
	}
	if _, ok := progress["startedAt"]; !ok {
		progress["startedAt"] = now
	}
	return progress
}
