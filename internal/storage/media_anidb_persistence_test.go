package storage

import "testing"

func TestAniDBAuxiliaryMetadataPersistsSuccessfullyWithTMDBPrimary(t *testing.T) {
	ctx, store := testDBStore(t)
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type: "serie", Title: "Primary series", Monitored: true,
		ExternalProvider: stringPtr("tmdb"), ExternalID: stringPtr("77"),
		MediaMetadataSnapshot: MediaMetadataSnapshot{Seasons: []MediaSeason{{
			Name: "Season 1", SeasonNumber: 1, Monitored: true,
			Episodes: []MediaEpisode{{Name: "Episode 1", EpisodeNumber: 1, Monitored: true}},
		}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	seasons, err := store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	episodeID := seasons[0].Episodes[0].ID
	anidb := "anidb"
	primaryPoster, primaryDate := "tmdb.jpg", "2024-01-01"
	source := map[string]any{
		"source": "anidb_http_api", "externalId": "123",
		"artwork": "https://cdn-eu.anidb.net/images/main/456.jpg", "startDate": "2024-02-03",
	}
	absolute := int32(1)
	updated, err := store.UpdateMediaItemMetadata(ctx, item.ID, MediaItemInput{
		Type: "serie", ContentKind: "anime", Title: item.Title, Monitored: true,
		ExternalProvider: item.ExternalProvider, ExternalID: item.ExternalID, PosterPath: &primaryPoster,
		MediaMetadataSnapshot: MediaMetadataSnapshot{
			ReleaseDate: &primaryDate,
			Facts: []MediaFact{
				{Label: "AniDB artwork", Value: source["artwork"].(string), Source: source},
				{Label: "AniDB start date", Value: source["startDate"].(string), Source: source},
			},
		},
		ProviderMappings: []MediaProviderMappingInput{
			{ProviderName: "tmdb", ProviderEntityType: "serie", ExternalID: *item.ExternalID, Canonical: true},
			{ProviderName: "anidb", ProviderEntityType: "serie", ExternalID: "123", Source: source},
		},
		Aliases: []MediaItemAliasInput{{Alias: "AniDB title", Kind: "romaji", ProviderName: &anidb, Source: source}},
		EpisodeNumbering: []MediaEpisodeNumberingInput{{
			EpisodeID: episodeID, ProviderName: "anidb", NumberingScheme: "absolute", AbsoluteNumber: &absolute, Source: source,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.ExternalProvider == nil || *updated.ExternalProvider != "tmdb" || updated.ContentKind != "anime" || updated.NumberingStrategy == nil || *updated.NumberingStrategy != "anidb_absolute" {
		t.Fatalf("primary/anime state = %#v", updated)
	}
	if updated.PosterPath == nil || *updated.PosterPath != primaryPoster || updated.ReleaseDate == nil || *updated.ReleaseDate != primaryDate {
		t.Fatalf("primary artwork/date changed: %#v", updated.MediaMetadataSnapshot)
	}
	if len(updated.Facts) != 2 || updated.Facts[1].Source["externalId"] != "123" {
		t.Fatalf("AniDB facts provenance = %#v", updated.Facts)
	}
	if len(updated.ProviderMappings) != 2 || updated.ProviderMappings[1].Source["artwork"] != source["artwork"] {
		t.Fatalf("AniDB mapping provenance = %#v", updated.ProviderMappings)
	}
	aliasPersisted := false
	for _, alias := range updated.Aliases {
		aliasPersisted = aliasPersisted || alias.Alias == "AniDB title" && alias.ProviderName != nil && *alias.ProviderName == "anidb"
	}
	if !aliasPersisted {
		t.Fatalf("AniDB alias = %#v", updated.Aliases)
	}
	if len(updated.EpisodeNumbering) != 1 || updated.EpisodeNumbering[0].AbsoluteNumber == nil || *updated.EpisodeNumbering[0].AbsoluteNumber != 1 {
		t.Fatalf("AniDB numbering = %#v", updated.EpisodeNumbering)
	}
}
