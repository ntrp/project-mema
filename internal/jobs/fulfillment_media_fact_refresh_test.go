package jobs

import (
	"strings"
	"testing"

	"media-manager/internal/delivery"
	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func TestLiveMediaFileFactInputMapsSubtitleTrackState(t *testing.T) {
	language := "eng"
	codec := "subrip"
	duration := 12.5
	filePath := "/library/movie.mkv"
	item := storage.MediaItem{
		ID: uuid.New(),
		FileFacts: []storage.MediaFileFact{{
			FilePath: filePath,
		}},
	}

	input := liveMediaFileFactInput(item, filePath, delivery.ProbeResult{
		DurationSeconds: &duration,
		Tracks: []delivery.Track{{
			Index:    int32Ptr(2),
			Type:     delivery.TrackSubtitle,
			Codec:    &codec,
			Language: &language,
			Duration: &duration,
		}},
	})

	if len(input.Tracks) != 1 {
		t.Fatalf("tracks = %#v", input.Tracks)
	}
	track := input.Tracks[0]
	if track.LanguageID == nil || *track.LanguageID != "eng" {
		t.Fatalf("language = %#v", track.LanguageID)
	}
	if track.Format == nil || *track.Format != "subrip" {
		t.Fatalf("format = %#v", track.Format)
	}
	if track.DurationMs == nil || *track.DurationMs != 12500 {
		t.Fatalf("duration = %#v", track.DurationMs)
	}
}

func TestPersistLiveMediaFileFactFailsWhenProbeFindsNoTracks(t *testing.T) {
	err := persistLiveMediaFileFact(t.Context(), nil, storage.MediaItem{ID: uuid.New()}, "/missing/movie.mkv")
	if err == nil || !strings.Contains(err.Error(), "probe returned no tracks") {
		t.Fatalf("error = %v", err)
	}
}
