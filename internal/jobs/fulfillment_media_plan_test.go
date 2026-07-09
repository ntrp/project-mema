package jobs

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestMediaFulfillmentScheduleIsSingleDisabledPlanner(t *testing.T) {
	definition, ok := fixedJobDefinitionByID("media_fulfillment")
	if !ok {
		t.Fatal("media fulfillment schedule missing")
	}
	if definition.Kind != (MediaFulfillmentArgs{}).Kind() || !definition.PausedByDefault {
		t.Fatalf("media fulfillment definition = %#v", definition.SystemJobScheduleDefinition)
	}
	for _, removed := range []string{"video_transcode", "audio_transcode", "container_remux", "subtitle_extract", "subtitle_convert", "subtitle_embed"} {
		if _, ok := fixedJobDefinitionByID(removed); ok {
			t.Fatalf("removed schedule still present: %s", removed)
		}
	}
}

func TestMediaFulfillmentProgressPercent(t *testing.T) {
	cases := []struct {
		processed int
		total     int
		want      int32
	}{
		{0, 0, 100},
		{0, 3, 0},
		{1, 3, 33},
		{2, 3, 66},
		{3, 3, 100},
		{4, 3, 100},
	}
	for _, tc := range cases {
		got := mediaFulfillmentProgressPercent(tc.processed, tc.total)
		if got != tc.want {
			t.Fatalf("mediaFulfillmentProgressPercent(%d, %d) = %d, want %d", tc.processed, tc.total, got, tc.want)
		}
	}
}

func TestEnqueueSubtitleFulfillmentJobsQueuesExtraction(t *testing.T) {
	language := "english"
	format := "srt"
	trackID := uuid.New()
	item := storage.MediaItem{
		ID:           uuid.New(),
		SubtitleMode: "external",
		SubtitleTargets: []storage.MediaProfileSubtitleTarget{{
			LanguageID: language,
			Formats:    []string{"srt"},
		}},
		FileFacts: []storage.MediaFileFact{{
			ID:       uuid.New(),
			FilePath: "/library/movie.mkv",
			Tracks: []storage.MediaFileTrackFact{{
				ID:         trackID,
				FilePath:   "/library/movie.mkv",
				TrackType:  "subtitle",
				LanguageID: &language,
				Format:     &format,
			}},
		}},
	}
	operations := []string{}
	enqueue := func(_ context.Context, operation string, args FulfillmentActionArgs) (int64, error) {
		operations = append(operations, operation)
		if args.FilePath != "/library/movie.mkv" || args.LanguageID != language {
			t.Fatalf("subtitle args = %#v", args)
		}
		return int64(len(operations)), nil
	}

	count, err := enqueueSubtitleFulfillmentJobs(context.Background(), nil, nil, enqueue, item)
	if err != nil {
		t.Fatalf("enqueue subtitle fulfillment: %v", err)
	}
	if count != 1 || operations[0] != "subtitle_extraction" {
		t.Fatalf("queued operations = %#v count=%d", operations, count)
	}
}
