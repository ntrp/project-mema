package jobs

import (
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func TestFulfillmentSchedulesUseUnifiedMediaPlanner(t *testing.T) {
	want := map[string]struct {
		kind  string
		queue string
	}{
		"media_fulfillment": {kind: MediaFulfillmentArgs{}.Kind(), queue: queueMediaAssembly},
		"audio_source":      {kind: AudioSourceArgs{}.Kind(), queue: queueMediaSearch},
		"subtitle_download": {kind: SubtitleDownloadArgs{}.Kind(), queue: queueMediaSearch},
	}
	removed := map[string]struct{}{
		"video_transcode":  {},
		"audio_transcode":  {},
		"container_remux":  {},
		"subtitle_embed":   {},
		"subtitle_extract": {},
		"subtitle_convert": {},
	}
	for _, definition := range fixedJobDefinitions() {
		if _, ok := removed[definition.ID]; ok {
			t.Fatalf("removed fulfillment schedule still registered: %s", definition.ID)
		}
		expected, ok := want[definition.ID]
		if !ok {
			continue
		}
		delete(want, definition.ID)
		if definition.Kind != expected.kind || definition.Queue != expected.queue {
			t.Fatalf("%s kind/queue = %s/%s", definition.ID, definition.Kind, definition.Queue)
		}
		if !definition.PausedByDefault || !definition.Automatic || !definition.ManualActionAvailable {
			t.Fatalf("%s flags = %#v", definition.ID, definition.SystemJobScheduleDefinition)
		}
	}
	if len(want) > 0 {
		t.Fatalf("missing fulfillment schedule definitions: %#v", want)
	}
}

func TestJobInsertOptsAreNonRetryable(t *testing.T) {
	metadata := []byte(`{"manual":true}`)
	opts := jobInsertOptsWithMetadataAndUnique(queueMediaSearch, metadata, river.UniqueOpts{
		ByState: []rivertype.JobState{rivertype.JobStateAvailable},
	})

	if opts.Queue != queueMediaSearch || opts.MaxAttempts != nonRetryableJobMaxAttempts {
		t.Fatalf("insert opts = %#v", opts)
	}
	if string(opts.Metadata) != string(metadata) || len(opts.UniqueOpts.ByState) != 1 {
		t.Fatalf("insert opts lost metadata/unique fields: %#v", opts)
	}
}

func TestFulfillmentTargetMatchesPartialAudioTranscode(t *testing.T) {
	target := targets.Target{
		Type:       targets.TypeAudio,
		State:      targets.StatePartial,
		LanguageID: "english",
	}
	args := FulfillmentActionArgs{TargetType: "audio", LanguageID: "eng"}

	if !fulfillmentTargetMatches(targets.OperationAudioTranscode, args, target) {
		t.Fatalf("partial audio target should match audio transcode")
	}
	args.LanguageID = "jpn"
	if fulfillmentTargetMatches(targets.OperationAudioTranscode, args, target) {
		t.Fatalf("language scoped request should not match another language")
	}
}

func TestFulfillmentTargetInRequestScopeFiltersExpectedMismatches(t *testing.T) {
	args := FulfillmentActionArgs{TargetType: "audio", LanguageID: "eng"}

	if fulfillmentTargetInRequestScope(args, targets.Target{Type: targets.TypeVideo, LanguageID: "eng"}) {
		t.Fatalf("video target should be outside audio request scope")
	}
	if !fulfillmentTargetInRequestScope(args, targets.Target{Type: targets.TypeAudio, LanguageID: "english"}) {
		t.Fatalf("english audio target should be inside eng request scope")
	}
}

func TestFulfillmentApplyTrackScopeNarrowsToOneTrack(t *testing.T) {
	trackID := uuid.New()
	language := "eng"
	item := storage.MediaItem{FileFacts: []storage.MediaFileFact{
		{
			FilePath: "/library/one.mkv",
			Tracks: []storage.MediaFileTrackFact{{
				ID:         uuid.New(),
				FilePath:   "/library/one.mkv",
				TrackType:  "audio",
				LanguageID: &language,
			}},
		},
		{
			FilePath: "/library/two.mkv",
			Tracks: []storage.MediaFileTrackFact{{
				ID:         trackID,
				FilePath:   "/library/two.mkv",
				TrackType:  "audio",
				LanguageID: &language,
			}},
		},
	}}

	scoped, args, track, err := fulfillmentApplyTrackScope(item, FulfillmentActionArgs{TrackID: trackID.String()})
	if err != nil {
		t.Fatalf("scope track: %v", err)
	}
	if track == nil || track.ID != trackID || len(scoped.FileFacts) != 1 || len(scoped.FileFacts[0].Tracks) != 1 {
		t.Fatalf("scoped item = %#v track=%#v", scoped, track)
	}
	if args.FilePath != "/library/two.mkv" || args.TargetType != "audio" || args.LanguageID != "eng" {
		t.Fatalf("scoped args = %#v", args)
	}
}

func TestFulfillmentActionDetailsIncludeManualScope(t *testing.T) {
	args := FulfillmentActionArgs{
		MediaItemID:        "media-1",
		FilePath:           "/library/movie.mkv",
		TargetType:         "audio",
		LanguageID:         "eng",
		TrackID:            "track-1",
		OtherFileID:        "other-1",
		ExternalSubtitleID: "subtitle-1",
	}
	got := fulfillmentActionDetails(targets.OperationAudioTranscode, args)

	for _, key := range []string{"operation", "mediaItemId", "filePath", "targetType", "languageId", "trackId", "otherFileId", "externalSubtitleId"} {
		if got[key] == "" {
			t.Fatalf("missing detail %s in %#v", key, got)
		}
	}
}

func TestAudioTrackTranscodeArgsTargetsSelectedAudioStream(t *testing.T) {
	args, err := audioTrackTranscodeArgs("/library/movie.mkv", "/library/.movie.tmp.mkv", 1, AudioConversionDecision{
		Allowed:           true,
		TargetCodec:       "eac3",
		TargetChannels:    "5.1",
		TargetBitrateKbps: 640,
	})
	if err != nil {
		t.Fatalf("audio track transcode args: %v", err)
	}
	if !hasArgPair(args, "-c:a:1", "eac3") {
		t.Fatalf("expected selected audio codec args, got %#v", args)
	}
	if !hasArgPair(args, "-ac:a:1", "6") {
		t.Fatalf("expected selected audio channel args, got %#v", args)
	}
	if !hasArgPair(args, "-b:a:1", "640k") {
		t.Fatalf("expected selected audio bitrate args, got %#v", args)
	}
	if !slices.Contains(args, "-map") || !slices.Contains(args, "-c") {
		t.Fatalf("expected full-file remux args, got %#v", args)
	}
}

func TestVideoTrackTranscodeArgsTargetsSelectedVideoStream(t *testing.T) {
	args, err := videoTrackTranscodeArgs("/library/movie.mkv", "/library/.movie.tmp.mkv", 1, VideoConversionDecision{
		Allowed:     true,
		TargetCodec: "hevc",
		TargetPixel: "yuv420p10le",
	})
	if err != nil {
		t.Fatalf("video track transcode args: %v", err)
	}
	if !hasArgPair(args, "-c:v:1", "libx265") {
		t.Fatalf("expected selected video codec args, got %#v", args)
	}
	if !hasArgPair(args, "-pix_fmt:v:1", "yuv420p10le") {
		t.Fatalf("expected selected video pixel args, got %#v", args)
	}
	if !slices.Contains(args, "-map") || !slices.Contains(args, "-c") {
		t.Fatalf("expected full-file remux args, got %#v", args)
	}
}

func TestContainerRemuxArgsCopyStreamsToTargetContainer(t *testing.T) {
	args, err := containerRemuxArgs("/library/movie.mkv", "/library/.movie.tmp.mp4", "mp4")
	if err != nil {
		t.Fatalf("container remux args: %v", err)
	}
	if !slices.Contains(args, "-map") || !slices.Contains(args, "-c") {
		t.Fatalf("expected full-file remux args, got %#v", args)
	}
	if !hasArgPair(args, "-movflags", "+faststart") {
		t.Fatalf("expected mp4 faststart args, got %#v", args)
	}
	if !hasArgPair(args, "-c:s", "mov_text") {
		t.Fatalf("expected mp4 subtitle conversion args, got %#v", args)
	}
	if args[len(args)-1] != "/library/.movie.tmp.mp4" {
		t.Fatalf("output arg = %q", args[len(args)-1])
	}
}

func TestAudioOrdinalFindsSelectedTrackAudioIndex(t *testing.T) {
	selected := uuid.New()
	item := storage.MediaItem{FileFacts: []storage.MediaFileFact{{
		FilePath: "/library/movie.mkv",
		Tracks: []storage.MediaFileTrackFact{
			{ID: uuid.New(), FilePath: "/library/movie.mkv", TrackType: "video"},
			{ID: uuid.New(), FilePath: "/library/movie.mkv", TrackType: "audio"},
			{ID: selected, FilePath: "/library/movie.mkv", TrackType: "audio"},
		},
	}}}
	got := audioOrdinal(item, storage.MediaFileTrackFact{ID: selected, FilePath: "/library/movie.mkv"})

	if got != 1 {
		t.Fatalf("audio ordinal = %d, want 1", got)
	}
}

func TestVideoOrdinalFindsSelectedTrackVideoIndex(t *testing.T) {
	selected := uuid.New()
	item := storage.MediaItem{FileFacts: []storage.MediaFileFact{{
		FilePath: "/library/movie.mkv",
		Tracks: []storage.MediaFileTrackFact{
			{ID: uuid.New(), FilePath: "/library/movie.mkv", TrackType: "video"},
			{ID: uuid.New(), FilePath: "/library/movie.mkv", TrackType: "audio"},
			{ID: selected, FilePath: "/library/movie.mkv", TrackType: "video"},
		},
	}}}
	got := videoOrdinal(item, storage.MediaFileTrackFact{ID: selected, FilePath: "/library/movie.mkv"})

	if got != 1 {
		t.Fatalf("video ordinal = %d, want 1", got)
	}
}

func hasArgPair(args []string, key string, value string) bool {
	for index := 0; index+1 < len(args); index++ {
		if args[index] == key && args[index+1] == value {
			return true
		}
	}
	return false
}
