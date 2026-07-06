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
		ContentKind:         optionalContentKind(request.ContentKind),
		Title:               title,
		Year:                request.Year,
		Monitored:           request.Monitored,
		ExternalProvider:    optionalTrimmedString(request.ExternalProvider),
		ExternalID:          optionalTrimmedString(request.ExternalId),
		Overview:            optionalTrimmedString(request.Overview),
		PosterPath:          optionalTrimmedString(request.PosterPath),
		MonitorMode:         string(request.MonitorMode),
		SeriesType:          optionalSeriesType(request.Type, request.SeriesType),
		NumberingStrategy:   optionalNumberingStrategy(request.NumberingStrategy),
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
	providerMappings := mediaProviderMappingResponses(item.ProviderMappings)
	aliases := mediaAliasResponses(item.Aliases)
	episodeNumbering := mediaEpisodeNumberingResponses(item.EpisodeNumbering)
	externalSubtitles := mediaSubtitleResponses(item.ExternalSubtitles)
	contentKind := MediaContentKind(item.ContentKind)
	return MediaItem{
		Id:                  openapi_types.UUID(item.ID),
		Type:                MediaType(item.Type),
		Title:               item.Title,
		Status:              MediaItemStatus(item.Status),
		Year:                item.Year,
		Monitored:           item.Monitored,
		ContentKind:         &contentKind,
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
		NumberingStrategy:   optionalOpenAPINumberingStrategy(item.NumberingStrategy),
		MinimumAvailability: MinimumAvailability(item.MinimumAvailability),
		QualityProfileId:    item.QualityProfileID,
		QualityProfileName:  item.QualityProfileName,
		LibraryFolderId:     optionalOpenAPIUUID(item.LibraryFolderID),
		LibraryFolderPath:   item.LibraryFolderPath,
		MediaFolderPath:     item.MediaFolderPath,
		FilePaths:           item.FilePaths,
		Files:               mediaFileInfoResponses(item.FilePaths, item.SubtitleLanguages, item.ExternalSubtitles),
		ExternalSubtitles:   &externalSubtitles,
		ProviderMappings:    &providerMappings,
		Aliases:             &aliases,
		EpisodeNumbering:    &episodeNumbering,
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

func optionalContentKind(value *MediaContentKind) string {
	if value == nil {
		return "standard"
	}
	return string(*value)
}

func optionalNumberingStrategy(value *MediaNumberingStrategy) *string {
	if value == nil {
		return nil
	}
	strategy := string(*value)
	return &strategy
}

func optionalOpenAPINumberingStrategy(value *string) *MediaNumberingStrategy {
	if value == nil {
		return nil
	}
	strategy := MediaNumberingStrategy(*value)
	return &strategy
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
