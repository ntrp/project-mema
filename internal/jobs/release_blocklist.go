package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/riverqueue/river"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

type ReleaseBlocklistCleanupArgs struct{}

func (ReleaseBlocklistCleanupArgs) Kind() string {
	return "release.blocklist_cleanup"
}

type ReleaseBlocklistCleanupWorker struct {
	river.WorkerDefaults[ReleaseBlocklistCleanupArgs]

	settings *storage.SettingsStore
	events   *events.Broker
}

func (w *ReleaseBlocklistCleanupWorker) Work(ctx context.Context, job *river.Job[ReleaseBlocklistCleanupArgs]) (err error) {
	publishJobUpdated(w.events, job.JobRow, "running")
	defer func() { publishJobFinished(w.events, job.JobRow, err) }()

	deleted, err := w.settings.CleanupExpiredReleaseBlocks(ctx)
	if err != nil {
		return fmt.Errorf("cleanup expired release blocks: %w", err)
	}
	if deleted > 0 {
		publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "downloads", "Expired release blocks cleaned up", map[string]any{"deletedCount": deleted})
	}
	return nil
}

func blockAutomaticRelease(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	release storage.ReleaseCandidateInput,
	reason string,
	source string,
) bool {
	expiresAt := automaticBlockExpiry(ctx, settings)
	block, err := settings.BlockReleaseCandidate(ctx, release, reason, source, &expiresAt)
	if err != nil {
		slog.Error("automatic release block failed", "mediaItemId", release.MediaItemID, "releaseTitle", release.Title, "error", err)
		publishSystemEvent(ctx, settings, eventBroker, jobEventError, "downloads", "Automatic release block failed", map[string]any{"mediaItemId": release.MediaItemID.String(), "releaseTitle": release.Title, "error": err.Error()})
		return false
	}
	publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "downloads", "Release temporarily blocklisted", map[string]any{"mediaItemId": release.MediaItemID.String(), "releaseTitle": release.Title, "blockId": block.ID.String(), "expiresAt": expiresAt})
	return true
}

func automaticBlockExpiry(ctx context.Context, settings *storage.SettingsStore) time.Time {
	config, err := settings.GetIndexerSearchSettings(ctx)
	if err != nil || config.AutomaticBlocklistExpiryDays < 1 {
		return time.Now().Add(7 * 24 * time.Hour)
	}
	return time.Now().Add(time.Duration(config.AutomaticBlocklistExpiryDays) * 24 * time.Hour)
}
