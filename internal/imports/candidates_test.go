package imports

import (
	"os"
	"path/filepath"
	"testing"

	"media-manager/internal/downloadclients"
	"media-manager/internal/storage"
)

const testMiB = 1024 * 1024

func TestSelectCompletedDownloadCandidatesPrefersLargestVideo(t *testing.T) {
	root := t.TempDir()
	small := writeSizedFile(t, root, "movie-sample.mkv", 700*testMiB)
	main := writeSizedFile(t, root, "movie-main.mkv", 900*testMiB)
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
	main := writeSizedFile(t, dir, "Feature.Movie.mkv", 900*testMiB)
	extra := writeSizedFile(t, dir, filepath.Join("extras", "behind-the-scenes.mp4"), 100*testMiB)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: dir, Complete: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection, main)
	assertRejectedReason(t, selection, extra, "sample_or_extra")
}

func TestSelectCompletedDownloadCandidatesUsesReportedSizeAndStableTieBreak(t *testing.T) {
	root := t.TempDir()
	alpha := writeSizedFile(t, root, "alpha.mp4", 1)
	zeta := writeSizedFile(t, root, "zeta.mp4", 1)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: zeta, SizeBytes: 500 * testMiB, Complete: true},
		{Path: alpha, SizeBytes: 500 * testMiB, Complete: true},
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
	source := writeSizedFile(t, appRoot, "movie.mkv", 500*testMiB)

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

func TestSelectCompletedDownloadCandidatesDoesNotRejectKeywordWithoutSizeEvidence(t *testing.T) {
	root := t.TempDir()
	source := writeSizedFile(t, root, "Sample.Collection.2026.mkv", 700*testMiB)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: source, Complete: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection, source)
}

func TestSelectCompletedDownloadCandidatesRejectsTinyFiles(t *testing.T) {
	root := t.TempDir()
	source := writeSizedFile(t, root, "tiny.mkv", 1*testMiB)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: source, Complete: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection)
	assertRejectedReason(t, selection, source, "tiny_file")
}

func TestSelectCompletedDownloadCandidatesRejectsRelativeTinyFiles(t *testing.T) {
	root := t.TempDir()
	main := writeSizedFile(t, root, "movie.mkv", 900*testMiB)
	clip := writeSizedFile(t, root, "clip.mkv", 60*testMiB)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: main, Complete: true},
		{Path: clip, Complete: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection, main)
	assertRejectedReason(t, selection, clip, "relative_tiny_file")
}

func TestSelectCompletedDownloadCandidatesRejectsMixedSampleAndExtras(t *testing.T) {
	root := t.TempDir()
	main := writeSizedFile(t, root, "Movie.2026.mkv", 900*testMiB)
	sample := writeSizedFile(t, root, "Movie.2026.sample.mkv", 100*testMiB)
	trailer := writeSizedFile(t, root, "Movie.2026.trailer.mp4", 100*testMiB)

	selection, err := selectCompletedDownloadCandidates([]downloadclients.StatusFile{
		{Path: sample, Complete: true},
		{Path: trailer, Complete: true},
		{Path: main, Complete: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	assertSelectedSources(t, selection, main)
	assertRejectedReason(t, selection, sample, "sample_or_extra")
	assertRejectedReason(t, selection, trailer, "sample_or_extra")
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
	if err := os.WriteFile(path, []byte{0}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Truncate(path, size); err != nil {
		t.Fatal(err)
	}
	return path
}
