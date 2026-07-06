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
