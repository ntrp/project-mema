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

type mediaFileMove struct {
	Source      string
	Destination string
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
	return s.moveMediaFileRecordsWithHistory(ctx, mediaItemID, []mediaFileMove{{
		Source: source, Destination: destination,
	}}, nil, history)
}

func (s *SettingsStore) moveMediaFileRecords(
	ctx context.Context,
	mediaItemID uuid.UUID,
	moves []mediaFileMove,
	metadataMoves []mediaFileMove,
) error {
	return s.moveMediaFileRecordsWithHistory(ctx, mediaItemID, moves, metadataMoves, mediaFileMoveHistory{
		Operation: "moved",
		ActorType: "user",
	})
}

func (s *SettingsStore) moveMediaFileRecordsWithHistory(
	ctx context.Context,
	mediaItemID uuid.UUID,
	moves []mediaFileMove,
	metadataMoves []mediaFileMove,
	history mediaFileMoveHistory,
) error {
	if len(moves) == 0 && len(metadataMoves) == 0 {
		return nil
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	queries := storagegen.New(tx)
	for _, move := range moves {
		if err := moveMediaFileRecordInTx(ctx, tx, queries, mediaItemID, move, history); err != nil {
			return err
		}
	}
	for _, move := range metadataMoves {
		if _, err := tx.Exec(ctx, `
update app.media_item_sidecars
set file_path = $3,
	updated_at = now()
where media_item_id = $1
	and file_path = $2`, mediaItemID, move.Source, move.Destination); err != nil {
			return err
		}
	}
	if err := queries.TouchMediaItem(ctx, mediaItemID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func moveMediaFileRecordInTx(
	ctx context.Context,
	q storagegen.DBTX,
	queries *storagegen.Queries,
	mediaItemID uuid.UUID,
	move mediaFileMove,
	history mediaFileMoveHistory,
) error {
	updated, err := queries.RenameMediaFileRecord(ctx, storagegen.RenameMediaFileRecordParams{
		MediaItemID:     &mediaItemID,
		SourcePath:      move.Source,
		DestinationPath: move.Destination,
		FileName:        filepath.Base(move.Destination),
	})
	if err != nil {
		return err
	}
	if updated == 0 {
		updated, err = renameRelativeMediaFileRecord(ctx, queries, mediaItemID, move.Source, move.Destination)
		if err != nil {
			return err
		}
	}
	if updated == 0 {
		return ErrInvalidInput
	}
	if err := moveMediaFileDerivedRows(ctx, q, mediaItemID, move.Source, move.Destination); err != nil {
		return err
	}
	_, err = createMediaFileHistory(ctx, q, MediaFileHistoryInput{
		MediaItemID:     &mediaItemID,
		FilePath:        move.Destination,
		SourcePath:      &move.Source,
		DestinationPath: &move.Destination,
		Operation:       history.Operation,
		Status:          "succeeded",
		ActorType:       history.ActorType,
	})
	return err
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
