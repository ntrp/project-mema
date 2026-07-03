package storage

import "context"

func (s *SettingsStore) BulkUpdateIndexers(ctx context.Context, input IndexerBulkUpdateInput) ([]Indexer, error) {
	if len(input.IDs) == 0 {
		return nil, ErrInvalidInput
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for _, id := range input.IDs {
		if _, err := tx.Exec(ctx, `
			update app.indexers
			set enabled = coalesce($2, enabled),
				app_profile_id = coalesce($3, app_profile_id),
				priority = coalesce($4, priority),
				minimum_seeders = coalesce($5, minimum_seeders),
				seed_ratio = coalesce($6, seed_ratio),
				seed_time = coalesce($7, seed_time),
				pack_seed_time = coalesce($8, pack_seed_time),
				prefer_magnet_url = coalesce($9, prefer_magnet_url),
				updated_at = now()
			where id = $1
		`, id, input.Enabled, input.AppProfileID, input.Priority, input.MinimumSeeders, input.SeedRatio,
			input.SeedTime, input.PackSeedTime, input.PreferMagnetURL); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.ListIndexers(ctx)
}
