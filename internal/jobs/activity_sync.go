package jobs

import (
	"context"
	"fmt"
	"strings"

	"github.com/riverqueue/river"

	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/imports"
	"media-manager/internal/storage"
)

type DownloadActivitySyncArgs struct{}

func (DownloadActivitySyncArgs) Kind() string {
	return "download.activity_sync"
}

type DownloadActivitySyncWorker struct {
	river.WorkerDefaults[DownloadActivitySyncArgs]

	settings        *storage.SettingsStore
	downloadClients *downloadclients.Service
	imports         *imports.Service
	events          *events.Broker
}

func (w *DownloadActivitySyncWorker) Work(ctx context.Context, _ *river.Job[DownloadActivitySyncArgs]) error {
	activities, err := w.settings.ListActiveDownloadActivity(ctx)
	if err != nil {
		return fmt.Errorf("list active download activity: %w", err)
	}
	clients, err := w.settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		return fmt.Errorf("list enabled download clients: %w", err)
	}
	clientsByName := map[string]storage.DownloadClient{}
	for _, client := range clients {
		clientsByName[client.Name] = client
	}

	var failures []string
	for _, activity := range activities {
		client, ok := clientsByName[activity.DownloadClientName]
		if !ok || activity.DownloadID == nil {
			continue
		}
		if err := w.syncActivity(ctx, activity, client); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %s", activity.ReleaseTitle, err.Error()))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("download activity sync failed for %d item(s): %s", len(failures), strings.Join(failures, "; "))
	}
	return nil
}

func (w *DownloadActivitySyncWorker) syncActivity(ctx context.Context, activity storage.DownloadActivity, client storage.DownloadClient) error {
	result := w.downloadClients.Status(ctx, downloadClientConfig(client), downloadclients.StatusRequest{
		DownloadID: *activity.DownloadID,
	})
	if !result.Success {
		message := strings.TrimSpace(result.Message)
		if message == "" {
			message = "Could not fetch download status"
		}
		updated, err := w.settings.UpdateDownloadActivityStatus(ctx, activity.ID, "failed", &message)
		if err == nil {
			w.publishActivity(updated, activity)
		}
		return err
	}
	if !result.Found {
		return nil
	}
	if result.Status == "completed" {
		if err := w.imports.ImportCompletedDownload(ctx, activity, result.Files); err != nil {
			return w.failActivity(ctx, activity, err.Error())
		}
	}
	message := (*string)(nil)
	if result.Status == "failed" {
		trimmed := strings.TrimSpace(result.Message)
		if trimmed != "" {
			message = &trimmed
		}
	}
	updated, err := w.settings.UpdateDownloadActivityStatus(ctx, activity.ID, result.Status, message)
	if err != nil {
		return err
	}
	w.publishActivity(updated, activity)
	return nil
}

func (w *DownloadActivitySyncWorker) failActivity(ctx context.Context, activity storage.DownloadActivity, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "Download import failed"
	}
	updated, err := w.settings.UpdateDownloadActivityStatus(ctx, activity.ID, "failed", &message)
	if err != nil {
		return err
	}
	w.publishActivity(updated, activity)
	return nil
}

func (w *DownloadActivitySyncWorker) publishActivity(updated storage.DownloadActivity, previous storage.DownloadActivity) {
	if updated.Status == previous.Status && sameStringPtr(updated.Error, previous.Error) {
		return
	}
	updated.MediaTitle = previous.MediaTitle
	updated.MediaType = previous.MediaType
	publishDownloadActivity(w.events, updated)
}

func sameStringPtr(left *string, right *string) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}

func publishDownloadActivity(broker *events.Broker, activity storage.DownloadActivity) {
	broker.Publish("activity.download.updated", map[string]interface{}{
		"id":                 activity.ID,
		"mediaItemId":        activity.MediaItemID,
		"mediaTitle":         activity.MediaTitle,
		"mediaType":          activity.MediaType,
		"releaseTitle":       activity.ReleaseTitle,
		"indexerName":        activity.IndexerName,
		"downloadClientName": activity.DownloadClientName,
		"downloadId":         activity.DownloadID,
		"downloadUrl":        activity.DownloadURL,
		"status":             activity.Status,
		"error":              activity.Error,
		"createdAt":          activity.CreatedAt,
		"updatedAt":          activity.UpdatedAt,
	})
}
