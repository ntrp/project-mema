package storage

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestMarshalCustomFormatSpecs(t *testing.T) {
	includeSpecs, excludeSpecs, err := marshalCustomFormatSpecs(CustomFormatInput{
		IncludeSpecs: []CustomFormatSpec{{
			ID:       "release-title",
			Name:     "Release Title",
			Type:     "releaseTitle",
			Value:    "WEB",
			Required: true,
		}},
		ExcludeSpecs: []CustomFormatSpec{{
			ID:       "not-cam",
			Name:     "Not CAM",
			Type:     "source",
			Value:    "CAM",
			Required: true,
		}},
	})
	if err != nil {
		t.Fatalf("marshal custom format specs: %v", err)
	}
	if string(includeSpecs) != `[{"id":"release-title","name":"Release Title","type":"releaseTitle","value":"WEB","required":true}]` {
		t.Fatalf("unexpected include specs JSON: %s", includeSpecs)
	}
	if string(excludeSpecs) != `[{"id":"not-cam","name":"Not CAM","type":"source","value":"CAM","required":true}]` {
		t.Fatalf("unexpected exclude specs JSON: %s", excludeSpecs)
	}
}

func TestCustomFormatsUseGeneratedQueries(t *testing.T) {
	ctx, store := testDBStore(t)
	id := uuid.New()

	created, err := store.CreateCustomFormat(ctx, CustomFormatInput{
		ID:                      id,
		Name:                    "Scenario WEB",
		IncludeInRenameTemplate: true,
		IncludeSpecs: []CustomFormatSpec{{
			ID:       "release-title",
			Name:     "Release Title",
			Type:     "releaseTitle",
			Value:    "WEB",
			Required: true,
		}},
		ExcludeSpecs: []CustomFormatSpec{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != id || created.Name != "Scenario WEB" || !created.IncludeInRenameTemplate {
		t.Fatalf("created custom format = %#v", created)
	}
	if len(created.IncludeSpecs) != 1 || created.IncludeSpecs[0].Value != "WEB" {
		t.Fatalf("created include specs = %#v", created.IncludeSpecs)
	}

	updated, err := store.UpdateCustomFormat(ctx, id, CustomFormatInput{
		Name:                    "Scenario Remux",
		IncludeInRenameTemplate: false,
		IncludeSpecs: []CustomFormatSpec{{
			ID:       "source",
			Name:     "Source",
			Type:     "source",
			Value:    "REMUX",
			Required: true,
		}},
		ExcludeSpecs: []CustomFormatSpec{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "Scenario Remux" || updated.IncludeInRenameTemplate {
		t.Fatalf("updated custom format = %#v", updated)
	}

	listed, err := store.ListCustomFormats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !customFormatListHasID(listed, id) {
		t.Fatalf("created custom format missing from list: %#v", listed)
	}

	if _, err := store.UpdateCustomFormat(ctx, uuid.New(), CustomFormatInput{Name: "Missing"}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("update missing custom format error = %v, want %v", err, ErrNotFound)
	}
	if err := store.DeleteCustomFormat(ctx, id); err != nil {
		t.Fatal(err)
	}
	if err := store.DeleteCustomFormat(ctx, id); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing custom format error = %v, want %v", err, ErrNotFound)
	}
}

func customFormatListHasID(formats []CustomFormat, id uuid.UUID) bool {
	for _, format := range formats {
		if format.ID == id {
			return true
		}
	}
	return false
}
