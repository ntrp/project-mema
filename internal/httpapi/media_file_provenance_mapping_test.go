package httpapi

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestHydrateTrackProvenanceMatchesSourceStreamID(t *testing.T) {
	itemID := uuid.New()
	streamID := int32(2)
	tracks := []MediaFileTrack{
		{Type: Audio, Index: provenanceInt32Ptr(1), Codec: provenanceStringPtr("aac")},
		{Type: Audio, Index: &streamID, Codec: provenanceStringPtr("dts")},
	}

	hydrateTrackProvenance("/library/Scenario.Release.mkv", tracks, []storage.MediaComponentProvenance{{
		ID:             uuid.New(),
		MediaItemID:    itemID,
		ComponentType:  "audio",
		ComponentKey:   "source-1",
		ReleaseGroup:   "ARR",
		ReleaseName:    "Scenario.Release",
		SourceStreamID: &streamID,
	}})

	if tracks[0].Provenance != nil {
		t.Fatalf("unexpected first track provenance: %#v", tracks[0].Provenance)
	}
	if tracks[1].Provenance == nil || tracks[1].Provenance.SourceStreamId == nil || *tracks[1].Provenance.SourceStreamId != streamID {
		t.Fatalf("second track provenance = %#v", tracks[1].Provenance)
	}
}

func TestHydrateTrackProvenanceSkipsAmbiguousTypeFallback(t *testing.T) {
	tracks := []MediaFileTrack{
		{Type: Subtitle, Index: provenanceInt32Ptr(3)},
		{Type: Subtitle, Index: provenanceInt32Ptr(4)},
	}

	hydrateTrackProvenance("/library/Scenario.Release.mkv", tracks, []storage.MediaComponentProvenance{{
		ID:            uuid.New(),
		MediaItemID:   uuid.New(),
		ComponentType: "subtitle",
		ComponentKey:  "source-1",
	}})

	if tracks[0].Provenance != nil || tracks[1].Provenance != nil {
		t.Fatalf("ambiguous provenance should not hydrate tracks: %#v", tracks)
	}
}

func TestHydrateTrackProvenanceFallsBackToImportedContainer(t *testing.T) {
	itemID := uuid.New()
	path := "/library/Scenario.Movie.2026.1080p-ARR.mkv"
	tracks := []MediaFileTrack{
		{Type: Video, Index: provenanceInt32Ptr(0), Codec: provenanceStringPtr("h264")},
		{Type: Audio, Index: provenanceInt32Ptr(1), Codec: provenanceStringPtr("dts")},
	}

	hydrateTrackProvenance(path, tracks, []storage.MediaComponentProvenance{{
		ID:             uuid.New(),
		MediaItemID:    itemID,
		ComponentType:  "container",
		ComponentKey:   "imported:" + path,
		ReleaseGroup:   "ARR",
		ReleaseName:    "Scenario.Movie.2026.1080p",
		SourceFilePath: &path,
	}})

	for _, track := range tracks {
		if track.Provenance == nil || track.Provenance.ReleaseGroup != "ARR" {
			t.Fatalf("track provenance = %#v", tracks)
		}
	}
}

func provenanceInt32Ptr(value int32) *int32 {
	return &value
}

func provenanceStringPtr(value string) *string {
	return &value
}
