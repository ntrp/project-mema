package httpapi

import (
	"os"
	"path/filepath"
	"testing"

	"media-manager/internal/storage"
)

func TestSCNMedia007StorageMediaSeasonsMapsNestedMonitorState(t *testing.T) {
	seasonMonitored := true
	episodeMonitored := false
	episodeCount := int32(2)
	airDate := "2026-01-02"
	values := []MediaMetadataSeason{{
		Name:         "Season 1",
		EpisodeCount: &episodeCount,
		AirDate:      &airDate,
		Monitored:    &seasonMonitored,
		Episodes: &[]MediaMetadataEpisode{{
			Name:          "Pilot",
			EpisodeNumber: 1,
			AirDate:       &airDate,
			Monitored:     &episodeMonitored,
		}},
	}}

	seasons := storageMediaSeasons(&values)

	if seasons == nil || len(*seasons) != 1 {
		t.Fatalf("seasons = %#v", seasons)
	}
	season := (*seasons)[0]
	if season.Name != "Season 1" || !season.Monitored || season.EpisodeCount == nil {
		t.Fatalf("season = %#v", season)
	}
	if len(season.Episodes) != 1 || season.Episodes[0].Monitored {
		t.Fatalf("episodes = %#v", season.Episodes)
	}
	if storageMediaSeasons(nil) != nil {
		t.Fatal("nil season pointer should stay nil")
	}
}

func TestSCNMedia007OptionalEnumAndBooleanValues(t *testing.T) {
	enabled := true
	disabled := false
	mode := FutureEpisodes
	availability := Released

	if !optionalBoolValue(&enabled) || optionalBoolValue(&disabled) || optionalBoolValue(nil) {
		t.Fatal("optionalBoolValue did not map pointers correctly")
	}
	if got := optionalMediaMonitorMode(&mode); got == nil || *got != "future_episodes" {
		t.Fatalf("optionalMediaMonitorMode = %v", got)
	}
	if got := optionalMinimumAvailability(&availability); got == nil || *got != "released" {
		t.Fatalf("optionalMinimumAvailability = %v", got)
	}
	if optionalMediaMonitorMode(nil) != nil || optionalMinimumAvailability(nil) != nil {
		t.Fatal("nil enum pointers should stay nil")
	}
}

func TestSCNMedia001MediaFileInfoResponsesExposeExistingFileSizes(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "Scenario.Movie.2026.mkv")
	missingPath := filepath.Join(dir, "missing.mkv")
	if err := os.WriteFile(filePath, []byte("scenario"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "poster.jpg"), []byte("poster"), 0o600); err != nil {
		t.Fatalf("write poster fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.bin"), []byte("notes"), 0o600); err != nil {
		t.Fatalf("write unknown fixture: %v", err)
	}
	missingPoster := filepath.Join(dir, "Scenario.Movie.2026-fanart.jpg")

	files := mediaFileInfoResponses(storage.MediaItem{
		FilePaths:       []string{filePath, missingPath},
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{LanguageID: "german", Formats: []string{"srt"}}},
		SubtitleMode:    "external",
		Sidecars: []storage.MediaItemSidecar{{
			MediaFilePath: filePath,
			FilePath:      missingPoster,
			SidecarType:   storage.MediaSidecarMetadata,
		}},
	})

	if files == nil || len(*files) != 2 {
		t.Fatalf("files = %#v", files)
	}
	if (*files)[0].Path != filePath || (*files)[0].Status != MediaFileInfoStatusAvailable || (*files)[0].SizeBytes == nil || *(*files)[0].SizeBytes != 8 {
		t.Fatalf("existing file response = %#v", (*files)[0])
	}
	if (*files)[1].Path != missingPath || (*files)[1].Status != MediaFileInfoStatusMissing || (*files)[1].SizeBytes != nil {
		t.Fatalf("missing file response = %#v", (*files)[1])
	}
	otherFiles := *(*files)[0].OtherFiles
	if !hasOtherFile(otherFiles, MediaFileOtherFileTypeSubtitle, filepath.Join(dir, "Scenario.Movie.2026.german.srt"), MediaFileOtherFileStatusMissing) {
		t.Fatalf("missing subtitle file not derived: %#v", otherFiles)
	}
	if !hasOtherFile(otherFiles, MediaFileOtherFileTypeMetadata, filepath.Join(dir, "poster.jpg"), MediaFileOtherFileStatusAvailable) {
		t.Fatalf("metadata file not listed: %#v", otherFiles)
	}
	if !hasOtherFile(otherFiles, MediaFileOtherFileTypeUnknown, filepath.Join(dir, "notes.bin"), MediaFileOtherFileStatusAvailable) {
		t.Fatalf("unknown file not listed: %#v", otherFiles)
	}
	if !hasOtherFile(otherFiles, MediaFileOtherFileTypeMetadata, missingPoster, MediaFileOtherFileStatusMissing) {
		t.Fatalf("missing metadata sidecar not listed: %#v", otherFiles)
	}
}

func TestSCNMedia001MediaFileProbePathConvertsRelativePaths(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	got := mediaFileProbePath(filepath.Join("relative", "Movie.mkv"))
	want := filepath.Join(dir, "relative", "Movie.mkv")
	if got != want {
		t.Fatalf("mediaFileProbePath = %q, want %q", got, want)
	}

	absolute := filepath.Join(dir, "absolute.mkv")
	if mediaFileProbePath(absolute) != absolute {
		t.Fatalf("absolute probe path changed")
	}
}

func hasOtherFile(
	files []MediaFileOtherFile,
	fileType MediaFileOtherFileType,
	path string,
	status MediaFileOtherFileStatus,
) bool {
	for _, file := range files {
		if file.Type == fileType && file.Path == path && file.Status == status {
			return true
		}
	}
	return false
}
