package httpapi

import (
	"strings"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func localPersonDetails(items []storage.MediaItem, provider, personID string) (metadata.PersonDetails, bool) {
	details := metadata.PersonDetails{ID: personID}
	found := false
	for _, item := range items {
		roles := []string{}
		for _, people := range [][]storage.MediaPerson{item.Cast, item.Crew} {
			for _, person := range people {
				if !samePerson(person, provider, personID) {
					continue
				}
				found = true
				if details.Name == "" {
					details.Name = strings.TrimSpace(person.Name)
				}
				if details.ProfilePath == nil && person.ProfilePath != nil {
					details.ProfilePath = person.ProfilePath
				}
				roles = appendUniqueRole(roles, person.Role)
			}
		}
		if len(roles) == 0 {
			continue
		}
		appearanceProvider, appearanceID, ok := mediaItemProviderIdentity(item)
		if !ok {
			continue
		}
		details.Appearances = append(details.Appearances, metadata.PersonAppearance{
			Title: item.Title, Type: item.Type, Year: item.Year,
			ExternalProvider: appearanceProvider, ExternalID: appearanceID,
			Overview: item.Overview, PosterPath: item.PosterPath, BackdropPath: item.BackdropPath,
			Role: optionalString(strings.Join(roles, ", ")), ReleaseDate: item.ReleaseDate,
		})
	}
	return details, found && details.Name != ""
}

func samePerson(person storage.MediaPerson, provider, personID string) bool {
	return person.ExternalProvider != nil && person.ExternalID != nil &&
		strings.EqualFold(strings.TrimSpace(*person.ExternalProvider), provider) &&
		strings.TrimSpace(*person.ExternalID) == personID
}

func appendUniqueRole(roles []string, role *string) []string {
	if role == nil || strings.TrimSpace(*role) == "" {
		return roles
	}
	value := strings.TrimSpace(*role)
	for _, existing := range roles {
		if existing == value {
			return roles
		}
	}
	return append(roles, value)
}

func mediaItemProviderIdentity(item storage.MediaItem) (string, string, bool) {
	if item.ExternalProvider != nil && item.ExternalID != nil && strings.TrimSpace(*item.ExternalID) != "" {
		return strings.TrimSpace(*item.ExternalProvider), strings.TrimSpace(*item.ExternalID), true
	}
	for _, mapping := range item.ProviderMappings {
		if mapping.EntityType == "media_item" && strings.TrimSpace(mapping.ExternalID) != "" {
			return strings.TrimSpace(mapping.ProviderName), strings.TrimSpace(mapping.ExternalID), true
		}
	}
	return "", "", false
}
