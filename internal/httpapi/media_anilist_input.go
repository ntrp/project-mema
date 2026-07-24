package httpapi

import (
	"strings"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func mediaProviderMappings(details metadata.Details) []storage.MediaProviderMappingInput {
	mappings := make([]storage.MediaProviderMappingInput, 0, max(1, len(details.ProviderMappings)))
	for _, mapping := range details.ProviderMappings {
		mappings = append(mappings, storage.MediaProviderMappingInput{
			EntityType:         "media_item",
			ProviderName:       mapping.ProviderName,
			ProviderEntityType: mapping.ProviderEntityType,
			ExternalID:         mapping.ExternalID,
			Canonical:          mapping.Canonical,
			Source:             mapping.Source,
		})
	}
	if len(mappings) > 0 || strings.TrimSpace(details.ExternalProvider) == "" || strings.TrimSpace(details.ExternalID) == "" {
		return mappings
	}
	source := map[string]any{"source": "primary_metadata"}
	if details.ExternalURL != nil && strings.TrimSpace(*details.ExternalURL) != "" {
		source["externalUrl"] = strings.TrimSpace(*details.ExternalURL)
	}
	return []storage.MediaProviderMappingInput{{EntityType: "media_item", ProviderName: details.ExternalProvider,
		ProviderEntityType: details.Type, ExternalID: details.ExternalID, Canonical: true, Source: source}}
}

func mediaAliases(details metadata.Details) []storage.MediaItemAliasInput {
	aliases := make([]storage.MediaItemAliasInput, 0, len(details.Aliases))
	for _, alias := range details.Aliases {
		provider := alias.ProviderName
		aliases = append(aliases, storage.MediaItemAliasInput{
			Alias: alias.Value, Language: alias.Language, Kind: alias.Kind, ProviderName: optionalString(provider), Source: alias.Source,
		})
	}
	return aliases
}

func mergeMediaProviderMappings(existing, next []storage.MediaProviderMappingInput, primaryProvider, primaryID, primaryType string) []storage.MediaProviderMappingInput {
	mappings := append([]storage.MediaProviderMappingInput(nil), existing...)
	if primaryProvider != "" && primaryID != "" && !hasMediaProviderMapping(existing, primaryProvider, primaryID) {
		mappings = append(mappings, storage.MediaProviderMappingInput{EntityType: "media_item", ProviderName: primaryProvider,
			ProviderEntityType: primaryType, ExternalID: primaryID, Canonical: true, Source: map[string]any{"source": "primary_metadata"}})
	}
	for _, mapping := range next {
		if primaryProvider != "" && mapping.ProviderName != primaryProvider {
			mapping.Canonical = false
		}
		mappings = append(mappings, mapping)
	}
	return mappings
}

func hasMediaProviderMapping(mappings []storage.MediaProviderMappingInput, provider, externalID string) bool {
	for _, mapping := range mappings {
		if (mapping.EntityType == "" || mapping.EntityType == "media_item") && mapping.ProviderName == provider && mapping.ExternalID == externalID {
			return true
		}
	}
	return false
}

func mergeMediaAliases(existing, next []storage.MediaItemAliasInput) []storage.MediaItemAliasInput {
	return append(append([]storage.MediaItemAliasInput(nil), existing...), next...)
}
