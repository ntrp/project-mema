package jobs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func fulfillmentItems(ctx context.Context, settings *storage.SettingsStore, mediaItemID string) ([]storage.MediaItem, error) {
	if strings.TrimSpace(mediaItemID) == "" {
		return settings.ListMediaItems(ctx)
	}
	id, err := uuid.Parse(mediaItemID)
	if err != nil {
		return nil, fmt.Errorf("parse media item id: %w", err)
	}
	item, err := settings.GetMediaItem(ctx, id)
	if err != nil {
		return nil, err
	}
	return []storage.MediaItem{item}, nil
}

func fulfillmentActionScoped(args FulfillmentActionArgs) bool {
	return strings.TrimSpace(args.MediaItemID) != "" ||
		strings.TrimSpace(args.FilePath) != "" ||
		strings.TrimSpace(args.TargetType) != "" ||
		strings.TrimSpace(args.LanguageID) != "" ||
		strings.TrimSpace(args.TrackID) != "" ||
		strings.TrimSpace(args.OtherFileID) != "" ||
		strings.TrimSpace(args.ExternalSubtitleID) != ""
}

func fmtScopedFulfillmentNotFound(args FulfillmentActionArgs) error {
	return fmt.Errorf("no fulfillment target found for scoped request: mediaItemId=%s filePath=%s trackId=%s otherFileId=%s",
		args.MediaItemID,
		args.FilePath,
		args.TrackID,
		args.OtherFileID,
	)
}
