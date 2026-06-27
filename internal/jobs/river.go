package jobs

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type SmokeArgs struct {
	Message string `json:"message"`
}

func (SmokeArgs) Kind() string {
	return "system.smoke"
}

type SmokeWorker struct {
	river.WorkerDefaults[SmokeArgs]
}

func (w *SmokeWorker) Work(ctx context.Context, job *river.Job[SmokeArgs]) error {
	return nil
}

func NewClient(pool *pgxpool.Pool) (*river.Client[pgx.Tx], error) {
	workers := river.NewWorkers()
	river.AddWorker(workers, &SmokeWorker{})

	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 4},
		},
		Workers: workers,
	})
}
