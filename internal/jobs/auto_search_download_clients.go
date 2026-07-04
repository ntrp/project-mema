package jobs

import (
	"context"
	"fmt"

	"media-manager/internal/downloadrouting"
	"media-manager/internal/events"
	"media-manager/internal/storage"
)

func autoSearchCandidatesWithDownloadClient(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	releases []storage.ReleaseCandidateInput,
	searchErrors []string,
) ([]storage.ReleaseCandidateInput, []storage.DownloadClient, bool, error) {
	clients, err := settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "downloads", "Automatic grab failed to list clients", map[string]any{"mediaItemId": item.ID.String(), "error": err.Error()})
		return nil, nil, false, fmt.Errorf("list enabled download clients: %w", err)
	}
	if len(clients) == 0 {
		message := downloadrouting.MissingClientMessage("")
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "downloads", "Automatic grab skipped because no download client is enabled", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title})
		return nil, nil, false, settings.ReplaceReleaseSearchResults(ctx, item.ID, releases, append(searchErrors, message))
	}
	releases, err = unblockedReleaseCandidates(ctx, settings, releases)
	if err != nil {
		return nil, nil, false, err
	}
	availableReleases := downloadrouting.ReleaseInputsForClients(releases, clients)
	if len(availableReleases) == 0 && len(releases) > 0 {
		message := "No enabled download client matches the release protocols"
		publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search found no compatible download client", map[string]any{"mediaItemId": item.ID.String(), "title": item.Title, "releaseCount": len(releases)})
		return nil, nil, false, settings.ReplaceReleaseSearchResults(ctx, item.ID, releases, append(searchErrors, message))
	}
	return availableReleases, clients, true, nil
}
