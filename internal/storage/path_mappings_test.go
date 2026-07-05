package storage

import (
	"errors"
	"testing"
)

func TestPathMappingsUseGeneratedQueries(t *testing.T) {
	ctx, store := testDBStore(t)

	created, err := store.CreatePathMapping(ctx, PathMappingInput{
		ClientPath: " /downloads/ ",
		AppPath:    " /media/ ",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ClientPath != "/downloads" || created.AppPath != "/media" {
		t.Fatalf("created path mapping = %#v", created)
	}

	updated, err := store.CreatePathMapping(ctx, PathMappingInput{
		ClientPath: "/downloads",
		AppPath:    "/library",
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.ID != created.ID || updated.AppPath != "/library" {
		t.Fatalf("updated path mapping = %#v, want id %s with /library", updated, created.ID)
	}

	listed, err := store.ListPathMappings(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0].ID != created.ID || listed[0].AppPath != "/library" {
		t.Fatalf("listed path mappings = %#v", listed)
	}

	if err := store.DeletePathMapping(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.DeletePathMapping(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing mapping error = %v, want %v", err, ErrNotFound)
	}
}
