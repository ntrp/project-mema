package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/downloadrouting"
	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

const maxAutomaticGrabAttempts = 3

var errRetryAlternativeRelease = errors.New("retry alternative release")

type AutoSearchDownloadArgs struct {
	MediaItemID string `json:"media_item_id" river:"unique"`
}

func (AutoSearchDownloadArgs) Kind() string {
	return "media.auto_search_download"
}

type MissingMediaRetryArgs struct{}

func (MissingMediaRetryArgs) Kind() string {
	return "media.missing_media_retry"
}

type AutoSearchDownloadWorker struct {
	river.WorkerDefaults[AutoSearchDownloadArgs]

	settings        *storage.SettingsStore
	indexers        *indexers.Service
	downloadClients *downloadclients.Service
	decisions       decisions.Engine
	events          *events.Broker
}

func (w *AutoSearchDownloadWorker) Work(ctx context.Context, job *river.Job[AutoSearchDownloadArgs]) (err error) {
	publishJobUpdated(w.events, job.JobRow, "running")
	defer func() { publishJobFinished(w.events, job.JobRow, err) }()

	mediaItemID, err := uuid.Parse(job.Args.MediaItemID)
	if err != nil {
		slog.Error("auto search download invalid media item id", "mediaItemId", job.Args.MediaItemID, "error", err)
		return fmt.Errorf("parse media item id: %w", err)
	}
	item, err := w.settings.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		slog.Error("auto search download media item load failed", "mediaItemId", mediaItemID, "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "jobs", "Automatic search failed to load media", map[string]any{"mediaItemId": mediaItemID.String(), "error": err.Error()})
		return fmt.Errorf("load media item: %w", err)
	}
	slog.Debug("auto search download started", "mediaItemId", item.ID, "title", item.Title, "status", item.Status)
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "jobs", "Automatic search started", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title})
	return autoSearchDownload(ctx, w.settings, w.indexers, w.downloadClients, w.decisions, w.events, item)
}

type MissingMediaRetryWorker struct {
	river.WorkerDefaults[MissingMediaRetryArgs]

	settings        *storage.SettingsStore
	indexers        *indexers.Service
	downloadClients *downloadclients.Service
	decisions       decisions.Engine
	events          *events.Broker
}

func (w *MissingMediaRetryWorker) Work(ctx context.Context, job *river.Job[MissingMediaRetryArgs]) (err error) {
	publishJobUpdated(w.events, job.JobRow, "running")
	defer func() { publishJobFinished(w.events, job.JobRow, err) }()

	items, err := w.settings.ListMissingMediaItems(ctx)
	if err != nil {
		slog.Error("missing media retry list failed", "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "jobs", "Missing media retry failed to list items", map[string]any{"error": err.Error()})
		return fmt.Errorf("list missing media: %w", err)
	}
	slog.Debug("missing media retry started", "itemCount", len(items))
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "jobs", "Missing media retry started", map[string]any{"itemCount": len(items)})
	var failures []string
	for _, item := range items {
		if err := autoSearchDownload(ctx, w.settings, w.indexers, w.downloadClients, w.decisions, w.events, item); err != nil {
			slog.Error("missing media retry item failed", "mediaItemId", item.ID, "title", item.Title, "error", err)
			failures = append(failures, fmt.Sprintf("%s: %s", item.Title, err.Error()))
		}
	}
	if len(failures) > 0 {
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "jobs", "Missing media retry finished with failures", map[string]any{"failureCount": len(failures)})
		return fmt.Errorf("missing media retry failed for %d item(s): %s", len(failures), strings.Join(failures, "; "))
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "jobs", "Missing media retry finished", map[string]any{"itemCount": len(items)})
	return nil
}

func autoSearchDownload(
	ctx context.Context,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	downloadClientService *downloadclients.Service,
	decisionEngine decisions.Engine,
	eventBroker *events.Broker,
	item storage.MediaItem,
) error {
	if item.Status == "downloaded" || item.Status == "downloading" {
		slog.Debug("auto search download skipped", "mediaItemId", item.ID, "title", item.Title, "status", item.Status)
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search skipped", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "status": item.Status})
		return nil
	}
	releases, searchErrors, err := searchReleases(ctx, settings, indexerService, item, decisions.SearchQueryForMediaItem(item), eventBroker, false)
	if err != nil {
		slog.Error("auto search release search failed", "mediaItemId", item.ID, "title", item.Title, "error", err)
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "jobs", "Automatic search failed", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "error": err.Error()})
		return err
	}
	if err := settings.ReplaceReleaseSearchResults(ctx, item.ID, releases, searchErrors); err != nil {
		slog.Error("auto search result storage failed", "mediaItemId", item.ID, "title", item.Title, "error", err)
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "jobs", "Automatic search result storage failed", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "error": err.Error()})
		return err
	}
	slog.Debug("auto search release search finished", "mediaItemId", item.ID, "title", item.Title, "releaseCount", len(releases), "errorCount", len(searchErrors))
	profile, formats, languages := releaseDecisionContext(ctx, settings, item)
	for attempt := 1; attempt <= maxAutomaticGrabAttempts; attempt++ {
		availableReleases, clients, ok, err := autoSearchCandidatesWithDownloadClient(ctx, settings, eventBroker, item, releases, searchErrors)
		if err != nil || !ok {
			return err
		}
		decision, ok := decisionEngine.ChooseReleaseWithProfileAndLanguages(
			item,
			profile,
			formats,
			languages,
			availableReleases,
		)
		if !ok {
			slog.Debug("auto search found no acceptable release", "mediaItemId", item.ID, "title", item.Title)
			publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search found no acceptable release", map[string]any{
				"mediaItemId":  item.ID.String(),
				"title":        item.Title,
				"releaseCount": len(availableReleases),
				"errorCount":   len(searchErrors),
				"reasons":      topDecisionRejections(item, profile, formats, languages, availableReleases),
			})
			return nil
		}
		if err := grabReleaseNow(ctx, settings, downloadClientService, eventBroker, clients, item, decision.Release); err != nil {
			slog.Error("auto search grab failed", "mediaItemId", item.ID, "title", item.Title, "releaseTitle", decision.Release.Title, "error", err)
			if shouldRetryAlternativeRelease(err, attempt) {
				publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search trying alternative release", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "failedReleaseTitle": decision.Release.Title, "attempt": attempt})
				continue
			}
			if automaticRetryLimitReached(err, attempt) {
				publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search retry limit reached", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "attempts": attempt})
				return nil
			} else {
				publishSystemEvent(ctx, settings, eventBroker, jobEventError, "jobs", "Automatic search grab failed", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseTitle": decision.Release.Title, "error": err.Error()})
			}
			return err
		}
		slog.Debug("auto search grab queued", "mediaItemId", item.ID, "title", item.Title, "releaseTitle", decision.Release.Title)
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "jobs", "Automatic search queued download", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseTitle": decision.Release.Title})
		return nil
	}
	return nil
}

func grabReleaseNow(
	ctx context.Context,
	settings *storage.SettingsStore,
	downloadClientService *downloadclients.Service,
	eventBroker *events.Broker,
	clients []storage.DownloadClient,
	item storage.MediaItem,
	release storage.ReleaseCandidateInput,
) error {
	if len(clients) == 0 {
		slog.Debug("grab release skipped because no download clients are enabled", "mediaItemId", item.ID, "title", item.Title, "releaseTitle", release.Title)
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "downloads", "Automatic grab skipped because no download client is enabled", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseTitle": release.Title})
		return settings.ReplaceReleaseSearchResults(ctx, item.ID, []storage.ReleaseCandidateInput{release}, []string{downloadrouting.MissingClientMessage("")})
	}

	client, ok := downloadrouting.ClientForProtocol(clients, release.IndexerProtocol)
	if !ok {
		message := downloadrouting.MissingClientMessage(release.IndexerProtocol)
		slog.Debug("grab release skipped because no compatible download client is enabled", "mediaItemId", item.ID, "title", item.Title, "releaseTitle", release.Title, "protocol", release.IndexerProtocol)
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "downloads", "Automatic grab skipped because no compatible download client is enabled", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseTitle": release.Title, "protocol": release.IndexerProtocol})
		return settings.ReplaceReleaseSearchResults(ctx, item.ID, []storage.ReleaseCandidateInput{release}, []string{message})
	}
	activity, err := settings.CreateDownloadActivity(ctx, storage.DownloadActivityInput{
		MediaItemID:        item.ID,
		ReleaseTitle:       release.Title,
		IndexerName:        release.IndexerName,
		DownloadClientName: client.Name,
		DownloadURL:        release.DownloadURL,
		Status:             "queued",
	})
	if err != nil {
		slog.Error("download activity create failed", "mediaItemId", item.ID, "title", item.Title, "releaseTitle", release.Title, "downloadClientName", client.Name, "error", err)
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "downloads", "Download activity creation failed", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseTitle": release.Title, "downloadClientName": client.Name, "error": err.Error()})
		return fmt.Errorf("record download activity: %w", err)
	}
	activity.MediaTitle = item.Title
	activity.MediaType = item.Type
	publishDownloadActivity(eventBroker, activity)
	result := downloadClientService.Add(ctx, downloadClientConfig(client), downloadclients.AddRequest{
		URL:      release.DownloadURL,
		Title:    release.Title,
		Category: client.Category,
	})
	if !result.Success {
		message := strings.TrimSpace(result.Message)
		if message == "" {
			message = "Download client rejected the release"
		}
		slog.Error("download client rejected release", "activityId", activity.ID, "mediaItemId", item.ID, "releaseTitle", release.Title, "downloadClientName", client.Name, "message", message)
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "downloads", "Download client rejected automatic release", map[string]any{"activityId": activity.ID.String(), "mediaItemId": item.ID.String(), "releaseTitle": release.Title, "downloadClientName": client.Name, "message": message})
		if blockAutomaticRelease(ctx, settings, eventBroker, release, message, "download_client_rejected") {
			if _, err := settings.FailDownloadActivity(ctx, activity.ID, &message, "download"); err != nil {
				return err
			}
			return fmt.Errorf("%w: %s", errRetryAlternativeRelease, release.Title)
		}
		_, err := settings.FailDownloadActivity(ctx, activity.ID, &message, "download")
		return err
	}
	slog.Debug("download client accepted release", "activityId", activity.ID, "mediaItemId", item.ID, "releaseTitle", release.Title, "downloadClientName", client.Name, "downloadId", result.DownloadID)
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "downloads", "Download client accepted automatic release", map[string]any{"activityId": activity.ID.String(), "mediaItemId": item.ID.String(), "releaseTitle": release.Title, "downloadClientName": client.Name, "downloadId": result.DownloadID})
	updated, err := settings.UpdateDownloadActivityClientState(ctx, activity.ID, "grabbed", optionalString(result.DownloadID), nil)
	if err == nil {
		updated.MediaTitle = item.Title
		updated.MediaType = item.Type
		publishDownloadActivity(eventBroker, updated)
	} else {
		slog.Error("download activity client state update failed", "activityId", activity.ID, "mediaItemId", item.ID, "error", err)
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "downloads", "Download activity state update failed", map[string]any{"activityId": activity.ID.String(), "mediaItemId": item.ID.String(), "error": err.Error()})
	}
	return err
}
