package httpapi

import (
	"context"
	"errors"
	"strings"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

var errMediaCollectionUnavailable = errors.New("media collection unavailable")

func (s *Server) createMediaForAdd(ctx context.Context, input storage.MediaItemInput) ([]storage.MediaItem, error) {
	inputs, err := s.mediaAddInputs(ctx, input)
	if err != nil {
		return nil, err
	}
	for index := range inputs {
		enriched, err := s.enrichMediaItemInput(ctx, inputs[index])
		if err != nil {
			return nil, err
		}
		inputs[index] = applySeriesMonitoring(enriched)
	}
	return s.createMediaInputs(ctx, inputs)
}

func (s *Server) createMediaInputs(ctx context.Context, inputs []storage.MediaItemInput) ([]storage.MediaItem, error) {
	items := make([]storage.MediaItem, 0, len(inputs))
	for _, nextInput := range inputs {
		item, err := s.settings.CreateMediaItem(ctx, nextInput)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func mediaInputFromRequest(
	request storage.MediaRequest,
	input storage.MediaRequestApprovalInput,
) storage.MediaItemInput {
	qualityProfileID := input.QualityProfileID
	libraryFolderID := input.LibraryFolderID
	return storage.MediaItemInput{
		Type:                request.Type,
		Title:               request.Title,
		Year:                request.Year,
		Monitored:           true,
		ExternalProvider:    request.ExternalProvider,
		ExternalID:          request.ExternalID,
		Overview:            request.Overview,
		PosterPath:          request.PosterPath,
		MonitorMode:         request.MonitorMode,
		SeriesType:          request.SeriesType,
		MinimumAvailability: request.MinimumAvailability,
		Tags:                request.Tags,
		QualityProfileID:    &qualityProfileID,
		LibraryFolderID:     &libraryFolderID,
	}
}

func (s *Server) enrichMediaItemInput(ctx context.Context, input storage.MediaItemInput) (storage.MediaItemInput, error) {
	if input.ExternalProvider == nil || input.ExternalID == nil || strings.TrimSpace(*input.ExternalID) == "" {
		return input, nil
	}
	providers, err := s.settings.ListMetadataProviders(ctx)
	if err != nil {
		return storage.MediaItemInput{}, err
	}
	provider, ok := metadataProviderByType(providers, strings.TrimSpace(*input.ExternalProvider))
	if !ok {
		return input, nil
	}
	details, err := s.metadataProviderDetails(ctx, provider, metadata.DetailsRequest{
		MediaType:  input.Type,
		ExternalID: strings.TrimSpace(*input.ExternalID),
	})
	if err != nil {
		return storage.MediaItemInput{}, err
	}
	return applyMediaDetails(input, details), nil
}

func applyMediaDetails(input storage.MediaItemInput, details metadata.Details) storage.MediaItemInput {
	input.Title = details.Title
	input.Type = details.Type
	input.Year = details.Year
	input.ExternalProvider = optionalString(details.ExternalProvider)
	input.ExternalID = optionalString(details.ExternalID)
	input.Overview = details.Overview
	input.PosterPath = details.PosterPath
	input.MediaMetadataSnapshot = storage.MediaMetadataSnapshot{
		CollectionID:     details.CollectionID,
		CollectionName:   details.CollectionName,
		BackdropPath:     details.BackdropPath,
		MetadataStatus:   details.Status,
		OriginalLanguage: details.OriginalLanguage,
		ReleaseDate:      details.ReleaseDate,
		FirstAirDate:     details.FirstAirDate,
		RuntimeMinutes:   details.RuntimeMinutes,
		SeasonCount:      details.SeasonCount,
		EpisodeCount:     details.EpisodeCount,
		VoteAverage:      details.VoteAverage,
		Genres:           append([]string(nil), details.Genres...),
		Facts:            mediaFacts(details.Facts),
		Seasons:          mediaSeasons(details.Seasons),
		Cast:             mediaCast(details.Cast),
	}
	return input
}

func mediaFacts(facts []metadata.Fact) []storage.MediaFact {
	items := make([]storage.MediaFact, 0, len(facts))
	for _, fact := range facts {
		items = append(items, storage.MediaFact{Label: fact.Label, Value: fact.Value})
	}
	return items
}

func mediaSeasons(seasons []metadata.Season) []storage.MediaSeason {
	items := make([]storage.MediaSeason, 0, len(seasons))
	for _, season := range seasons {
		items = append(items, storage.MediaSeason{
			Name:         season.Name,
			EpisodeCount: season.EpisodeCount,
			AirDate:      season.AirDate,
			PosterPath:   season.PosterPath,
			Episodes:     mediaEpisodes(season.Episodes),
		})
	}
	return items
}

func mediaEpisodes(episodes []metadata.Episode) []storage.MediaEpisode {
	items := make([]storage.MediaEpisode, 0, len(episodes))
	for _, episode := range episodes {
		items = append(items, storage.MediaEpisode{
			Name:          episode.Name,
			EpisodeNumber: episode.EpisodeNumber,
			Overview:      episode.Overview,
			AirDate:       episode.AirDate,
			StillPath:     episode.StillPath,
		})
	}
	return items
}

func mediaCast(cast []metadata.Person) []storage.MediaPerson {
	items := make([]storage.MediaPerson, 0, len(cast))
	for _, person := range cast {
		items = append(items, storage.MediaPerson{
			Name:        person.Name,
			Role:        person.Role,
			ProfilePath: person.ProfilePath,
		})
	}
	return items
}

func (s *Server) mediaAddInputs(ctx context.Context, input storage.MediaItemInput) ([]storage.MediaItemInput, error) {
	if input.Type != "movie" || input.MonitorMode != "collection" {
		return []storage.MediaItemInput{input}, nil
	}
	collection, err := s.mediaCollection(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(collection.Parts) == 0 {
		return []storage.MediaItemInput{input}, nil
	}
	inputs := make([]storage.MediaItemInput, 0, len(collection.Parts))
	for _, part := range collection.Parts {
		nextInput := input
		nextInput.Type = part.Type
		nextInput.Title = part.Title
		nextInput.Year = part.Year
		nextInput.ExternalProvider = optionalString(part.ExternalProvider)
		nextInput.ExternalID = optionalString(part.ExternalID)
		nextInput.Overview = part.Overview
		nextInput.PosterPath = part.PosterPath
		inputs = append(inputs, nextInput)
	}
	return inputs, nil
}

func (s *Server) mediaCollection(ctx context.Context, input storage.MediaItemInput) (metadata.Collection, error) {
	if input.Type != "movie" || input.ExternalID == nil || strings.TrimSpace(*input.ExternalID) == "" {
		return metadata.Collection{}, errMediaCollectionUnavailable
	}
	provider, ok, err := s.tmdbProvider(ctx)
	if err != nil {
		return metadata.Collection{}, err
	}
	if !ok {
		return metadata.Collection{}, errMediaCollectionUnavailable
	}
	details, err := s.metadata.Details(ctx, metadataProviderConfig(provider), metadata.DetailsRequest{
		MediaType:  input.Type,
		ExternalID: strings.TrimSpace(*input.ExternalID),
	})
	if err != nil {
		return metadata.Collection{}, err
	}
	if details.CollectionID == nil || strings.TrimSpace(*details.CollectionID) == "" {
		return metadata.Collection{}, errMediaCollectionUnavailable
	}
	return s.metadata.Collection(ctx, metadataProviderConfig(provider), *details.CollectionID)
}

func (s *Server) tmdbProvider(ctx context.Context) (storage.MetadataProvider, bool, error) {
	providers, err := s.settings.ListEnabledMetadataProviders(ctx, "movie")
	if err != nil {
		return storage.MetadataProvider{}, false, err
	}
	for _, provider := range providers {
		if provider.Type == "tmdb" {
			return provider, true, nil
		}
	}
	return storage.MetadataProvider{}, false, nil
}

func (s *Server) enqueueAutomaticSearch(ctx context.Context, items []storage.MediaItem) {
	for _, item := range items {
		_, _ = s.jobs.EnqueueAutoSearchDownload(ctx, item.ID)
	}
}
