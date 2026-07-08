package storage

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

type MockSubtitleProviderRow struct {
	ID         uuid.UUID
	ProviderID uuid.UUID
	Title      string
	LanguageID string
	Format     string
	SortOrder  int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MockSubtitleProviderRowInput struct {
	Title      string
	LanguageID string
	Format     string
}

func subtitleProviderWithRows(
	ctx context.Context,
	q storagegen.DBTX,
	row storagegen.AppSubtitleProvider,
) (SubtitleProvider, error) {
	provider := subtitleProviderFromRow(row)
	mockRows, err := storagegen.New(q).ListMockSubtitleProviderRows(ctx, provider.ID)
	if err != nil {
		return SubtitleProvider{}, err
	}
	provider.MockSubtitles = mockSubtitleProviderRowsFromRows(mockRows)
	return provider, nil
}

func replaceMockSubtitleProviderRows(
	ctx context.Context,
	q *storagegen.Queries,
	providerID uuid.UUID,
	inputs []MockSubtitleProviderRowInput,
) ([]MockSubtitleProviderRow, error) {
	if err := q.DeleteMockSubtitleProviderRows(ctx, providerID); err != nil {
		return nil, err
	}
	rows := make([]MockSubtitleProviderRow, 0, len(inputs))
	for index, input := range inputs {
		row, err := q.CreateMockSubtitleProviderRow(ctx, storagegen.CreateMockSubtitleProviderRowParams{
			ID:         uuid.New(),
			ProviderID: providerID,
			Title:      strings.TrimSpace(input.Title),
			LanguageID: strings.TrimSpace(input.LanguageID),
			Format:     normalizeMockSubtitleFormat(input.Format),
			SortOrder:  int32(index),
		})
		if err != nil {
			return nil, err
		}
		rows = append(rows, mockSubtitleProviderRowFromRow(row))
	}
	return rows, nil
}

func mockSubtitleProviderRowsFromRows(rows []storagegen.AppMockSubtitleProviderRow) []MockSubtitleProviderRow {
	items := make([]MockSubtitleProviderRow, 0, len(rows))
	for _, row := range rows {
		items = append(items, mockSubtitleProviderRowFromRow(row))
	}
	return items
}

func mockSubtitleProviderRowFromRow(row storagegen.AppMockSubtitleProviderRow) MockSubtitleProviderRow {
	return MockSubtitleProviderRow{
		ID:         row.ID,
		ProviderID: row.ProviderID,
		Title:      row.Title,
		LanguageID: row.LanguageID,
		Format:     row.Format,
		SortOrder:  row.SortOrder,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}

func normalizeMockSubtitleFormat(value string) string {
	switch strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".") {
	case "", "srt", "subrip":
		return "subrip"
	case "vtt", "webvtt":
		return "vtt"
	case "ass":
		return "ass"
	case "ssa":
		return "ssa"
	default:
		return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".")
	}
}
