package storage

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestDeleteLibraryFolderDeletesDetectedMediaAndRelatedData(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Deleted Movie",
		Year:            int32Ptr(2026),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	mediaFile := filepath.Join(*item.MediaFolderPath, "Deleted.Movie.2026.mkv")
	if err := store.RecordImportedMediaFileWithHistory(ctx, item, "/downloads/deleted.mkv", mediaFile, "hardlink"); err != nil {
		t.Fatal(err)
	}
	activityID := uuid.New()
	if _, err := store.CreateImportAttempt(ctx, ImportAttemptInput{
		ActivityID:             activityID,
		MediaItemID:            item.ID,
		TargetPath:             &mediaFile,
		Status:                 "succeeded",
		InsertedMediaFilePaths: []string{mediaFile},
	}); err != nil {
		t.Fatal(err)
	}
	other := createLibraryMediaItem(t, ctx, store, "Other Movie")

	if err := store.DeleteLibraryFolder(ctx, folder.ID); err != nil {
		t.Fatalf("delete library folder: %v", err)
	}

	if _, err := store.GetMediaItem(ctx, item.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted folder media still exists, err=%v", err)
	}
	if _, err := store.GetMediaItem(ctx, other.ID); err != nil {
		t.Fatalf("other folder media was deleted: %v", err)
	}
	if countLibraryDeleteRows(t, ctx, store, "select count(*) from app.media_file_history where file_path = $1", mediaFile) != 0 {
		t.Fatal("media file history was not deleted")
	}
	if countLibraryDeleteRows(t, ctx, store, "select count(*) from app.import_attempts where activity_id = $1", activityID) != 0 {
		t.Fatal("import attempts were not deleted")
	}
	if countLibraryDeleteRows(t, ctx, store, "select count(*) from app.media_component_provenance where media_item_id = $1", item.ID) != 0 {
		t.Fatal("component provenance was not deleted")
	}
}

func createLibraryMediaItem(t *testing.T, ctx context.Context, store *SettingsStore, title string) MediaItem {
	t.Helper()
	folder, err := store.CreateLibraryFolder(ctx, filepath.Join(t.TempDir(), title), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           title,
		Year:            int32Ptr(2026),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func countLibraryDeleteRows(
	t *testing.T,
	ctx context.Context,
	store *SettingsStore,
	query string,
	args ...any,
) int {
	t.Helper()
	var count int
	if err := store.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		t.Fatal(err)
	}
	return count
}
