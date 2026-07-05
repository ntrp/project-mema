package storage

import (
	"context"
	"strings"

	storagegen "media-manager/internal/storage/generated"

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
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return MediaItem{}, err
	}
	seasons, updateSeasons := mediaItemUpdateSeasons(item.Seasons, input)
	seasonsPayload, err := mediaItemSeasonsPayload(seasons)
	if err != nil {
		return MediaItem{}, err
	}
	rows, err := storagegen.New(s.pool).UpdateMediaItemOptionsRecord(ctx, storagegen.UpdateMediaItemOptionsRecordParams{
		QualityProfileID:    textValue(qualityProfileID),
		MinimumAvailability: textValue(minimumAvailability),
		Monitored:           boolValue(input.Monitored),
		MonitorMode:         textValue(monitorMode),
		UpdateSeasons:       updateSeasons,
		Seasons:             seasonsPayload,
		LibraryFolderID:     input.LibraryFolderID,
		MediaFolderPath:     textValue(mediaFolderPath),
		ID:                  id,
	})
	if err != nil {
		return MediaItem{}, err
	}
	if rows == 0 {
		return MediaItem{}, ErrNotFound
	}
	if updateSeasons && seasons != nil {
		if err := applyMediaSeriesMonitorSnapshot(ctx, s.pool, id, *seasons); err != nil {
			return MediaItem{}, err
		}
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

func mediaItemUpdateSeasons(current []MediaSeason, input MediaItemOptionsInput) (*[]MediaSeason, bool) {
	if input.Seasons != nil {
		return input.Seasons, true
	}
	if input.MonitorSeasonName == nil {
		return nil, false
	}
	seasons := make([]MediaSeason, len(current))
	copy(seasons, current)
	for index := range seasons {
		seasons[index].Episodes = append([]MediaEpisode(nil), seasons[index].Episodes...)
	}
	for index := range seasons {
		if seasons[index].Name != *input.MonitorSeasonName {
			continue
		}
		applySeasonMonitorPatch(&seasons[index], input)
		return &seasons, true
	}
	return nil, false
}

func applySeasonMonitorPatch(season *MediaSeason, input MediaItemOptionsInput) {
	if input.SeasonMonitored != nil {
		season.Monitored = *input.SeasonMonitored
		for index := range season.Episodes {
			season.Episodes[index].Monitored = *input.SeasonMonitored
		}
		return
	}
	if input.MonitorEpisodeNumber == nil || input.EpisodeMonitored == nil {
		return
	}
	for index := range season.Episodes {
		if season.Episodes[index].EpisodeNumber == *input.MonitorEpisodeNumber {
			season.Episodes[index].Monitored = *input.EpisodeMonitored
			break
		}
	}
	season.Monitored = mediaSeasonHasMonitoredEpisode(*season)
}

func mediaSeasonHasMonitoredEpisode(season MediaSeason) bool {
	for _, episode := range season.Episodes {
		if episode.Monitored {
			return true
		}
	}
	return false
}

func mediaItemSeasonsPayload(seasons *[]MediaSeason) ([]byte, error) {
	if seasons == nil {
		return []byte("[]"), nil
	}
	payload, err := marshalJSONArray(*seasons)
	return payload, err
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
