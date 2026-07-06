package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/imports"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

type workerDependencies struct {
	settings        *storage.SettingsStore
	indexers        *indexers.Service
	downloadClients *downloadclients.Service
	decisions       decisions.Engine
	imports         *imports.Service
	subtitles       *subtitles.Service
	events          *events.Broker
}

func addWorkers(workers *river.Workers, deps workerDependencies) {
	river.AddWorker(workers, &ReleaseSearchWorker{settings: deps.settings, indexers: deps.indexers, events: deps.events})
	river.AddWorker(workers, &AutoSearchDownloadWorker{
		settings:        deps.settings,
		indexers:        deps.indexers,
		downloadClients: deps.downloadClients,
		decisions:       deps.decisions,
		events:          deps.events,
	})
	river.AddWorker(workers, &RSSSyncWorker{
		settings:        deps.settings,
		indexers:        deps.indexers,
		downloadClients: deps.downloadClients,
		decisions:       deps.decisions,
		events:          deps.events,
	})
	river.AddWorker(workers, &GrabReleaseWorker{settings: deps.settings, downloadClients: deps.downloadClients, events: deps.events})
	river.AddWorker(workers, &DownloadActivitySyncWorker{
		settings:        deps.settings,
		indexers:        deps.indexers,
		downloadClients: deps.downloadClients,
		decisions:       deps.decisions,
		imports:         deps.imports,
		events:          deps.events,
	})
	river.AddWorker(workers, &ReleaseBlocklistCleanupWorker{settings: deps.settings, events: deps.events})
	river.AddWorker(workers, &SubtitleSearchWorker{settings: deps.settings, subtitles: deps.subtitles, events: deps.events})
	river.AddWorker(workers, &SubtitleRetryWorker{settings: deps.settings, subtitles: deps.subtitles, events: deps.events})
	river.AddWorker(workers, &ComponentExtractionWorker{settings: deps.settings, events: deps.events})
}

func periodicJobs() []*river.PeriodicJob {
	return []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(15*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return RSSSyncArgs{}, &river.InsertOpts{Queue: queueMediaSearch}
			},
			&river.PeriodicJobOpts{ID: "rss_sync"},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(10*time.Second),
			func() (river.JobArgs, *river.InsertOpts) {
				return DownloadActivitySyncArgs{}, &river.InsertOpts{Queue: queueDownloads}
			},
			&river.PeriodicJobOpts{ID: "download_activity_sync"},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(1*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return ReleaseBlocklistCleanupArgs{}, &river.InsertOpts{Queue: queueDownloads}
			},
			&river.PeriodicJobOpts{ID: "release_blocklist_cleanup"},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(6*time.Hour),
			func() (river.JobArgs, *river.InsertOpts) {
				return SubtitleRetryArgs{}, &river.InsertOpts{Queue: queueMediaSearch}
			},
			&river.PeriodicJobOpts{ID: "subtitle_retry"},
		),
	}
}

func cleanupLegacyJobs(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, `
		delete from river_job
		where kind in ('media.wanted_rss_sync', 'media.missing_media_retry')
			and state in ('available', 'scheduled', 'retryable')
	`)
	if err != nil {
		slog.Debug("legacy RSS job cleanup skipped", "error", err)
	}
}
