package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestMediaSeriesRowsUpsertAndMonitorState(t *testing.T) {
	ctx, store := testDBStore(t)
	item := createSeriesItem(t, ctx, store)

	seasons, err := store.UpsertMediaSeriesSeasons(ctx, item.ID, []MediaSeriesSeasonInput{
		{
			ExternalProvider: stringPtr("tmdb"),
			ExternalID:       stringPtr("season-1"),
			SeasonNumber:     1,
			Name:             "Season 1",
			EpisodeCount:     int32Ptr(2),
			Monitored:        true,
			Source:           map[string]any{"provider": "tmdb"},
			Episodes: []MediaSeriesEpisodeInput{
				{ExternalID: stringPtr("e1"), EpisodeNumber: 1, Name: "Pilot", Monitored: true},
				{ExternalID: stringPtr("e2"), EpisodeNumber: 2, Name: "Second", Monitored: true},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(seasons) != 1 || len(seasons[0].Episodes) != 2 {
		t.Fatalf("expected one season with two episodes, got %#v", seasons)
	}
	seasonID := seasons[0].ID
	episodeID := seasons[0].Episodes[0].ID

	seasons, err = store.UpsertMediaSeriesSeasons(ctx, item.ID, []MediaSeriesSeasonInput{
		{
			SeasonNumber: 1,
			Name:         "Updated",
			Monitored:    true,
			Episodes: []MediaSeriesEpisodeInput{
				{EpisodeNumber: 1, Name: "Pilot Updated", Monitored: true},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if seasons[0].ID != seasonID || seasons[0].Episodes[0].ID != episodeID {
		t.Fatalf("expected stable ids, got season=%s episode=%s", seasons[0].ID, seasons[0].Episodes[0].ID)
	}
	if seasons[0].Name != "Updated" || seasons[0].Episodes[0].Name != "Pilot Updated" {
		t.Fatalf("expected updated rows, got %#v", seasons[0])
	}

	if _, err := store.SetMediaSeriesSeasonMonitored(ctx, seasonID, false); err != nil {
		t.Fatal(err)
	}
	seasons, err = store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if seasons[0].Monitored || seasons[0].Episodes[0].Monitored || seasons[0].Episodes[1].Monitored {
		t.Fatalf("expected season and episodes unmonitored, got %#v", seasons[0])
	}

	if _, err := store.SetMediaSeriesEpisodeMonitored(ctx, episodeID, true); err != nil {
		t.Fatal(err)
	}
	seasons, err = store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !seasons[0].Monitored || !seasons[0].Episodes[0].Monitored || seasons[0].Episodes[1].Monitored {
		t.Fatalf("expected parent to follow monitored episode, got %#v", seasons[0])
	}
}

func TestMediaSeriesRowsCascadeAndJsonCompatibility(t *testing.T) {
	ctx, store := testDBStore(t)
	item := createSeriesItem(t, ctx, store)
	if _, err := store.UpsertMediaSeriesSeasons(ctx, item.ID, []MediaSeriesSeasonInput{
		{SeasonNumber: 1, Name: "Season 1", Monitored: true},
	}); err != nil {
		t.Fatal(err)
	}

	got, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Seasons) != 1 || got.Seasons[0].Name != "Season 1" {
		t.Fatalf("expected relational seasons to drive response, got %#v", got.Seasons)
	}
	if err := store.DeleteMediaItem(ctx, item.ID, true); err != nil {
		t.Fatal(err)
	}
	seasons, err := store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(seasons) != 0 {
		t.Fatalf("expected cascade delete, got %#v", seasons)
	}
}

func TestMediaSeriesSnapshotRefreshPreservesMonitorState(t *testing.T) {
	ctx, store := testDBStore(t)
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:      "serie",
		Title:     "Refresh Series " + uuid.NewString(),
		Monitored: true,
		MediaMetadataSnapshot: MediaMetadataSnapshot{
			Seasons: []MediaSeason{{
				Name:         "Season 1",
				SeasonNumber: 1,
				Monitored:    true,
				Episodes: []MediaEpisode{
					{Name: "Pilot", EpisodeNumber: 1, Monitored: true},
					{Name: "Second", EpisodeNumber: 2, Monitored: true},
				},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	seasons, err := store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(seasons) != 1 || len(seasons[0].Episodes) != 2 {
		t.Fatalf("expected create to materialize rows, got %#v", seasons)
	}
	if _, err := store.SetMediaSeriesEpisodeMonitored(ctx, seasons[0].Episodes[0].ID, false); err != nil {
		t.Fatal(err)
	}

	refreshed, err := store.UpdateMediaItemMetadata(ctx, item.ID, MediaItemInput{
		Type:      "serie",
		Title:     item.Title,
		Monitored: true,
		MediaMetadataSnapshot: MediaMetadataSnapshot{
			Seasons: []MediaSeason{{
				Name:         "Season One",
				SeasonNumber: 1,
				Monitored:    true,
				Episodes: []MediaEpisode{
					{Name: "Pilot Updated", EpisodeNumber: 1, Monitored: true},
					{Name: "Second Updated", EpisodeNumber: 2, Monitored: true},
					{Name: "Third", EpisodeNumber: 3, Monitored: true},
				},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.Seasons[0].Name != "Season One" {
		t.Fatalf("expected season metadata update, got %#v", refreshed.Seasons[0])
	}
	episodes := refreshed.Seasons[0].Episodes
	if len(episodes) != 3 {
		t.Fatalf("expected added episode, got %#v", episodes)
	}
	if episodes[0].Name != "Pilot Updated" || episodes[0].Monitored {
		t.Fatalf("expected renamed episode with preserved monitor=false, got %#v", episodes[0])
	}
	if !episodes[1].Monitored || !episodes[2].Monitored {
		t.Fatalf("expected existing true and new true monitor state, got %#v", episodes)
	}
}

func createSeriesItem(t *testing.T, ctx context.Context, store *SettingsStore) MediaItem {
	t.Helper()
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:      "serie",
		Title:     "Series " + uuid.NewString(),
		Monitored: true,
		MediaMetadataSnapshot: MediaMetadataSnapshot{
			Seasons: []MediaSeason{{Name: "Legacy Season", Monitored: true}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return item
}
