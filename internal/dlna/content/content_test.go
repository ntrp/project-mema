package content

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"media-manager/internal/storage"

	"github.com/google/uuid"
)

type fakeSource struct {
	items []storage.MediaItem
}

func (s fakeSource) ListMediaItems(context.Context) ([]storage.MediaItem, error) {
	return s.items, nil
}

func TestObjectIDIsOpaqueAndStable(t *testing.T) {
	mediaID := uuid.New()
	path := "/library/Movies/Scenario Movie/Scenario Movie.mkv"

	id := EncodeID(FileRef(mediaID, path))
	again := EncodeID(FileRef(mediaID, path))

	if id != again {
		t.Fatalf("id not stable: %q != %q", id, again)
	}
	if strings.Contains(id, path) || strings.Contains(id, "Scenario") || strings.Contains(id, "library") {
		t.Fatalf("id exposes file path: %q", id)
	}
	ref, err := DecodeID(id)
	if err != nil {
		t.Fatal(err)
	}
	if ref.Kind != "file" || ref.Key != mediaID.String() || ref.Aux == "" {
		t.Fatalf("decoded ref = %#v", ref)
	}
}

func TestRootBrowseReturnsExpectedContainersAndSkipsMissingFiles(t *testing.T) {
	ctx := context.Background()
	collection := "Scenario Collection"
	year := int32(2026)
	moviePath := "/media/Scenario.Movie.2026.mkv"
	missingPath := "/media/Missing.Movie.2026.mkv"
	items := []storage.MediaItem{
		{
			ID:        uuid.New(),
			Type:      "movie",
			Title:     "Scenario Movie",
			Year:      &year,
			FilePaths: []string{moviePath},
			MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
				CollectionName: &collection,
				Genres:         []string{"Drama"},
			},
		},
		{
			ID:        uuid.New(),
			Type:      "movie",
			Title:     "Missing Movie",
			Year:      &year,
			FilePaths: []string{missingPath},
		},
	}
	tree := NewTree(fakeSource{items: items}).WithStat(fakeStat(moviePath))

	root, err := tree.BrowseChildren(ctx, RootID)
	if err != nil {
		t.Fatal(err)
	}
	assertContainer(t, root, "Movies", 1)
	assertContainer(t, root, "Collections", 1)
	assertContainer(t, root, "Genres", 1)
	assertContainer(t, root, "Years", 1)

	movies, err := tree.BrowseChildren(ctx, EncodeID(RootContainerRef("movies")))
	if err != nil {
		t.Fatal(err)
	}
	if len(movies) != 1 || movies[0].Title != "Scenario Movie (2026)" || movies[0].ChildCount != 1 {
		t.Fatalf("movie browse = %#v", movies)
	}
	files, err := tree.BrowseChildren(ctx, movies[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].FilePath != moviePath || strings.Contains(files[0].ID, moviePath) {
		t.Fatalf("file browse = %#v", files)
	}
}

func TestSeriesBrowseMapsSeasonsEpisodesAndFiles(t *testing.T) {
	ctx := context.Background()
	year := int32(2026)
	seasonID := uuid.New()
	episodeID := uuid.New()
	path := "/media/Scenario.Show/S01E02.mkv"
	item := storage.MediaItem{
		ID:        uuid.New(),
		Type:      "serie",
		Title:     "Scenario Show",
		Year:      &year,
		FilePaths: []string{path},
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{{
				ID:           &seasonID,
				Name:         "Season 1",
				SeasonNumber: 1,
				Episodes: []storage.MediaEpisode{{
					ID:            &episodeID,
					Name:          "Second",
					EpisodeNumber: 2,
				}},
			}},
		},
	}
	tree := NewTree(fakeSource{items: []storage.MediaItem{item}}).WithStat(fakeStat(path))

	shows, err := tree.BrowseChildren(ctx, EncodeID(RootContainerRef("series")))
	if err != nil {
		t.Fatal(err)
	}
	if len(shows) != 1 || shows[0].ChildCount != 1 {
		t.Fatalf("shows = %#v", shows)
	}
	seasons, err := tree.BrowseChildren(ctx, shows[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(seasons) != 1 || seasons[0].Title != "Season 1" || seasons[0].ChildCount != 1 {
		t.Fatalf("seasons = %#v", seasons)
	}
	episodes, err := tree.BrowseChildren(ctx, seasons[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(episodes) != 1 || episodes[0].Title != "Second" || episodes[0].ChildCount != 1 {
		t.Fatalf("episodes = %#v", episodes)
	}
	files, err := tree.BrowseChildren(ctx, episodes[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Title != filepath.Base(path) {
		t.Fatalf("files = %#v", files)
	}
}

func assertContainer(t *testing.T, objects []Object, title string, childCount int) {
	t.Helper()
	for _, object := range objects {
		if object.Title == title {
			if object.Kind != ObjectContainer || object.ChildCount != childCount {
				t.Fatalf("%s container = %#v", title, object)
			}
			return
		}
	}
	t.Fatalf("missing container %q in %#v", title, objects)
}

func fakeStat(paths ...string) FileStatFunc {
	available := map[string]struct{}{}
	for _, path := range paths {
		available[path] = struct{}{}
	}
	return func(path string) (os.FileInfo, error) {
		if _, ok := available[path]; !ok {
			return nil, os.ErrNotExist
		}
		return fileInfo{name: filepath.Base(path), size: 100, isDir: false}, nil
	}
}
