package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

func TestSCNSystem008IndexerSearchHistoryPublishesObservablePayload(t *testing.T) {
	broker := events.NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updates := broker.Subscribe(ctx)
	message := "rate limited"
	createdAt := time.Date(2026, 7, 3, 4, 5, 6, 0, time.UTC)

	publishIndexerSearchHistoryCreated(broker, storage.IndexerSearchHistoryEntry{
		IndexerName:     "Local Torznab",
		IndexerProtocol: "torrent",
		MediaType:       "movie",
		Query:           "Scenario Movie 2026",
		CacheHit:        false,
		Success:         false,
		ResultCount:     0,
		Error:           &message,
		Response:        `{"error":"rate limited"}`,
		CreatedAt:       createdAt,
	})

	event := readEvent(t, updates, "indexer.search.history.created")
	payload, ok := event.Data.(indexerSearchHistoryEntryEvent)
	if !ok {
		t.Fatalf("event payload = %#v", event.Data)
	}
	if payload.IndexerName != "Local Torznab" || payload.Query != "Scenario Movie 2026" {
		t.Fatalf("unexpected history payload: %#v", payload)
	}
	if payload.Error == nil || *payload.Error != message || payload.Success {
		t.Fatalf("unexpected error payload: %#v", payload)
	}
}

func TestSCNSystem008IndexerSearchCachePayloadPreservesEntryState(t *testing.T) {
	expiresAt := time.Date(2026, 7, 3, 5, 0, 0, 0, time.UTC)
	entry := storage.IndexerSearchCacheEntry{
		IndexerID:       uuid.New(),
		IndexerName:     "Local Torznab",
		IndexerProtocol: "torrent",
		MediaType:       "series",
		Query:           "Scenario Series S01",
		ResultCount:     12,
		ExpiresAt:       expiresAt,
		CreatedAt:       expiresAt.Add(-time.Hour),
		UpdatedAt:       expiresAt.Add(-time.Minute),
		Expired:         true,
	}

	payload := indexerSearchCacheEntryPayload(entry)

	if payload.IndexerName != entry.IndexerName || payload.ResultCount != entry.ResultCount {
		t.Fatalf("payload = %#v, entry = %#v", payload, entry)
	}
	if !payload.Expired || !payload.ExpiresAt.Equal(entry.ExpiresAt) {
		t.Fatalf("unexpected cache lifetime payload: %#v", payload)
	}
}

func TestSCNMedia009DownloadActivityPublishesObservablePayload(t *testing.T) {
	broker := events.NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updates := broker.Subscribe(ctx)
	progress := 42
	errorMessage := "client disconnected"
	failureType := "download"

	publishDownloadActivity(broker, storage.DownloadActivity{
		ID:                 uuid.New(),
		MediaItemID:        uuid.New(),
		MediaTitle:         "Scenario Movie",
		MediaType:          "movie",
		ReleaseTitle:       "Scenario.Movie.2026.1080p",
		IndexerName:        "Local Torznab",
		DownloadClientName: "Local SABnzbd",
		DownloadURL:        "https://downloads.invalid/scenario",
		Status:             "failed",
		ProgressPercent:    &progress,
		Error:              &errorMessage,
		FailureType:        &failureType,
		CreatedAt:          time.Date(2026, 7, 3, 1, 0, 0, 0, time.UTC),
		UpdatedAt:          time.Date(2026, 7, 3, 1, 1, 0, 0, time.UTC),
	})

	event := readEvent(t, updates, "activity.download.updated")
	payload, ok := event.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("event payload = %#v", event.Data)
	}
	if payload["mediaTitle"] != "Scenario Movie" || payload["status"] != "failed" {
		t.Fatalf("unexpected activity payload: %#v", payload)
	}
	if payload["progressPercent"] != &progress || payload["error"] != &errorMessage {
		t.Fatalf("unexpected activity progress payload: %#v", payload)
	}
}

func readEvent(t *testing.T, updates <-chan events.Event, eventType string) events.Event {
	t.Helper()
	select {
	case event := <-updates:
		if event.Type != eventType {
			t.Fatalf("event type = %q, want %q", event.Type, eventType)
		}
		return event
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for %s", eventType)
		return events.Event{}
	}
}
