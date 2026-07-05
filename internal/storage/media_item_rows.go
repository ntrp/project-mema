package storage

import (
	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func mediaItemRecordParams(id uuid.UUID, input MediaItemInput, payloads mediaMetadataPayloads, mediaFolderPath *string) storagegen.CreateMediaItemRecordParams {
	return storagegen.CreateMediaItemRecordParams{
		ID:                  id,
		MediaType:           input.Type,
		Title:               input.Title,
		Year:                int4Value(input.Year),
		Monitored:           input.Monitored,
		ExternalProvider:    textValue(input.ExternalProvider),
		ExternalID:          textValue(input.ExternalID),
		Overview:            textValue(input.Overview),
		PosterPath:          textValue(input.PosterPath),
		CollectionID:        textValue(input.CollectionID),
		CollectionName:      textValue(input.CollectionName),
		BackdropPath:        textValue(input.BackdropPath),
		MetadataStatus:      textValue(input.MetadataStatus),
		OriginalLanguage:    textValue(input.OriginalLanguage),
		SeriesType:          textValue(input.SeriesType),
		ReleaseDate:         textValue(input.ReleaseDate),
		FirstAirDate:        textValue(input.FirstAirDate),
		RuntimeMinutes:      int4Value(input.RuntimeMinutes),
		SeasonCount:         int4Value(input.SeasonCount),
		EpisodeCount:        int4Value(input.EpisodeCount),
		VoteAverage:         float8Value(input.VoteAverage),
		Genres:              payloads.genres,
		Keywords:            payloads.keywords,
		Facts:               payloads.facts,
		Seasons:             payloads.seasons,
		CastMembers:         payloads.cast,
		CrewMembers:         payloads.crew,
		Recommendations:     payloads.recommendations,
		SimilarMedia:        payloads.similar,
		MonitorMode:         input.MonitorMode,
		MinimumAvailability: input.MinimumAvailability,
		QualityProfileID:    textValue(input.QualityProfileID),
		LibraryFolderID:     input.LibraryFolderID,
		MediaFolderPath:     textValue(mediaFolderPath),
	}
}

func mediaItemFromListRow(row storagegen.ListMediaItemsRow) MediaItem {
	return mediaItemFromGetRow(storagegen.GetMediaItemRow(row))
}

func mediaItemFromSearchRow(row storagegen.SearchMediaItemsRow) MediaItem {
	return mediaItemFromGetRow(storagegen.GetMediaItemRow(row))
}

func mediaItemFromMissingRow(row storagegen.ListMissingMediaItemsRow) MediaItem {
	return mediaItemFromGetRow(storagegen.GetMediaItemRow(row))
}

func mediaItemFromMatchRow(row storagegen.FindMonitoredMediaMatchRow) MediaItem {
	return mediaItemFromGetRow(storagegen.GetMediaItemRow(row))
}

func mediaItemFromGetRow(row storagegen.GetMediaItemRow) MediaItem {
	item := MediaItem{
		ID:                  row.ID,
		Type:                row.MediaType,
		Title:               row.Title,
		Year:                int4Ptr(row.Year),
		Monitored:           row.Monitored,
		ExternalProvider:    textPtr(row.ExternalProvider),
		ExternalID:          textPtr(row.ExternalID),
		Overview:            textPtr(row.Overview),
		PosterPath:          textPtr(row.PosterPath),
		MonitorMode:         row.MonitorMode,
		SeriesType:          textPtr(row.SeriesType),
		MinimumAvailability: row.MinimumAvailability,
		QualityProfileID:    textPtr(row.QualityProfileID),
		QualityProfileName:  textPtr(row.QualityProfileName),
		Status:              row.Status,
		LibraryFolderID:     row.LibraryFolderID,
		MediaFolderPath:     textPtr(row.MediaFolderPath),
		LibraryFolderPath:   emptyStringPtr(row.LibraryFolderPath),
		FilePaths:           row.FilePaths,
		Tags:                row.Tags,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
	item.CollectionID = textPtr(row.CollectionID)
	item.CollectionName = textPtr(row.CollectionName)
	item.BackdropPath = textPtr(row.BackdropPath)
	item.MetadataStatus = textPtr(row.MetadataStatus)
	item.OriginalLanguage = textPtr(row.OriginalLanguage)
	item.ReleaseDate = textPtr(row.ReleaseDate)
	item.FirstAirDate = textPtr(row.FirstAirDate)
	item.RuntimeMinutes = int4Ptr(row.RuntimeMinutes)
	item.SeasonCount = int4Ptr(row.SeasonCount)
	item.EpisodeCount = int4Ptr(row.EpisodeCount)
	item.VoteAverage = float8Ptr(row.VoteAverage)
	scanMediaMetadata(&item.MediaMetadataSnapshot, row.Genres, row.Keywords, row.Facts, row.Seasons, row.CastMembers, row.CrewMembers, row.Recommendations, row.SimilarMedia)
	item.MetadataFilePaths = collectMetadataFilePaths(item.FilePaths)
	return item
}

func emptyStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
