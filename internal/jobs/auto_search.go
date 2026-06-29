package jobs

import (
	"context"
	"fmt"
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
		return fmt.Errorf("parse media item id: %w", err)
	}
	item, err := w.settings.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return fmt.Errorf("load media item: %w", err)
	}
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
		return fmt.Errorf("list missing media: %w", err)
	}
	var failures []string
	for _, item := range items {
		if err := autoSearchDownload(ctx, w.settings, w.indexers, w.downloadClients, w.decisions, w.events, item); err != nil {
			failures = append(failures, fmt.Sprintf("%s: %s", item.Title, err.Error()))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("missing media retry failed for %d item(s): %s", len(failures), strings.Join(failures, "; "))
	}
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
		return nil
	}
	releases, searchErrors, err := searchReleases(ctx, settings, indexerService, item)
	if err != nil {
		return err
	}
	if err := settings.ReplaceReleaseSearchResults(ctx, item.ID, releases, searchErrors); err != nil {
		return err
	}
	decision, ok := decisionEngine.ChooseRelease(releases)
	if !ok {
		return nil
	}
	if err := grabReleaseNow(ctx, settings, downloadClientService, eventBroker, item, decision.Release); err != nil {
		return err
	}
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
		return nil, []string{"No enabled indexer is configured"}, nil
	}

	releases := []storage.ReleaseCandidateInput{}
	searchErrors := []string{}
	for _, config := range configs {
		found, err := indexerService.Search(ctx, indexerConfig(config), item.Title, item.Type)
		if err != nil {
			searchErrors = append(searchErrors, fmt.Sprintf("%s: %s", config.Name, err.Error()))
			continue
		}
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
		return fmt.Errorf("list enabled download clients: %w", err)
	}
	if len(clients) == 0 {
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
		_, err := settings.UpdateDownloadActivityStatus(ctx, activity.ID, "failed", &message)
		return err
	}
	updated, err := settings.UpdateDownloadActivityClientState(ctx, activity.ID, "grabbed", optionalString(result.DownloadID), nil)
	if err == nil {
		updated.MediaTitle = item.Title
		updated.MediaType = item.Type
		publishDownloadActivity(eventBroker, updated)
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
