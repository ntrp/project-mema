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

	upserted, err := store.CreateLibraryFolder(ctx, path, "series")
	if err != nil {
		t.Fatal(err)
	}
	if upserted.ID != created.ID {
		t.Fatalf("upserted folder id = %s, want %s", upserted.ID, created.ID)
	}
	if upserted.Kind != "series" {
		t.Fatalf("upserted folder kind = %s, want series", upserted.Kind)
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

func libraryFolderListHasID(folders []LibraryFolder, id uuid.UUID) bool {
	for _, folder := range folders {
		if folder.ID == id {
			return true
		}
	}
	return false
}
