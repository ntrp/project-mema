package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestMediaAnimeMetadataHydratesMappingsAliasesAndNumbering(t *testing.T) {
	ctx, store := testDBStore(t)
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:             "serie",
		ContentKind:      "anime",
		Title:            "Frieren: Beyond Journey's End",
		Monitored:        true,
		ExternalProvider: stringPtr("tmdb"),
		ExternalID:       stringPtr("209867"),
		Aliases: []MediaItemAliasInput{
			{Alias: "Sousou no Frieren", Kind: "romaji", ProviderName: stringPtr("anilist")},
			{Alias: "Frieren", Kind: "release_title", ProviderName: stringPtr("anidb")},
		},
		MediaMetadataSnapshot: MediaMetadataSnapshot{
			Seasons: []MediaSeason{{
				Name:         "Season 1",
				SeasonNumber: 1,
				Monitored:    true,
				Episodes: []MediaEpisode{{
					Name:          "The Journey's End",
					EpisodeNumber: 1,
					Monitored:     true,
				}},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.ContentKind != "anime" || item.NumberingStrategy == nil || *item.NumberingStrategy != "anidb_absolute" {
		t.Fatalf("expected anime absolute state, got %#v", item)
	}
	if len(item.ProviderMappings) != 1 || item.ProviderMappings[0].ProviderName != "tmdb" || !item.ProviderMappings[0].Canonical {
		t.Fatalf("expected canonical tmdb mapping, got %#v", item.ProviderMappings)
	}
	if len(item.Aliases) != 2 || item.Aliases[0].NormalizedAlias == "" {
		t.Fatalf("expected provider aliases, got %#v", item.Aliases)
	}
	seasons, err := store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	episode := seasons[0].Episodes[0]
	if err := upsertAnimeMetadata(ctx, store.pool, item.ID, MediaItemInput{
		EpisodeNumbering: []MediaEpisodeNumberingInput{{
			SeasonID:        &episode.SeasonID,
			EpisodeID:       episode.ID,
			ProviderName:    "anidb",
			NumberingScheme: "absolute",
			AbsoluteNumber:  int32Ptr(1),
		}},
	}); err != nil {
		t.Fatal(err)
	}
	item, err = store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(item.EpisodeNumbering) != 1 || *item.EpisodeNumbering[0].AbsoluteNumber != 1 {
		t.Fatalf("expected absolute numbering, got %#v", item.EpisodeNumbering)
	}
}

func TestAniListMetadataPersistsAlongsideTMDBPrimaryAtomically(t *testing.T) {
	ctx, store := testDBStore(t)
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type: "serie", Title: "Frieren: Beyond Journey's End", Monitored: true,
		ExternalProvider: stringPtr("tmdb"), ExternalID: stringPtr("209867"),
	})
	if err != nil {
		t.Fatal(err)
	}
	input := MediaItemInput{
		Type: "serie", ContentKind: "anime", Title: item.Title, Monitored: true,
		ExternalProvider: item.ExternalProvider, ExternalID: item.ExternalID,
		ProviderMappings: []MediaProviderMappingInput{
			{EntityType: "media_item", ProviderName: "tmdb", ProviderEntityType: "serie", ExternalID: "209867", Canonical: true},
			{EntityType: "media_item", ProviderName: "anilist", ProviderEntityType: "serie", ExternalID: "154587", Source: map[string]any{"source": "anilist", "raw": map[string]any{"id": float64(154587)}}},
		},
		Aliases: []MediaItemAliasInput{{Alias: "Sousou no Frieren", Language: stringPtr("ja-Latn"), Kind: "romaji", ProviderName: stringPtr("anilist"), Source: map[string]any{"field": "title.romaji"}}},
	}
	updated, err := store.UpdateMediaItemMetadata(ctx, item.ID, input)
	if err != nil {
		t.Fatal(err)
	}
	if updated.ExternalProvider == nil || *updated.ExternalProvider != "tmdb" || len(updated.ProviderMappings) != 2 || len(updated.Aliases) != 2 {
		t.Fatalf("updated AniList metadata = %#v", updated)
	}
	if updated.ProviderMappings[1].ProviderName != "anilist" || updated.ProviderMappings[1].Source["raw"] == nil {
		t.Fatalf("AniList provenance = %#v", updated.ProviderMappings)
	}

	input.Title = "Should Roll Back"
	input.ProviderMappings = append(input.ProviderMappings, MediaProviderMappingInput{
		EntityType: "media_item", ProviderName: "unsupported", ProviderEntityType: "serie", ExternalID: "1",
	})
	if _, err := store.UpdateMediaItemMetadata(ctx, item.ID, input); err == nil {
		t.Fatal("expected invalid provider mapping to fail")
	}
	afterFailure, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if afterFailure.Title != item.Title || len(afterFailure.ProviderMappings) != 2 {
		t.Fatalf("failed update partially persisted: %#v", afterFailure)
	}
}

func TestUpdateMediaItemMetadataRollsBackAniDBArtifactsTogether(t *testing.T) {
	ctx, store := testDBStore(t)
	item := createSeriesItem(t, ctx, store)
	provider := "anidb"
	_, err := store.UpdateMediaItemMetadata(ctx, item.ID, MediaItemInput{
		Type: item.Type, Title: item.Title, Year: item.Year, Monitored: item.Monitored,
		ExternalProvider: item.ExternalProvider, ExternalID: item.ExternalID,
		MediaMetadataSnapshot: item.MediaMetadataSnapshot,
		ProviderMappings: []MediaProviderMappingInput{{
			ProviderName: "anidb", ProviderEntityType: "serie", ExternalID: "123",
			Source: map[string]any{"source": "anidb_http_api"},
		}},
		Aliases: []MediaItemAliasInput{{
			Alias: "Uncommitted alias", Kind: "synonym", ProviderName: &provider,
		}},
		EpisodeNumbering: []MediaEpisodeNumberingInput{{
			EpisodeID: uuid.New(), ProviderName: "anidb", NumberingScheme: "absolute", AbsoluteNumber: int32Ptr(1),
		}},
	})
	if err == nil {
		t.Fatal("expected invalid episode mapping to fail")
	}
	got, getErr := store.GetMediaItem(ctx, item.ID)
	if getErr != nil {
		t.Fatal(getErr)
	}
	for _, mapping := range got.ProviderMappings {
		if mapping.ProviderName == "anidb" {
			t.Fatalf("partial mapping persisted: %#v", mapping)
		}
	}
	for _, alias := range got.Aliases {
		if alias.Alias == "Uncommitted alias" {
			t.Fatalf("partial alias persisted: %#v", alias)
		}
	}
}

func TestAnimeMetadataProviderDefaultsAndMediaTypeFiltering(t *testing.T) {
	ctx, store := testDBStore(t)
	if err := store.EnsureDefaultMetadataProviders(ctx); err != nil {
		t.Fatal(err)
	}
	providers, err := store.ListMetadataProviders(ctx)
	if err != nil {
		t.Fatal(err)
	}
	byType := make(map[string]MetadataProvider, len(providers))
	for _, provider := range providers {
		byType[provider.Type] = provider
	}
	if !byType["anilist"].Enabled || byType["anidb"].Enabled {
		t.Fatalf("anime provider defaults = %#v", byType)
	}
	client, version := "client1", "1"
	anidbProvider := byType["anidb"]
	if _, err := store.UpdateMetadataProvider(ctx, anidbProvider.ID, MetadataProviderInput{
		Name: anidbProvider.Name, Type: anidbProvider.Type, BaseURL: anidbProvider.BaseURL,
		APIKey: &client, PIN: &version, Enabled: true, Priority: anidbProvider.Priority,
	}); err != nil {
		t.Fatal(err)
	}
	seriesProviders, err := store.ListEnabledMetadataProviders(ctx, "serie")
	if err != nil {
		t.Fatal(err)
	}
	movieProviders, err := store.ListEnabledMetadataProviders(ctx, "movie")
	if err != nil {
		t.Fatal(err)
	}
	if !hasMetadataProviderType(seriesProviders, "anilist") || !hasMetadataProviderType(seriesProviders, "anidb") {
		t.Fatalf("series providers = %#v", seriesProviders)
	}
	if !hasMetadataProviderType(movieProviders, "anilist") || hasMetadataProviderType(movieProviders, "anidb") {
		t.Fatalf("movie providers = %#v", movieProviders)
	}
}

func hasMetadataProviderType(providers []MetadataProvider, providerType string) bool {
	for _, provider := range providers {
		if provider.Type == providerType {
			return true
		}
	}
	return false
}

func TestUnconfiguredAniDBProviderIsDisabled(t *testing.T) {
	name, badVersion, version := "client1", "invalid", "2"
	invalidNames := []string{"Client", "client name", "client_name", " client"}
	inputs := []MetadataProviderInput{
		{Type: "anidb", Enabled: true},
		{Type: "anidb", APIKey: &name, PIN: &badVersion, Enabled: true},
	}
	for index := range invalidNames {
		inputs = append(inputs, MetadataProviderInput{Type: "anidb", APIKey: &invalidNames[index], PIN: &version, Enabled: true})
	}
	for _, input := range inputs {
		if disableUnconfiguredAniDB(input).Enabled {
			t.Fatalf("unconfigured provider remained enabled: %#v", input)
		}
	}
	input := MetadataProviderInput{Type: "anidb", APIKey: &name, PIN: &version, Enabled: true}
	if !disableUnconfiguredAniDB(input).Enabled {
		t.Fatal("configured AniDB provider was disabled")
	}
}

func TestAniDBProviderPersistenceDisablesInvalidClientName(t *testing.T) {
	ctx, store := testDBStore(t)
	invalidName, validName, version := "InvalidClient", "validclient1", "1"
	provider, err := store.CreateMetadataProvider(ctx, MetadataProviderInput{
		Name: "AniDB validation " + uuid.NewString(), Type: "anidb", BaseURL: "http://api.anidb.net:9001",
		APIKey: &invalidName, PIN: &version, Enabled: true, Priority: 120,
	})
	if err != nil {
		t.Fatal(err)
	}
	if provider.Enabled {
		t.Fatal("provider with uppercase client name was persisted as enabled")
	}
	provider, err = store.UpdateMetadataProvider(ctx, provider.ID, MetadataProviderInput{
		Name: provider.Name, Type: "anidb", BaseURL: provider.BaseURL,
		APIKey: &validName, PIN: &version, Enabled: true, Priority: provider.Priority,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !provider.Enabled {
		t.Fatal("provider with valid client name was not enabled")
	}
}

func TestMediaKindToMediaTypeAndContentMapsAnimeHints(t *testing.T) {
	mediaType, contentKind, ok := mediaKindToMediaTypeAndContent("anime_series")
	if !ok || mediaType != "serie" || contentKind != "anime" {
		t.Fatalf("anime series mapping = %q %q %v", mediaType, contentKind, ok)
	}
	mediaType, contentKind, ok = mediaKindToMediaTypeAndContent("anime_movie")
	if !ok || mediaType != "movie" || contentKind != "anime" {
		t.Fatalf("anime movie mapping = %q %q %v", mediaType, contentKind, ok)
	}
}

func createAnimeEpisodeNumberingFixture(t *testing.T, ctx context.Context, store *SettingsStore) (MediaItem, uuid.UUID) {
	t.Helper()
	item := createSeriesItem(t, ctx, store)
	seasons, err := store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	return item, seasons[0].Episodes[0].ID
}
