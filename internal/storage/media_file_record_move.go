package storage

import (
	"context"
	"path/filepath"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

type mediaFileMoveHistory struct {
	Operation string
	ActorType string
}

func (s *SettingsStore) RecordContainerRemuxedMediaFile(
	ctx context.Context,
	mediaItemID uuid.UUID,
	source string,
	destination string,
) error {
	return s.moveMediaFileRecord(ctx, mediaItemID, source, destination, mediaFileMoveHistory{
		Operation: "container_remux",
		ActorType: "system",
	})
}

func (s *SettingsStore) moveMediaFileRecord(
	ctx context.Context,
	mediaItemID uuid.UUID,
	source string,
	destination string,
	history mediaFileMoveHistory,
) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	queries := storagegen.New(s.pool).WithTx(tx)
	updated, err := queries.RenameMediaFileRecord(ctx, storagegen.RenameMediaFileRecordParams{
		MediaItemID:     &mediaItemID,
		SourcePath:      source,
		DestinationPath: destination,
		FileName:        filepath.Base(destination),
	})
	if err != nil {
		return err
	}
	if updated == 0 {
		updated, err = renameRelativeMediaFileRecord(ctx, queries, mediaItemID, source, destination)
		if err != nil {
			return err
		}
	}
	if updated == 0 {
		return ErrInvalidInput
	}
	if err := moveMediaFileDerivedRows(ctx, tx, mediaItemID, source, destination); err != nil {
		return err
	}
	if _, err := createMediaFileHistory(ctx, tx, MediaFileHistoryInput{
		MediaItemID:     &mediaItemID,
		FilePath:        destination,
		SourcePath:      &source,
		DestinationPath: &destination,
		Operation:       history.Operation,
		Status:          "succeeded",
		ActorType:       history.ActorType,
	}); err != nil {
		return err
	}
	if err := queries.TouchMediaItem(ctx, mediaItemID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func moveMediaFileDerivedRows(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
	source string,
	destination string,
) error {
	if source == destination {
		return nil
	}
	_, err := q.Exec(ctx, `
delete from app.media_file_facts source_fact
where source_fact.media_item_id = $1
	and source_fact.file_path = $2
	and exists (
		select 1
		from app.media_file_facts destination_fact
		where destination_fact.media_item_id = $1
			and destination_fact.file_path = $3
	)`, mediaItemID, source, destination)
	if err != nil {
		return err
	}
	if _, err := q.Exec(ctx, `
update app.media_file_facts
set file_path = $3,
	updated_at = now()
where media_item_id = $1
	and file_path = $2`, mediaItemID, source, destination); err != nil {
		return err
	}
	if _, err := q.Exec(ctx, `
update app.media_file_tracks
set file_path = $3,
	updated_at = now()
where media_item_id = $1
	and file_path = $2`, mediaItemID, source, destination); err != nil {
		return err
	}
	_, err = q.Exec(ctx, `
update app.media_item_sidecars
set media_file_path = $3,
	updated_at = now()
where media_item_id = $1
	and media_file_path = $2`, mediaItemID, source, destination)
	return err
}
