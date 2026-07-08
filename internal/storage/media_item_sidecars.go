package storage

import (
	"context"
	"strings"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

func listMediaItemSidecars(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
) ([]MediaItemSidecar, error) {
	rows, err := storagegen.New(q).ListMediaItemSidecars(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	items := make([]MediaItemSidecar, 0, len(rows))
	for _, row := range rows {
		items = append(items, mediaItemSidecarFromRow(row))
	}
	return items, nil
}

func upsertMediaItemSidecar(
	ctx context.Context,
	q storagegen.DBTX,
	input MediaItemSidecarInput,
) (MediaItemSidecar, error) {
	row, err := storagegen.New(q).UpsertMediaItemSidecar(ctx, storagegen.UpsertMediaItemSidecarParams{
		ID:            uuid.New(),
		MediaItemID:   input.MediaItemID,
		MediaFilePath: strings.TrimSpace(input.MediaFilePath),
		FilePath:      strings.TrimSpace(input.FilePath),
		SidecarType:   string(input.SidecarType),
		Subtype:       textValue(optionalTrimmedValue(input.Subtype)),
		LanguageID:    textValue(optionalTrimmedValue(input.LanguageID)),
		Format:        textValue(optionalTrimmedValue(input.Format)),
	})
	return mediaItemSidecarFromRow(row), err
}

func mediaItemSidecarFromRow(row storagegen.AppMediaItemSidecar) MediaItemSidecar {
	return MediaItemSidecar{
		ID:            row.ID,
		MediaItemID:   row.MediaItemID,
		MediaFilePath: row.MediaFilePath,
		FilePath:      row.FilePath,
		SidecarType:   MediaSidecarType(row.SidecarType),
		Subtype:       textPtr(row.Subtype),
		LanguageID:    textPtr(row.LanguageID),
		Format:        textPtr(row.Format),
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func optionalTrimmedValue(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
