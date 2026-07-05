package storage

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) CreateMediaFileHistory(
	ctx context.Context,
	input MediaFileHistoryInput,
) (MediaFileHistoryEntry, error) {
	return createMediaFileHistory(ctx, s.pool, input)
}

func (s *SettingsStore) ListMediaFileHistory(
	ctx context.Context,
	mediaItemID uuid.UUID,
) ([]MediaFileHistoryEntry, error) {
	rows, err := storagegen.New(s.pool).ListMediaFileHistory(ctx, &mediaItemID)
	if err != nil {
		return nil, err
	}
	items := make([]MediaFileHistoryEntry, 0, len(rows))
	for _, row := range rows {
		entry, err := mediaFileHistoryFromRow(row)
		if err != nil {
			return nil, err
		}
		items = append(items, entry)
	}
	return items, nil
}

func createMediaFileHistory(
	ctx context.Context,
	q storagegen.DBTX,
	input MediaFileHistoryInput,
) (MediaFileHistoryEntry, error) {
	details := input.Details
	if details == nil {
		details = map[string]any{}
	}
	payload, err := json.Marshal(details)
	if err != nil {
		return MediaFileHistoryEntry{}, err
	}
	params := storagegen.CreateMediaFileHistoryParams{
		ID:              uuid.New(),
		MediaItemID:     input.MediaItemID,
		FilePath:        input.FilePath,
		SourcePath:      textValue(input.SourcePath),
		DestinationPath: textValue(input.DestinationPath),
		Operation:       input.Operation,
		Status:          input.Status,
		ActorType:       defaultActorType(input.ActorType),
		ActorID:         textValue(input.ActorID),
		JobID:           textValue(input.JobID),
		Details:         payload,
		FailureDetails:  textValue(input.FailureDetails),
	}
	row, err := storagegen.New(q).CreateMediaFileHistory(ctx, params)
	if err != nil {
		return MediaFileHistoryEntry{}, err
	}
	return mediaFileHistoryFromRow(row)
}

func mediaFileHistoryFromRow(row storagegen.AppMediaFileHistory) (MediaFileHistoryEntry, error) {
	details := map[string]any{}
	if len(row.Details) > 0 {
		if err := json.Unmarshal(row.Details, &details); err != nil {
			return MediaFileHistoryEntry{}, err
		}
	}
	return MediaFileHistoryEntry{
		ID:              row.ID,
		MediaItemID:     row.MediaItemID,
		FilePath:        row.FilePath,
		SourcePath:      textPtr(row.SourcePath),
		DestinationPath: textPtr(row.DestinationPath),
		Operation:       row.Operation,
		Status:          row.Status,
		ActorType:       row.ActorType,
		ActorID:         textPtr(row.ActorID),
		JobID:           textPtr(row.JobID),
		Details:         details,
		FailureDetails:  textPtr(row.FailureDetails),
		CreatedAt:       row.CreatedAt,
	}, nil
}

func defaultActorType(value string) string {
	if value == "" {
		return "system"
	}
	return value
}
