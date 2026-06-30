package jobs

import (
	"context"
	"fmt"
	"log/slog"
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
		slog.Error("download activity sync list failed", "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "jobs", "Download activity sync failed to list activity", map[string]any{"error": err.Error()})
		return fmt.Errorf("list active download activity: %w", err)
	}
	clients, err := w.settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		slog.Error("download activity sync client list failed", "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "jobs", "Download activity sync failed to list clients", map[string]any{"error": err.Error()})
		return fmt.Errorf("list enabled download clients: %w", err)
	}
	slog.Debug("download activity sync started", "activityCount", len(activities), "clientCount", len(clients))
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
			slog.Error("download activity sync item failed", "activityId", activity.ID, "mediaItemId", activity.MediaItemID, "releaseTitle", activity.ReleaseTitle, "error", err)
			failures = append(failures, fmt.Sprintf("%s: %s", activity.ReleaseTitle, err.Error()))
		}
	}
	if len(failures) > 0 {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "jobs", "Download activity sync finished with failures", map[string]any{"failureCount": len(failures)})
		return fmt.Errorf("download activity sync failed for %d item(s): %s", len(failures), strings.Join(failures, "; "))
	}
	if len(activities) > 0 {
		publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "jobs", "Download activity sync finished", map[string]any{"activityCount": len(activities), "clientCount": len(clients)})
	}
	return nil
}

func (w *DownloadActivitySyncWorker) syncActivity(ctx context.Context, activity storage.DownloadActivity, client storage.DownloadClient) error {
	slog.Debug("checking download activity status", "activityId", activity.ID, "mediaItemId", activity.MediaItemID, "downloadClientName", activity.DownloadClientName, "downloadId", activity.DownloadID)
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
			publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download activity status check failed", map[string]any{"activityId": activity.ID.String(), "message": message})
		} else {
			slog.Error("failed to mark download activity failed", "activityId", activity.ID, "error", err)
			publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download activity failure update failed", map[string]any{"activityId": activity.ID.String(), "error": err.Error()})
		}
		return err
	}
	if !result.Found {
		slog.Debug("download activity not found in client", "activityId", activity.ID, "downloadId", activity.DownloadID)
		return nil
	}
	slog.Debug("download activity status received", "activityId", activity.ID, "status", result.Status, "progressPercent", result.ProgressPercent, "fileCount", len(result.Files))
	if result.Status == "completed" {
		if err := w.imports.ImportCompletedDownload(ctx, activity, result.Files); err != nil {
			slog.Error("completed download import failed", "activityId", activity.ID, "mediaItemId", activity.MediaItemID, "error", err)
			publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Completed download import failed", map[string]any{"activityId": activity.ID.String(), "mediaItemId": activity.MediaItemID.String(), "error": err.Error()})
			return w.failActivity(ctx, activity, err.Error())
		}
		slog.Debug("completed download imported", "activityId", activity.ID, "mediaItemId", activity.MediaItemID, "fileCount", len(result.Files))
		publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "downloads", "Completed download imported", map[string]any{"activityId": activity.ID.String(), "mediaItemId": activity.MediaItemID.String(), "fileCount": len(result.Files)})
	}
	message := (*string)(nil)
	if result.Status == "failed" {
		trimmed := strings.TrimSpace(result.Message)
		if trimmed != "" {
			message = &trimmed
		}
	}
	updated, err := w.settings.UpdateDownloadActivityProgress(ctx, activity.ID, result.Status, result.ProgressPercent, message)
	if err != nil {
		slog.Error("failed to update download activity progress", "activityId", activity.ID, "status", result.Status, "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download activity progress update failed", map[string]any{"activityId": activity.ID.String(), "status": result.Status, "error": err.Error()})
		return err
	}
	if result.Status == "failed" {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download activity failed", map[string]any{"activityId": activity.ID.String(), "status": result.Status})
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
	if updated.Status == previous.Status && sameStringPtr(updated.Error, previous.Error) && sameIntPtr(updated.ProgressPercent, previous.ProgressPercent) {
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

func sameIntPtr(left *int, right *int) bool {
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
		"mediaYear":          activity.MediaYear,
		"releaseTitle":       activity.ReleaseTitle,
		"indexerName":        activity.IndexerName,
		"downloadClientName": activity.DownloadClientName,
		"downloadId":         activity.DownloadID,
		"downloadUrl":        activity.DownloadURL,
		"status":             activity.Status,
		"progressPercent":    activity.ProgressPercent,
		"error":              activity.Error,
		"createdAt":          activity.CreatedAt,
		"updatedAt":          activity.UpdatedAt,
	})
}
