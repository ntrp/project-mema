package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

type MissingMediaRetryWorker struct {
	river.WorkerDefaults[MissingMediaRetryArgs]

	settings        *storage.SettingsStore
	indexers        *indexers.Service
	downloadClients *downloadclients.Service
	decisions       decisions.Engine
	events          *events.Broker
}

func (w *MissingMediaRetryWorker) Work(ctx context.Context, job *river.Job[MissingMediaRetryArgs]) (err error) {
	return runWantedRSSSyncWorker(ctx, job.JobRow, w.settings, w.indexers, w.downloadClients, w.decisions, w.events)
}

type WantedRSSSyncWorker struct {
	river.WorkerDefaults[WantedRSSSyncArgs]

	settings        *storage.SettingsStore
	indexers        *indexers.Service
	downloadClients *downloadclients.Service
	decisions       decisions.Engine
	events          *events.Broker
}

func (w *WantedRSSSyncWorker) Work(ctx context.Context, job *river.Job[WantedRSSSyncArgs]) (err error) {
	return runWantedRSSSyncWorker(ctx, job.JobRow, w.settings, w.indexers, w.downloadClients, w.decisions, w.events)
}

func runWantedRSSSyncWorker(
	ctx context.Context,
	jobRow *rivertype.JobRow,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	downloadClientService *downloadclients.Service,
	decisionEngine decisions.Engine,
	eventBroker *events.Broker,
) (err error) {
	publishJobUpdated(eventBroker, jobRow, "running")
	defer func() { publishJobFinished(eventBroker, jobRow, err) }()

	items, err := settings.ListMissingMediaItems(ctx)
	if err != nil {
		slog.Error("wanted rss sync list failed", "error", err)
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "jobs", "Wanted RSS sync failed to list items", map[string]any{"error": err.Error()})
		return fmt.Errorf("list missing media: %w", err)
	}
	slog.Debug("wanted rss sync started", "itemCount", len(items))
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "jobs", "Wanted RSS sync started", map[string]any{"itemCount": len(items)})
	var failures []string
	for _, item := range items {
		err := autoSearchDownload(ctx, settings, indexerService, downloadClientService, decisionEngine, eventBroker, item)
		if err != nil {
			slog.Error("wanted rss sync item failed", "mediaItemId", item.ID, "title", item.Title, "error", err)
			failures = append(failures, fmt.Sprintf("%s: %s", item.Title, err.Error()))
		}
	}
	if len(failures) > 0 {
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "jobs", "Wanted RSS sync finished with failures", map[string]any{"failureCount": len(failures)})
		return fmt.Errorf("wanted rss sync failed for %d item(s): %s", len(failures), strings.Join(failures, "; "))
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "jobs", "Wanted RSS sync finished", map[string]any{"itemCount": len(items)})
	return nil
}
