package storage

import (
	"context"
	"time"

	storagegen "media-manager/internal/storage/generated"
)

type SystemJobFilters struct {
	States []string
	Queue  string
	Kind   string
	Query  string
	Limit  int32
}

type SystemJob struct {
	ID          int64
	State       string
	Kind        string
	Queue       string
	Attempt     int32
	MaxAttempts int32
	Priority    int32
	Args        string
	Metadata    string
	Errors      string
	InfoMessage string
	ScheduledAt time.Time
	CreatedAt   time.Time
	AttemptedAt *time.Time
	FinalizedAt *time.Time
}

func (s *SettingsStore) ListSystemJobs(ctx context.Context, filters SystemJobFilters) ([]SystemJob, error) {
	rows, err := storagegen.New(s.pool).ListSystemJobs(ctx, storagegen.ListSystemJobsParams{
		States:      filters.States,
		Queue:       filters.Queue,
		Kind:        filters.Kind,
		SearchQuery: filters.Query,
		RowLimit:    systemJobLimit(filters.Limit),
	})
	if err != nil {
		return nil, err
	}

	jobs := make([]SystemJob, 0, len(rows))
	for _, row := range rows {
		jobs = append(jobs, systemJobFromListRow(row))
	}
	return jobs, nil
}

func (s *SettingsStore) GetSystemJob(ctx context.Context, id int64) (SystemJob, error) {
	row, err := storagegen.New(s.pool).GetSystemJob(ctx, id)
	return systemJobFromGetRow(row), err
}

func systemJobFromListRow(row storagegen.ListSystemJobsRow) SystemJob {
	return systemJobFromGetRow(storagegen.GetSystemJobRow(row))
}

func systemJobFromGetRow(row storagegen.GetSystemJobRow) SystemJob {
	return SystemJob{
		ID:          row.ID,
		State:       row.State,
		Kind:        row.Kind,
		Queue:       row.Queue,
		Attempt:     row.Attempt,
		MaxAttempts: row.MaxAttempts,
		Priority:    row.Priority,
		Args:        row.Args,
		Metadata:    row.Metadata,
		Errors:      row.Errors,
		InfoMessage: row.InfoMessage,
		ScheduledAt: row.ScheduledAt,
		CreatedAt:   row.CreatedAt,
		AttemptedAt: row.AttemptedAt,
		FinalizedAt: row.FinalizedAt,
	}
}

func systemJobLimit(limit int32) int32 {
	if limit <= 0 {
		return 100
	}
	if limit > 500 {
		return 500
	}
	return limit
}
