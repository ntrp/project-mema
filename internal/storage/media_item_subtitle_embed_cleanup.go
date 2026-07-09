package storage

import (
	"context"
	"os"

	"github.com/google/uuid"
)

func (s *SettingsStore) RemoveExternalSubtitleAfterEmbed(
	ctx context.Context,
	mediaItemID uuid.UUID,
	subtitleID uuid.UUID,
	filePath string,
) error {
	item, err := s.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return err
	}
	target, err := mediaItemSubtitleTarget(item, filePath)
	if err != nil {
		return err
	}
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, `
delete from app.media_item_subtitles
where media_item_id = $1
	and (id = $2 or file_path = $3)`, mediaItemID, subtitleID, filePath); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
delete from app.media_item_sidecars
where media_item_id = $1
	and file_path = $2`, mediaItemID, filePath); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
