package jobs

import (
	"context"
	"log/slog"
	"time"

	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func recordIndexerSearchFailure(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	indexer storage.Indexer,
	err error,
) {
	statusCode := indexers.StatusCode(err)
	permanent := indexers.IsPermanentFailure(statusCode)
	retryUntil := retryUntilFromIndexerError(err)
	updated, updateErr := settings.RecordIndexerFailure(
		ctx,
		indexer.ID,
		statusCode,
		err.Error(),
		permanent,
		retryUntil,
	)
	if updateErr != nil {
		slog.Error("indexer failure state update failed", "indexerName", indexer.Name, "error", updateErr)
		return
	}
	if updated.HealthStatus == "disabled" {
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "indexers", "Indexer disabled", map[string]any{
			"indexerId":   indexer.ID.String(),
			"indexerName": indexer.Name,
			"statusCode":  statusCode,
			"message":     err.Error(),
		})
	}
}

func retryUntilFromIndexerError(err error) *time.Time {
	delay := indexers.RetryAfter(err)
	if delay <= 0 {
		return nil
	}
	when := time.Now().Add(delay)
	return &when
}
