package httpapi

import (
	"testing"
	"time"

	"media-manager/internal/storage"
)

func TestSCNSystem006SystemJobResponsePreservesRuntimeFields(t *testing.T) {
	scheduledAt := time.Date(2026, time.July, 3, 5, 0, 0, 0, time.UTC)
	createdAt := scheduledAt.Add(-time.Minute)
	attemptedAt := scheduledAt.Add(time.Minute)
	finalizedAt := scheduledAt.Add(2 * time.Minute)

	response := systemJobResponse(storage.SystemJob{
		ID:          42,
		State:       "retryable",
		Kind:        "auto_search",
		Queue:       "default",
		Attempt:     2,
		MaxAttempts: 5,
		Priority:    3,
		Args:        `{"mediaItemId":"movie-1"}`,
		Metadata:    `{"source":"test"}`,
		Errors:      `last error`,
		InfoMessage: "Waiting for indexer",
		ScheduledAt: scheduledAt,
		CreatedAt:   createdAt,
		AttemptedAt: &attemptedAt,
		FinalizedAt: &finalizedAt,
	})

	if response.Id != 42 || response.Status != "retryable" || response.Kind != "auto_search" {
		t.Fatalf("identity/status = %#v", response)
	}
	if response.Queue != "default" || response.Attempt != 2 || response.MaxAttempts != 5 {
		t.Fatalf("queue/attempts = %#v", response)
	}
	if response.Args == "" || response.Metadata == "" || response.Errors == "" || response.InfoMessage == "" {
		t.Fatalf("expected job detail strings to be preserved: %#v", response)
	}
	if !response.ScheduledAt.Equal(scheduledAt) || response.AttemptedAt == nil || response.FinalizedAt == nil {
		t.Fatalf("timestamps = %#v", response)
	}
}

func TestSCNSystem006SystemJobFilterParamsAreNormalized(t *testing.T) {
	statuses := []string{" available ", "", "running", " "}
	query := "  movie "
	limit := int32(25)

	if got := stringList(&statuses); len(got) != 2 || got[0] != "available" || got[1] != "running" {
		t.Fatalf("stringList = %#v", got)
	}
	if stringList(nil) != nil {
		t.Fatal("nil string list should stay nil")
	}
	if got := optionalStringParam(&query); got != "movie" {
		t.Fatalf("optionalStringParam = %q, want movie", got)
	}
	if got := optionalStringParam(nil); got != "" {
		t.Fatalf("optionalStringParam(nil) = %q, want empty", got)
	}
	if got := optionalInt32(&limit); got != 25 {
		t.Fatalf("optionalInt32 = %d, want 25", got)
	}
	if got := optionalInt32(nil); got != 0 {
		t.Fatalf("optionalInt32(nil) = %d, want 0", got)
	}
}

func TestSCNSystem006SystemJobOverviewResponsesPreserveProgressAndSettings(t *testing.T) {
	now := time.Date(2026, time.July, 8, 10, 0, 0, 0, time.UTC)
	activeID := int64(101)
	progress := int32(50)
	schedule := systemJobScheduleResponse(storage.SystemJobSchedule{
		ID:                    "rss_sync",
		Name:                  "RSS sync",
		Kind:                  "media.rss_sync",
		Queue:                 "media_search",
		IntervalSeconds:       900,
		ActiveRiverJobID:      &activeID,
		ActiveStatus:          "running",
		ActiveProgressPercent: &progress,
		ActiveProgressLabel:   "Checking indexers",
		CreatedAt:             now,
		UpdatedAt:             now,
	})

	if schedule.ActiveRiverJobId == nil || *schedule.ActiveRiverJobId != activeID {
		t.Fatalf("active job = %#v", schedule)
	}
	if schedule.ActiveProgressPercent == nil || *schedule.ActiveProgressPercent != progress {
		t.Fatalf("progress = %#v", schedule.ActiveProgressPercent)
	}
	if settings := systemJobHistorySettingsResponse(storage.SystemJobHistorySettings{RetentionDays: 45}); settings.RetentionDays != 45 {
		t.Fatalf("settings = %#v", settings)
	}
}

func TestSCNSystem006SystemJobExecutionResponsePreservesLogsContract(t *testing.T) {
	now := time.Date(2026, time.July, 8, 10, 0, 0, 0, time.UTC)
	progress := int32(75)
	response := systemJobExecutionResponse(storage.SystemJobExecution{
		RiverJobID:      77,
		ScheduleID:      "download_activity_sync",
		Classification:  "fixed",
		Status:          "running",
		Kind:            "download.activity_sync",
		Queue:           "downloads",
		ProgressPercent: &progress,
		ProgressLabel:   "Checking downloads",
		ScheduledAt:     now,
		CreatedAt:       now,
		UpdatedAt:       now,
	})

	if response.ScheduleId == nil || *response.ScheduleId != "download_activity_sync" {
		t.Fatalf("schedule id = %#v", response.ScheduleId)
	}
	if response.Classification != Fixed || response.ProgressPercent == nil || *response.ProgressPercent != progress {
		t.Fatalf("execution response = %#v", response)
	}
	logs := systemJobExecutionLogResponses([]storage.SystemJobExecutionLog{{
		ID: 1, RiverJobID: 77, Severity: "info", Message: "Started", Data: map[string]any{"queue": "downloads"}, CreatedAt: now,
	}})
	if len(logs) != 1 || logs[0].Data["queue"] != "downloads" {
		t.Fatalf("logs = %#v", logs)
	}
}

func TestSCNSystem006HistoryLimitKeepsOneRowForHasMore(t *testing.T) {
	limit := int32(500)
	if got := historyLimit(&limit); got != 499 {
		t.Fatalf("historyLimit(500) = %d, want 499", got)
	}
	if got := historyLimit(nil); got != 100 {
		t.Fatalf("historyLimit(nil) = %d, want 100", got)
	}
}
