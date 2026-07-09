package httpapi

import (
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
