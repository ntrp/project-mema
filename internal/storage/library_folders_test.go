package storage

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestLibraryFoldersUseGeneratedQueries(t *testing.T) {
	ctx, store := testDBStore(t)
	path := filepath.Join(t.TempDir(), "library")

	created, err := store.CreateLibraryFolder(ctx, path, "movie")
	if err != nil {
		t.Fatal(err)
	}
	if created.Path != path || created.Kind != "movie" {
		t.Fatalf("created folder = %#v", created)
	}

	upserted, err := store.CreateLibraryFolderWithManual(ctx, path, "series", true)
	if err != nil {
		t.Fatal(err)
	}
	if upserted.ID != created.ID {
		t.Fatalf("upserted folder id = %s, want %s", upserted.ID, created.ID)
	}
	if upserted.Kind != "series" || !upserted.Manual {
		t.Fatalf("upserted folder = %#v, want manual series", upserted)
	}
	allowed, err := store.MediaItemAllowsAutomation(ctx, MediaItem{LibraryFolderID: &created.ID}, "")
	if err != nil || allowed {
		t.Fatalf("manual folder automation = %v, %v; want false, nil", allowed, err)
	}
	nestedPath := filepath.Join(path, "Nested", "Movie.mkv")
	allowed, err = store.MediaItemAllowsAutomation(ctx, MediaItem{}, nestedPath)
	if err != nil || allowed {
		t.Fatalf("nested manual path automation = %v, %v; want false, nil", allowed, err)
	}
	outsidePath := filepath.Join(t.TempDir(), "Other.mkv")
	allowed, err = store.MediaItemAllowsAutomation(ctx, MediaItem{LibraryFolderID: &created.ID}, outsidePath)
	if err != nil || !allowed {
		t.Fatalf("outside path automation = %v, %v; want true, nil", allowed, err)
	}
	facts, err := store.AutomatedMediaFileFacts(ctx, MediaItem{FileFacts: []MediaFileFact{
		{FilePath: nestedPath}, {FilePath: outsidePath},
	}})
	if err != nil || len(facts) != 1 || facts[0].FilePath != outsidePath {
		t.Fatalf("automated file facts = %#v, %v; want outside path only", facts, err)
	}
	updated, err := store.UpdateLibraryFolderManual(ctx, created.ID, false)
	if err != nil || updated.Manual {
		t.Fatalf("updated folder = %#v, %v; want automation enabled", updated, err)
	}

	fetched, err := store.GetLibraryFolder(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetched.ID != created.ID || fetched.Path != path {
		t.Fatalf("fetched folder = %#v", fetched)
	}

	exists, err := store.LibraryFolderExists(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected created folder to exist")
	}

	listed, err := store.ListLibraryFolders(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !libraryFolderListHasID(listed, created.ID) {
		t.Fatalf("created folder missing from list: %#v", listed)
	}

	if err := store.DeleteLibraryFolder(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetLibraryFolder(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing folder error = %v, want %v", err, ErrNotFound)
	}
	if err := store.DeleteLibraryFolder(ctx, uuid.New()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing folder error = %v, want %v", err, ErrNotFound)
	}
}

func TestCreateLibraryFolderStoresAbsolutePath(t *testing.T) {
	ctx, store := testDBStore(t)
	relativePath := filepath.Join(".data", "relative-library")

	created, err := store.CreateLibraryFolder(ctx, relativePath, "movie")
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.Abs(relativePath)
	if err != nil {
		t.Fatal(err)
	}
	if created.Path != want {
		t.Fatalf("created path = %q, want %q", created.Path, want)
	}
	if !filepath.IsAbs(created.Path) {
		t.Fatalf("created path is not absolute: %q", created.Path)
	}
}

func libraryFolderListHasID(folders []LibraryFolder, id uuid.UUID) bool {
	for _, folder := range folders {
		if folder.ID == id {
			return true
		}
	}
	return false
}
