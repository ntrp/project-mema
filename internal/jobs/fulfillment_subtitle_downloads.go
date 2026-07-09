package jobs

import (
	"context"

	"media-manager/internal/events"
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

func runSubtitleDownloads(
	ctx context.Context,
	settings *storage.SettingsStore,
	service *subtitles.Service,
	eventBroker *events.Broker,
	item storage.MediaItem,
	args FulfillmentActionArgs,
) int {
	count := 0
	for _, target := range item.SubtitleTargets {
		if args.LanguageID != "" && target.LanguageID != args.LanguageID {
			continue
		}
		if err := subtitleSearchDownload(ctx, settings, service, eventBroker, item, SubtitleSearchArgs{
			MediaItemID: item.ID.String(),
			LanguageID:  target.LanguageID,
			FilePath:    args.FilePath,
		}); err == nil {
			count++
		}
	}
	return count
}
