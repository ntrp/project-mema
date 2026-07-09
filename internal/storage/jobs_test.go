package storage

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	storagegen "media-manager/internal/storage/generated"
)

func TestSCNSystem006SystemJobMapperPreservesRiverFields(t *testing.T) {
	scheduledAt := time.Date(2026, time.July, 3, 5, 0, 0, 0, time.UTC)
	createdAt := scheduledAt.Add(-time.Minute)
	attemptedAt := scheduledAt.Add(time.Minute)
	finalizedAt := scheduledAt.Add(2 * time.Minute)

	job := systemJobFromGetRow(storagegen.GetSystemJobRow{
		ID:          42,
		State:       "running",
		Kind:        "media.release_search",
		Queue:       "media_search",
		Attempt:     2,
		MaxAttempts: 5,
		Priority:    3,
		Args:        `{"media_item_id":"movie-1"}`,
		Metadata:    `{"source":"test"}`,
		Errors:      `[{"error":"timeout"}]`,
		InfoMessage: "timeout",
		ScheduledAt: scheduledAt,
		CreatedAt:   createdAt,
		AttemptedAt: &attemptedAt,
		FinalizedAt: &finalizedAt,
	})

	if job.ID != 42 || job.State != "running" || job.Kind != "media.release_search" {
		t.Fatalf("identity/status = %#v", job)
	}
	if job.Queue != "media_search" || job.Attempt != 2 || job.MaxAttempts != 5 || job.Priority != 3 {
		t.Fatalf("queue/attempts = %#v", job)
	}
	if job.Args == "" || job.Metadata == "" || job.Errors == "" || job.InfoMessage != "timeout" {
		t.Fatalf("details = %#v", job)
	}
	if !job.ScheduledAt.Equal(scheduledAt) || job.AttemptedAt == nil || job.FinalizedAt == nil {
		t.Fatalf("timestamps = %#v", job)
	}
}

func TestSCNSystem006SystemJobScheduleOverviewMapsActiveAndLastRuns(t *testing.T) {
	createdAt := time.Date(2026, time.July, 8, 10, 0, 0, 0, time.UTC)
	finalizedAt := createdAt.Add(2 * time.Minute)
	progress := int32(42)

	schedule := systemJobScheduleFromRow(storagegen.ListSystemJobSchedulesRow{
		ID:                    "rss_sync",
		Name:                  "RSS sync",
		Category:              "release_search",
		Description:           "Checks feeds",
		Kind:                  "media.rss_sync",
		Queue:                 "media_search",
		IntervalSeconds:       900,
		HistoryPolicy:         "routine",
		IntervalConfigurable:  true,
		Automatic:             true,
		ManualActionAvailable: true,
		CreatedAt:             createdAt.Add(-time.Hour),
		UpdatedAt:             createdAt,
		ActiveRiverJobID:      101,
		ActiveStatus:          "running",
		ActiveProgressPercent: pgtype.Int4{Int32: progress, Valid: true},
		ActiveProgressLabel:   "Checking indexers",
		ActiveProgressData:    []byte(`{"phase":"release_search","mediaItemId":"media-1"}`),
		LastRiverJobID:        100,
		LastStatus:            "completed",
		LastCreatedAt:         createdAt,
		LastFinalizedAt:       &finalizedAt,
	})

	if schedule.ActiveRiverJobID == nil || *schedule.ActiveRiverJobID != 101 {
		t.Fatalf("active execution = %#v", schedule.ActiveRiverJobID)
	}
	if schedule.ActiveProgressPercent == nil || *schedule.ActiveProgressPercent != progress {
		t.Fatalf("progress = %#v", schedule.ActiveProgressPercent)
	}
	if schedule.ActiveProgressData["phase"] != "release_search" {
		t.Fatalf("progress data = %#v", schedule.ActiveProgressData)
	}
	if schedule.NextRunAt == nil || !schedule.NextRunAt.Equal(createdAt.Add(15*time.Minute)) {
		t.Fatalf("next run = %#v", schedule.NextRunAt)
	}
	if schedule.HistoryPolicy != "routine" || !schedule.IntervalConfigurable {
		t.Fatalf("schedule policy/configurable = %#v", schedule)
	}
	if schedule.Category != "release_search" || !schedule.Automatic || !schedule.ManualActionAvailable || !schedule.Enabled {
		t.Fatalf("schedule fulfillment settings = %#v", schedule)
	}
	if schedule.LastFinalizedAt == nil || !schedule.LastFinalizedAt.Equal(finalizedAt) {
		t.Fatalf("last finalized = %#v", schedule.LastFinalizedAt)
	}
}

func TestSCNSystem006ConfigurableSystemJobScheduleInterval(t *testing.T) {
	ctx, store := testDBStore(t)
	err := store.SyncSystemJobSchedules(ctx, []SystemJobScheduleDefinition{
		{ID: "download_activity_sync", Name: "Download activity sync", Kind: "download.activity_sync", Queue: "downloads", IntervalSeconds: 15, IntervalConfigurable: true, HistoryPolicy: "routine"},
		{ID: "rss_sync", Name: "RSS sync", Kind: "media.rss_sync", Queue: "media_search", IntervalSeconds: 900},
	})
	if err != nil {
		t.Fatalf("sync schedules: %v", err)
	}

	if _, err := store.SetSystemJobScheduleInterval(ctx, "download_activity_sync", 10); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("interval below minimum error = %v, want invalid input", err)
	}
	if _, err := store.SetSystemJobScheduleInterval(ctx, "rss_sync", 30); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("non-configurable interval error = %v, want no rows", err)
	}
	schedule, err := store.SetSystemJobScheduleInterval(ctx, "download_activity_sync", 30)
	if err != nil {
		t.Fatalf("set configurable interval: %v", err)
	}
	if schedule.IntervalSeconds != 30 {
		t.Fatalf("interval = %d, want 30", schedule.IntervalSeconds)
	}
	if err := store.SyncSystemJobSchedules(ctx, []SystemJobScheduleDefinition{
		{ID: "download_activity_sync", Name: "Download activity sync", Kind: "download.activity_sync", Queue: "downloads", IntervalSeconds: 15, IntervalConfigurable: true, HistoryPolicy: "routine"},
	}); err != nil {
		t.Fatalf("resync schedules: %v", err)
	}
	schedules, err := store.ListSystemJobSchedules(ctx)
	if err != nil {
		t.Fatalf("list schedules: %v", err)
	}
	if scheduleByID(schedules, "download_activity_sync").IntervalSeconds != 30 {
		t.Fatalf("resync overwrote configured interval: %#v", schedules)
	}
}

func TestSCNSystem006ScheduleSyncPreservesDisabledFulfillmentDefault(t *testing.T) {
	ctx, store := testDBStore(t)
	definition := SystemJobScheduleDefinition{
		ID:                    "subtitle_embed",
		Name:                  "Subtitle embed",
		Kind:                  "media.fulfillment.subtitle_embed",
		Queue:                 "media_assembly",
		IntervalSeconds:       3600,
		IntervalConfigurable:  true,
		Automatic:             true,
		ManualActionAvailable: true,
		PausedByDefault:       true,
	}
	if err := store.SyncSystemJobSchedules(ctx, []SystemJobScheduleDefinition{definition}); err != nil {
		t.Fatalf("sync schedule: %v", err)
	}
	schedules, err := store.ListSystemJobSchedules(ctx)
	if err != nil {
		t.Fatalf("list schedules: %v", err)
	}
	if schedule := scheduleByID(schedules, "subtitle_embed"); schedule.Enabled || !schedule.Paused {
		t.Fatalf("new fulfillment schedule should be disabled: %#v", schedule)
	}
	if _, err := store.SetSystemJobSchedulePaused(ctx, "subtitle_embed", false); err != nil {
		t.Fatalf("enable schedule: %v", err)
	}
	if err := store.SyncSystemJobSchedules(ctx, []SystemJobScheduleDefinition{definition}); err != nil {
		t.Fatalf("resync schedule: %v", err)
	}
	schedules, err = store.ListSystemJobSchedules(ctx)
	if err != nil {
		t.Fatalf("list schedules after resync: %v", err)
	}
	if schedule := scheduleByID(schedules, "subtitle_embed"); !schedule.Enabled || schedule.Paused {
		t.Fatalf("resync overwrote enabled state: %#v", schedule)
	}
}

func TestSCNSystem006RoutineExecutionsAreHiddenFromDefaultHistory(t *testing.T) {
	ctx, store := testDBStore(t)
	err := store.SyncSystemJobSchedules(ctx, []SystemJobScheduleDefinition{
		{ID: "download_activity_sync", Name: "Download activity sync", Kind: "download.activity_sync", Queue: "downloads", IntervalSeconds: 15, IntervalConfigurable: true, HistoryPolicy: "routine"},
		{ID: "rss_sync", Name: "RSS sync", Kind: "media.rss_sync", Queue: "media_search", IntervalSeconds: 900},
	})
	if err != nil {
		t.Fatalf("sync schedules: %v", err)
	}
	now := time.Now().UTC()
	finalized := now.Add(time.Minute)
	_, _ = store.UpsertSystemJobExecution(ctx, executionInput(201, "download_activity_sync", "download.activity_sync", "completed", now, &finalized))
	_, _ = store.UpsertSystemJobExecution(ctx, executionInput(202, "download_activity_sync", "download.activity_sync", "retryable", now.Add(time.Second), nil))
	_, _ = store.UpsertSystemJobExecution(ctx, executionInput(203, "rss_sync", "media.rss_sync", "completed", now.Add(2*time.Second), &finalized))

	defaultRows, err := store.ListSystemJobExecutions(ctx, SystemJobExecutionFilters{Limit: 10})
	if err != nil {
		t.Fatalf("list default history: %v", err)
	}
	if ids(defaultRows) != "203" {
		t.Fatalf("default history ids = %s, want 203", ids(defaultRows))
	}
	withRoutine, err := store.ListSystemJobExecutions(ctx, SystemJobExecutionFilters{Limit: 10, IncludeRoutine: true})
	if err != nil {
		t.Fatalf("list routine history: %v", err)
	}
	if ids(withRoutine) != "203,201" {
		t.Fatalf("routine history ids = %s, want 203,201", ids(withRoutine))
	}
}

func TestSCNSystem006RoutineSuccessfulExecutionsPruneEarly(t *testing.T) {
	ctx, store := testDBStore(t)
	err := store.SyncSystemJobSchedules(ctx, []SystemJobScheduleDefinition{
		{ID: "download_activity_sync", Name: "Download activity sync", Kind: "download.activity_sync", Queue: "downloads", IntervalSeconds: 15, IntervalConfigurable: true, HistoryPolicy: "routine"},
		{ID: "rss_sync", Name: "RSS sync", Kind: "media.rss_sync", Queue: "media_search", IntervalSeconds: 900},
	})
	if err != nil {
		t.Fatalf("sync schedules: %v", err)
	}
	old := time.Now().UTC().Add(-2 * time.Hour)
	_, _ = store.UpsertSystemJobExecution(ctx, executionInput(301, "download_activity_sync", "download.activity_sync", "completed", old, &old))
	_, _ = store.UpsertSystemJobExecution(ctx, executionInput(302, "download_activity_sync", "download.activity_sync", "retryable", old, &old))
	_, _ = store.UpsertSystemJobExecution(ctx, executionInput(303, "rss_sync", "media.rss_sync", "completed", old, &old))

	if _, err := store.UpdateSystemJobHistorySettings(ctx, SystemJobHistorySettings{RetentionDays: 1, RoutineRetentionHours: 1}); err != nil {
		t.Fatalf("update history settings: %v", err)
	}
	if _, err := store.GetSystemJobExecution(ctx, 301); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("routine success load error = %v, want no rows", err)
	}
	if _, err := store.GetSystemJobExecution(ctx, 302); err != nil {
		t.Fatalf("routine failure should remain: %v", err)
	}
	if _, err := store.GetSystemJobExecution(ctx, 303); err != nil {
		t.Fatalf("standard execution should remain: %v", err)
	}
}

func executionInput(id int64, scheduleID string, kind string, status string, createdAt time.Time, finalizedAt *time.Time) SystemJobExecutionInput {
	return SystemJobExecutionInput{
		RiverJobID:     id,
		ScheduleID:     scheduleID,
		Classification: "fixed",
		Status:         status,
		Kind:           kind,
		Queue:          "downloads",
		Attempt:        1,
		MaxAttempts:    3,
		Priority:       1,
		Args:           []byte("{}"),
		Metadata:       []byte("{}"),
		Errors:         []byte("[]"),
		InfoMessage:    status,
		ScheduledAt:    createdAt,
		CreatedAt:      createdAt,
		FinalizedAt:    finalizedAt,
	}
}

func ids(executions []SystemJobExecution) string {
	values := make([]string, 0, len(executions))
	for _, execution := range executions {
		values = append(values, strconv.FormatInt(execution.RiverJobID, 10))
	}
	return strings.Join(values, ",")
}

func scheduleByID(schedules []SystemJobSchedule, id string) SystemJobSchedule {
	for _, schedule := range schedules {
		if schedule.ID == id {
			return schedule
		}
	}
	return SystemJobSchedule{}
}

func TestSCNSystem006PausedSystemJobScheduleHasNoNextRun(t *testing.T) {
	updatedAt := time.Date(2026, time.July, 8, 10, 0, 0, 0, time.UTC)
	schedule := systemJobScheduleFromRow(storagegen.ListSystemJobSchedulesRow{
		ID:              "subtitle_retry",
		Name:            "Subtitle retry",
		Kind:            "subtitle.retry",
		Queue:           "media_search",
		IntervalSeconds: 60,
		Paused:          true,
		CreatedAt:       updatedAt,
		UpdatedAt:       updatedAt,
	})

	if schedule.NextRunAt != nil {
		t.Fatalf("paused schedule next run = %#v", schedule.NextRunAt)
	}
}

func TestSCNSystem006SystemJobLimitUsesSafeBounds(t *testing.T) {
	cases := []struct {
		name  string
		limit int32
		want  int32
	}{
		{name: "default", limit: 0, want: 100},
		{name: "negative", limit: -1, want: 100},
		{name: "requested", limit: 25, want: 25},
		{name: "maximum", limit: 999, want: 500},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := systemJobLimit(tc.limit); got != tc.want {
				t.Fatalf("systemJobLimit(%d) = %d, want %d", tc.limit, got, tc.want)
			}
		})
	}
}
