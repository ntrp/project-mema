package httpapi

import (
	"context"
	"strings"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

const anidbContractSource = "https://wiki.anidb.net/HTTP_API_Definition"

func (s *Server) enrichAniDBFacts(
	ctx context.Context,
	item storage.MediaItem,
	providers []storage.MetadataProvider,
	input storage.MediaItemInput,
	primary metadata.Details,
) (storage.MediaItemInput, error) {
	if primary.ExternalProvider == "anidb" {
		return appendAniDBArtifacts(item, input, primary, true), nil
	}
	externalID, ok := anidbMediaMapping(item.ProviderMappings)
	if !ok {
		return input, nil
	}
	provider, ok := metadataProviderByType(providers, "anidb")
	if !ok {
		return input, nil
	}
	details, err := s.freshMetadataProviderDetails(ctx, provider, metadata.DetailsRequest{
		MediaType: "serie", ExternalID: externalID,
	})
	if err != nil {
		return storage.MediaItemInput{}, err
	}
	return appendAniDBArtifacts(item, input, details, false), nil
}

func anidbMediaMapping(mappings []storage.MediaProviderMapping) (string, bool) {
	for _, mapping := range mappings {
		if mapping.ProviderName == "anidb" && mapping.EntityType == "media_item" && strings.TrimSpace(mapping.ExternalID) != "" {
			return strings.TrimSpace(mapping.ExternalID), true
		}
	}
	return "", false
}

func appendAniDBArtifacts(
	item storage.MediaItem,
	input storage.MediaItemInput,
	details metadata.Details,
	canonical bool,
) storage.MediaItemInput {
	mapping := mediaProviderMappings(details)
	if len(mapping) > 0 {
		mapping[0].Canonical = canonical
		mapping[0].Source = anidbArtifactSource(details)
		input.ProviderMappings = append(input.ProviderMappings, mapping[0])
	}
	for _, alias := range mediaAliases(details) {
		alias.Source = anidbSource(details.ExternalID)
		input.Aliases = append(input.Aliases, alias)
	}
	if !canonical {
		input.Facts = append(input.Facts, anidbMediaFacts(details)...)
	}
	input.EpisodeNumbering = append(input.EpisodeNumbering, resolvedAniDBNumbering(item.ProviderMappings, details.EpisodeFacts)...)
	return input
}

func resolvedAniDBNumbering(
	mappings []storage.MediaProviderMapping,
	facts []metadata.EpisodeFact,
) []storage.MediaEpisodeNumberingInput {
	byID := make(map[string]metadata.EpisodeFact, len(facts))
	for _, fact := range facts {
		byID[fact.ExternalID] = fact
	}
	items := []storage.MediaEpisodeNumberingInput{}
	for _, mapping := range mappings {
		if mapping.ProviderName != "anidb" || mapping.EntityType != "episode" || mapping.EpisodeID == nil {
			continue
		}
		fact, ok := byID[mapping.ExternalID]
		if !ok || fact.AbsoluteNumber == nil {
			continue
		}
		items = append(items, storage.MediaEpisodeNumberingInput{
			SeasonID: mapping.SeasonID, EpisodeID: *mapping.EpisodeID,
			ProviderName: "anidb", NumberingScheme: "absolute", AbsoluteNumber: fact.AbsoluteNumber,
			Source: map[string]any{
				"source": "anidb_http_api", "contract": anidbContractSource,
				"externalId": fact.ExternalID, "originalType": fact.OriginalType,
				"originalNumber": fact.OriginalNumber, "resolvedBy": "provider_mapping",
			},
		})
	}
	return items
}

func anidbSource(externalID string) map[string]any {
	return map[string]any{
		"source": "anidb_http_api", "contract": anidbContractSource, "externalId": externalID,
	}
}

func anidbArtifactSource(details metadata.Details) map[string]any {
	source := anidbSource(details.ExternalID)
	if details.PosterPath != nil {
		source["artwork"] = *details.PosterPath
	}
	if details.ReleaseDate != nil {
		source["startDate"] = *details.ReleaseDate
	}
	return source
}

func anidbMediaFacts(details metadata.Details) []storage.MediaFact {
	source := anidbSource(details.ExternalID)
	facts := make([]storage.MediaFact, 0, len(details.Facts)+2)
	for _, fact := range details.Facts {
		facts = append(facts, storage.MediaFact{Label: fact.Label, Value: fact.Value, Source: source})
	}
	if details.PosterPath != nil {
		facts = append(facts, storage.MediaFact{Label: "AniDB artwork", Value: *details.PosterPath, Source: source})
	}
	if details.ReleaseDate != nil {
		facts = append(facts, storage.MediaFact{Label: "AniDB start date", Value: *details.ReleaseDate, Source: source})
	}
	return facts
}
