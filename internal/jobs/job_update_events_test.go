package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/riverqueue/river/rivertype"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

func TestSCNSystem008JobUpdatedPublishesObservablePayload(t *testing.T) {
	broker := events.NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updates := broker.Subscribe(ctx)

	scheduledAt := time.Date(2026, 7, 3, 1, 2, 3, 0, time.UTC)
	createdAt := scheduledAt.Add(-time.Hour)
	attemptedAt := scheduledAt.Add(-time.Minute)
	row := &rivertype.JobRow{
		ID:          42,
		State:       rivertype.JobStateAvailable,
		Kind:        "media.release_search",
		Queue:       "media_search",
		Attempt:     2,
		MaxAttempts: 3,
		Priority:    1,
		EncodedArgs: []byte(`{"mediaItemId":"media-1"}`),
		Metadata:    []byte(`{"source":"test"}`),
		Errors: []rivertype.AttemptError{{
			At:      scheduledAt,
			Attempt: 1,
			Error:   " indexer timed out ",
		}},
		ScheduledAt: scheduledAt,
		CreatedAt:   createdAt,
		AttemptedAt: &attemptedAt,
	}

	publishJobUpdated(broker, row, "running")

	update := readJobUpdate(t, updates)
	if update.ID != row.ID || update.Status != "running" || update.Kind != row.Kind {
		t.Fatalf("unexpected update identity: %#v", update)
	}
	if update.Args != `{"mediaItemId":"media-1"}` || update.Metadata != `{"source":"test"}` {
		t.Fatalf("unexpected update json fields: %#v", update)
	}
	if update.InfoMessage != "indexer timed out" {
		t.Fatalf("info message = %q", update.InfoMessage)
	}
	if update.FinalizedAt != nil {
		t.Fatalf("running update finalizedAt = %v, want nil", update.FinalizedAt)
	}
}

func TestSCNSystem008JobFinishedPublishesRetryableAndCompletedStates(t *testing.T) {
	broker := events.NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updates := broker.Subscribe(ctx)
	row := &rivertype.JobRow{
		ID:          43,
		State:       rivertype.JobStateAvailable,
		Kind:        "download.activity_sync",
		Queue:       "downloads",
		MaxAttempts: 3,
		ScheduledAt: time.Date(2026, 7, 3, 2, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, 7, 3, 1, 0, 0, 0, time.UTC),
	}

	publishJobFinished(broker, row, errors.New("download client unavailable"))
	if update := readJobUpdate(t, updates); update.Status != "retryable" || update.FinalizedAt != nil {
		t.Fatalf("retryable update = %#v", update)
	}

	publishJobFinished(broker, row, nil)
	update := readJobUpdate(t, updates)
	if update.Status != "completed" {
		t.Fatalf("completed status = %q", update.Status)
	}
	if update.FinalizedAt == nil {
		t.Fatal("completed update did not include finalizedAt")
	}
	if update.Args != "{}" || update.Metadata != "{}" || update.Errors != "[]" {
		t.Fatalf("completed update fallbacks = %#v", update)
	}
}

func TestSCNSystem008JobExecutionInputClassifiesPeriodicJobs(t *testing.T) {
	row := &rivertype.JobRow{
		ID:          44,
		State:       rivertype.JobStateAvailable,
		Kind:        "media.rss_sync",
		Queue:       "media_search",
		MaxAttempts: 3,
		Metadata:    []byte(`{"river:periodic_job_id":"rss_sync"}`),
		ScheduledAt: time.Date(2026, 7, 8, 10, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2026, 7, 8, 9, 59, 0, 0, time.UTC),
	}

	input := jobExecutionInputFromRow(row, "running")
	if input.ScheduleID != "rss_sync" || input.Classification != "fixed" {
		t.Fatalf("periodic classification = %#v", input)
	}

	row.Metadata = []byte(`{"source":"manual"}`)
	input = jobExecutionInputFromRow(row, "")
	if input.ScheduleID != "" || input.Classification != "one_shot" || input.Status != "available" {
		t.Fatalf("one-shot classification = %#v", input)
	}
}

func TestSCNSystem008JobExecutionEventPreservesProgress(t *testing.T) {
	progress := int32(67)
	execution := jobExecutionEventFromStorage(storage.SystemJobExecution{
		RiverJobID:      45,
		Classification:  "one_shot",
		HistoryPolicy:   "routine",
		Status:          "running",
		Kind:            "media.release_search",
		Queue:           "media_search",
		ProgressPercent: &progress,
		ProgressLabel:   "Searching indexers",
	})

	if execution.RiverJobID != 45 || execution.ProgressPercent == nil || *execution.ProgressPercent != progress {
		t.Fatalf("execution event = %#v", execution)
	}
	if execution.HistoryPolicy != "routine" {
		t.Fatalf("history policy = %q", execution.HistoryPolicy)
	}
	if execution.ProgressLabel != "Searching indexers" {
		t.Fatalf("progress label = %q", execution.ProgressLabel)
	}
}

func readJobUpdate(t *testing.T, updates <-chan events.Event) jobUpdateEvent {
	t.Helper()
	select {
	case event := <-updates:
		if event.Type != "system.job.updated" {
			t.Fatalf("event type = %q", event.Type)
		}
		update, ok := event.Data.(jobUpdateEvent)
		if !ok {
			t.Fatalf("event data = %#v", event.Data)
		}
		return update
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for job update")
		return jobUpdateEvent{}
	}
}
