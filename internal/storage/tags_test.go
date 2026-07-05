package storage

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNormalizeTagNames(t *testing.T) {
	tags := normalizeTagNames([]string{
		"  Anime  ",
		"anime",
		"",
		"  4K   Preferred ",
		"4k preferred",
		"Documentary",
	})

	expectStrings(t, tags, []string{"Anime", "4K Preferred", "Documentary"})
}

func TestNormalizeTagName(t *testing.T) {
	if got := normalizeTagName("  Family   Movies  "); got != "Family Movies" {
		t.Fatalf("expected compacted tag name, got %q", got)
	}
}

func TestTagsUseGeneratedQueries(t *testing.T) {
	ctx, store := testDBStore(t)

	created, err := store.SaveTag(ctx, nil, " Scenario   Tag ")
	if err != nil {
		t.Fatal(err)
	}
	if created.Name != "Scenario Tag" {
		t.Fatalf("created tag = %#v", created)
	}

	upserted, err := store.SaveTag(ctx, nil, "scenario tag")
	if err != nil {
		t.Fatal(err)
	}
	if upserted.ID != created.ID || upserted.Name != "scenario tag" {
		t.Fatalf("upserted tag = %#v, want id %s and updated name", upserted, created.ID)
	}

	renamed, err := store.SaveTag(ctx, &created.ID, " Renamed   Tag ")
	if err != nil {
		t.Fatal(err)
	}
	if renamed.ID != created.ID || renamed.Name != "Renamed Tag" {
		t.Fatalf("renamed tag = %#v", renamed)
	}

	listed, err := store.ListTags(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !tagListHasID(listed, created.ID) {
		t.Fatalf("created tag missing from list: %#v", listed)
	}

	if err := store.DeleteTag(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.DeleteTag(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing tag error = %v, want %v", err, ErrNotFound)
	}
}

func tagListHasID(tags []Tag, id uuid.UUID) bool {
	for _, tag := range tags {
		if tag.ID == id {
			return true
		}
	}
	return false
}
