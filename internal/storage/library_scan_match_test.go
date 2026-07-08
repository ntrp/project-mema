package storage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestImportLibraryScanItemRecordsOriginalFileNameProvenance(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	fileName := "Linked.Movie.2026.1080p-GRP.mkv"
	path := filepath.Join(folder.Path, fileName)
	if err := os.WriteFile(path, []byte("movie"), 0o644); err != nil {
		t.Fatal(err)
	}
	title := "Linked Movie"
	kind := "movie"
	scan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{{
		Path:              path,
		FileName:          fileName,
		SizeBytes:         1024,
		DetectedTitle:     title,
		DetectedYear:      int32Ptr(2026),
		DetectedMediaKind: kind,
	}})
	if err != nil {
		t.Fatal(err)
	}

	_, item, err := store.ImportLibraryScanItem(ctx, scan.ID, scan.Items[0].ID, LibraryMatchInput{
		MediaKind:           kind,
		Title:               title,
		Year:                int32Ptr(2026),
		Monitored:           true,
		QualityProfileID:    "any",
		MonitorMode:         "movie",
		MinimumAvailability: "released",
	})
	if err != nil {
		t.Fatal(err)
	}
	provenance, err := store.ListMediaComponentProvenance(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(provenance) != 1 {
		t.Fatalf("provenance = %#v", provenance)
	}
	chain := provenance[0].TransformationChain
	if len(chain) != 1 || chain[0]["originalFileName"] != fileName {
		t.Fatalf("transformation chain = %#v, want originalFileName %q", chain, fileName)
	}
}

func TestCreateLibraryScanCanonicalizesItemPath(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(folder.Path, "Nested", "..", "Linked.Movie.2026.mkv")

	scan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{{
		Path:              path,
		FileName:          "stale-name.mkv",
		DetectedTitle:     "Linked Movie",
		DetectedMediaKind: "movie",
	}})
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(folder.Path, "Linked.Movie.2026.mkv")
	if scan.Items[0].Path != expected || scan.Items[0].FileName != filepath.Base(expected) {
		t.Fatalf("scan item path = %q, file name = %q", scan.Items[0].Path, scan.Items[0].FileName)
	}
}

func TestResetLibraryScanItemImportClearsLinkedLibraryRows(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Linked Movie",
		Year:            int32Ptr(2026),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	title := "Linked Movie"
	kind := "movie"
	source := "library"
	input := LibraryScanItemInput{
		Path:              filepath.Join(folder.Path, "Linked.Movie.2026.mkv"),
		FileName:          "Linked.Movie.2026.mkv",
		SizeBytes:         1024,
		DetectedTitle:     title,
		DetectedYear:      int32Ptr(2026),
		DetectedMediaKind: kind,
		Imported:          true,
		MatchedTitle:      &title,
		MatchedYear:       int32Ptr(2026),
		MatchedMediaKind:  &kind,
		MatchSource:       &source,
		MediaItemID:       &item.ID,
	}
	previousScan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{input})
	if err != nil {
		t.Fatal(err)
	}
	scan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{{
		Path:              input.Path,
		FileName:          "Linked.Movie.2026.mkv",
		SizeBytes:         1024,
		DetectedTitle:     title,
		DetectedYear:      int32Ptr(2026),
		DetectedMediaKind: kind,
		Imported:          true,
		MatchedTitle:      &title,
		MatchedYear:       int32Ptr(2026),
		MatchedMediaKind:  &kind,
		MatchSource:       &source,
		MediaItemID:       &item.ID,
	}})
	if err != nil {
		t.Fatal(err)
	}

	reset, err := store.ResetLibraryScanItemImport(ctx, scan.ID, scan.Items[0].ID)
	if err != nil {
		t.Fatalf("reset linked scan item: %v", err)
	}
	if reset.Item.Imported || reset.Item.MediaItemID != nil || reset.Item.MatchSource != nil {
		t.Fatalf("reset item = %#v", reset.Item)
	}
	if reset.RemovedMediaItemID == nil || *reset.RemovedMediaItemID != item.ID {
		t.Fatalf("removed media item = %v, want %s", reset.RemovedMediaItemID, item.ID)
	}
	previous, err := store.GetLibraryScan(ctx, previousScan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if previous.Items[0].Imported || previous.Items[0].MediaItemID != nil {
		t.Fatalf("previous scan item still linked = %#v", previous.Items[0])
	}
	if _, err := store.GetMediaItem(ctx, item.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("media item still exists, err=%v", err)
	}
}
