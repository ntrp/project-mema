package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
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
	Mux(ctx context.Context, command string, args []string) (string, error)
}

type mkvMergeRunner struct{}

func (mkvMergeRunner) Mux(ctx context.Context, command string, args []string) (string, error) {
	output, err := mediatools.RunOutput(ctx, mediatools.CommandSpec{
		Name:           command,
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
	command, args, err := ComponentMuxCommand(run.OutputPath, run.Inputs)
	if err != nil {
		return w.failRun(ctx, runID, "", err)
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "media", "Component assembly started", map[string]any{
		"runId":       run.ID.String(),
		"mediaItemId": run.MediaItemID.String(),
	})
	summary, err := w.runnerOrDefault().Mux(ctx, command, args)
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

func ComponentMuxCommand(outputPath string, inputs []storage.MediaComponentAssemblyInput) (string, []string, error) {
	if strings.EqualFold(filepath.Ext(outputPath), ".mp4") {
		args, err := FfmpegMP4RemuxArgs(outputPath, inputs)
		return "ffmpeg", args, err
	}
	args, err := MkvMergeArgs(outputPath, assemblyInputPaths(inputs))
	return "mkvmerge", args, err
}

func FfmpegMP4RemuxArgs(outputPath string, inputs []storage.MediaComponentAssemblyInput) ([]string, error) {
	if len(inputs) < 1 {
		return nil, fmt.Errorf("at least one input is required")
	}
	if err := mediatools.SafePathArg(outputPath); err != nil {
		return nil, err
	}
	args := []string{"-y"}
	for _, input := range inputs {
		if input.StreamType == "subtitle" {
			return nil, fmt.Errorf("mp4 assembly does not support subtitle input %s", input.InputPath)
		}
		if err := mediatools.SafePathArg(input.InputPath); err != nil {
			return nil, err
		}
		args = append(args, "-i", input.InputPath)
	}
	args = append(args, "-map", "0:v:0?")
	for index, input := range inputs {
		if input.StreamType == "audio" {
			args = append(args, "-map", strconv.Itoa(index)+":a:0?")
		}
	}
	return append(args, "-c", "copy", outputPath), nil
}

func assemblyInputPaths(inputs []storage.MediaComponentAssemblyInput) []string {
	paths := make([]string, 0, len(inputs))
	for _, input := range inputs {
		paths = append(paths, input.InputPath)
	}
	return paths
}
