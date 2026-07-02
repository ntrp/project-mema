package jobs

import (
	"testing"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func TestShouldFallbackEpisodeToSeasonWhenNoReleases(t *testing.T) {
	item, criteria := episodeFallbackFixture()
	if !shouldFallbackEpisodeToSeason(item, criteria, nil, nil, nil, nil) {
		t.Fatal("expected episode search with no releases to fall back to season search")
	}
}

func TestShouldFallbackEpisodeToSeasonWhenOnlyRejectedReleases(t *testing.T) {
	item, criteria := episodeFallbackFixture()
	releases := []storage.ReleaseCandidateInput{{
		Title:            "The.Show.S01E02.1080p.WEB-DL",
		SearchKind:       "episode",
		RequestedSeason:  criteria.SeasonNumber,
		RequestedEpisode: criteria.EpisodeNumber,
	}}
	if !shouldFallbackEpisodeToSeason(item, criteria, nil, nil, nil, releases) {
		t.Fatal("expected rejected episode releases to fall back to season search")
	}
}

func TestShouldNotFallbackEpisodeToSeasonWhenEpisodeMatches(t *testing.T) {
	item, criteria := episodeFallbackFixture()
	releases := []storage.ReleaseCandidateInput{{
		Title:            "The.Show.S01E01.1080p.WEB-DL",
		SearchKind:       "episode",
		RequestedSeason:  criteria.SeasonNumber,
		RequestedEpisode: criteria.EpisodeNumber,
	}}
	if shouldFallbackEpisodeToSeason(item, criteria, nil, nil, nil, releases) {
		t.Fatal("did not expect matching episode release to fall back to season search")
	}
}

func TestSeasonFallbackCriteriaDropsEpisode(t *testing.T) {
	_, criteria := episodeFallbackFixture()
	seasonCriteria := seasonFallbackCriteria(criteria)
	if seasonCriteria.Kind != "season" {
		t.Fatalf("kind = %q, want season", seasonCriteria.Kind)
	}
	if seasonCriteria.EpisodeNumber != nil {
		t.Fatalf("episode number = %d, want nil", *seasonCriteria.EpisodeNumber)
	}
	if seasonCriteria.SeasonNumber == nil || *seasonCriteria.SeasonNumber != 1 {
		t.Fatalf("season number = %#v, want 1", seasonCriteria.SeasonNumber)
	}
}

func episodeFallbackFixture() (storage.MediaItem, decisions.ReleaseSearchCriteria) {
	season := int32(1)
	episode := int32(1)
	item := storage.MediaItem{Type: "series", Title: "The Show"}
	criteria := decisions.ReleaseSearchCriteria{
		Kind:          "episode",
		Title:         "The Show",
		SeasonNumber:  &season,
		EpisodeNumber: &episode,
	}
	return item, criteria
}
