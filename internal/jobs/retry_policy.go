package jobs

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func reconcileActiveJobsNonRetryable(ctx context.Context, pool *pgxpool.Pool) {
	tag, err := pool.Exec(ctx, `
		update river_job
		set max_attempts = $1
		where max_attempts <> $1
			and state in ('available', 'pending', 'retryable', 'running', 'scheduled')
	`, nonRetryableJobMaxAttempts)
	if err != nil {
		slog.Warn("reconcile non-retryable jobs failed", "error", err)
		return
	}
	if tag.RowsAffected() > 0 {
		slog.Info("reconciled active jobs as non-retryable", "count", tag.RowsAffected())
	}
}
