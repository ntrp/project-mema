package content

import (
	"context"
	"errors"
	"testing"

	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func TestParseCriteriaSupportsKnownFields(t *testing.T) {
	criteria, err := ParseCriteria(`dc:title contains "scenario" and upnp:class derivedfrom "object.item.videoItem"`)
	if err != nil {
		t.Fatal(err)
	}
	if len(criteria) != 2 || criteria[0].Field != "dc:title" || criteria[1].Op != "derivedfrom" {
		t.Fatalf("criteria = %#v", criteria)
	}
}

func TestParseCriteriaRejectsUnsupportedFields(t *testing.T) {
	_, err := ParseCriteria(`res@duration > "0:01:00"`)
	if err == nil {
		t.Fatal("expected unsupported criteria error")
	}
}

func TestSearchReturnsMatchingMoviesAndEpisodes(t *testing.T) {
	ctx := context.Background()
	moviePath := "/media/Scenario.Movie.mkv"
	episodePath := "/media/Scenario.Show/S01E02.mkv"
	airDate := "2026-07-07"
	episodeID := uuid.New()
	seasonID := uuid.New()
	tree := NewTree(fakeSource{items: []storage.MediaItem{
		{
			ID:        uuid.New(),
			Type:      "movie",
			Title:     "Scenario Movie",
			FilePaths: []string{moviePath},
			MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
				Genres: []string{"Drama"},
				Cast:   []storage.MediaPerson{{Name: "Actor One"}},
			},
		},
		{
			ID:        uuid.New(),
			Type:      "serie",
			Title:     "Scenario Show",
			FilePaths: []string{episodePath},
			MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
				Seasons: []storage.MediaSeason{{
					ID:           &seasonID,
					Name:         "Season 1",
					SeasonNumber: 1,
					Episodes: []storage.MediaEpisode{{
						ID:            &episodeID,
						Name:          "Pilot",
						EpisodeNumber: 2,
						AirDate:       &airDate,
					}},
				}},
			},
		},
	}}).WithStat(fakeStat(moviePath, episodePath))

	movies, err := tree.Search(ctx, SearchRequest{
		ContainerID:    RootID,
		SearchCriteria: `upnp:genre = "Drama" and dc:creator contains "actor"`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if movies.TotalMatches != 1 || movies.Objects[0].Title != "Scenario Movie" {
		t.Fatalf("movie search = %#v", movies)
	}
	episodes, err := tree.Search(ctx, SearchRequest{
		ContainerID:    RootID,
		SearchCriteria: `upnp:class derivedfrom "object.item.videoItem.episode" and dc:date >= "2026"`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if episodes.TotalMatches != 1 || episodes.Objects[0].Title != "Pilot" {
		t.Fatalf("episode search = %#v", episodes)
	}
}

func TestSearchInvalidContainerReturnsNotFound(t *testing.T) {
	tree := NewTree(fakeSource{})

	_, err := tree.Search(context.Background(), SearchRequest{
		ContainerID:    "bad-container",
		SearchCriteria: "*",
	})
	if !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("err = %v", err)
	}
}
