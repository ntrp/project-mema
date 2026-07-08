package httpapi

import (
	"testing"

	"media-manager/internal/storage"
)

func TestMediaFileRequirementStatesComputeTrackSidecarAndMissingRows(t *testing.T) {
	minBitrate := int32(384)
	item := storage.MediaItem{
		VideoTarget:             storage.MediaProfileVideoTarget{Codecs: []string{"h264"}},
		AudioTargets:            []storage.MediaProfileAudioTarget{{LanguageID: "german", MinimumBitrateKbps: &minBitrate}},
		SubtitleTargets:         []storage.MediaProfileSubtitleTarget{{LanguageID: "english", Formats: []string{"srt"}}},
		SubtitleMode:            "mixed",
		RemoveUnwantedSubtitles: true,
	}
	file := MediaFileInfo{
		Path:   "/media/movie.mkv",
		Status: MediaFileInfoStatusAvailable,
		SubtitleSatisfaction: &MediaFileSubtitleSatisfaction{
			Mode:             MediaProfileSubtitleModeMixed,
			State:            MediaFileSubtitleSatisfactionStateMissing,
			WantedLanguages:  []string{"english"},
			MatchedLanguages: []string{},
			MissingLanguages: []string{"english"},
		},
	}
	tracks := []MediaFileTrack{
		{Type: MediaFileTrackTypeVideo, Codec: stringPtr("h264")},
		{Type: MediaFileTrackTypeAudio, Language: stringPtr("english"), Codec: stringPtr("aac")},
		{Type: MediaFileTrackTypeSubtitle, Language: stringPtr("english"), Codec: stringPtr("subrip")},
	}
	otherFiles := []MediaFileOtherFile{
		{Type: MediaFileOtherFileTypeSubtitle, Path: "/media/movie.spanish.srt", Status: MediaFileOtherFileStatusAvailable, Language: stringPtr("spanish")},
	}

	applyMediaFileRequirementStates(&file, item, tracks, otherFiles)

	if file.Requirements.Video.State != MediaFileRequirementStateSatisfied {
		t.Fatalf("video requirement = %#v", file.Requirements.Video)
	}
	if file.Requirements.Audio.State != MediaFileRequirementStateMissing {
		t.Fatalf("audio requirement = %#v", file.Requirements.Audio)
	}
	if file.Requirements.Subtitles.State != MediaFileRequirementStatePending {
		t.Fatalf("subtitle requirement = %#v", file.Requirements.Subtitles)
	}
	if tracks[2].State == nil || tracks[2].State.OperationLabel == nil || *tracks[2].State.OperationLabel != "Convert subtitle" {
		t.Fatalf("subtitle track state = %#v", tracks[2].State)
	}
	if otherFiles[0].State == nil || otherFiles[0].State.VisualState != MediaFileDetailVisualStateUnwanted {
		t.Fatalf("sidecar state = %#v", otherFiles[0].State)
	}
	if file.MissingTracks == nil || len(*file.MissingTracks) != 1 || (*file.MissingTracks)[0].Type != MediaFileMissingTrackTypeAudio {
		t.Fatalf("missing tracks = %#v", file.MissingTracks)
	}
}

func TestMediaFileRequirementStatesMarkMissingExternalSubtitle(t *testing.T) {
	item := storage.MediaItem{
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{LanguageID: "japanese", Formats: []string{"srt"}}},
		SubtitleMode:    "external",
	}
	file := MediaFileInfo{
		Path:   "/media/movie.mkv",
		Status: MediaFileInfoStatusAvailable,
		SubtitleSatisfaction: &MediaFileSubtitleSatisfaction{
			Mode:             MediaProfileSubtitleModeExternal,
			State:            MediaFileSubtitleSatisfactionStateMissing,
			WantedLanguages:  []string{"japanese"},
			MatchedLanguages: []string{},
			MissingLanguages: []string{"japanese"},
		},
	}
	otherFiles := []MediaFileOtherFile{
		{Type: MediaFileOtherFileTypeSubtitle, Path: "/media/movie.japanese.srt", Status: MediaFileOtherFileStatusMissing, Language: stringPtr("japanese")},
	}

	applyMediaFileRequirementStates(&file, item, nil, otherFiles)

	if otherFiles[0].State == nil || otherFiles[0].State.VisualState != MediaFileDetailVisualStateMissingPlaceholder {
		t.Fatalf("missing external state = %#v", otherFiles[0].State)
	}
	if file.MissingTracks != nil {
		t.Fatalf("missing subtitle sidecar should be the candidate row, got %#v", file.MissingTracks)
	}
}
