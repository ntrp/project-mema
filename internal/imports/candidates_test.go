package imports

import (
	"os"
	"path/filepath"
	"testing"

	"media-manager/internal/downloadclients"
	"media-manager/internal/storage"
)

func TestSelectCompletedDownloadCandidatesPrefersLargestVideo(t *testing.T) {
	root := t.TempDir()
	small := writeSizedFile(t, root, "movie-sample.mkv", 10)
	main := writeSizedFile(t, root, "movie-main.mkv", 100)
	notes := writeSizedFile(t, root, "notes.txt", 1)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: small, Complete: true},
		{Path: main, Complete: true},
		{Path: notes, Complete: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection, main)
	assertRejectedReason(t, selection, small, "lower_scoring_candidate")
	assertRejectedReason(t, selection, notes, "not_video_file")
}

func TestSelectCompletedDownloadCandidatesFindsVideosInDirectory(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "release")
	main := writeSizedFile(t, dir, "Feature.Movie.mkv", 200)
	extra := writeSizedFile(t, dir, filepath.Join("extras", "behind-the-scenes.mp4"), 20)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: dir, Complete: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection, main)
	assertRejectedReason(t, selection, extra, "lower_scoring_candidate")
}

func TestSelectCompletedDownloadCandidatesUsesReportedSizeAndStableTieBreak(t *testing.T) {
	root := t.TempDir()
	alpha := writeSizedFile(t, root, "alpha.mp4", 1)
	zeta := writeSizedFile(t, root, "zeta.mp4", 1)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: zeta, SizeBytes: 500, Complete: true},
		{Path: alpha, SizeBytes: 500, Complete: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection, alpha)
	assertRejectedReason(t, selection, zeta, "lower_scoring_candidate")
}

func TestSelectCompletedDownloadCandidatesMapsPathsAndRejectsIncomplete(t *testing.T) {
	root := t.TempDir()
	appRoot := filepath.Join(root, "client")
	source := writeSizedFile(t, appRoot, "movie.mkv", 100)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: "/client/movie.mkv", Complete: true},
		{Path: "/client/partial.mkv", Complete: false},
	}, []storage.PathMapping{{ClientPath: "/client", AppPath: appRoot}})
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection, source)
	assertRejectedReason(t, selection, "/client/partial.mkv", "incomplete")
}

func assertSelectedSources(t *testing.T, selection completedDownloadSelection, paths ...string) {
	t.Helper()
	if len(selection.SelectedSources) != len(paths) {
		t.Fatalf("selected = %#v, want %#v", selection.SelectedSources, paths)
	}
	for index, path := range paths {
		if selection.SelectedSources[index] != path {
			t.Fatalf("selected = %#v, want %#v", selection.SelectedSources, paths)
		}
	}
}

func assertRejectedReason(t *testing.T, selection completedDownloadSelection, path string, reason string) {
	t.Helper()
	for _, rejected := range selection.RejectedCandidates {
		if rejected.SourcePath == path && rejected.Reason == reason {
			return
		}
	}
	t.Fatalf("rejection %s/%s not found in %#v", path, reason, selection.RejectedCandidates)
}

func writeSizedFile(t *testing.T, root string, relativePath string, size int64) string {
	t.Helper()
	path := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
