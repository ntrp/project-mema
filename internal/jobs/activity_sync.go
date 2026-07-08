package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/riverqueue/river"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/imports"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

type DownloadActivitySyncArgs struct{}

func (DownloadActivitySyncArgs) Kind() string {
	return "download.activity_sync"
}

type DownloadActivitySyncWorker struct {
	river.WorkerDefaults[DownloadActivitySyncArgs]

	settings        *storage.SettingsStore
	indexers        *indexers.Service
	downloadClients *downloadclients.Service
	decisions       decisions.Engine
	imports         *imports.Service
	events          *events.Broker
}

func (w *DownloadActivitySyncWorker) Work(ctx context.Context, job *river.Job[DownloadActivitySyncArgs]) (err error) {
	ctx = withJobExecution(ctx, job.JobRow.ID)
	recordJobUpdated(ctx, w.settings, w.events, job.JobRow, "running")
	defer func() { recordJobFinished(ctx, w.settings, w.events, job.JobRow, err) }()

	activities, err := w.settings.ListActiveDownloadActivity(ctx)
	if err != nil {
		slog.Error("download activity sync list failed", "error", err)
		return fmt.Errorf("list active download activity: %w", err)
	}
	clients, err := w.settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		slog.Error("download activity sync client list failed", "error", err)
		return fmt.Errorf("list enabled download clients: %w", err)
	}
	slog.Debug("download activity sync started", "activityCount", len(activities), "clientCount", len(clients))
	recordJobProgressData(ctx, w.settings, w.events, nil, fmt.Sprintf("Checking %d download activity item(s)", len(activities)), map[string]any{
		"phase":     "download_activity_sync",
		"unitTotal": len(activities),
	})
	clientsByName := map[string]storage.DownloadClient{}
	for _, client := range clients {
		clientsByName[client.Name] = client
	}

	var failures []string
	for index, activity := range activities {
		recordJobProgressData(ctx, w.settings, w.events, nil, fmt.Sprintf("Checking download activity %d of %d", index+1, len(activities)), map[string]any{
			"mediaItemId":      activity.MediaItemID.String(),
			"mediaTitle":       activity.MediaTitle,
			"phase":            "download_activity_status",
			"unitCurrent":      index + 1,
			"unitTotal":        len(activities),
			"pendingOperation": activity.DownloadClientName,
		})
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
		return fmt.Errorf("download activity sync failed for %d item(s): %s", len(failures), strings.Join(failures, "; "))
	}
	done := int32(100)
	recordJobProgressData(ctx, w.settings, w.events, &done, "Download activity sync finished", map[string]any{
		"phase":       "download_activity_sync",
		"unitCurrent": len(activities),
		"unitTotal":   len(activities),
	})
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
		updated, err := w.settings.FailDownloadActivity(ctx, activity.ID, &message, "download")
		if err == nil {
			blocked := w.blockActivity(ctx, activity, message, "download_status_unavailable")
			w.publishActivity(updated, activity)
			publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download activity status check failed", map[string]any{"activityId": activity.ID.String(), "message": message})
			if blocked {
				w.retryAlternativeRelease(ctx, activity, "download_status_unavailable")
			}
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
	blocked := false
	if result.Status == "failed" {
		trimmed := strings.TrimSpace(result.Message)
		if trimmed != "" {
			message = &trimmed
		}
		blockMessage := "Download failed"
		if message != nil {
			blockMessage = *message
		}
		blocked = w.blockActivity(ctx, activity, blockMessage, "download_failed")
	}
	updated, err := w.settings.UpdateDownloadActivityProgress(ctx, activity.ID, result.Status, result.ProgressPercent, message)
	if err != nil {
		slog.Error("failed to update download activity progress", "activityId", activity.ID, "status", result.Status, "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download activity progress update failed", map[string]any{"activityId": activity.ID.String(), "status": result.Status, "error": err.Error()})
		return err
	}
	if result.Status == "failed" {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download activity failed", map[string]any{"activityId": activity.ID.String(), "status": result.Status})
		if blocked {
			w.retryAlternativeRelease(ctx, activity, "download_failed")
		}
	}
	w.publishActivity(updated, activity)
	return nil
}

func (w *DownloadActivitySyncWorker) failActivity(ctx context.Context, activity storage.DownloadActivity, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "Download import failed"
	}
	updated, err := w.settings.FailDownloadActivity(ctx, activity.ID, &message, "import")
	if err != nil {
		return err
	}
	w.publishActivity(updated, activity)
	if w.blockActivity(ctx, activity, message, "import_failed") {
		w.retryAlternativeRelease(ctx, activity, "import_failed")
	}
	return nil
}

func (w *DownloadActivitySyncWorker) blockActivity(ctx context.Context, activity storage.DownloadActivity, reason string, source string) bool {
	expiresAt := automaticBlockExpiry(ctx, w.settings)
	if _, err := w.settings.BlockReleaseActivity(ctx, activity, reason, source, &expiresAt); err != nil {
		slog.Error("release block failed", "activityId", activity.ID, "releaseTitle", activity.ReleaseTitle, "source", source, "error", err)
		return false
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventWarning, "downloads", "Release temporarily blocklisted", map[string]any{"mediaItemId": activity.MediaItemID.String(), "activityId": activity.ID.String(), "releaseTitle": activity.ReleaseTitle, "source": source, "expiresAt": expiresAt})
	return true
}

func (w *DownloadActivitySyncWorker) retryAlternativeRelease(ctx context.Context, activity storage.DownloadActivity, source string) {
	if w.indexers == nil || w.downloadClients == nil {
		return
	}
	item, err := w.settings.GetMediaItem(ctx, activity.MediaItemID)
	if err != nil {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Alternative release retry failed to load media", map[string]any{"mediaItemId": activity.MediaItemID.String(), "activityId": activity.ID.String(), "error": err.Error()})
		return
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "downloads", "Trying alternative release after blocklist", map[string]any{"mediaItemId": item.ID.String(), "activityId": activity.ID.String(), "title": item.Title, "source": source})
	if err := autoSearchDownload(ctx, w.settings, w.indexers, w.downloadClients, w.decisions, w.events, item); err != nil {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Alternative release retry failed", map[string]any{"mediaItemId": item.ID.String(), "activityId": activity.ID.String(), "title": item.Title, "error": err.Error()})
	}
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
		"failureType":        activity.FailureType,
		"createdAt":          activity.CreatedAt,
		"updatedAt":          activity.UpdatedAt,
	})
}
