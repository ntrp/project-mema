package storage

import (
	"context"
	"fmt"
	"strings"
	"time"
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
	clauses := []string{"true"}
	args := []any{}
	if len(filters.States) > 0 {
		args = append(args, filters.States)
		clauses = append(clauses, fmt.Sprintf("state::text = any($%d)", len(args)))
	}
	if filters.Queue != "" {
		args = append(args, filters.Queue)
		clauses = append(clauses, fmt.Sprintf("queue = $%d", len(args)))
	}
	if filters.Kind != "" {
		args = append(args, filters.Kind)
		clauses = append(clauses, fmt.Sprintf("kind = $%d", len(args)))
	}
	if filters.Query != "" {
		args = append(args, "%"+filters.Query+"%")
		index := len(args)
		clauses = append(clauses, fmt.Sprintf("(kind ilike $%d or queue ilike $%d or args::text ilike $%d or errors::text ilike $%d)", index, index, index, index))
	}
	args = append(args, systemJobLimit(filters.Limit))
	query := fmt.Sprintf(`
		select id, state::text, kind, queue, attempt::int, max_attempts::int, priority::int,
			args::text, metadata::text, coalesce(array_to_json(errors), '[]'::json)::text,
			coalesce(errors[array_length(errors, 1)]->>'error', errors[array_length(errors, 1)]->>'message', state::text),
			scheduled_at, created_at, attempted_at, finalized_at
		from river_job
		where %s
		order by coalesce(finalized_at, attempted_at, scheduled_at, created_at) desc, id desc
		limit $%d
	`, strings.Join(clauses, " and "), len(args))

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := []SystemJob{}
	for rows.Next() {
		job, err := scanSystemJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

func (s *SettingsStore) GetSystemJob(ctx context.Context, id int64) (SystemJob, error) {
	return scanSystemJob(s.pool.QueryRow(ctx, `
		select id, state::text, kind, queue, attempt::int, max_attempts::int, priority::int,
			args::text, metadata::text, coalesce(array_to_json(errors), '[]'::json)::text,
			coalesce(errors[array_length(errors, 1)]->>'error', errors[array_length(errors, 1)]->>'message', state::text),
			scheduled_at, created_at, attempted_at, finalized_at
		from river_job
		where id = $1
	`, id))
}

type systemJobScanner interface {
	Scan(dest ...any) error
}

func scanSystemJob(row systemJobScanner) (SystemJob, error) {
	var job SystemJob
	err := row.Scan(
		&job.ID,
		&job.State,
		&job.Kind,
		&job.Queue,
		&job.Attempt,
		&job.MaxAttempts,
		&job.Priority,
		&job.Args,
		&job.Metadata,
		&job.Errors,
		&job.InfoMessage,
		&job.ScheduledAt,
		&job.CreatedAt,
		&job.AttemptedAt,
		&job.FinalizedAt,
	)
	return job, err
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
