package storage

import "testing"

func TestNormalizeMediaItemOptions(t *testing.T) {
	seriesType := "weird"
	input := normalizeMediaItemOptions(MediaItemInput{
		Type:                "serie",
		MonitorMode:         "invalid",
		SeriesType:          &seriesType,
		MinimumAvailability: "soon",
	})

	if input.MonitorMode != "all_episodes" || !input.Monitored {
		t.Fatalf("expected series monitor fallback with monitored=true, got mode=%q monitored=%v", input.MonitorMode, input.Monitored)
	}
	if input.SeriesType == nil || *input.SeriesType != "standard" {
		t.Fatalf("expected standard series type fallback, got %#v", input.SeriesType)
	}
	if input.MinimumAvailability != "released" {
		t.Fatalf("expected released availability fallback, got %q", input.MinimumAvailability)
	}
}

func TestNormalizeMediaItemOptionsForUnmonitoredMovie(t *testing.T) {
	input := normalizeMediaItemOptions(MediaItemInput{
		Type:        "movie",
		MonitorMode: "none",
	})

	if input.MonitorMode != "none" || input.Monitored {
		t.Fatalf("expected movie none monitor mode with monitored=false, got mode=%q monitored=%v", input.MonitorMode, input.Monitored)
	}
	if input.SeriesType != nil {
		t.Fatalf("expected movie series type to be cleared, got %#v", input.SeriesType)
	}
}

func TestNormalizeMediaRequestOptions(t *testing.T) {
	seriesType := "daily"
	input := normalizeMediaRequestOptions(MediaRequestInput{
		Type:                "serie",
		MonitorMode:         "future_episodes",
		SeriesType:          &seriesType,
		MinimumAvailability: "in_cinema",
	})

	if input.MonitorMode != "future_episodes" {
		t.Fatalf("expected request monitor mode to be preserved, got %q", input.MonitorMode)
	}
	if input.SeriesType == nil || *input.SeriesType != "daily" {
		t.Fatalf("expected daily series type, got %#v", input.SeriesType)
	}
	if input.MinimumAvailability != "in_cinema" {
		t.Fatalf("expected in_cinema availability, got %q", input.MinimumAvailability)
	}
}

func TestSCNMedia007SeasonMonitorPatchUpdatesEpisodes(t *testing.T) {
	season := MediaSeason{
		Name:      "Season 1",
		Monitored: true,
		Episodes: []MediaEpisode{
			{Name: "Pilot", EpisodeNumber: 1, Monitored: true},
			{Name: "Second", EpisodeNumber: 2, Monitored: true},
		},
	}

	monitored := false
	applySeasonMonitorPatch(&season, MediaItemOptionsInput{SeasonMonitored: &monitored})
	if season.Monitored {
		t.Fatal("expected season to be unmonitored")
	}
	for _, episode := range season.Episodes {
		if episode.Monitored {
			t.Fatalf("expected all episodes to be unmonitored, got %#v", season.Episodes)
		}
	}
}

func TestSCNMedia007EpisodeMonitorPatchUpdatesSeasonState(t *testing.T) {
	seasons := []MediaSeason{
		{
			Name:      "Season 1",
			Monitored: true,
			Episodes: []MediaEpisode{
				{Name: "Pilot", EpisodeNumber: 1, Monitored: true},
				{Name: "Second", EpisodeNumber: 2, Monitored: false},
			},
		},
	}

	seasonName := "Season 1"
	episodeNumber := int32(1)
	episodeMonitored := false
	next, updated := mediaItemUpdateSeasons(seasons, MediaItemOptionsInput{
		MonitorSeasonName:    &seasonName,
		MonitorEpisodeNumber: &episodeNumber,
		EpisodeMonitored:     &episodeMonitored,
	})

	if !updated || next == nil {
		t.Fatal("expected matching season to update")
	}
	if (*next)[0].Monitored || (*next)[0].Episodes[0].Monitored {
		t.Fatalf("expected season and first episode to be unmonitored, got %#v", (*next)[0])
	}
	if (*next)[0].Episodes[1].Monitored {
		t.Fatalf("expected second episode to remain unmonitored, got %#v", (*next)[0])
	}
	if seasons[0].Episodes[0].Monitored != true {
		t.Fatalf("expected original season slice to remain unchanged, got %#v", seasons[0])
	}
}

func TestSCNMedia007SeasonPayloadAndOptionalInputs(t *testing.T) {
	if payload, err := mediaItemSeasonsPayload(nil); err != nil || string(payload) != "[]" {
		t.Fatalf("nil seasons payload = %q err=%v", payload, err)
	}

	seasons := []MediaSeason{{Name: "Season 1", Monitored: true}}
	payload, err := mediaItemSeasonsPayload(&seasons)
	if err != nil {
		t.Fatalf("season payload: %v", err)
	}
	if string(payload) != `[{"name":"Season 1","monitored":true}]` {
		t.Fatalf("season payload = %s", payload)
	}

	blank := "  "
	if optionalTrimmed(&blank) != nil {
		t.Fatal("expected blank optional value to be omitted")
	}
	availability := "soon"
	if got := optionalMinimumAvailability(&availability); got == nil || *got != "released" {
		t.Fatalf("expected invalid availability fallback, got %#v", got)
	}
}
