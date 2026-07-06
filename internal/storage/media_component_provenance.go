package storage

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) ListMediaComponentProvenance(
	ctx context.Context,
	mediaItemID uuid.UUID,
) ([]MediaComponentProvenance, error) {
	rows, err := storagegen.New(s.pool).ListMediaComponentProvenance(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	items := make([]MediaComponentProvenance, 0, len(rows))
	for _, row := range rows {
		items = append(items, mediaComponentProvenanceFromRow(row))
	}
	return items, nil
}

func (s *SettingsStore) UpsertMediaComponentProvenance(
	ctx context.Context,
	input MediaComponentProvenanceInput,
) (MediaComponentProvenance, error) {
	row, err := storagegen.New(s.pool).UpsertMediaComponentProvenance(ctx, mediaComponentProvenanceParams(input))
	return mediaComponentProvenanceFromRow(row), err
}

func mediaComponentProvenanceParams(
	input MediaComponentProvenanceInput,
) storagegen.UpsertMediaComponentProvenanceParams {
	return storagegen.UpsertMediaComponentProvenanceParams{
		ID:                  uuid.New(),
		MediaItemID:         input.MediaItemID,
		ComponentType:       strings.TrimSpace(input.ComponentType),
		ComponentKey:        strings.TrimSpace(input.ComponentKey),
		ReleaseGroup:        strings.TrimSpace(input.ReleaseGroup),
		ReleaseName:         strings.TrimSpace(input.ReleaseName),
		ReleaseID:           textValue(input.ReleaseID),
		SourceProvider:      textValue(input.SourceProvider),
		SourceFilePath:      textValue(input.SourceFilePath),
		RetainedSourceID:    input.RetainedSourceID,
		SourceStreamID:      int4Value(input.SourceStreamID),
		TransformationChain: jsonArray(input.TransformationChain),
	}
}

func mediaComponentProvenanceFromRow(row storagegen.AppMediaComponentProvenance) MediaComponentProvenance {
	return MediaComponentProvenance{
		ID:                  row.ID,
		MediaItemID:         row.MediaItemID,
		ComponentType:       row.ComponentType,
		ComponentKey:        row.ComponentKey,
		ReleaseGroup:        row.ReleaseGroup,
		ReleaseName:         row.ReleaseName,
		ReleaseID:           textPtr(row.ReleaseID),
		SourceProvider:      textPtr(row.SourceProvider),
		SourceFilePath:      textPtr(row.SourceFilePath),
		RetainedSourceID:    row.RetainedSourceID,
		SourceStreamID:      int4Ptr(row.SourceStreamID),
		TransformationChain: jsonArrayMap(row.TransformationChain),
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
}

func jsonArray(source []map[string]any) []byte {
	if source == nil {
		source = []map[string]any{}
	}
	payload, _ := json.Marshal(source)
	return payload
}

func jsonArrayMap(payload []byte) []map[string]any {
	items := []map[string]any{}
	if len(payload) > 0 {
		_ = json.Unmarshal(payload, &items)
	}
	return items
}
