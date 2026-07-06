package storage

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

var aliasNormalizePattern = regexp.MustCompile(`[^a-z0-9]+`)

func hydrateMediaItemAnime(ctx context.Context, q storagegen.DBTX, item MediaItem) (MediaItem, error) {
	queries := storagegen.New(q)
	state, err := queries.GetMediaItemAnimeState(ctx, item.ID)
	if err != nil {
		return item, err
	}
	item.ContentKind = state.ContentKind
	item.NumberingStrategy = textPtr(state.NumberingStrategy)
	mappings, err := queries.ListMediaProviderMappings(ctx, item.ID)
	if err != nil {
		return item, err
	}
	for _, row := range mappings {
		item.ProviderMappings = append(item.ProviderMappings, mediaProviderMappingFromRow(row))
	}
	aliases, err := queries.ListMediaItemAliases(ctx, item.ID)
	if err != nil {
		return item, err
	}
	for _, row := range aliases {
		item.Aliases = append(item.Aliases, mediaItemAliasFromRow(row))
	}
	numbering, err := queries.ListMediaEpisodeNumbering(ctx, item.ID)
	if err != nil {
		return item, err
	}
	for _, row := range numbering {
		item.EpisodeNumbering = append(item.EpisodeNumbering, mediaEpisodeNumberingFromRow(row))
	}
	return item, nil
}

func upsertAnimeMetadata(ctx context.Context, q storagegen.DBTX, mediaItemID uuid.UUID, input MediaItemInput) error {
	input.ProviderMappings = defaultProviderMappings(input)
	input.Aliases = defaultMediaAliases(input)
	queries := storagegen.New(q)
	for _, mapping := range input.ProviderMappings {
		if strings.TrimSpace(mapping.ProviderName) == "" || strings.TrimSpace(mapping.ExternalID) == "" {
			continue
		}
		if _, err := queries.UpsertMediaProviderMapping(ctx, providerMappingParams(mediaItemID, mapping)); err != nil {
			return err
		}
	}
	for _, alias := range input.Aliases {
		if strings.TrimSpace(alias.Alias) == "" {
			continue
		}
		if _, err := queries.UpsertMediaItemAlias(ctx, aliasParams(mediaItemID, alias)); err != nil {
			return err
		}
	}
	for _, numbering := range input.EpisodeNumbering {
		if _, err := queries.UpsertMediaEpisodeNumbering(ctx, numberingParams(mediaItemID, numbering)); err != nil {
			return err
		}
	}
	return nil
}

func defaultProviderMappings(input MediaItemInput) []MediaProviderMappingInput {
	if len(input.ProviderMappings) > 0 || input.ExternalProvider == nil || input.ExternalID == nil {
		return input.ProviderMappings
	}
	return []MediaProviderMappingInput{{
		EntityType:         "media_item",
		ProviderName:       strings.TrimSpace(*input.ExternalProvider),
		ProviderEntityType: input.Type,
		ExternalID:         strings.TrimSpace(*input.ExternalID),
		Canonical:          true,
		Source:             map[string]any{"source": "primary_metadata"},
	}}
}

func defaultMediaAliases(input MediaItemInput) []MediaItemAliasInput {
	if len(input.Aliases) > 0 || strings.TrimSpace(input.Title) == "" {
		return input.Aliases
	}
	alias := MediaItemAliasInput{Alias: input.Title, Kind: "canonical", Source: map[string]any{"source": "primary_title"}}
	if input.ExternalProvider != nil {
		provider := strings.TrimSpace(*input.ExternalProvider)
		alias.ProviderName = &provider
	}
	return []MediaItemAliasInput{alias}
}

func providerMappingParams(mediaItemID uuid.UUID, input MediaProviderMappingInput) storagegen.UpsertMediaProviderMappingParams {
	if input.EntityType == "" {
		input.EntityType = "media_item"
	}
	if input.ProviderEntityType == "" {
		input.ProviderEntityType = input.EntityType
	}
	return storagegen.UpsertMediaProviderMappingParams{
		ID:                 uuid.New(),
		MediaItemID:        mediaItemID,
		SeasonID:           input.SeasonID,
		EpisodeID:          input.EpisodeID,
		EntityType:         input.EntityType,
		ProviderName:       input.ProviderName,
		ProviderEntityType: input.ProviderEntityType,
		ExternalID:         input.ExternalID,
		Canonical:          input.Canonical,
		Confidence:         float8Value(input.Confidence),
		Source:             jsonObject(input.Source),
	}
}

func aliasParams(mediaItemID uuid.UUID, input MediaItemAliasInput) storagegen.UpsertMediaItemAliasParams {
	if input.Kind == "" {
		input.Kind = "canonical"
	}
	return storagegen.UpsertMediaItemAliasParams{
		ID:                uuid.New(),
		MediaItemID:       mediaItemID,
		Alias:             strings.TrimSpace(input.Alias),
		NormalizedAlias:   normalizedAlias(input.Alias),
		Language:          textValue(input.Language),
		AliasKind:         input.Kind,
		ProviderName:      textValue(input.ProviderName),
		ProviderMappingID: input.ProviderMappingID,
		Source:            jsonObject(input.Source),
	}
}

func numberingParams(mediaItemID uuid.UUID, input MediaEpisodeNumberingInput) storagegen.UpsertMediaEpisodeNumberingParams {
	return storagegen.UpsertMediaEpisodeNumberingParams{
		ID:              uuid.New(),
		MediaItemID:     mediaItemID,
		SeasonID:        input.SeasonID,
		EpisodeID:       input.EpisodeID,
		ProviderName:    input.ProviderName,
		NumberingScheme: input.NumberingScheme,
		SeasonNumber:    int4Value(input.SeasonNumber),
		EpisodeNumber:   int4Value(input.EpisodeNumber),
		AbsoluteNumber:  int4Value(input.AbsoluteNumber),
		Source:          jsonObject(input.Source),
	}
}

func mediaProviderMappingFromRow(row storagegen.AppMediaProviderMapping) MediaProviderMapping {
	return MediaProviderMapping{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		SeasonID:           row.SeasonID,
		EpisodeID:          row.EpisodeID,
		EntityType:         row.EntityType,
		ProviderName:       row.ProviderName,
		ProviderEntityType: row.ProviderEntityType,
		ExternalID:         row.ExternalID,
		Canonical:          row.Canonical,
		Confidence:         float8Ptr(row.Confidence),
		Source:             jsonMap(row.Source),
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}

func mediaItemAliasFromRow(row storagegen.AppMediaItemAlias) MediaItemAlias {
	return MediaItemAlias{
		ID:                row.ID,
		MediaItemID:       row.MediaItemID,
		Alias:             row.Alias,
		NormalizedAlias:   row.NormalizedAlias,
		Language:          textPtr(row.Language),
		Kind:              row.AliasKind,
		ProviderName:      textPtr(row.ProviderName),
		ProviderMappingID: row.ProviderMappingID,
		Source:            jsonMap(row.Source),
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}

func mediaEpisodeNumberingFromRow(row storagegen.AppMediaEpisodeNumbering) MediaEpisodeNumbering {
	return MediaEpisodeNumbering{
		ID:              row.ID,
		MediaItemID:     row.MediaItemID,
		SeasonID:        row.SeasonID,
		EpisodeID:       row.EpisodeID,
		ProviderName:    row.ProviderName,
		NumberingScheme: row.NumberingScheme,
		SeasonNumber:    int4Ptr(row.SeasonNumber),
		EpisodeNumber:   int4Ptr(row.EpisodeNumber),
		AbsoluteNumber:  int4Ptr(row.AbsoluteNumber),
		Source:          jsonMap(row.Source),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func normalizedAlias(value string) string {
	return aliasNormalizePattern.ReplaceAllString(strings.ToLower(strings.TrimSpace(value)), "")
}

func jsonObject(source map[string]any) []byte {
	if source == nil {
		source = map[string]any{}
	}
	payload, _ := json.Marshal(source)
	return payload
}

func jsonMap(payload []byte) map[string]any {
	source := map[string]any{}
	if len(payload) > 0 {
		_ = json.Unmarshal(payload, &source)
	}
	return source
}
