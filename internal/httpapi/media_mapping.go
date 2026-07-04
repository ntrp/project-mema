package httpapi

import (
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func mediaItemInput(request MediaItemCreateRequest) (storage.MediaItemInput, bool) {
	title := strings.TrimSpace(request.Title)
	if title == "" || !request.Type.Valid() {
		return storage.MediaItemInput{}, false
	}
	return storage.MediaItemInput{
		Type:                string(request.Type),
		Title:               title,
		Year:                request.Year,
		Monitored:           request.Monitored,
		ExternalProvider:    optionalTrimmedString(request.ExternalProvider),
		ExternalID:          optionalTrimmedString(request.ExternalId),
		Overview:            optionalTrimmedString(request.Overview),
		PosterPath:          optionalTrimmedString(request.PosterPath),
		MonitorMode:         string(request.MonitorMode),
		SeriesType:          optionalSeriesType(request.Type, request.SeriesType),
		MinimumAvailability: string(request.MinimumAvailability),
		QualityProfileID:    optionalTrimmedString(request.QualityProfileId),
		LibraryFolderID:     optionalUUID(request.LibraryFolderId),
		Tags:                optionalStringSlice(request.Tags),
	}, true
}

func mediaItemResponse(item storage.MediaItem) MediaItem {
	genres := append([]string(nil), item.Genres...)
	keywords := append([]string(nil), item.Keywords...)
	facts := mediaFactResponses(item.Facts)
	seasons := mediaSeasonResponses(item.Seasons)
	cast := mediaPersonResponses(item.Cast)
	crew := mediaPersonResponses(item.Crew)
	recommendations := mediaRelatedResponses(item.Recommendations)
	similar := mediaRelatedResponses(item.Similar)
	return MediaItem{
		Id:                  openapi_types.UUID(item.ID),
		Type:                MediaType(item.Type),
		Title:               item.Title,
		Status:              MediaItemStatus(item.Status),
		Year:                item.Year,
		Monitored:           item.Monitored,
		ExternalProvider:    item.ExternalProvider,
		ExternalId:          item.ExternalID,
		Overview:            item.Overview,
		PosterPath:          item.PosterPath,
		CollectionId:        item.CollectionID,
		CollectionName:      item.CollectionName,
		BackdropPath:        item.BackdropPath,
		MetadataStatus:      item.MetadataStatus,
		OriginalLanguage:    item.OriginalLanguage,
		ReleaseDate:         item.ReleaseDate,
		FirstAirDate:        item.FirstAirDate,
		RuntimeMinutes:      item.RuntimeMinutes,
		SeasonCount:         item.SeasonCount,
		EpisodeCount:        item.EpisodeCount,
		VoteAverage:         item.VoteAverage,
		Genres:              &genres,
		Keywords:            &keywords,
		Facts:               &facts,
		Seasons:             &seasons,
		Cast:                &cast,
		Crew:                &crew,
		Recommendations:     &recommendations,
		Similar:             &similar,
		MonitorMode:         MediaMonitorMode(item.MonitorMode),
		SeriesType:          optionalOpenAPISeriesType(item.SeriesType),
		MinimumAvailability: MinimumAvailability(item.MinimumAvailability),
		QualityProfileId:    item.QualityProfileID,
		QualityProfileName:  item.QualityProfileName,
		LibraryFolderId:     optionalOpenAPIUUID(item.LibraryFolderID),
		LibraryFolderPath:   item.LibraryFolderPath,
		MediaFolderPath:     item.MediaFolderPath,
		FilePaths:           item.FilePaths,
		Files:               mediaFileInfoResponses(item.FilePaths),
		MetadataFilePaths:   item.MetadataFilePaths,
		Tags:                &item.Tags,
		CreatedAt:           item.CreatedAt,
		UpdatedAt:           item.UpdatedAt,
	}
}

func mediaRequestInput(request MediaRequestCreateRequest, requestedByUserID uuid.UUID) (storage.MediaRequestInput, bool) {
	title := strings.TrimSpace(request.Title)
	if title == "" || !request.Type.Valid() {
		return storage.MediaRequestInput{}, false
	}
	return storage.MediaRequestInput{
		RequestedByUserID:   requestedByUserID,
		Type:                string(request.Type),
		Title:               title,
		Year:                request.Year,
		ExternalProvider:    optionalTrimmedString(request.ExternalProvider),
		ExternalID:          optionalTrimmedString(request.ExternalId),
		Overview:            optionalTrimmedString(request.Overview),
		PosterPath:          optionalTrimmedString(request.PosterPath),
		MonitorMode:         string(request.MonitorMode),
		SeriesType:          optionalSeriesType(request.Type, request.SeriesType),
		MinimumAvailability: string(request.MinimumAvailability),
		Tags:                optionalStringSlice(request.Tags),
	}, true
}

func mediaRequestResponse(request storage.MediaRequest) MediaRequest {
	return MediaRequest{
		Id:                  openapi_types.UUID(request.ID),
		RequestedByUserId:   openapi_types.UUID(request.RequestedByUserID),
		RequestedByUsername: request.RequestedByUsername,
		Type:                MediaType(request.Type),
		Title:               request.Title,
		Year:                request.Year,
		ExternalProvider:    request.ExternalProvider,
		ExternalId:          request.ExternalID,
		Overview:            request.Overview,
		PosterPath:          request.PosterPath,
		MonitorMode:         MediaMonitorMode(request.MonitorMode),
		SeriesType:          optionalOpenAPISeriesType(request.SeriesType),
		MinimumAvailability: MinimumAvailability(request.MinimumAvailability),
		Tags:                &request.Tags,
		Status:              MediaRequestStatus(request.Status),
		QualityProfileId:    request.QualityProfileID,
		LibraryFolderId:     optionalOpenAPIUUID(request.LibraryFolderID),
		MediaItemId:         optionalOpenAPIUUID(request.MediaItemID),
		DecidedAt:           request.DecidedAt,
		CreatedAt:           request.CreatedAt,
		UpdatedAt:           request.UpdatedAt,
	}
}

func optionalStringSlice(values *[]string) []string {
	if values == nil {
		return nil
	}
	return append([]string(nil), (*values)...)
}

func mediaFactResponses(values []storage.MediaFact) []MediaMetadataFact {
	items := make([]MediaMetadataFact, 0, len(values))
	for _, value := range values {
		items = append(items, MediaMetadataFact{Label: value.Label, Value: value.Value})
	}
	return items
}

func mediaSeasonResponses(values []storage.MediaSeason) []MediaMetadataSeason {
	items := make([]MediaMetadataSeason, 0, len(values))
	for _, value := range values {
		episodes := mediaEpisodeResponses(value.Episodes)
		items = append(items, MediaMetadataSeason{
			Name:         value.Name,
			EpisodeCount: value.EpisodeCount,
			AirDate:      value.AirDate,
			PosterPath:   value.PosterPath,
			Monitored:    &value.Monitored,
			Episodes:     &episodes,
		})
	}
	return items
}

func mediaEpisodeResponses(values []storage.MediaEpisode) []MediaMetadataEpisode {
	items := make([]MediaMetadataEpisode, 0, len(values))
	for _, value := range values {
		items = append(items, MediaMetadataEpisode{
			Name:          value.Name,
			EpisodeNumber: value.EpisodeNumber,
			Overview:      value.Overview,
			AirDate:       value.AirDate,
			StillPath:     value.StillPath,
			Monitored:     &value.Monitored,
		})
	}
	return items
}

func optionalSeriesType(mediaType MediaType, value *SeriesType) *string {
	if mediaType != MediaTypeSerie || value == nil {
		return nil
	}
	seriesType := string(*value)
	return &seriesType
}

func optionalOpenAPISeriesType(value *string) *SeriesType {
	if value == nil {
		return nil
	}
	seriesType := SeriesType(*value)
	return &seriesType
}

func mediaPersonResponses(values []storage.MediaPerson) []MediaMetadataPerson {
	items := make([]MediaMetadataPerson, 0, len(values))
	for _, value := range values {
		items = append(items, MediaMetadataPerson{
			ExternalProvider: metadataProviderType(value.ExternalProvider),
			ExternalId:       value.ExternalID,
			Name:             value.Name,
			Role:             value.Role,
			ProfilePath:      value.ProfilePath,
		})
	}
	return items
}

func mediaRelatedResponses(values []storage.MediaRelatedItem) []MediaSearchResult {
	items := make([]MediaSearchResult, 0, len(values))
	for _, value := range values {
		items = append(items, MediaSearchResult{
			Title:            value.Title,
			Type:             MediaType(value.Type),
			Year:             value.Year,
			ExternalProvider: &value.ExternalProvider,
			ExternalId:       &value.ExternalID,
			Overview:         value.Overview,
			PosterPath:       value.PosterPath,
		})
	}
	return items
}

func downloadActivityResponse(activity storage.DownloadActivity) DownloadActivity {
	return DownloadActivity{
		Id:                 openapi_types.UUID(activity.ID),
		MediaItemId:        openapi_types.UUID(activity.MediaItemID),
		MediaTitle:         activity.MediaTitle,
		MediaType:          MediaType(activity.MediaType),
		MediaYear:          activity.MediaYear,
		ReleaseTitle:       activity.ReleaseTitle,
		IndexerName:        activity.IndexerName,
		DownloadClientName: activity.DownloadClientName,
		DownloadId:         activity.DownloadID,
		DownloadUrl:        activity.DownloadURL,
		Status:             DownloadActivityStatus(activity.Status),
		ProgressPercent:    activity.ProgressPercent,
		Error:              activity.Error,
		FailureType:        downloadActivityFailureType(activity.FailureType),
		CreatedAt:          activity.CreatedAt,
		UpdatedAt:          activity.UpdatedAt,
	}
}

func downloadActivityFailureType(value *string) *DownloadActivityFailureType {
	if value == nil {
		return nil
	}
	failureType := DownloadActivityFailureType(*value)
	return &failureType
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func optionalUUID(value *openapi_types.UUID) *uuid.UUID {
	if value == nil {
		return nil
	}
	converted := uuid.UUID(*value)
	return &converted
}

func optionalOpenAPIUUID(value *uuid.UUID) *openapi_types.UUID {
	if value == nil {
		return nil
	}
	converted := openapi_types.UUID(*value)
	return &converted
}
