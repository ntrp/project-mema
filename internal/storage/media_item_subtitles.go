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
	return upsertMediaItemSubtitle(ctx, s.pool, input)
}

func upsertMediaItemSubtitle(
	ctx context.Context,
	q storagegen.DBTX,
	input MediaItemSubtitleInput,
) (MediaItemSubtitle, error) {
	row, err := storagegen.New(q).UpsertMediaItemSubtitle(ctx, subtitleParams(input))
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

func (s *SettingsStore) UpdateMediaItemSubtitleSelection(
	ctx context.Context,
	mediaItemID uuid.UUID,
	subtitleID uuid.UUID,
	input MediaItemSubtitleSelectionInput,
) (MediaItem, error) {
	mode, err := normalizeSubtitleRetention(input.RetentionMode)
	if err != nil {
		return MediaItem{}, err
	}
	q := storagegen.New(s.pool)
	subtitle, err := q.GetMediaItemSubtitle(ctx, storagegen.GetMediaItemSubtitleParams{
		MediaItemID: mediaItemID,
		ID:          subtitleID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, ErrNotFound
	}
	if err != nil {
		return MediaItem{}, err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaItem{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	txq := storagegen.New(tx)
	if input.Selected {
		err = txq.ClearSelectedMediaItemSubtitles(ctx, storagegen.ClearSelectedMediaItemSubtitlesParams{
			MediaItemID: mediaItemID,
			LanguageID:  subtitle.LanguageID,
			FilePath:    subtitle.FilePath,
			ID:          subtitleID,
		})
		if err != nil {
			return MediaItem{}, err
		}
	}
	_, err = txq.UpdateMediaItemSubtitleSelection(ctx, storagegen.UpdateMediaItemSubtitleSelectionParams{
		MediaItemID:   mediaItemID,
		ID:            subtitleID,
		Selected:      input.Selected,
		RetentionMode: string(mode),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, ErrNotFound
	}
	if err != nil {
		return MediaItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaItem{}, err
	}
	return s.GetMediaItem(ctx, mediaItemID)
}

func (s *SettingsStore) ListSelectedSubtitleArtifacts(
	ctx context.Context,
	mediaItemID uuid.UUID,
) ([]SubtitleAssemblyArtifact, error) {
	item, err := s.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	return SelectedSubtitleArtifacts(item), nil
}

func SelectedSubtitleArtifacts(item MediaItem) []SubtitleAssemblyArtifact {
	artifacts := []SubtitleAssemblyArtifact{}
	for _, subtitle := range item.ExternalSubtitles {
		if !subtitle.Selected || subtitle.RetentionMode != SubtitleRetentionMux {
			continue
		}
		artifacts = append(artifacts, SubtitleAssemblyArtifact{
			ID:                 subtitle.ID,
			MediaItemID:        subtitle.MediaItemID,
			LanguageID:         subtitle.LanguageID,
			Format:             subtitle.Format,
			FilePath:           subtitle.FilePath,
			RetentionMode:      subtitle.RetentionMode,
			ProviderName:       subtitle.ProviderName,
			SourceURL:          subtitle.SourceURL,
			SourceRef:          subtitle.SourceRef,
			ProviderSubtitleID: subtitle.ProviderSubtitleID,
			Checksum:           subtitle.Checksum,
			SizeBytes:          subtitle.SizeBytes,
			DownloadedAt:       subtitle.DownloadedAt,
		})
	}
	return artifacts
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
	selected := true
	if input.Selected != nil {
		selected = *input.Selected
	}
	retentionMode, err := normalizeSubtitleRetention(input.RetentionMode)
	if err != nil {
		retentionMode = SubtitleRetentionExternal
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
		Selected:           selected,
		RetentionMode:      string(retentionMode),
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
		Selected:           row.Selected,
		RetentionMode:      SubtitleRetentionMode(row.RetentionMode),
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}

func normalizeSubtitleRetention(value SubtitleRetentionMode) (SubtitleRetentionMode, error) {
	if value == "" {
		return SubtitleRetentionExternal, nil
	}
	switch value {
	case SubtitleRetentionExternal, SubtitleRetentionMux, SubtitleRetentionIgnore:
		return value, nil
	default:
		return "", ErrInvalidInput
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
