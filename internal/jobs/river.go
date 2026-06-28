package jobs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

const (
	queueMediaSearch = "media_search"
	queueDownloads   = "downloads"
)

type Client struct {
	river *river.Client[pgx.Tx]
}

type ReleaseSearchArgs struct {
	MediaItemID string `json:"media_item_id" river:"unique"`
}

func (ReleaseSearchArgs) Kind() string {
	return "media.release_search"
}

type GrabReleaseArgs struct {
	ActivityID  string `json:"activity_id" river:"unique"`
	MediaItemID string `json:"media_item_id"`
	Title       string `json:"title"`
	DownloadURL string `json:"download_url"`
	IndexerName string `json:"indexer_name"`
}

func (GrabReleaseArgs) Kind() string {
	return "media.grab_release"
}

type ReleaseSearchWorker struct {
	river.WorkerDefaults[ReleaseSearchArgs]

	settings *storage.SettingsStore
	indexers *indexers.Service
}

func (w *ReleaseSearchWorker) Work(ctx context.Context, job *river.Job[ReleaseSearchArgs]) error {
	mediaItemID, err := uuid.Parse(job.Args.MediaItemID)
	if err != nil {
		return fmt.Errorf("parse media item id: %w", err)
	}
	item, err := w.settings.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return fmt.Errorf("load media item: %w", err)
	}
	configs, err := w.settings.ListEnabledIndexers(ctx)
	if err != nil {
		return fmt.Errorf("list enabled indexers: %w", err)
	}
	if len(configs) == 0 {
		return w.settings.ReplaceReleaseSearchResults(ctx, mediaItemID, nil, []string{"No enabled indexer is configured"})
	}

	releases := []storage.ReleaseCandidateInput{}
	searchErrors := []string{}
	for _, config := range configs {
		found, err := w.indexers.Search(ctx, indexerConfig(config), item.Title, item.Type)
		if err != nil {
			searchErrors = append(searchErrors, fmt.Sprintf("%s: %s", config.Name, err.Error()))
			continue
		}
		for _, release := range found {
			releases = append(releases, releaseCandidateInput(mediaItemID, release))
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
	return w.settings.ReplaceReleaseSearchResults(ctx, mediaItemID, releases, searchErrors)
}

type GrabReleaseWorker struct {
	river.WorkerDefaults[GrabReleaseArgs]

	settings        *storage.SettingsStore
	downloadClients *downloadclients.Service
}

func (w *GrabReleaseWorker) Work(ctx context.Context, job *river.Job[GrabReleaseArgs]) error {
	activityID, err := uuid.Parse(job.Args.ActivityID)
	if err != nil {
		return fmt.Errorf("parse activity id: %w", err)
	}
	clients, err := w.settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		return fmt.Errorf("list enabled download clients: %w", err)
	}
	if len(clients) == 0 {
		return w.markGrabFailed(ctx, activityID, "No enabled download client is configured")
	}

	client := clients[0]
	result := w.downloadClients.Add(ctx, downloadClientConfig(client), downloadclients.AddRequest{
		URL:      job.Args.DownloadURL,
		Title:    job.Args.Title,
		Category: client.Category,
	})
	if !result.Success {
		return w.markGrabFailed(ctx, activityID, result.Message)
	}
	_, err = w.settings.UpdateDownloadActivityStatus(ctx, activityID, "grabbed", nil)
	return err
}

func (w *GrabReleaseWorker) markGrabFailed(ctx context.Context, activityID uuid.UUID, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "Download client rejected the release"
	}
	_, err := w.settings.UpdateDownloadActivityStatus(ctx, activityID, "failed", &message)
	return err
}

func NewClient(pool *pgxpool.Pool, settings *storage.SettingsStore, indexerService *indexers.Service, downloadClientService *downloadclients.Service) (*Client, error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, &ReleaseSearchWorker{settings: settings, indexers: indexerService})
	river.AddWorker(workers, &GrabReleaseWorker{settings: settings, downloadClients: downloadClientService})

	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			queueMediaSearch: {MaxWorkers: 2},
			queueDownloads:   {MaxWorkers: 2},
		},
		SoftStopTimeout: 10 * time.Second,
		Workers:         workers,
	})
	if err != nil {
		return nil, err
	}
	return &Client{river: riverClient}, nil
}

func (c *Client) Start(ctx context.Context) error {
	return c.river.Start(ctx)
}

func (c *Client) Stop(ctx context.Context) error {
	return c.river.Stop(ctx)
}

func (c *Client) EnqueueReleaseSearch(ctx context.Context, mediaItemID uuid.UUID) (int64, error) {
	result, err := c.river.Insert(ctx, ReleaseSearchArgs{MediaItemID: mediaItemID.String()}, &river.InsertOpts{
		Queue: queueMediaSearch,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	if err != nil {
		return 0, err
	}
	return result.Job.ID, nil
}

func (c *Client) EnqueueGrabRelease(ctx context.Context, args GrabReleaseArgs) (int64, error) {
	result, err := c.river.Insert(ctx, args, &river.InsertOpts{
		Queue: queueDownloads,
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
		},
	})
	if err != nil {
		return 0, err
	}
	return result.Job.ID, nil
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
