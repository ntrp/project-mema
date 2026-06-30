package storage

import (
	"context"

	"github.com/google/uuid"
)

func (s *SettingsStore) UpdateMediaItemManual(ctx context.Context, id uuid.UUID, manual bool) (MediaItem, error) {
	tag, err := s.pool.Exec(ctx, `
		update app.media_items
		set manual = $2,
			updated_at = now()
		where id = $1
	`, id, manual)
	if err != nil {
		return MediaItem{}, err
	}
	if tag.RowsAffected() == 0 {
		return MediaItem{}, ErrNotFound
	}
	return s.GetMediaItem(ctx, id)
}
