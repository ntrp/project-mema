package storage

import (
	"context"

	"github.com/google/uuid"

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
	return storagegen.UpsertMediaItemSubtitleParams{
		ID:           uuid.New(),
		MediaItemID:  input.MediaItemID,
		SeasonID:     input.SeasonID,
		EpisodeID:    input.EpisodeID,
		ProviderID:   input.ProviderID,
		ProviderName: input.ProviderName,
		LanguageID:   input.LanguageID,
		FilePath:     input.FilePath,
		SourceUrl:    textValue(input.SourceURL),
		ReleaseName:  textValue(input.ReleaseName),
	}
}

func mediaItemSubtitleFromRow(row storagegen.AppMediaItemSubtitle) MediaItemSubtitle {
	return MediaItemSubtitle{
		ID:           row.ID,
		MediaItemID:  row.MediaItemID,
		SeasonID:     row.SeasonID,
		EpisodeID:    row.EpisodeID,
		ProviderID:   row.ProviderID,
		ProviderName: row.ProviderName,
		LanguageID:   row.LanguageID,
		FilePath:     row.FilePath,
		SourceURL:    textPtr(row.SourceUrl),
		ReleaseName:  textPtr(row.ReleaseName),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func subtitleFilePaths(subtitles []MediaItemSubtitle) []string {
	paths := make([]string, 0, len(subtitles))
	for _, subtitle := range subtitles {
		paths = append(paths, subtitle.FilePath)
	}
	return paths
}
