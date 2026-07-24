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
