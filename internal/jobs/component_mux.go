package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"media-manager/internal/events"
	"media-manager/internal/storage"
	mediatools "media-manager/internal/tools"
)

type MediaComponentMuxArgs struct {
	RunID string `json:"run_id" river:"unique"`
}

func (MediaComponentMuxArgs) Kind() string {
	return "media.component_mux"
}

type ComponentMuxWorker struct {
	river.WorkerDefaults[MediaComponentMuxArgs]

	settings *storage.SettingsStore
	events   *events.Broker
	runner   componentMuxRunner
}

type componentMuxRunner interface {
	Mux(ctx context.Context, args []string) (string, error)
}

type mkvMergeRunner struct{}

func (mkvMergeRunner) Mux(ctx context.Context, args []string) (string, error) {
	output, err := mediatools.RunOutput(ctx, mediatools.CommandSpec{
		Name:           "mkvmerge",
		Args:           args,
		Timeout:        30 * time.Minute,
		MaxOutputBytes: 64 * 1024,
		MaxStderrBytes: 64 * 1024,
	})
	return string(output), err
}

func (w *ComponentMuxWorker) Work(ctx context.Context, job *river.Job[MediaComponentMuxArgs]) (err error) {
	publishJobUpdated(w.events, job.JobRow, "running")
	defer func() { publishJobFinished(w.events, job.JobRow, err) }()

	runID, err := uuid.Parse(job.Args.RunID)
	if err != nil {
		return fmt.Errorf("parse assembly run id: %w", err)
	}
	run, err := w.settings.StartMediaComponentAssemblyRun(ctx, runID)
	if err != nil {
		return fmt.Errorf("start component assembly: %w", err)
	}
	args, err := MkvMergeArgs(run.OutputPath, assemblyInputPaths(run.Inputs))
	if err != nil {
		return w.failRun(ctx, runID, "", err)
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "media", "Component assembly started", map[string]any{
		"runId":       run.ID.String(),
		"mediaItemId": run.MediaItemID.String(),
	})
	summary, err := w.runnerOrDefault().Mux(ctx, args)
	if err != nil {
		slog.Error("component assembly failed", "runId", runID, "error", err)
		return w.failRun(ctx, runID, summary, err)
	}
	run, err = w.settings.CompleteMediaComponentAssemblyRun(ctx, runID, summary)
	if err != nil {
		return fmt.Errorf("complete component assembly: %w", err)
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "media", "Component assembly finished", map[string]any{
		"runId":       run.ID.String(),
		"mediaItemId": run.MediaItemID.String(),
		"outputPath":  run.OutputPath,
	})
	return nil
}

func (w *ComponentMuxWorker) runnerOrDefault() componentMuxRunner {
	if w.runner != nil {
		return w.runner
	}
	return mkvMergeRunner{}
}

func (w *ComponentMuxWorker) failRun(ctx context.Context, runID uuid.UUID, summary string, err error) error {
	message := strings.TrimSpace(err.Error())
	run, markErr := w.settings.FailMediaComponentAssemblyRun(ctx, runID, summary, message)
	if markErr != nil {
		return fmt.Errorf("%w; mark failed: %v", err, markErr)
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventError, "media", "Component assembly failed", map[string]any{
		"runId":       run.ID.String(),
		"mediaItemId": run.MediaItemID.String(),
		"error":       message,
	})
	return err
}

func MkvMergeArgs(outputPath string, inputPaths []string) ([]string, error) {
	if len(inputPaths) < 2 {
		return nil, fmt.Errorf("at least two inputs are required")
	}
	if err := mediatools.SafePathArg(outputPath); err != nil {
		return nil, err
	}
	args := []string{"-o", outputPath}
	for _, inputPath := range inputPaths {
		if err := mediatools.SafePathArg(inputPath); err != nil {
			return nil, err
		}
		args = append(args, inputPath)
	}
	return args, nil
}

func assemblyInputPaths(inputs []storage.MediaComponentAssemblyInput) []string {
	paths := make([]string, 0, len(inputs))
	for _, input := range inputs {
		paths = append(paths, input.InputPath)
	}
	return paths
}
