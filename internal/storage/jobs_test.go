package storage

import (
	"testing"
	"time"
)

func TestSCNSystem006ScanSystemJobPreservesRiverFields(t *testing.T) {
	scheduledAt := time.Date(2026, time.July, 3, 5, 0, 0, 0, time.UTC)
	createdAt := scheduledAt.Add(-time.Minute)
	attemptedAt := scheduledAt.Add(time.Minute)
	finalizedAt := scheduledAt.Add(2 * time.Minute)

	job, err := scanSystemJob(fakeSystemJobRow{
		values: []any{
			int64(42),
			"running",
			"media.release_search",
			"media_search",
			int32(2),
			int32(5),
			int32(3),
			`{"media_item_id":"movie-1"}`,
			`{"source":"test"}`,
			`[{"error":"timeout"}]`,
			"timeout",
			scheduledAt,
			createdAt,
			&attemptedAt,
			&finalizedAt,
		},
	})
	if err != nil {
		t.Fatalf("scanSystemJob returned error: %v", err)
	}

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

type fakeSystemJobRow struct {
	values []any
}

func (r fakeSystemJobRow) Scan(dest ...any) error {
	for index, value := range r.values {
		switch target := dest[index].(type) {
		case *int64:
			*target = value.(int64)
		case *int32:
			*target = value.(int32)
		case *string:
			*target = value.(string)
		case *time.Time:
			*target = value.(time.Time)
		case **time.Time:
			*target = value.(*time.Time)
		}
	}
	return nil
}
