package storage

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

func (s *SettingsStore) UpdateMediaItemOptions(ctx context.Context, id uuid.UUID, input MediaItemOptionsInput) (MediaItem, error) {
	qualityProfileID := optionalTrimmed(input.QualityProfileID)
	minimumAvailability := optionalMinimumAvailability(input.MinimumAvailability)
	monitorMode := optionalTrimmed(input.MonitorMode)
	mediaFolderPath, err := s.updatedMediaFolderPath(ctx, id, input.LibraryFolderID)
	if err != nil {
		return MediaItem{}, err
	}
	seasonsPayload, updateSeasons, err := mediaItemSeasonsPayload(input.Seasons)
	if err != nil {
		return MediaItem{}, err
	}
	tag, err := s.pool.Exec(ctx, `
		update app.media_items
		set quality_profile_id = coalesce($2, quality_profile_id),
			minimum_availability = coalesce($3, minimum_availability),
			monitored = coalesce($4, monitored),
			monitor_mode = coalesce($5, monitor_mode),
			seasons = case when $6 then $7::jsonb else seasons end,
			library_folder_id = coalesce($8, library_folder_id),
			media_folder_path = coalesce($9, media_folder_path),
			updated_at = now()
		where id = $1
	`, id, qualityProfileID, minimumAvailability, input.Monitored, monitorMode, updateSeasons, seasonsPayload,
		input.LibraryFolderID, mediaFolderPath)
	if err != nil {
		return MediaItem{}, err
	}
	if tag.RowsAffected() == 0 {
		return MediaItem{}, ErrNotFound
	}
	if input.LibraryFolderID != nil {
		return s.RescanMediaItemFiles(ctx, id)
	}
	return s.GetMediaItem(ctx, id)
}

func (s *SettingsStore) updatedMediaFolderPath(
	ctx context.Context,
	id uuid.UUID,
	libraryFolderID *uuid.UUID,
) (*string, error) {
	if libraryFolderID == nil {
		return nil, nil
	}
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return nil, err
	}
	mediaFolderPath, err := ensureMediaMainFolder(ctx, s.pool, mediaItemInputForFolder(item, libraryFolderID))
	if err != nil || mediaFolderPath == nil {
		return mediaFolderPath, err
	}
	if err := moveMediaItemFiles(item, *mediaFolderPath); err != nil {
		return nil, err
	}
	return mediaFolderPath, nil
}

func mediaItemInputForFolder(item MediaItem, libraryFolderID *uuid.UUID) MediaItemInput {
	return MediaItemInput{
		Type:                  item.Type,
		Title:                 item.Title,
		Year:                  item.Year,
		MediaMetadataSnapshot: item.MediaMetadataSnapshot,
		LibraryFolderID:       libraryFolderID,
	}
}

func optionalMinimumAvailability(value *string) *string {
	if value == nil {
		return nil
	}
	normalized := normalizeMinimumAvailability(*value)
	return &normalized
}

func mediaItemSeasonsPayload(seasons *[]MediaSeason) ([]byte, bool, error) {
	if seasons == nil {
		return []byte("[]"), false, nil
	}
	payload, err := marshalJSONArray(*seasons)
	return payload, true, err
}

func optionalTrimmed(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
