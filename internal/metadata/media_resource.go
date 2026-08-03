package metadata

import "strings"

// MergeMediaResources fills fields missing from primary without replacing provider data.
func MergeMediaResources(primary, fallback MediaResource) MediaResource {
	if strings.TrimSpace(primary.Title) == "" {
		primary.Title = fallback.Title
	}
	if strings.TrimSpace(primary.Type) == "" {
		primary.Type = fallback.Type
	}
	fillPointer(&primary.Year, fallback.Year)
	fillPointer(&primary.ExternalURL, fallback.ExternalURL)
	fillPointer(&primary.Overview, fallback.Overview)
	fillPointer(&primary.PosterPath, fallback.PosterPath)
	fillPointer(&primary.CollectionID, fallback.CollectionID)
	fillPointer(&primary.CollectionName, fallback.CollectionName)
	fillPointer(&primary.BackdropPath, fallback.BackdropPath)
	fillPointer(&primary.TrailerURL, fallback.TrailerURL)
	fillPointer(&primary.Status, fallback.Status)
	fillPointer(&primary.OriginalLanguage, fallback.OriginalLanguage)
	fillPointer(&primary.ReleaseDate, fallback.ReleaseDate)
	fillPointer(&primary.FirstAirDate, fallback.FirstAirDate)
	fillPointer(&primary.RuntimeMinutes, fallback.RuntimeMinutes)
	fillPointer(&primary.SeasonCount, fallback.SeasonCount)
	fillPointer(&primary.EpisodeCount, fallback.EpisodeCount)
	fillPointer(&primary.VoteAverage, fallback.VoteAverage)
	fillSlice(&primary.Aliases, fallback.Aliases)
	fillSlice(&primary.Genres, fallback.Genres)
	fillSlice(&primary.Keywords, fallback.Keywords)
	fillSlice(&primary.Facts, fallback.Facts)
	fillSlice(&primary.EpisodeFacts, fallback.EpisodeFacts)
	fillSlice(&primary.Seasons, fallback.Seasons)
	fillSlice(&primary.Cast, fallback.Cast)
	fillSlice(&primary.Crew, fallback.Crew)
	fillSlice(&primary.Recommendations, fallback.Recommendations)
	fillSlice(&primary.Similar, fallback.Similar)
	primary.ProviderMappings = mergeResourceMappings(primary.ProviderMappings, fallback)
	return primary
}

func fillPointer[T any](target **T, fallback *T) {
	if *target == nil {
		*target = fallback
	}
}

func fillSlice[T any](target *[]T, fallback []T) {
	if len(*target) == 0 {
		*target = fallback
	}
}

func mergeResourceMappings(current []ProviderMapping, fallback MediaResource) []ProviderMapping {
	mappings := append([]ProviderMapping{}, current...)
	mappings = appendUniqueMappings(mappings, fallback.ProviderMappings)
	if strings.TrimSpace(fallback.ExternalProvider) != "" && strings.TrimSpace(fallback.ExternalID) != "" {
		mappings = appendUniqueMappings(mappings, []ProviderMapping{{
			ProviderName: fallback.ExternalProvider, ProviderEntityType: fallback.Type,
			ExternalID: fallback.ExternalID, Source: map[string]any{"source": "metadata_fallback"},
		}})
	}
	return mappings
}

func appendUniqueMappings(current, additions []ProviderMapping) []ProviderMapping {
	for _, addition := range additions {
		found := false
		for _, existing := range current {
			if existing.ProviderName == addition.ProviderName && existing.ProviderEntityType == addition.ProviderEntityType && existing.ExternalID == addition.ExternalID {
				found = true
				break
			}
		}
		if !found {
			current = append(current, addition)
		}
	}
	return current
}
