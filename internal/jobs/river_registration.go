package jobs

import (
	"context"
	"log/slog"
	"strings"
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
	settings           *storage.SettingsStore
	indexers           *indexers.Service
	downloadClients    *downloadclients.Service
	decisions          decisions.Engine
	imports            *imports.Service
	subtitles          *subtitles.Service
	events             *events.Broker
	enqueueFulfillment fulfillmentEnqueueFunc
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
	river.AddWorker(workers, &VideoTranscodeWorker{settings: deps.settings, events: deps.events, enqueueFulfillment: deps.enqueueFulfillment})
	river.AddWorker(workers, &AudioTranscodeWorker{settings: deps.settings, events: deps.events, enqueueFulfillment: deps.enqueueFulfillment})
	river.AddWorker(workers, &AudioSourceWorker{settings: deps.settings, events: deps.events})
	river.AddWorker(workers, &ContainerRemuxWorker{settings: deps.settings, events: deps.events})
	river.AddWorker(workers, &SubtitleDownloadWorker{settings: deps.settings, subtitles: deps.subtitles, events: deps.events})
	river.AddWorker(workers, &SubtitleEmbedWorker{settings: deps.settings, events: deps.events})
	river.AddWorker(workers, &SubtitleExtractWorker{settings: deps.settings, events: deps.events})
	river.AddWorker(workers, &SubtitleConvertWorker{settings: deps.settings, events: deps.events})
}

func fixedJobDefinitions() []fixedJobDefinition {
	return []fixedJobDefinition{
		{
			SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
				ID:                    "rss_sync",
				Name:                  "RSS sync",
				Category:              "release_search",
				Description:           "Checks indexer feeds and grabs matching wanted releases.",
				Kind:                  RSSSyncArgs{}.Kind(),
				Queue:                 queueMediaSearch,
				IntervalSeconds:       int32((15 * time.Minute).Seconds()),
				IntervalConfigurable:  true,
				Automatic:             true,
				ManualActionAvailable: true,
			},
			args: func() river.JobArgs { return RSSSyncArgs{} },
		},
		{
			SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
				ID:                    "download_activity_sync",
				Name:                  "Download activity sync",
				Category:              "download_import",
				Description:           "Checks download clients and imports completed media.",
				Kind:                  DownloadActivitySyncArgs{}.Kind(),
				Queue:                 queueDownloads,
				IntervalSeconds:       storage.MinSystemJobScheduleIntervalSeconds,
				IntervalConfigurable:  true,
				HistoryPolicy:         "routine",
				Automatic:             true,
				ManualActionAvailable: true,
			},
			args: func() river.JobArgs { return DownloadActivitySyncArgs{} },
		},
		{
			SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
				ID:                    "release_blocklist_cleanup",
				Name:                  "Release blocklist cleanup",
				Category:              "maintenance",
				Description:           "Expires temporary failed-release blocks.",
				Kind:                  ReleaseBlocklistCleanupArgs{}.Kind(),
				Queue:                 queueDownloads,
				IntervalSeconds:       int32(time.Hour.Seconds()),
				IntervalConfigurable:  true,
				Automatic:             true,
				ManualActionAvailable: true,
			},
			args: func() river.JobArgs { return ReleaseBlocklistCleanupArgs{} },
		},
		{
			SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
				ID:                    "subtitle_retry",
				Name:                  "Subtitle retry",
				Category:              "subtitle_fulfillment",
				Description:           "Retries queued subtitle downloads and imports.",
				Kind:                  SubtitleRetryArgs{}.Kind(),
				Queue:                 queueMediaSearch,
				IntervalSeconds:       int32((6 * time.Hour).Seconds()),
				IntervalConfigurable:  true,
				Automatic:             true,
				ManualActionAvailable: true,
			},
			args: func() river.JobArgs { return SubtitleRetryArgs{} },
		},
		fulfillmentSchedule("video_transcode", "Video transcoding", "Transforms supported video codec or pixel format mismatches.", queueMediaAssembly, VideoTranscodeArgs{}),
		fulfillmentSchedule("audio_transcode", "Audio transcoding", "Transforms audio codec, channels, or bitrate when policy allows.", queueMediaAssembly, AudioTranscodeArgs{}),
		fulfillmentSchedule("audio_source", "Audio sourcing", "Sources desired audio tracks from alternate releases.", queueMediaSearch, AudioSourceArgs{}),
		fulfillmentSchedule("container_remux", "Container remuxing", "Moves selected streams into the target container.", queueMediaAssembly, ContainerRemuxArgs{}),
		fulfillmentSchedule("subtitle_download", "Subtitle download", "Downloads stored external subtitles for missing subtitle targets.", queueMediaSearch, SubtitleDownloadArgs{}),
		fulfillmentSchedule("subtitle_embed", "Subtitle merge/embed", "Embeds external subtitles when embedded mode requires it.", queueMediaAssembly, SubtitleEmbedArgs{}),
		fulfillmentSchedule("subtitle_extract", "Subtitle extraction", "Extracts embedded subtitles when external subtitles are required.", queueMediaAssembly, SubtitleExtractArgs{}),
		fulfillmentSchedule("subtitle_convert", "Subtitle conversion", "Converts subtitle format when tooling supports it.", queueMediaAssembly, SubtitleConvertArgs{}),
	}
}

func fulfillmentSchedule(id string, name string, description string, queue string, args river.JobArgs) fixedJobDefinition {
	return fixedJobDefinition{
		SystemJobScheduleDefinition: storage.SystemJobScheduleDefinition{
			ID:                    id,
			Name:                  name,
			Category:              "fulfillment",
			Description:           description,
			Kind:                  args.Kind(),
			Queue:                 queue,
			IntervalSeconds:       int32((6 * time.Hour).Seconds()),
			IntervalConfigurable:  true,
			Automatic:             true,
			ManualActionAvailable: true,
			PausedByDefault:       true,
		},
		args: func() river.JobArgs { return args },
	}
}

func fixedJobDefinitionByID(id string) (fixedJobDefinition, bool) {
	id = strings.TrimSpace(id)
	for _, definition := range fixedJobDefinitions() {
		if definition.ID == id {
			return definition, true
		}
	}
	return fixedJobDefinition{}, false
}

func periodicJobs(settings *storage.SettingsStore) []*river.PeriodicJob {
	definitions := fixedJobDefinitions()
	jobs := make([]*river.PeriodicJob, 0, len(definitions))
	for _, definition := range definitions {
		definition := definition
		interval := time.Duration(definition.IntervalSeconds) * time.Second
		if definition.IntervalConfigurable {
			interval = time.Duration(storage.MinSystemJobScheduleIntervalSeconds) * time.Second
		}
		jobs = append(jobs, river.NewPeriodicJob(
			river.PeriodicInterval(interval),
			func() (river.JobArgs, *river.InsertOpts) {
				if settings != nil {
					if definition.IntervalConfigurable {
						if !settings.SystemJobScheduleReady(context.Background(), definition.ID) {
							return nil, nil
						}
					} else if settings.SystemJobSchedulePaused(context.Background(), definition.ID) {
						return nil, nil
					}
				}
				return definition.args(), jobInsertOpts(definition.Queue)
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
