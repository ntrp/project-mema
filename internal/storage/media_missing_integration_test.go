package storage

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestScenarioSCNMedia013StorageListsWantedMedia(t *testing.T) {
	requireStorageScenario(t, "SCN-MEDIA-013")
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	prefix := "Wanted Contract " + suffix

	missing, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:      "movie",
		Title:     prefix + " Missing",
		Year:      int32Ptr(2026),
		Monitored: true,
	})
	if err != nil {
		t.Fatalf("create missing media: %v", err)
	}
	if _, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:        "movie",
		Title:       prefix + " Unmonitored",
		Year:        int32Ptr(2026),
		Monitored:   false,
		MonitorMode: "none",
	}); err != nil {
		t.Fatalf("create unmonitored media: %v", err)
	}
	downloading, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:      "movie",
		Title:     prefix + " Downloading",
		Year:      int32Ptr(2026),
		Monitored: true,
	})
	if err != nil {
		t.Fatalf("create downloading media: %v", err)
	}
	if _, err := store.CreateDownloadActivity(ctx, DownloadActivityInput{
		MediaItemID:        downloading.ID,
		ReleaseTitle:       prefix + ".2026.1080p.WEB-DL",
		IndexerName:        "Wanted Indexer",
		DownloadClientName: "Wanted Client",
		DownloadID:         stringPtr("download-" + suffix),
		DownloadURL:        "http://download.test/" + suffix,
		Status:             "downloading",
	}); err != nil {
		t.Fatalf("create download activity: %v", err)
	}

	items, err := store.ListMissingMediaItems(ctx)
	if err != nil {
		t.Fatalf("list missing media: %v", err)
	}
	var matching []MediaItem
	for _, item := range items {
		if strings.HasPrefix(item.Title, prefix) {
			matching = append(matching, item)
		}
	}
	if len(matching) != 1 || matching[0].ID != missing.ID || matching[0].Status != "missing" {
		t.Fatalf("wanted media for %q = %#v", prefix, matching)
	}
}
