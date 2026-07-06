package httpapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func mediaProviderMappingResponses(values []storage.MediaProviderMapping) []MediaProviderMapping {
	items := make([]MediaProviderMapping, 0, len(values))
	for _, value := range values {
		items = append(items, MediaProviderMapping{
			Id:                 openapi_types.UUID(value.ID),
			SeasonId:           optionalOpenAPIUUID(value.SeasonID),
			EpisodeId:          optionalOpenAPIUUID(value.EpisodeID),
			EntityType:         MediaProviderMappingEntityType(value.EntityType),
			ProviderName:       MetadataProviderType(value.ProviderName),
			ProviderEntityType: value.ProviderEntityType,
			ExternalId:         value.ExternalID,
			Canonical:          value.Canonical,
			Confidence:         value.Confidence,
		})
	}
	return items
}

func mediaAliasResponses(values []storage.MediaItemAlias) []MediaItemAlias {
	items := make([]MediaItemAlias, 0, len(values))
	for _, value := range values {
		items = append(items, MediaItemAlias{
			Id:              openapi_types.UUID(value.ID),
			Alias:           value.Alias,
			NormalizedAlias: value.NormalizedAlias,
			Language:        value.Language,
			Kind:            MediaItemAliasKind(value.Kind),
			ProviderName:    optionalMetadataProviderType(value.ProviderName),
		})
	}
	return items
}

func mediaEpisodeNumberingResponses(values []storage.MediaEpisodeNumbering) []MediaEpisodeNumbering {
	items := make([]MediaEpisodeNumbering, 0, len(values))
	for _, value := range values {
		items = append(items, MediaEpisodeNumbering{
			Id:              openapi_types.UUID(value.ID),
			SeasonId:        optionalOpenAPIUUID(value.SeasonID),
			EpisodeId:       openapi_types.UUID(value.EpisodeID),
			ProviderName:    MetadataProviderType(value.ProviderName),
			NumberingScheme: MediaEpisodeNumberingNumberingScheme(value.NumberingScheme),
			SeasonNumber:    value.SeasonNumber,
			EpisodeNumber:   value.EpisodeNumber,
			AbsoluteNumber:  value.AbsoluteNumber,
		})
	}
	return items
}

func optionalMetadataProviderType(value *string) *MetadataProviderType {
	if value == nil {
		return nil
	}
	provider := MetadataProviderType(*value)
	return &provider
}
