package jobs

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestVideoTranscodePlansCreateNonManualTrackJobs(t *testing.T) {
	mediaID := uuid.New()
	trackID := uuid.New()
	sourceCodec := "h264"
	targetCodec := "hevc"
	sourcePixel := "yuv420p"
	targetPixel := "yuv420p10le"
	item := storage.MediaItem{
		ID: mediaID,
		VideoTarget: storage.MediaProfileVideoTarget{
			Codecs:       []string{targetCodec},
			PixelFormats: []string{targetPixel},
		},
		FileFacts: []storage.MediaFileFact{{
			ID:       uuid.New(),
			FilePath: "/library/movie.mkv",
			Tracks: []storage.MediaFileTrackFact{{
				ID:          trackID,
				FilePath:    "/library/movie.mkv",
				TrackType:   "video",
				Codec:       &sourceCodec,
				PixelFormat: &sourcePixel,
			}},
		}},
	}

	plans := videoTranscodePlans(item, FulfillmentActionArgs{})
	if len(plans) != 1 {
		t.Fatalf("plans = %#v", plans)
	}
	got := plans[0].args
	if got.MediaItemID != mediaID.String() || got.FilePath != "/library/movie.mkv" ||
		got.TargetType != "video" || got.TrackID != trackID.String() || got.Manual {
		t.Fatalf("planned args = %#v", got)
	}
}

func TestVideoTranscodePlansSkipUnsupportedHdrOnlyTracks(t *testing.T) {
	sourceHDR := "sdr"
	sourceCodec := "hevc"
	item := storage.MediaItem{
		ID: uuid.New(),
		VideoTarget: storage.MediaProfileVideoTarget{
			Codecs:     []string{"hevc"},
			HDRFormats: []string{"hdr10"},
		},
		FileFacts: []storage.MediaFileFact{{
			FilePath: "/library/movie.mkv",
			Tracks: []storage.MediaFileTrackFact{{
				ID:        uuid.New(),
				FilePath:  "/library/movie.mkv",
				TrackType: "video",
				Codec:     &sourceCodec,
				HDRFormat: &sourceHDR,
			}},
		}},
	}

	if plans := videoTranscodePlans(item, FulfillmentActionArgs{}); len(plans) != 0 {
		t.Fatalf("plans = %#v", plans)
	}
}
