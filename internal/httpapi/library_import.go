package httpapi

import (
	"context"

	"media-manager/internal/storage"
)

func (s *Server) enrichLibraryImportMatch(ctx context.Context, input storage.LibraryMatchInput) (storage.LibraryMatchInput, error) {
	if input.MediaItemID != nil {
		return input, nil
	}
	mediaType, ok := scanInputMediaType(input.MediaKind)
	if !ok {
		return input, nil
	}
	enriched, err := s.enrichMediaItemInput(ctx, storage.MediaItemInput{
		Type:                mediaType,
		Title:               input.Title,
		Year:                input.Year,
		Monitored:           input.Monitored,
		ExternalProvider:    input.ExternalProvider,
		ExternalID:          input.ExternalID,
		Overview:            input.Overview,
		PosterPath:          input.PosterPath,
		MonitorMode:         input.MonitorMode,
		SeriesType:          input.SeriesType,
		MinimumAvailability: input.MinimumAvailability,
		QualityProfileID:    &input.QualityProfileID,
	})
	if err != nil {
		return storage.LibraryMatchInput{}, err
	}
	input.Title = enriched.Title
	input.Year = enriched.Year
	input.ExternalProvider = enriched.ExternalProvider
	input.ExternalID = enriched.ExternalID
	input.Overview = enriched.Overview
	input.PosterPath = enriched.PosterPath
	input.MediaMetadataSnapshot = enriched.MediaMetadataSnapshot
	return input, nil
}
