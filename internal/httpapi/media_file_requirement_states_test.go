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
		{Type: MediaFileTrackTypeSubtitle, Language: stringPtr("english"), Codec: stringPtr("ass")},
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

func TestMediaFileRequirementStatesKeepAudioSummaryPartialForExistingBadTrack(t *testing.T) {
	minBitrate := int32(384)
	item := storage.MediaItem{
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:         "english",
			TargetCodec:        stringPtr("aac"),
			TargetChannels:     []string{"5.1"},
			MinimumBitrateKbps: &minBitrate,
		}},
	}
	file := MediaFileInfo{
		Path:   "/media/movie.mkv",
		Status: MediaFileInfoStatusAvailable,
		SubtitleSatisfaction: &MediaFileSubtitleSatisfaction{
			Mode:  MediaProfileSubtitleModeMixed,
			State: MediaFileSubtitleSatisfactionStateIgnored,
		},
	}
	tracks := []MediaFileTrack{
		{
			Type:          MediaFileTrackTypeAudio,
			Language:      stringPtr("eng"),
			Codec:         stringPtr("ac3"),
			Channels:      int32Ptr(2),
			ChannelLayout: stringPtr("stereo"),
			BitRate:       stringPtr("192000"),
		},
	}

	applyMediaFileRequirementStates(&file, item, tracks, nil)

	if file.Requirements.Audio.State != MediaFileRequirementStatePartial {
		t.Fatalf("audio requirement = %#v", file.Requirements.Audio)
	}
	if file.MissingTracks != nil {
		t.Fatalf("partial audio track should not create missing rows: %#v", file.MissingTracks)
	}
	if tracks[0].State == nil || tracks[0].State.VisualState != MediaFileDetailVisualStatePartial {
		t.Fatalf("audio track state = %#v", tracks[0].State)
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

func TestMediaFileRequirementStatesCheckVideoResolutionAgainstQuality(t *testing.T) {
	qualityID := "webdl-1080p"
	item := storage.MediaItem{
		FileFacts: []storage.MediaFileFact{{
			FilePath:  "/media/movie.mkv",
			QualityID: &qualityID,
		}},
	}
	file := MediaFileInfo{
		Path:   "/media/movie.mkv",
		Status: MediaFileInfoStatusAvailable,
		SubtitleSatisfaction: &MediaFileSubtitleSatisfaction{
			Mode:  MediaProfileSubtitleModeMixed,
			State: MediaFileSubtitleSatisfactionStateIgnored,
		},
	}
	tracks := []MediaFileTrack{
		{Type: MediaFileTrackTypeVideo, Width: int32Ptr(1280), Height: int32Ptr(536)},
	}

	applyMediaFileRequirementStates(&file, item, tracks, nil)

	if file.Requirements.Video.State != MediaFileRequirementStatePartial {
		t.Fatalf("video requirement = %#v", file.Requirements.Video)
	}
	if tracks[0].State == nil || tracks[0].State.VisualState != MediaFileDetailVisualStatePartial {
		t.Fatalf("video track state = %#v", tracks[0].State)
	}
	if tracks[0].State.Details[0] != "Video resolution is below selected quality 1080p" {
		t.Fatalf("video details = %#v", tracks[0].State.Details)
	}
}

func TestMediaFileRequirementStatesAcceptCroppedVideoAtQualityWidth(t *testing.T) {
	qualityID := "webdl-1080p"
	item := storage.MediaItem{
		FileFacts: []storage.MediaFileFact{{
			FilePath:  "/media/movie.mkv",
			QualityID: &qualityID,
		}},
	}
	file := MediaFileInfo{
		Path:   "/media/movie.mkv",
		Status: MediaFileInfoStatusAvailable,
		SubtitleSatisfaction: &MediaFileSubtitleSatisfaction{
			Mode:  MediaProfileSubtitleModeMixed,
			State: MediaFileSubtitleSatisfactionStateIgnored,
		},
	}
	tracks := []MediaFileTrack{
		{Type: MediaFileTrackTypeVideo, Width: int32Ptr(1920), Height: int32Ptr(800)},
	}

	applyMediaFileRequirementStates(&file, item, tracks, nil)

	if file.Requirements.Video.State != MediaFileRequirementStateSatisfied {
		t.Fatalf("video requirement = %#v", file.Requirements.Video)
	}
}

func TestMediaFileRequirementStatesMatchSubtitleFormatAliases(t *testing.T) {
	item := storage.MediaItem{
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{LanguageID: "english", Formats: []string{"srt"}}},
		SubtitleMode:    "mixed",
	}
	file := MediaFileInfo{
		Path:   "/media/movie.mkv",
		Status: MediaFileInfoStatusAvailable,
		SubtitleSatisfaction: &MediaFileSubtitleSatisfaction{
			Mode:             MediaProfileSubtitleModeMixed,
			State:            MediaFileSubtitleSatisfactionStateSatisfied,
			WantedLanguages:  []string{"english"},
			MatchedLanguages: []string{"english"},
			MissingLanguages: []string{},
		},
	}
	tracks := []MediaFileTrack{
		{Type: MediaFileTrackTypeSubtitle, Language: stringPtr("english"), Codec: stringPtr("subrip")},
	}

	applyMediaFileRequirementStates(&file, item, tracks, nil)

	if tracks[0].State == nil || tracks[0].State.VisualState != MediaFileDetailVisualStateMatching {
		t.Fatalf("subtitle track state = %#v", tracks[0].State)
	}
}

func TestMediaFileRequirementStatesDoNotOfferBitmapSubtitleConversion(t *testing.T) {
	item := storage.MediaItem{
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{LanguageID: "english", Formats: []string{"srt"}}},
		SubtitleMode:    "mixed",
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
		{Type: MediaFileTrackTypeSubtitle, Language: stringPtr("english"), Codec: stringPtr("pgs")},
	}

	applyMediaFileRequirementStates(&file, item, tracks, nil)

	if tracks[0].State == nil || tracks[0].State.VisualState != MediaFileDetailVisualStatePartial {
		t.Fatalf("subtitle track state = %#v", tracks[0].State)
	}
	if tracks[0].State.OperationLabel != nil {
		t.Fatalf("unexpected operation = %#v", tracks[0].State.OperationLabel)
	}
}

func TestMediaFileRequirementStatesMarkOffTargetSubtitleUnwanted(t *testing.T) {
	item := storage.MediaItem{
		SubtitleTargets:         []storage.MediaProfileSubtitleTarget{{LanguageID: "italian", Formats: []string{"srt"}}},
		SubtitleMode:            "mixed",
		RemoveUnwantedSubtitles: false,
	}
	file := MediaFileInfo{
		Path:   "/media/movie.mkv",
		Status: MediaFileInfoStatusAvailable,
		SubtitleSatisfaction: &MediaFileSubtitleSatisfaction{
			Mode:             MediaProfileSubtitleModeMixed,
			State:            MediaFileSubtitleSatisfactionStateSatisfied,
			WantedLanguages:  []string{"italian"},
			MatchedLanguages: []string{"italian"},
			MissingLanguages: []string{},
		},
	}
	tracks := []MediaFileTrack{
		{Type: MediaFileTrackTypeSubtitle, Language: stringPtr("eng"), Codec: stringPtr("subrip")},
	}
	otherFiles := []MediaFileOtherFile{
		{Type: MediaFileOtherFileTypeSubtitle, Path: "/media/movie.ita.srt", Status: MediaFileOtherFileStatusAvailable, Language: stringPtr("italian")},
	}

	applyMediaFileRequirementStates(&file, item, tracks, otherFiles)

	if tracks[0].State == nil || tracks[0].State.VisualState != MediaFileDetailVisualStateUnwanted {
		t.Fatalf("off-target embedded subtitle state = %#v", tracks[0].State)
	}
	if otherFiles[0].State == nil || otherFiles[0].State.VisualState != MediaFileDetailVisualStateMatching {
		t.Fatalf("target external subtitle state = %#v", otherFiles[0].State)
	}
	if file.Requirements.Subtitles.State != MediaFileRequirementStateSatisfied {
		t.Fatalf("subtitle requirement = %#v", file.Requirements.Subtitles)
	}
}
