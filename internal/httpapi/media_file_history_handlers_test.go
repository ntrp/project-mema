package httpapi

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestMediaFileHistoryResponseMapsProvenance(t *testing.T) {
	mediaItemID := uuid.New()
	source := "/downloads/movie.mkv"
	destination := "/library/movie.mkv"
	failure := "permission denied"
	createdAt := time.Now().UTC()

	response := mediaFileHistoryResponse(storage.MediaFileHistoryEntry{
		ID:              uuid.New(),
		MediaItemID:     &mediaItemID,
		FilePath:        destination,
		SourcePath:      &source,
		DestinationPath: &destination,
		Operation:       "imported",
		Status:          "failed",
		ActorType:       "job",
		Details:         map[string]any{"importMode": "move"},
		FailureDetails:  &failure,
		CreatedAt:       createdAt,
	})

	if response.MediaItemId == nil || uuid.UUID(*response.MediaItemId) != mediaItemID {
		t.Fatalf("media item id = %#v", response.MediaItemId)
	}
	if response.Operation != Imported ||
		response.Status != Failed ||
		response.ActorType != MediaFileHistoryEntryActorTypeJob {
		t.Fatalf("response = %#v", response)
	}
	if response.SourcePath == nil || *response.SourcePath != source {
		t.Fatalf("source = %#v", response.SourcePath)
	}
	if response.FailureDetails == nil || *response.FailureDetails != failure {
		t.Fatalf("failure = %#v", response.FailureDetails)
	}
	if response.Details["importMode"] != "move" || !response.CreatedAt.Equal(createdAt) {
		t.Fatalf("details/time = %#v %s", response.Details, response.CreatedAt)
	}
}
