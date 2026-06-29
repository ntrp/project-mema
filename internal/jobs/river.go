package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/imports"
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
	releases, searchErrors, err := searchReleases(ctx, w.settings, w.indexers, item)
	if err != nil {
		return err
	}
	return w.settings.ReplaceReleaseSearchResults(ctx, mediaItemID, releases, searchErrors)
}

type GrabReleaseWorker struct {
	river.WorkerDefaults[GrabReleaseArgs]

	settings        *storage.SettingsStore
	downloadClients *downloadclients.Service
	events          *events.Broker
}

func (w *GrabReleaseWorker) Work(ctx context.Context, job *river.Job[GrabReleaseArgs]) error {
	activityID, err := uuid.Parse(job.Args.ActivityID)
	if err != nil {
		return fmt.Errorf("parse activity id: %w", err)
	}
	activity, err := w.settings.GetDownloadActivity(ctx, activityID)
	if err != nil {
		return fmt.Errorf("load download activity: %w", err)
	}
	if activity.Status == "cancelled" {
		return nil
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
	downloadID := optionalString(result.DownloadID)
	activity, err = w.settings.UpdateDownloadActivityClientState(ctx, activityID, "grabbed", downloadID, nil)
	if err == nil {
		publishDownloadActivity(w.events, activity)
	}
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

func NewClient(pool *pgxpool.Pool, settings *storage.SettingsStore, indexerService *indexers.Service, downloadClientService *downloadclients.Service, eventBroker *events.Broker) (*Client, error) {
	workers := river.NewWorkers()
	decisionEngine := decisions.NewEngine()
	importService := imports.NewService(settings)
	if eventBroker == nil {
		eventBroker = events.NewBroker()
	}
	river.AddWorker(workers, &ReleaseSearchWorker{settings: settings, indexers: indexerService})
	river.AddWorker(workers, &AutoSearchDownloadWorker{
		settings:        settings,
		indexers:        indexerService,
		downloadClients: downloadClientService,
		decisions:       decisionEngine,
		events:          eventBroker,
	})
	river.AddWorker(workers, &MissingMediaRetryWorker{
		settings:        settings,
		indexers:        indexerService,
		downloadClients: downloadClientService,
		decisions:       decisionEngine,
		events:          eventBroker,
	})
	river.AddWorker(workers, &GrabReleaseWorker{settings: settings, downloadClients: downloadClientService, events: eventBroker})
	river.AddWorker(workers, &DownloadActivitySyncWorker{
		settings:        settings,
		downloadClients: downloadClientService,
		imports:         importService,
		events:          eventBroker,
	})

	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			queueMediaSearch: {MaxWorkers: 2},
			queueDownloads:   {MaxWorkers: 2},
		},
		PeriodicJobs: []*river.PeriodicJob{
			river.NewPeriodicJob(
				river.PeriodicInterval(6*time.Hour),
				func() (river.JobArgs, *river.InsertOpts) {
					return MissingMediaRetryArgs{}, &river.InsertOpts{Queue: queueMediaSearch}
				},
				&river.PeriodicJobOpts{ID: "missing_media_retry"},
			),
			river.NewPeriodicJob(
				river.PeriodicInterval(30*time.Second),
				func() (river.JobArgs, *river.InsertOpts) {
					return DownloadActivitySyncArgs{}, &river.InsertOpts{Queue: queueDownloads}
				},
				&river.PeriodicJobOpts{ID: "download_activity_sync"},
			),
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

func (c *Client) EnqueueAutoSearchDownload(ctx context.Context, mediaItemID uuid.UUID) (int64, error) {
	result, err := c.river.Insert(ctx, AutoSearchDownloadArgs{MediaItemID: mediaItemID.String()}, &river.InsertOpts{
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
