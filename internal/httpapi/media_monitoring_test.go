package httpapi

import (
	"testing"
	"time"

	"media-manager/internal/storage"
)

func TestSCNMedia007SeriesMonitoringDefaultsToAllEpisodes(t *testing.T) {
	past := "2026-01-01"
	input := storage.MediaItemInput{
		Type: "serie",
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{{
				Name: "Season 1",
				Episodes: []storage.MediaEpisode{{
					Name:          "Pilot",
					EpisodeNumber: 1,
					AirDate:       &past,
				}},
			}},
		},
	}

	got := applySeriesMonitoring(input)

	if got.MonitorMode != "all_episodes" {
		t.Fatalf("MonitorMode = %q, want all_episodes", got.MonitorMode)
	}
	if !got.Seasons[0].Monitored || !got.Seasons[0].Episodes[0].Monitored {
		t.Fatalf("expected season and episode to be monitored: %#v", got.Seasons[0])
	}
}

func TestSCNMedia007SeriesMonitoringFutureEpisodesSkipsPastAndUnknown(t *testing.T) {
	past := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	future := time.Now().AddDate(0, 0, 2).Format("2006-01-02")
	input := storage.MediaItemInput{
		Type:        "serie",
		MonitorMode: "future_episodes",
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{{
				Name: "Season 1",
				Episodes: []storage.MediaEpisode{
					{Name: "Past", EpisodeNumber: 1, AirDate: &past},
					{Name: "Future", EpisodeNumber: 2, AirDate: &future},
					{Name: "Unknown", EpisodeNumber: 3},
				},
			}},
		},
	}

	got := applySeriesMonitoring(input)

	if !got.Seasons[0].Monitored {
		t.Fatalf("expected season monitored when at least one episode is future: %#v", got.Seasons[0])
	}
	if got.Seasons[0].Episodes[0].Monitored {
		t.Fatal("expected past episode to be unmonitored")
	}
	if !got.Seasons[0].Episodes[1].Monitored {
		t.Fatal("expected future episode to be monitored")
	}
	if got.Seasons[0].Episodes[2].Monitored {
		t.Fatal("expected unknown-air-date episode to be unmonitored")
	}
}

func TestSCNMedia007NoSpecialsLeavesSpecialSeasonsUnmonitored(t *testing.T) {
	input := storage.MediaItemInput{
		Type:        "serie",
		MonitorMode: "no_specials",
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{
				{Name: "Specials", Episodes: []storage.MediaEpisode{{Name: "Bonus", EpisodeNumber: 1}}},
				{Name: "Season 1", Episodes: []storage.MediaEpisode{{Name: "Pilot", EpisodeNumber: 1}}},
			},
		},
	}

	got := applySeriesMonitoring(input)

	if got.Seasons[0].Monitored || got.Seasons[0].Episodes[0].Monitored {
		t.Fatalf("expected specials to stay unmonitored: %#v", got.Seasons[0])
	}
	if !got.Seasons[1].Monitored || !got.Seasons[1].Episodes[0].Monitored {
		t.Fatalf("expected regular season to be monitored: %#v", got.Seasons[1])
	}
}
