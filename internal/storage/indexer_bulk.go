package storage

import (
	"context"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) BulkUpdateIndexers(ctx context.Context, input IndexerBulkUpdateInput) ([]Indexer, error) {
	if len(input.IDs) == 0 {
		return nil, ErrInvalidInput
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	queries := storagegen.New(s.pool).WithTx(tx)
	for _, id := range input.IDs {
		if err := queries.BulkUpdateIndexer(ctx, bulkUpdateIndexerParams(id, input)); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.ListIndexers(ctx)
}

func bulkUpdateIndexerParams(id uuid.UUID, input IndexerBulkUpdateInput) storagegen.BulkUpdateIndexerParams {
	return storagegen.BulkUpdateIndexerParams{
		ID:              id,
		Enabled:         boolValue(input.Enabled),
		AppProfileID:    textValue(input.AppProfileID),
		Priority:        int4Value(input.Priority),
		MinimumSeeders:  int4Value(input.MinimumSeeders),
		SeedRatio:       input.SeedRatio,
		SeedTime:        int4Value(input.SeedTime),
		PackSeedTime:    int4Value(input.PackSeedTime),
		PreferMagnetUrl: boolValue(input.PreferMagnetURL),
	}
}
