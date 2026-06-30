package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

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

func (w *AutoSearchDownloadWorker) Work(ctx context.Context, job *river.Job[AutoSearchDownloadArgs]) error {
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

func (w *MissingMediaRetryWorker) Work(ctx context.Context, _ *river.Job[MissingMediaRetryArgs]) error {
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
	releases, searchErrors, err := searchReleases(ctx, settings, indexerService, item)
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
	decision, ok := decisionEngine.ChooseRelease(releases)
	if !ok {
		slog.Debug("auto search found no acceptable release", "mediaItemId", item.ID, "title", item.Title)
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search found no acceptable release", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseCount": len(releases), "errorCount": len(searchErrors)})
		return nil
	}
	if err := grabReleaseNow(ctx, settings, downloadClientService, eventBroker, item, decision.Release); err != nil {
		slog.Error("auto search grab failed", "mediaItemId", item.ID, "title", item.Title, "releaseTitle", decision.Release.Title, "error", err)
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "jobs", "Automatic search grab failed", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseTitle": decision.Release.Title, "error": err.Error()})
		return err
	}
	slog.Debug("auto search grab queued", "mediaItemId", item.ID, "title", item.Title, "releaseTitle", decision.Release.Title)
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "jobs", "Automatic search queued download", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseTitle": decision.Release.Title})
	return nil
}

func searchReleases(
	ctx context.Context,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	item storage.MediaItem,
) ([]storage.ReleaseCandidateInput, []string, error) {
	configs, err := settings.ListEnabledIndexers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list enabled indexers: %w", err)
	}
	if len(configs) == 0 {
		slog.Debug("release search skipped because no indexers are enabled", "mediaItemId", item.ID, "title", item.Title)
		return nil, []string{"No enabled indexer is configured"}, nil
	}

	releases := []storage.ReleaseCandidateInput{}
	searchErrors := []string{}
	for _, config := range configs {
		found, err := indexerService.Search(ctx, indexerConfig(config), item.Title, item.Type)
		if err != nil {
			recordIndexerSearchFailure(ctx, settings, config, err)
			slog.Error("indexer release search failed", "mediaItemId", item.ID, "title", item.Title, "indexerName", config.Name, "error", err)
			searchErrors = append(searchErrors, fmt.Sprintf("%s: %s", config.Name, err.Error()))
			continue
		}
		if _, err := settings.RecordIndexerSuccess(ctx, config.ID); err != nil {
			slog.Error("indexer success state update failed", "indexerName", config.Name, "error", err)
		}
		slog.Debug("indexer release search finished", "mediaItemId", item.ID, "title", item.Title, "indexerName", config.Name, "releaseCount", len(found))
		for _, release := range found {
			releases = append(releases, releaseCandidateInput(item.ID, release))
		}
	}
	sort.SliceStable(releases, func(i, j int) bool {
		left := releases[i]
		right := releases[j]
		if left.Seeders != nil && right.Seeders != nil && *left.Seeders != *right.Seeders {
			return *left.Seeders > *right.Seeders
		}
		return left.SizeBytes > right.SizeBytes
	})
	if len(releases) == 0 && len(searchErrors) == 0 {
		searchErrors = append(searchErrors, "No releases found")
	}
	return releases, searchErrors, nil
}

func grabReleaseNow(
	ctx context.Context,
	settings *storage.SettingsStore,
	downloadClientService *downloadclients.Service,
	eventBroker *events.Broker,
	item storage.MediaItem,
	release storage.ReleaseCandidateInput,
) error {
	clients, err := settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "downloads", "Automatic grab failed to list clients", map[string]any{"mediaItemId": item.ID.String(), "error": err.Error()})
		return fmt.Errorf("list enabled download clients: %w", err)
	}
	if len(clients) == 0 {
		slog.Debug("grab release skipped because no download clients are enabled", "mediaItemId", item.ID, "title", item.Title, "releaseTitle", release.Title)
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "downloads", "Automatic grab skipped because no download client is enabled", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseTitle": release.Title})
		return settings.ReplaceReleaseSearchResults(ctx, item.ID, []storage.ReleaseCandidateInput{release}, []string{"No enabled download client is configured"})
	}

	client := clients[0]
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

func indexerConfig(indexer storage.Indexer) indexers.Config {
	return indexers.Config{
		ID:         indexer.ID.String(),
		Name:       indexer.Name,
		Type:       indexer.Type,
		BaseURL:    indexer.BaseURL,
		APIKey:     indexer.APIKey,
		Categories: indexer.Categories,
	}
}

func downloadClientConfig(client storage.DownloadClient) downloadclients.Config {
	return downloadclients.Config{
		Name:     client.Name,
		Type:     client.Type,
		BaseURL:  client.BaseURL,
		Username: client.Username,
		Password: client.Password,
		APIKey:   client.APIKey,
		Category: client.Category,
	}
}

func releaseCandidateInput(mediaItemID uuid.UUID, release indexers.Release) storage.ReleaseCandidateInput {
	var indexerID *uuid.UUID
	if parsed, err := uuid.Parse(release.IndexerID); err == nil {
		indexerID = &parsed
	}
	return storage.ReleaseCandidateInput{
		MediaItemID: mediaItemID,
		IndexerID:   indexerID,
		IndexerName: release.IndexerName,
		IndexerType: release.IndexerType,
		Title:       release.Title,
		DownloadURL: release.DownloadURL,
		InfoURL:     optionalString(release.InfoURL),
		GUID:        optionalString(release.GUID),
		SizeBytes:   release.SizeBytes,
		Seeders:     release.Seeders,
		Peers:       release.Peers,
	}
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
