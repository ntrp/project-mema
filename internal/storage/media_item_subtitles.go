package storage

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) ListMediaItemSubtitles(
	ctx context.Context,
	mediaItemID uuid.UUID,
) ([]MediaItemSubtitle, error) {
	return listMediaItemSubtitles(ctx, s.pool, mediaItemID)
}

func (s *SettingsStore) UpsertMediaItemSubtitle(
	ctx context.Context,
	input MediaItemSubtitleInput,
) (MediaItemSubtitle, error) {
	row, err := storagegen.New(s.pool).UpsertMediaItemSubtitle(ctx, subtitleParams(input))
	return mediaItemSubtitleFromRow(row), err
}

func (s *SettingsStore) DeleteMediaItemSubtitle(
	ctx context.Context,
	mediaItemID uuid.UUID,
	subtitleID uuid.UUID,
) (MediaItem, error) {
	item, err := s.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return MediaItem{}, err
	}
	subtitle, err := storagegen.New(s.pool).GetMediaItemSubtitle(ctx, storagegen.GetMediaItemSubtitleParams{
		MediaItemID: mediaItemID,
		ID:          subtitleID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, ErrNotFound
	}
	if err != nil {
		return MediaItem{}, err
	}
	target, err := mediaItemSubtitleTarget(item, subtitle.FilePath)
	if err != nil {
		return MediaItem{}, err
	}
	result := s.applyFileDeletePolicy(ctx, item, target)
	if err := s.recordFileDeletePolicy(ctx, item.ID, result); err != nil {
		return MediaItem{}, err
	}
	if result.Status == "failed" {
		return MediaItem{}, ErrInvalidInput
	}
	if _, err := storagegen.New(s.pool).DeleteMediaItemSubtitle(ctx, storagegen.DeleteMediaItemSubtitleParams{
		MediaItemID: mediaItemID,
		ID:          subtitleID,
	}); err != nil {
		return MediaItem{}, err
	}
	return s.GetMediaItem(ctx, mediaItemID)
}

func listMediaItemSubtitles(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
) ([]MediaItemSubtitle, error) {
	rows, err := storagegen.New(q).ListMediaItemSubtitles(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	items := make([]MediaItemSubtitle, 0, len(rows))
	for _, row := range rows {
		items = append(items, mediaItemSubtitleFromRow(row))
	}
	return items, nil
}

func subtitleParams(input MediaItemSubtitleInput) storagegen.UpsertMediaItemSubtitleParams {
	format := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(input.Format)), ".")
	if format == "" {
		format = strings.TrimPrefix(strings.ToLower(filepath.Ext(input.FilePath)), ".")
	}
	if format == "" {
		format = "srt"
	}
	downloadedAt := input.DownloadedAt
	if downloadedAt.IsZero() {
		downloadedAt = time.Now().UTC()
	}
	return storagegen.UpsertMediaItemSubtitleParams{
		ID:                 uuid.New(),
		MediaItemID:        input.MediaItemID,
		SeasonID:           input.SeasonID,
		EpisodeID:          input.EpisodeID,
		ProviderID:         input.ProviderID,
		ProviderName:       input.ProviderName,
		LanguageID:         input.LanguageID,
		Format:             format,
		FilePath:           input.FilePath,
		SourceUrl:          textValue(input.SourceURL),
		SourceReference:    textValue(input.SourceRef),
		ReleaseName:        textValue(input.ReleaseName),
		ProviderSubtitleID: textValue(input.ProviderSubtitleID),
		Checksum:           textValue(input.Checksum),
		SizeBytes:          int8Value(input.SizeBytes),
		DownloadedAt:       downloadedAt,
	}
}

func mediaItemSubtitleFromRow(row storagegen.AppMediaItemSubtitle) MediaItemSubtitle {
	return MediaItemSubtitle{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		SeasonID:           row.SeasonID,
		EpisodeID:          row.EpisodeID,
		ProviderID:         row.ProviderID,
		ProviderName:       row.ProviderName,
		LanguageID:         row.LanguageID,
		Format:             row.Format,
		FilePath:           row.FilePath,
		SourceURL:          textPtr(row.SourceUrl),
		SourceRef:          textPtr(row.SourceReference),
		ReleaseName:        textPtr(row.ReleaseName),
		ProviderSubtitleID: textPtr(row.ProviderSubtitleID),
		Checksum:           textPtr(row.Checksum),
		SizeBytes:          int8Ptr(row.SizeBytes),
		DownloadedAt:       row.DownloadedAt,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}

func subtitleFilePaths(subtitles []MediaItemSubtitle) []string {
	paths := make([]string, 0, len(subtitles))
	for _, subtitle := range subtitles {
		paths = append(paths, subtitle.FilePath)
	}
	return paths
}

func mediaItemSubtitleTarget(item MediaItem, filePath string) (string, error) {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return "", ErrInvalidInput
	}
	return safePathUnderRoot(*item.MediaFolderPath, filePath, false)
}
