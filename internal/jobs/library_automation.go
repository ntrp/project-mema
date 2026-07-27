package jobs

import (
	"context"

	"media-manager/internal/storage"
)

func libraryOperationAllowed(ctx context.Context, settings *storage.SettingsStore, item storage.MediaItem, manual bool, path string) (bool, error) {
	if manual {
		return true, nil
	}
	return settings.MediaItemAllowsAutomation(ctx, item, path)
}
