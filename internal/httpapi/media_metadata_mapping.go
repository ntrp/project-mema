package httpapi

import "media-manager/internal/metadata"

func metadataDetailsResponse(details metadata.Details) MediaMetadataDetails {
	genres := append([]string(nil), details.Genres...)
	keywords := append([]string(nil), details.Keywords...)
	facts := make([]MediaMetadataFact, 0, len(details.Facts))
	for _, fact := range details.Facts {
		facts = append(facts, MediaMetadataFact{
			Label: fact.Label,
			Value: fact.Value,
		})
	}
	seasons := make([]MediaMetadataSeason, 0, len(details.Seasons))
	for _, season := range details.Seasons {
		episodes := make([]MediaMetadataEpisode, 0, len(season.Episodes))
		for _, episode := range season.Episodes {
			episodes = append(episodes, MediaMetadataEpisode{
				Name:          episode.Name,
				EpisodeNumber: episode.EpisodeNumber,
				Overview:      episode.Overview,
				AirDate:       episode.AirDate,
				StillPath:     episode.StillPath,
			})
		}
		seasons = append(seasons, MediaMetadataSeason{
			Name:         season.Name,
			EpisodeCount: season.EpisodeCount,
			AirDate:      season.AirDate,
			PosterPath:   season.PosterPath,
			Episodes:     &episodes,
		})
	}
	cast := make([]MediaMetadataPerson, 0, len(details.Cast))
	for _, person := range details.Cast {
		cast = append(cast, MediaMetadataPerson{
			Name:        person.Name,
			Role:        person.Role,
			ProfilePath: person.ProfilePath,
		})
	}
	recommendations := make([]MediaSearchResult, 0, len(details.Recommendations))
	for _, result := range details.Recommendations {
		recommendations = append(recommendations, metadataSearchResultResponse(result))
	}
	similar := make([]MediaSearchResult, 0, len(details.Similar))
	for _, result := range details.Similar {
		similar = append(similar, metadataSearchResultResponse(result))
	}
	return MediaMetadataDetails{
		Title:            details.Title,
		Type:             MediaType(details.Type),
		Year:             details.Year,
		ExternalProvider: details.ExternalProvider,
		ExternalId:       details.ExternalID,
		Overview:         details.Overview,
		PosterPath:       details.PosterPath,
		CollectionId:     details.CollectionID,
		CollectionName:   details.CollectionName,
		BackdropPath:     details.BackdropPath,
		TrailerUrl:       details.TrailerURL,
		Status:           details.Status,
		OriginalLanguage: details.OriginalLanguage,
		ReleaseDate:      details.ReleaseDate,
		FirstAirDate:     details.FirstAirDate,
		RuntimeMinutes:   details.RuntimeMinutes,
		SeasonCount:      details.SeasonCount,
		EpisodeCount:     details.EpisodeCount,
		VoteAverage:      details.VoteAverage,
		Genres:           &genres,
		Keywords:         &keywords,
		Facts:            &facts,
		Seasons:          &seasons,
		Cast:             &cast,
		Recommendations:  &recommendations,
		Similar:          &similar,
	}
}
