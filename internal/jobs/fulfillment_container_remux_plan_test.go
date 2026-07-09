package jobs

import (
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestContainerRemuxPlansCreateFileScopedJobs(t *testing.T) {
	mediaID := uuid.New()
	container := "mkv"
	item := storage.MediaItem{
		ID:             mediaID,
		FinalContainer: "mp4",
		FileFacts: []storage.MediaFileFact{{
			ID:              uuid.New(),
			FilePath:        "/library/movie.mkv",
			ContainerFormat: &container,
		}},
	}

	plans := containerRemuxPlans(item, FulfillmentActionArgs{})
	if len(plans) != 1 {
		t.Fatalf("plans = %#v", plans)
	}
	got := plans[0].args
	if got.MediaItemID != mediaID.String() || got.FilePath != "/library/movie.mkv" ||
		got.TargetType != "video" || got.TrackID != "" || got.Manual {
		t.Fatalf("planned args = %#v", got)
	}
}

func TestContainerRemuxPlansSkipMatchingContainer(t *testing.T) {
	container := "mp4"
	item := storage.MediaItem{
		ID:             uuid.New(),
		FinalContainer: "mp4",
		FileFacts: []storage.MediaFileFact{{
			FilePath:        "/library/movie.mp4",
			ContainerFormat: &container,
		}},
	}

	if plans := containerRemuxPlans(item, FulfillmentActionArgs{}); len(plans) != 0 {
		t.Fatalf("plans = %#v", plans)
	}
}
