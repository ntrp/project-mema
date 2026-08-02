package httpapi

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestLibraryScanDuplicateKeyScopesSeriesMatchesToEpisodes(t *testing.T) {
	mediaID := uuid.New()
	season := int32(2)
	episode11 := int32(11)
	episode21 := int32(21)
	base := storage.LibraryScanItemInput{
		DetectedTitle:     "The Apothecary Diaries",
		DetectedMediaKind: "series",
		MediaItemID:       &mediaID,
		SeasonNumber:      &season,
	}

	first := base
	first.EpisodeNumber = &episode11
	second := base
	second.EpisodeNumber = &episode21
	if libraryScanDuplicateKey(first) == libraryScanDuplicateKey(second) {
		t.Fatal("different series episodes must not share a duplicate key")
	}

	unknown := base
	if key := libraryScanDuplicateKey(unknown); key != "" {
		t.Fatalf("unidentified series episode duplicate key = %q, want empty", key)
	}
}
