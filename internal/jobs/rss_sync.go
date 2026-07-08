package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/downloadrouting"
	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

type RSSSyncWorker struct {
	river.WorkerDefaults[RSSSyncArgs]

	settings        *storage.SettingsStore
	indexers        *indexers.Service
	downloadClients *downloadclients.Service
	decisions       decisions.Engine
	events          *events.Broker
}

func (w *RSSSyncWorker) Work(ctx context.Context, job *river.Job[RSSSyncArgs]) (err error) {
	return runRSSSyncWorker(ctx, job.JobRow, w.settings, w.indexers, w.downloadClients, w.decisions, w.events)
}

func runRSSSyncWorker(
	ctx context.Context,
	jobRow *rivertype.JobRow,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	downloadClientService *downloadclients.Service,
	decisionEngine decisions.Engine,
	eventBroker *events.Broker,
) (err error) {
	ctx = withJobExecution(ctx, jobRow.ID)
	recordJobUpdated(ctx, settings, eventBroker, jobRow, "running")
	defer func() { recordJobFinished(ctx, settings, eventBroker, jobRow, err) }()

	items, err := settings.ListMissingMediaItems(ctx)
	if err != nil {
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "jobs", "RSS sync failed to list media", map[string]any{"error": err.Error()})
		return fmt.Errorf("list monitored missing media: %w", err)
	}
	configs, err := settings.ListRSSEnabledIndexers(ctx)
	if err != nil {
		return fmt.Errorf("list RSS indexers: %w", err)
	}
	if len(configs) == 0 {
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "RSS sync skipped because no RSS indexers are enabled", nil)
		return nil
	}

	slog.Debug("rss sync started", "indexerCount", len(configs), "candidateMediaCount", len(items))
	recordJobProgress(ctx, settings, eventBroker, nil, fmt.Sprintf("Checking %d RSS indexer(s)", len(configs)))
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "jobs", "RSS sync started", map[string]any{"indexerCount": len(configs), "candidateMediaCount": len(items)})
	matchesByMedia, failures := fetchRSSMatches(ctx, settings, indexerService, eventBroker, configs, items)
	if err := processRSSMatches(ctx, settings, downloadClientService, decisionEngine, eventBroker, matchesByMedia); err != nil {
		return err
	}
	if len(failures) > 0 {
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "RSS sync finished with indexer failures", map[string]any{"failureCount": len(failures)})
		return fmt.Errorf("rss sync failed for %d indexer(s): %s", len(failures), strings.Join(failures, "; "))
	}
	done := int32(100)
	recordJobProgress(ctx, settings, eventBroker, &done, "RSS sync finished")
	publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "jobs", "RSS sync finished", map[string]any{"matchedMediaCount": len(matchesByMedia)})
	return nil
}

func fetchRSSMatches(
	ctx context.Context,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	eventBroker *events.Broker,
	configs []storage.Indexer,
	items []storage.MediaItem,
) (map[uuid.UUID][]storage.ReleaseCandidateInput, []string) {
	matchesByMedia := map[uuid.UUID][]storage.ReleaseCandidateInput{}
	failures := []string{}
	for _, config := range configs {
		releases, err := indexerService.Recent(ctx, indexerConfig(config))
		if err != nil {
			recordIndexerSearchFailure(ctx, settings, eventBroker, config, err)
			failures = append(failures, fmt.Sprintf("%s: %s", config.Name, err.Error()))
			continue
		}
		recordIndexerSearchSuccess(ctx, settings, config)
		unseen, markerCovered := unseenRSSReleases(config, releases)
		if !markerCovered {
			publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "indexers", "RSS marker was not found in the returned feed page", map[string]any{"indexerId": config.ID.String(), "indexerName": config.Name})
		}
		for _, release := range unseen {
			for _, item := range items {
				candidate, ok := rssReleaseCandidate(ctx, settings, item, release)
				if ok {
					matchesByMedia[item.ID] = append(matchesByMedia[item.ID], candidate)
				}
			}
		}
		if marker := newestRSSMarker(releases); marker.DownloadURL != nil || marker.GUID != nil || marker.PublishedAt != nil {
			if err := settings.UpdateIndexerRSSMarker(ctx, config.ID, marker); err != nil {
				slog.Error("rss marker update failed", "indexerName", config.Name, "error", err)
			}
		}
	}
	return matchesByMedia, failures
}

func processRSSMatches(
	ctx context.Context,
	settings *storage.SettingsStore,
	downloadClientService *downloadclients.Service,
	decisionEngine decisions.Engine,
	eventBroker *events.Broker,
	matchesByMedia map[uuid.UUID][]storage.ReleaseCandidateInput,
) error {
	clients, err := settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		return fmt.Errorf("list enabled download clients: %w", err)
	}
	for mediaItemID, releases := range matchesByMedia {
		item, err := settings.GetMediaItem(ctx, mediaItemID)
		if err != nil {
			return fmt.Errorf("load media item: %w", err)
		}
		profile, formats, languages := releaseDecisionContext(ctx, settings, item)
		releases = dedupeReleaseCandidates(item, profile, formats, languages, releases)
		sortReleaseCandidates(item, profile, formats, languages, releases)
		if err := settings.MergeReleaseCandidates(ctx, mediaItemID, releases); err != nil {
			return fmt.Errorf("merge RSS release candidates: %w", err)
		}
		releases, err = unblockedReleaseCandidates(ctx, settings, releases)
		if err != nil {
			return err
		}
		available := downloadrouting.ReleaseInputsForClients(releases, clients)
		decision, ok := decisionEngine.ChooseReleaseWithProfileAndLanguages(item, profile, formats, languages, available)
		if !ok {
			continue
		}
		if err := grabReleaseNow(ctx, settings, downloadClientService, eventBroker, clients, item, decision.Release); err != nil {
			return err
		}
	}
	return nil
}
