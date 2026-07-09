package httpapi

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestFulfillmentActionScopeHelpers(t *testing.T) {
	trackID := uuid.New()
	subtitleID := uuid.New()
	item := storage.MediaItem{
		ExternalSubtitles: []storage.MediaItemSubtitle{{ID: subtitleID}},
		FileFacts: []storage.MediaFileFact{
			{
				FilePath: "/library/movie/movie.mkv",
				Tracks:   []storage.MediaFileTrackFact{{ID: trackID}},
			},
		},
	}

	if mediaItemTrack(item, "", trackID.String()) == nil {
		t.Fatalf("track should match item scope")
	}
	if mediaItemTrack(item, "/library/movie/movie.mkv", trackID.String()) == nil {
		t.Fatalf("track should match file scope")
	}
	if mediaItemTrack(item, "/library/movie/other.mkv", trackID.String()) != nil {
		t.Fatalf("track should not match another file")
	}
	if mediaItemTrack(item, "", uuid.New().String()) != nil {
		t.Fatalf("unknown track should not match")
	}
	if !mediaItemHasSubtitle(item, subtitleID.String()) {
		t.Fatalf("subtitle should match item scope")
	}
	if mediaItemHasSubtitle(item, uuid.New().String()) {
		t.Fatalf("unknown subtitle should not match")
	}
}

func TestValidTrackScopedFulfillmentActionRequiresMatchingTrack(t *testing.T) {
	audioTrack := &storage.MediaFileTrackFact{TrackType: "audio"}
	videoTrack := &storage.MediaFileTrackFact{TrackType: "video"}

	if !validTrackScopedFulfillmentAction(httptest.NewRecorder(), "audio_transcode", audioTrack) {
		t.Fatalf("audio track should be valid for audio transcode")
	}
	if !validTrackScopedFulfillmentAction(httptest.NewRecorder(), "video_transcode", videoTrack) {
		t.Fatalf("video track should be valid for video transcode")
	}
	if validTrackScopedFulfillmentAction(httptest.NewRecorder(), "video_transcode", audioTrack) {
		t.Fatalf("audio track should not be valid for video transcode")
	}
	if validTrackScopedFulfillmentAction(httptest.NewRecorder(), "audio_transcode", nil) {
		t.Fatalf("missing track should not be valid for audio transcode")
	}
}
