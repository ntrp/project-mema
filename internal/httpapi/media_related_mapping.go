package httpapi

import (
	"strings"

	"media-manager/internal/storage"
)

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
			ExternalUrl:      value.ExternalURL,
			Overview:         value.Overview,
			PosterPath:       value.PosterPath,
		})
	}
	return items
}

func mediaItemExternalURL(item storage.MediaItem) *string {
	for _, mapping := range item.ProviderMappings {
		if !mapping.Canonical || !strings.EqualFold(mapping.EntityType, "media_item") {
			continue
		}
		if value, ok := mapping.Source["externalUrl"].(string); ok && strings.TrimSpace(value) != "" {
			return optionalString(strings.TrimSpace(value))
		}
	}
	if item.ExternalProvider == nil || item.ExternalID == nil {
		return nil
	}
	externalID := strings.TrimSpace(*item.ExternalID)
	if externalID == "" {
		return nil
	}
	if strings.EqualFold(*item.ExternalProvider, "tmdb") {
		path := "movie"
		if strings.EqualFold(item.Type, "serie") {
			path = "tv"
		}
		return optionalString("https://www.themoviedb.org/" + path + "/" + externalID)
	}
	if strings.EqualFold(*item.ExternalProvider, "tvdb") {
		path := "movie"
		if strings.EqualFold(item.Type, "serie") {
			path = "series"
		}
		return optionalString("https://thetvdb.com/dereferrer/" + path + "/" + externalID)
	}
	return nil
}
