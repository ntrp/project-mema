package jobs

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

type MediaRefreshArgs struct{}

func (MediaRefreshArgs) Kind() string { return "media.refresh" }

type MediaRefreshWorker struct {
	river.WorkerDefaults[MediaRefreshArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

func (w *MediaRefreshWorker) Work(ctx context.Context, job *river.Job[MediaRefreshArgs]) (err error) {
	ctx = withJobExecution(ctx, job.JobRow.ID)
	recordJobUpdated(ctx, w.settings, w.events, job.JobRow, "running")
	defer func() { recordJobFinished(ctx, w.settings, w.events, job.JobRow, err) }()
	return refreshAllMediaFiles(ctx, w.settings, w.events)
}

func refreshAllMediaFiles(ctx context.Context, settings *storage.SettingsStore, eventBroker *events.Broker) error {
	if settings == nil {
		return fmt.Errorf("settings store is not configured")
	}
	recordJobProgressData(ctx, settings, eventBroker, nil, "Loading media items", nil)
	items, err := settings.ListMediaItems(ctx)
	if err != nil {
		return err
	}
	total := len(items)
	failed := 0
	recordMediaRefreshProgress(ctx, settings, eventBroker, 0, total, failed)
	for index, item := range items {
		if _, err := settings.RescanMediaItemFiles(ctx, item.ID); err != nil {
			failed++
			publishSystemEvent(ctx, settings, eventBroker, jobEventError, "media", "Media refresh item failed", map[string]any{
				"mediaItemId": item.ID.String(),
				"title":       item.Title,
				"error":       err.Error(),
			})
		}
		recordMediaRefreshProgress(ctx, settings, eventBroker, index+1, total, failed)
	}
	if failed > 0 {
		return fmt.Errorf("media refresh failed for %d of %d media item(s)", failed, total)
	}
	return nil
}

func recordMediaRefreshProgress(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	processed int,
	total int,
	failed int,
) {
	percent := mediaRefreshProgressPercent(processed, total)
	recordJobProgressData(ctx, settings, eventBroker, &percent, fmt.Sprintf("Refreshing media %d/%d", processed, total), map[string]any{
		"processedMediaItemCount": processed,
		"mediaItemCount":          total,
		"failedMediaItemCount":    failed,
	})
}

func mediaRefreshProgressPercent(processed int, total int) int32 {
	if total <= 0 {
		return 100
	}
	if processed <= 0 {
		return 0
	}
	if processed >= total {
		return 100
	}
	return int32(processed * 100 / total)
}
