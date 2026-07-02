package jobs

import (
	"context"
	"strings"

	"media-manager/internal/decisions"
	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func SearchManualReleases(
	ctx context.Context,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	item storage.MediaItem,
	query string,
	eventBroker *events.Broker,
	progress ReleaseSearchProgress,
) ([]storage.ReleaseCandidateInput, []string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		query = decisions.SearchQueryForMediaItem(item)
	}
	publishReleaseSearchProgress(progress, "Starting release search for %s", item.Title)
	return searchReleasesWithProgress(ctx, settings, indexerService, item, query, eventBroker, true, progress)
}
