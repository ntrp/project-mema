package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/downloadrouting"
	"media-manager/internal/events"
	"media-manager/internal/imports"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

const (
	queueMediaSearch   = "media_search"
	queueDownloads     = "downloads"
	queueMediaAssembly = "media_assembly"
)

type Client struct {
	river    *river.Client[pgx.Tx]
	settings *storage.SettingsStore
	events   *events.Broker
}

type ReleaseSearchArgs struct {
	MediaItemID string `json:"media_item_id" river:"unique"`
	Query       string `json:"query,omitempty" river:"unique"`
}

func (ReleaseSearchArgs) Kind() string {
	return "media.release_search"
}

type ReleaseSearchWorker struct {
	river.WorkerDefaults[ReleaseSearchArgs]

	settings *storage.SettingsStore
	indexers *indexers.Service
	events   *events.Broker
}

func (w *ReleaseSearchWorker) Work(ctx context.Context, job *river.Job[ReleaseSearchArgs]) (err error) {
	ctx = withJobExecution(ctx, job.JobRow.ID)
	recordJobUpdated(ctx, w.settings, w.events, job.JobRow, "running")
	defer func() { recordJobFinished(ctx, w.settings, w.events, job.JobRow, err) }()

	mediaItemID, err := uuid.Parse(job.Args.MediaItemID)
	if err != nil {
		slog.Error("release search invalid media item id", "mediaItemId", job.Args.MediaItemID, "error", err)
		return fmt.Errorf("parse media item id: %w", err)
	}
	item, err := w.settings.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		slog.Error("release search media item load failed", "mediaItemId", mediaItemID, "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "jobs", "Release search failed to load media", map[string]any{"mediaItemId": mediaItemID.String(), "error": err.Error()})
		return fmt.Errorf("load media item: %w", err)
	}
	query := strings.TrimSpace(job.Args.Query)
	if query == "" {
		query = decisions.SearchQueryForMediaItem(item)
	}
	slog.Debug("release search started", "mediaItemId", item.ID, "title", item.Title, "query", query)
	recordJobProgress(ctx, w.settings, w.events, nil, "Searching releases")
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "jobs", "Release search started", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "query": query})
	releases, searchErrors, err := searchReleases(ctx, w.settings, w.indexers, item, query, w.events, true)
	if err != nil {
		slog.Error("release search failed", "mediaItemId", item.ID, "title", item.Title, "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "jobs", "Release search failed", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "error": err.Error()})
		return err
	}
	slog.Debug("release search finished", "mediaItemId", item.ID, "title", item.Title, "releaseCount", len(releases), "errorCount", len(searchErrors))
	severity := jobEventInfo
	if len(searchErrors) > 0 {
		severity = jobEventWarning
	}
	publishSystemEvent(ctx, w.settings, w.events, severity, "jobs", "Release search finished", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseCount": len(releases), "errorCount": len(searchErrors)})
	return w.settings.ReplaceReleaseSearchResults(ctx, mediaItemID, releases, searchErrors)
}

type GrabReleaseWorker struct {
	river.WorkerDefaults[GrabReleaseArgs]

	settings        *storage.SettingsStore
	downloadClients *downloadclients.Service
	events          *events.Broker
}

func (w *GrabReleaseWorker) Work(ctx context.Context, job *river.Job[GrabReleaseArgs]) (err error) {
	ctx = withJobExecution(ctx, job.JobRow.ID)
	recordJobUpdated(ctx, w.settings, w.events, job.JobRow, "running")
	defer func() { recordJobFinished(ctx, w.settings, w.events, job.JobRow, err) }()

	activityID, err := uuid.Parse(job.Args.ActivityID)
	if err != nil {
		slog.Error("grab release invalid activity id", "activityId", job.Args.ActivityID, "error", err)
		return fmt.Errorf("parse activity id: %w", err)
	}
	activity, err := w.settings.GetDownloadActivity(ctx, activityID)
	if err != nil {
		slog.Error("grab release activity load failed", "activityId", activityID, "error", err)
		return fmt.Errorf("load download activity: %w", err)
	}
	if activity.Status == "cancelled" {
		slog.Debug("grab release skipped cancelled activity", "activityId", activity.ID)
		publishSystemEvent(ctx, w.settings, w.events, jobEventWarning, "downloads", "Download job skipped because activity was cancelled", map[string]any{"activityId": activity.ID.String()})
		return nil
	}
	slog.Debug("grab release started", "activityId", activity.ID, "mediaItemId", activity.MediaItemID, "releaseTitle", job.Args.Title)
	progress := int32(20)
	recordJobProgress(ctx, w.settings, w.events, &progress, "Sending release to download client")
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "downloads", "Download job started", map[string]any{"activityId": activity.ID.String(), "mediaItemId": activity.MediaItemID.String(), "releaseTitle": job.Args.Title})
	clients, err := w.settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		slog.Error("grab release download client list failed", "activityId", activity.ID, "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download job failed to load clients", map[string]any{"activityId": activity.ID.String(), "error": err.Error()})
		return fmt.Errorf("list enabled download clients: %w", err)
	}
	if len(clients) == 0 {
		return w.markGrabFailed(ctx, activityID, "No enabled download client is configured")
	}

	client, ok := downloadrouting.NamedClientForProtocol(clients, activity.DownloadClientName, job.Args.Protocol)
	if !ok {
		return w.markGrabFailed(ctx, activityID, downloadrouting.MissingClientMessage(job.Args.Protocol))
	}
	result := w.downloadClients.Add(ctx, downloadClientConfig(client), downloadclients.AddRequest{
		URL:      job.Args.DownloadURL,
		Title:    job.Args.Title,
		Category: client.Category,
	})
	if !result.Success {
		slog.Error("grab release download client rejected", "activityId", activity.ID, "downloadClientName", client.Name, "message", result.Message)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download client rejected release", map[string]any{"activityId": activity.ID.String(), "downloadClientName": client.Name, "message": result.Message})
		return w.markGrabFailed(ctx, activityID, result.Message)
	}
	downloadID := optionalString(result.DownloadID)
	activity, err = w.settings.UpdateDownloadActivityClientState(ctx, activityID, "grabbed", downloadID, nil)
	if err == nil {
		done := int32(100)
		recordJobProgress(ctx, w.settings, w.events, &done, "Download accepted")
		publishDownloadActivity(w.events, activity)
		slog.Debug("grab release finished", "activityId", activity.ID, "downloadClientName", client.Name, "downloadId", result.DownloadID)
		publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "downloads", "Download sent to client", map[string]any{"activityId": activity.ID.String(), "downloadClientName": client.Name, "downloadId": result.DownloadID})
	} else {
		slog.Error("grab release activity update failed", "activityId", activityID, "error", err)
		publishSystemEvent(ctx, w.settings, w.events, jobEventError, "downloads", "Download activity update failed", map[string]any{"activityId": activityID.String(), "error": err.Error()})
	}
	return err
}

func (w *GrabReleaseWorker) markGrabFailed(ctx context.Context, activityID uuid.UUID, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "Download client rejected the release"
	}
	slog.Error("marking grab release failed", "activityId", activityID, "message", message)
	activity, loadErr := w.settings.GetDownloadActivity(ctx, activityID)
	if loadErr == nil {
		expiresAt := automaticBlockExpiry(ctx, w.settings)
		if _, err := w.settings.BlockReleaseActivity(ctx, activity, message, "download_client_rejected", &expiresAt); err != nil {
			slog.Error("manual release block failed", "activityId", activityID, "releaseTitle", activity.ReleaseTitle, "error", err)
		}
	}
	_, err := w.settings.FailDownloadActivity(ctx, activityID, &message, "download")
	return err
}

func NewClient(pool *pgxpool.Pool, settings *storage.SettingsStore, indexerService *indexers.Service, downloadClientService *downloadclients.Service, eventBroker *events.Broker) (*Client, error) {
	workers := river.NewWorkers()
	decisionEngine := decisions.NewEngine()
	importService := imports.NewService(settings)
	subtitleService := subtitles.NewService(nil)
	if eventBroker == nil {
		eventBroker = events.NewBroker()
	}
	if settings != nil {
		if err := settings.SyncSystemJobSchedules(context.Background(), fixedScheduleDefinitions()); err != nil {
			return nil, err
		}
	}
	addWorkers(workers, workerDependencies{
		settings:        settings,
		indexers:        indexerService,
		downloadClients: downloadClientService,
		decisions:       decisionEngine,
		imports:         importService,
		subtitles:       subtitleService,
		events:          eventBroker,
	})
	cleanupLegacyJobs(context.Background(), pool)

	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			queueMediaSearch:   {MaxWorkers: 2},
			queueDownloads:     {MaxWorkers: 2},
			queueMediaAssembly: {MaxWorkers: 2},
		},
		PeriodicJobs:    periodicJobs(settings),
		SoftStopTimeout: 10 * time.Second,
		Workers:         workers,
	})
	if err != nil {
		return nil, err
	}
	return &Client{river: riverClient, settings: settings, events: eventBroker}, nil
}
