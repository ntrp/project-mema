package jobs

import (
	"context"
	"log/slog"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/storage"
)

func grabAutomaticTargets(
	ctx context.Context,
	settings *storage.SettingsStore,
	downloadClientService *downloadclients.Service,
	decisionEngine decisions.Engine,
	eventBroker *events.Broker,
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	languages []storage.Language,
	clients []storage.DownloadClient,
	targets []automaticReleaseTarget,
) error {
	grabbed := 0
	for _, target := range targets {
		candidates := append([]storage.ReleaseCandidateInput(nil), target.candidates...)
		for attempt := 1; attempt <= maxAutomaticGrabAttempts; attempt++ {
			decision, ok := decisionEngine.ChooseReleaseWithProfileAndLanguages(item, profile, formats, languages, candidates)
			if !ok {
				publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search found no acceptable release for target", map[string]any{
					"mediaItemId": item.ID.String(), "title": item.Title, "target": target.label,
				})
				break
			}
			err := grabReleaseNow(ctx, settings, downloadClientService, eventBroker, clients, item, decision.Release)
			if err == nil {
				grabbed++
				slog.Debug("auto search grab queued", "mediaItemId", item.ID, "title", item.Title, "target", target.label, "releaseTitle", decision.Release.Title)
				publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "jobs", "Automatic search queued download", map[string]any{
					"mediaItemId": item.ID.String(), "title": item.Title, "target": target.label, "releaseTitle": decision.Release.Title,
				})
				break
			}
			slog.Error("auto search grab failed", "mediaItemId", item.ID, "title", item.Title, "target", target.label, "releaseTitle", decision.Release.Title, "error", err)
			if shouldRetryAlternativeRelease(err, attempt) {
				candidates = removeAutomaticCandidate(candidates, decision.Release)
				publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search trying alternative release", map[string]any{
					"mediaItemId": item.ID.String(), "title": item.Title, "target": target.label, "failedReleaseTitle": decision.Release.Title, "attempt": attempt,
				})
				continue
			}
			if automaticRetryLimitReached(err, attempt) {
				publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "jobs", "Automatic search retry limit reached", map[string]any{
					"mediaItemId": item.ID.String(), "title": item.Title, "target": target.label, "attempts": attempt,
				})
				break
			}
			publishSystemEvent(ctx, settings, eventBroker, jobEventError, "jobs", "Automatic search grab failed", map[string]any{
				"mediaItemId": item.ID.String(), "title": item.Title, "target": target.label, "releaseTitle": decision.Release.Title, "error": err.Error(),
			})
			return err
		}
	}
	if grabbed > 0 {
		publishSystemEvent(ctx, settings, eventBroker, jobEventInfo, "jobs", "Automatic search queued downloads", map[string]any{
			"mediaItemId": item.ID.String(), "title": item.Title, "downloadCount": grabbed,
		})
	}
	return nil
}

func removeAutomaticCandidate(candidates []storage.ReleaseCandidateInput, failed storage.ReleaseCandidateInput) []storage.ReleaseCandidateInput {
	filtered := make([]storage.ReleaseCandidateInput, 0, len(candidates)-1)
	for _, candidate := range candidates {
		if candidate.Title == failed.Title && candidate.DownloadURL == failed.DownloadURL && candidate.IndexerName == failed.IndexerName {
			continue
		}
		filtered = append(filtered, candidate)
	}
	return filtered
}
