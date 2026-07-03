package storage

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestScenarioSCNSettings015StorageSystemEventsLifecycle(t *testing.T) {
	requireStorageScenario(t, "SCN-SETTINGS-015")
	ctx, store := testDBStore(t)
	if err := store.ClearSystemEvents(ctx); err != nil {
		t.Fatalf("clear events: %v", err)
	}

	settings, err := store.UpdateSystemEventSettings(ctx, SystemEventSettingsInput{RetentionDays: 14})
	if err != nil {
		t.Fatalf("update event settings: %v", err)
	}
	if settings.RetentionDays != 14 {
		t.Fatalf("settings = %#v", settings)
	}
	if _, err := store.UpdateSystemEventSettings(ctx, SystemEventSettingsInput{RetentionDays: 366}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid retention error, got %v", err)
	}

	event, err := store.CreateSystemEvent(ctx, SystemEventInput{
		Severity: "warning",
		Category: "indexer",
		Message:  "Indexer degraded",
		Data:     map[string]any{"indexer": "scenario"},
	})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	if event.Data["indexer"] != "scenario" {
		t.Fatalf("event data = %#v", event.Data)
	}

	events, err := store.ListSystemEvents(ctx, 1, nil)
	if err != nil {
		t.Fatalf("list events: %v", err)
	}
	if len(events) != 1 || events[0].ID != event.ID {
		t.Fatalf("events = %#v", events)
	}

	if err := store.DeleteSystemEvent(ctx, event.ID); err != nil {
		t.Fatalf("delete event: %v", err)
	}
	if err := store.DeleteSystemEvent(ctx, uuid.New()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected missing event error, got %v", err)
	}
}
