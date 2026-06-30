package jobs

import (
	"context"
	"log/slog"

	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func recordIndexerSearchFailure(
	ctx context.Context,
	settings *storage.SettingsStore,
	indexer storage.Indexer,
	err error,
) {
	statusCode := indexers.StatusCode(err)
	permanent := indexers.IsPermanentFailure(statusCode)
	if _, updateErr := settings.RecordIndexerFailure(
		ctx,
		indexer.ID,
		statusCode,
		err.Error(),
		permanent,
	); updateErr != nil {
		slog.Error("indexer failure state update failed", "indexerName", indexer.Name, "error", updateErr)
	}
}
