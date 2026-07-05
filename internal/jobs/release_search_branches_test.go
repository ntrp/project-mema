package jobs

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func TestReleaseSearchBranchesExpandsMonitoredSeasons(t *testing.T) {
	item := storage.MediaItem{
		Type:  "serie",
		Title: "The Show",
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{
				{Name: "Season 1", Monitored: true},
				{Name: "Season 2", Monitored: false},
			},
		},
	}
	branches := releaseSearchBranches(item, decisions.SearchCriteriaForQuery(item, ""), "")
	if len(branches) != 1 {
		t.Fatalf("len(branches) = %d, want 1", len(branches))
	}
	if branches[0].criteria.Kind != "season" || *branches[0].criteria.SeasonNumber != 1 {
		t.Fatalf("branch criteria = %#v, want season 1", branches[0].criteria)
	}
}

func TestReleaseSearchBranchesCarriesPersistedSeasonAndEpisodeIDs(t *testing.T) {
	seasonID := uuid.New()
	episodeID := uuid.New()
	item := storage.MediaItem{
		Type:  "serie",
		Title: "The Show",
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{{
				ID:           &seasonID,
				Name:         "Season 1",
				SeasonNumber: 1,
				Monitored:    false,
				Episodes: []storage.MediaEpisode{{
					ID:            &episodeID,
					Name:          "Pilot",
					EpisodeNumber: 1,
					Monitored:     true,
				}},
			}},
		},
	}
	branches := releaseSearchBranches(item, decisions.SearchCriteriaForQuery(item, ""), "")
	if len(branches) != 1 {
		t.Fatalf("len(branches) = %d, want 1", len(branches))
	}
	criteria := branches[0].criteria
	if criteria.SeasonID == nil || *criteria.SeasonID != seasonID || criteria.EpisodeID == nil || *criteria.EpisodeID != episodeID {
		t.Fatalf("expected persisted ids in criteria, got %#v", criteria)
	}
}

func TestReleaseSearchBranchesExpandsMonitoredEpisodesWithoutSeasonMonitor(t *testing.T) {
	airDate := "2026-01-01"
	item := storage.MediaItem{
		Type:  "serie",
		Title: "The Show",
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{{
				Name: "Season 1",
				Episodes: []storage.MediaEpisode{
					{Name: "Pilot", EpisodeNumber: 1, AirDate: &airDate, Monitored: true},
					{Name: "Future", EpisodeNumber: 2, Monitored: false},
				},
			}},
		},
	}
	branches := releaseSearchBranches(item, decisions.SearchCriteriaForQuery(item, ""), "")
	if len(branches) != 1 {
		t.Fatalf("len(branches) = %d, want 1", len(branches))
	}
	if branches[0].criteria.Kind != "episode" || *branches[0].criteria.EpisodeNumber != 1 {
		t.Fatalf("branch criteria = %#v, want episode 1", branches[0].criteria)
	}
}
