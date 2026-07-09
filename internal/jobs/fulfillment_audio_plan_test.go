package jobs

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestAudioTranscodePlansCreateNonManualTrackJobs(t *testing.T) {
	mediaID := uuid.New()
	fileID := uuid.New()
	trackID := uuid.New()
	language := "eng"
	sourceCodec := "aac"
	targetCodec := "eac3"
	item := storage.MediaItem{
		ID: mediaID,
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:  "english",
			TargetCodec: &targetCodec,
		}},
		FileFacts: []storage.MediaFileFact{{
			ID:       fileID,
			FilePath: "/library/movie.mkv",
			Tracks: []storage.MediaFileTrackFact{{
				ID:         trackID,
				FilePath:   "/library/movie.mkv",
				TrackType:  "audio",
				LanguageID: &language,
				Codec:      &sourceCodec,
			}},
		}},
	}

	plans := audioTranscodePlansForPolicy("lossyToLossy", item, FulfillmentActionArgs{})
	if len(plans) != 1 {
		t.Fatalf("plans = %#v", plans)
	}
	got := plans[0].args
	if got.MediaItemID != mediaID.String() || got.FilePath != "/library/movie.mkv" || got.TargetType != "audio" ||
		got.LanguageID != "english" || got.TrackID != trackID.String() || got.Manual {
		t.Fatalf("planned args = %#v", got)
	}
}

func TestAudioTranscodePlansSkipPolicyBlockedTracks(t *testing.T) {
	language := "eng"
	sourceCodec := "aac"
	targetCodec := "eac3"
	item := storage.MediaItem{
		ID: uuid.New(),
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:  "eng",
			TargetCodec: &targetCodec,
		}},
		FileFacts: []storage.MediaFileFact{{
			ID:       uuid.New(),
			FilePath: "/library/movie.mkv",
			Tracks: []storage.MediaFileTrackFact{{
				ID:         uuid.New(),
				FilePath:   "/library/movie.mkv",
				TrackType:  "audio",
				LanguageID: &language,
				Codec:      &sourceCodec,
			}},
		}},
	}

	if plans := audioTranscodePlansForPolicy("disabled", item, FulfillmentActionArgs{}); len(plans) != 0 {
		t.Fatalf("plans = %#v", plans)
	}
}

func TestAudioTranscodePlansSkipInvalidNoOpChannelTarget(t *testing.T) {
	language := "eng"
	sourceCodec := "aac"
	channels := "2.0"
	item := storage.MediaItem{
		ID: uuid.New(),
		AudioTargets: []storage.MediaProfileAudioTarget{{
			LanguageID:     "eng",
			TargetChannels: []string{"objectaudio"},
		}},
		FileFacts: []storage.MediaFileFact{{
			ID:       uuid.New(),
			FilePath: "/library/movie.mkv",
			Tracks: []storage.MediaFileTrackFact{{
				ID:         uuid.New(),
				FilePath:   "/library/movie.mkv",
				TrackType:  "audio",
				LanguageID: &language,
				Codec:      &sourceCodec,
				Channels:   &channels,
			}},
		}},
	}

	if plans := audioTranscodePlansForPolicy("lossyToLossy", item, FulfillmentActionArgs{}); len(plans) != 0 {
		t.Fatalf("plans = %#v", plans)
	}
}
