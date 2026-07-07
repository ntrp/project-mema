package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestApplySelectedMediaItemRenameSkipsUncheckedRows(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.SaveFileNamingSettings(ctx, defaultFileNamingSettingsInput()); err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Scenario Movie",
		Year:            int32Ptr(2026),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	firstPath := filepath.Join(*item.MediaFolderPath, "First.Name.mkv")
	secondPath := filepath.Join(*item.MediaFolderPath, "Second.Name.mkv")
	writePreviewFile(t, firstPath)
	writePreviewFile(t, secondPath)
	if err := store.RecordImportedMediaFile(ctx, item, firstPath); err != nil {
		t.Fatal(err)
	}
	if err := store.RecordImportedMediaFile(ctx, item, secondPath); err != nil {
		t.Fatal(err)
	}

	result, err := store.ApplySelectedMediaItemRename(ctx, item.ID, []string{firstPath})
	if err != nil {
		t.Fatal(err)
	}
	if result.AppliedCount != 1 || result.SkippedCount != 1 || result.FailedCount != 0 {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Stat(firstPath); !os.IsNotExist(err) {
		t.Fatalf("first source stat err = %v, want missing", err)
	}
	if _, err := os.Stat(secondPath); err != nil {
		t.Fatalf("second source should remain: %v", err)
	}
}
