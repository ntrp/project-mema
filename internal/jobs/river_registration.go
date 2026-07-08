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

type fixedJobDefinition struct {
	storage.SystemJobScheduleDefinition
	args func() river.JobArgs
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
	river.AddWorker(workers, &ComponentMuxWorker{settings: deps.settings, events: deps.events})
}

func fixedJobDefinitions() []fixedJobDefinition {
	return []fixedJobDefinition{
		{
			SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
				ID:              "rss_sync",
				Name:            "RSS sync",
				Kind:            RSSSyncArgs{}.Kind(),
				Queue:           queueMediaSearch,
				IntervalSeconds: int32((15 * time.Minute).Seconds()),
			},
			args: func() river.JobArgs { return RSSSyncArgs{} },
		},
		{
			SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
				ID:              "download_activity_sync",
				Name:            "Download activity sync",
				Kind:            DownloadActivitySyncArgs{}.Kind(),
				Queue:           queueDownloads,
				IntervalSeconds: int32((10 * time.Second).Seconds()),
			},
			args: func() river.JobArgs { return DownloadActivitySyncArgs{} },
		},
		{
			SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
				ID:              "release_blocklist_cleanup",
				Name:            "Release blocklist cleanup",
				Kind:            ReleaseBlocklistCleanupArgs{}.Kind(),
				Queue:           queueDownloads,
				IntervalSeconds: int32(time.Hour.Seconds()),
			},
			args: func() river.JobArgs { return ReleaseBlocklistCleanupArgs{} },
		},
		{
			SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
				ID:              "subtitle_retry",
				Name:            "Subtitle retry",
				Kind:            SubtitleRetryArgs{}.Kind(),
				Queue:           queueMediaSearch,
				IntervalSeconds: int32((6 * time.Hour).Seconds()),
			},
			args: func() river.JobArgs { return SubtitleRetryArgs{} },
		},
	}
}

func periodicJobs(settings *storage.SettingsStore) []*river.PeriodicJob {
	definitions := fixedJobDefinitions()
	jobs := make([]*river.PeriodicJob, 0, len(definitions))
	for _, definition := range definitions {
		definition := definition
		jobs = append(jobs, river.NewPeriodicJob(
			river.PeriodicInterval(time.Duration(definition.IntervalSeconds)*time.Second),
			func() (river.JobArgs, *river.InsertOpts) {
				if settings != nil && settings.SystemJobSchedulePaused(context.Background(), definition.ID) {
					return nil, nil
				}
				return definition.args(), &river.InsertOpts{Queue: definition.Queue}
			},
			&river.PeriodicJobOpts{ID: definition.ID},
		))
	}
	return jobs
}

func fixedScheduleDefinitions() []storage.SystemJobScheduleDefinition {
	definitions := fixedJobDefinitions()
	schedules := make([]storage.SystemJobScheduleDefinition, 0, len(definitions))
	for _, definition := range definitions {
		schedules = append(schedules, definition.SystemJobScheduleDefinition)
	}
	return schedules
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
