package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) ListMediaComponentSources(
	ctx context.Context,
	mediaItemID uuid.UUID,
) ([]MediaComponentSource, error) {
	return listMediaComponentSources(ctx, s.pool, mediaItemID)
}

func (s *SettingsStore) GetMediaComponentSource(
	ctx context.Context,
	mediaItemID uuid.UUID,
	sourceID uuid.UUID,
) (MediaComponentSource, error) {
	row, err := storagegen.New(s.pool).GetMediaComponentSource(ctx, storagegen.GetMediaComponentSourceParams{
		MediaItemID: mediaItemID,
		ID:          sourceID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaComponentSource{}, ErrNotFound
	}
	if err != nil {
		return MediaComponentSource{}, err
	}
	return mediaComponentSourceFromRow(row), nil
}

func (s *SettingsStore) RetainMediaComponentSource(
	ctx context.Context,
	mediaItemID uuid.UUID,
	input MediaComponentSourceInput,
) (MediaComponentSource, error) {
	item, err := s.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return MediaComponentSource{}, err
	}
	source, target, id, size, err := retainedComponentSourceTarget(item, input.SourceFilePath)
	if err != nil {
		return MediaComponentSource{}, err
	}
	if err := copyFile(source, target); err != nil {
		return MediaComponentSource{}, err
	}
	row, err := storagegen.New(s.pool).CreateMediaComponentSource(ctx, storagegen.CreateMediaComponentSourceParams{
		ID:              id,
		MediaItemID:     mediaItemID,
		SourceRole:      normalizeComponentSourceRole(input.SourceRole),
		SourceFilePath:  source,
		RetainedPath:    target,
		ReleaseTitle:    textValue(input.ReleaseTitle),
		ReleaseGroup:    textValue(input.ReleaseGroup),
		ReleaseName:     textValue(input.ReleaseName),
		ReleaseID:       textValue(input.ReleaseID),
		SourceMetadata:  textValue(input.SourceMetadata),
		StreamInventory: strings.TrimSpace(input.StreamInventory),
		Checksum:        textValue(input.Checksum),
		SizeBytes:       int8Value(&size),
	})
	if err != nil {
		_ = os.Remove(target)
		return MediaComponentSource{}, err
	}
	return mediaComponentSourceFromRow(row), nil
}

func (s *SettingsStore) ReleaseMediaComponentSource(
	ctx context.Context,
	mediaItemID uuid.UUID,
	sourceID uuid.UUID,
) (MediaComponentSource, error) {
	item, err := s.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return MediaComponentSource{}, err
	}
	source, err := s.GetMediaComponentSource(ctx, mediaItemID, sourceID)
	if err != nil {
		return MediaComponentSource{}, err
	}
	target, err := mediaComponentSourceTarget(item, source.RetainedPath)
	if err != nil {
		return MediaComponentSource{}, err
	}
	result := s.applyFileDeletePolicy(ctx, item, target)
	if err := s.recordFileDeletePolicy(ctx, item.ID, result); err != nil {
		return MediaComponentSource{}, err
	}
	if result.Status == "failed" {
		return MediaComponentSource{}, ErrInvalidInput
	}
	row, err := storagegen.New(s.pool).ReleaseMediaComponentSource(ctx, storagegen.ReleaseMediaComponentSourceParams{
		MediaItemID: mediaItemID,
		ID:          sourceID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaComponentSource{}, ErrNotFound
	}
	if err != nil {
		return MediaComponentSource{}, err
	}
	return mediaComponentSourceFromRow(row), nil
}

func listMediaComponentSources(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
) ([]MediaComponentSource, error) {
	rows, err := storagegen.New(q).ListMediaComponentSources(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	sources := make([]MediaComponentSource, 0, len(rows))
	for _, row := range rows {
		sources = append(sources, mediaComponentSourceFromRow(row))
	}
	return sources, nil
}

func retainedComponentSourceTarget(
	item MediaItem,
	sourcePath string,
) (string, string, uuid.UUID, int64, error) {
	source, err := mediaComponentSourceTarget(item, sourcePath)
	if err != nil {
		return "", "", uuid.Nil, 0, err
	}
	info, err := os.Stat(source)
	if err != nil || info.IsDir() {
		return "", "", uuid.Nil, 0, ErrInvalidInput
	}
	id := uuid.New()
	target, err := mediaComponentSourceTarget(
		item,
		filepath.Join(".mema", "component-sources", id.String(), filepath.Base(source)),
	)
	if err != nil {
		return "", "", uuid.Nil, 0, err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", "", uuid.Nil, 0, err
	}
	return source, target, id, info.Size(), nil
}

func mediaComponentSourceTarget(item MediaItem, path string) (string, error) {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return "", ErrInvalidInput
	}
	return safePathUnderRoot(*item.MediaFolderPath, path, false)
}

func normalizeComponentSourceRole(value string) string {
	switch strings.TrimSpace(value) {
	case "baseVideo", "audio", "subtitle":
		return strings.TrimSpace(value)
	default:
		return "other"
	}
}

func mediaComponentSourceFromRow(row storagegen.AppMediaComponentSource) MediaComponentSource {
	return MediaComponentSource{
		ID:              row.ID,
		MediaItemID:     row.MediaItemID,
		SourceRole:      row.SourceRole,
		SourceFilePath:  row.SourceFilePath,
		RetainedPath:    row.RetainedPath,
		ReleaseTitle:    textPtr(row.ReleaseTitle),
		ReleaseGroup:    textPtr(row.ReleaseGroup),
		ReleaseName:     textPtr(row.ReleaseName),
		ReleaseID:       textPtr(row.ReleaseID),
		SourceMetadata:  textPtr(row.SourceMetadata),
		StreamInventory: row.StreamInventory,
		Checksum:        textPtr(row.Checksum),
		SizeBytes:       int8Ptr(row.SizeBytes),
		RetentionState:  row.RetentionState,
		RetainedAt:      row.RetainedAt,
		ReleasedAt:      row.ReleasedAt,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}
