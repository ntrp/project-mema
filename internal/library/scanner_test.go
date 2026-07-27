package library

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSCNMedia010DiscoverClassifiesLibraryFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "Movies/Scenario.Movie.2026.mkv")
	writeFile(t, root, "Shows/Scenario Series/Season 1/Scenario.Series.S01E02.mkv")
	writeFile(t, root, "Unsafe/Other.Movie.2026.1080p.WEBDL.mkv")
	writeFile(t, root, ".hidden/Hidden.Movie.2026.mkv")
	writeFile(t, root, "Notes/readme.txt")

	files, err := Discover(root)
	if err != nil {
		t.Fatalf("discover failed: %v", err)
	}
	if len(files) != 3 {
		t.Fatalf("discovered %d files, want 3: %#v", len(files), files)
	}

	byPath := map[string]DiscoveredFile{}
	for _, file := range files {
		byPath[file.Path] = file
	}

	movie := byPath["Movies/Scenario.Movie.2026.mkv"]
	if movie.DetectedKind != MediaKindMovie || movie.DetectedTitle != "Scenario Movie" {
		t.Fatalf("movie discovery = %#v", movie)
	}
	if movie.DetectedYear == nil || *movie.DetectedYear != 2026 {
		t.Fatalf("movie year = %#v, want 2026", movie.DetectedYear)
	}
	if !movie.SafeMatch {
		t.Fatalf("movie should be a safe match: %#v", movie)
	}
	if movie.SizeBytes != int64(len("fixture")) {
		t.Fatalf("movie size = %d, want fixture size", movie.SizeBytes)
	}

	episode := byPath["Shows/Scenario Series/Season 1/Scenario.Series.S01E02.mkv"]
	if episode.DetectedKind != MediaKindSeries || episode.DetectedTitle != "Scenario Series" {
		t.Fatalf("episode discovery = %#v", episode)
	}
	if !episode.SafeMatch {
		t.Fatalf("episode should be a safe match: %#v", episode)
	}
	if episode.SeasonNumber == nil || *episode.SeasonNumber != 1 || episode.EpisodeNumber == nil || *episode.EpisodeNumber != 2 {
		t.Fatalf("episode numbers = S%vE%v, want S1E2", episode.SeasonNumber, episode.EpisodeNumber)
	}

	unsafe := byPath["Unsafe/Other.Movie.2026.1080p.WEBDL.mkv"]
	if unsafe.SafeMatch {
		t.Fatalf("release-token file should not be a safe match: %#v", unsafe)
	}
}

func TestDiscoverSeriesGroupsNestedEpisodesByTopLevelFolder(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "Scenario Series (2026)/Season 1/Scenario.Series.S01E02.mkv")
	writeFile(t, root, "Scenario Series (2026)/Season 2/Nested/Scenario.Series.S02E03.mkv")
	writeFile(t, root, "Scenario Series (2026)/Season 3/04 - Fourth Episode.mkv")
	writeFile(t, root, "Other Series/Other.Series.S01E01.mkv")

	files, err := DiscoverSeries(root)
	if err != nil {
		t.Fatalf("discover series failed: %v", err)
	}
	if len(files) != 4 {
		t.Fatalf("discovered %d files, want 4: %#v", len(files), files)
	}
	byPath := map[string]DiscoveredFile{}
	for _, file := range files {
		byPath[file.Path] = file
	}
	for _, path := range []string{
		"Scenario Series (2026)/Season 1/Scenario.Series.S01E02.mkv",
		"Scenario Series (2026)/Season 2/Nested/Scenario.Series.S02E03.mkv",
		"Scenario Series (2026)/Season 3/04 - Fourth Episode.mkv",
	} {
		file := byPath[path]
		if file.SeriesFolderPath != "Scenario Series (2026)" || file.DetectedTitle != "Scenario Series" {
			t.Fatalf("series folder discovery = %#v", file)
		}
		if file.DetectedYear == nil || *file.DetectedYear != 2026 || file.DetectedKind != MediaKindSeries {
			t.Fatalf("series metadata = %#v", file)
		}
	}
	nested := byPath["Scenario Series (2026)/Season 2/Nested/Scenario.Series.S02E03.mkv"]
	if nested.SeasonNumber == nil || *nested.SeasonNumber != 2 || nested.EpisodeNumber == nil || *nested.EpisodeNumber != 3 {
		t.Fatalf("nested episode = %#v, want S02E03", nested)
	}
	folderEpisode := byPath["Scenario Series (2026)/Season 3/04 - Fourth Episode.mkv"]
	if folderEpisode.SeasonNumber == nil || *folderEpisode.SeasonNumber != 3 || folderEpisode.EpisodeNumber == nil || *folderEpisode.EpisodeNumber != 4 {
		t.Fatalf("folder episode = %#v, want S03E04", folderEpisode)
	}
}

func TestDiscoverUsesFolderTitlesWithoutExternalIDs(t *testing.T) {
	movieRoot := t.TempDir()
	writeFile(t, movieRoot, "The Matrix (1999) [imdbid-tt0133093]/video.1080p.mkv")
	movies, err := DiscoverMovies(movieRoot)
	if err != nil || len(movies) != 1 {
		t.Fatalf("discover movies = %#v, %v", movies, err)
	}
	if movies[0].DetectedTitle != "The Matrix" || movies[0].DetectedYear == nil || *movies[0].DetectedYear != 1999 {
		t.Fatalf("movie metadata = %#v", movies[0])
	}

	seriesRoot := t.TempDir()
	writeFile(t, seriesRoot, "Dark (2017) {tmdb-70523}/Season 1/01.mkv")
	series, err := DiscoverSeries(seriesRoot)
	if err != nil || len(series) != 1 {
		t.Fatalf("discover series = %#v, %v", series, err)
	}
	if series[0].DetectedTitle != "Dark" || series[0].DetectedYear == nil || *series[0].DetectedYear != 2017 {
		t.Fatalf("series metadata = %#v", series[0])
	}
}

func TestDiscoverMoviesTreatsEpisodeNamedFilesAsMovies(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "Movies/Scenario.Movie.S01E02.2026.mkv")

	files, err := DiscoverMovies(root)
	if err != nil {
		t.Fatalf("discover movies failed: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("discovered %d files, want 1: %#v", len(files), files)
	}
	file := files[0]
	if file.DetectedKind != MediaKindMovie {
		t.Fatalf("detected kind = %s, want movie: %#v", file.DetectedKind, file)
	}
	if file.SeasonNumber != nil || file.EpisodeNumber != nil {
		t.Fatalf("movie scan should not carry episode numbers: %#v", file)
	}
}

func TestDiscoverStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := DiscoverMoviesContext(ctx, t.TempDir())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("DiscoverMoviesContext error = %v, want context canceled", err)
	}
}

func writeFile(t *testing.T, root string, relativePath string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %q: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte("fixture"), 0o644); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}
