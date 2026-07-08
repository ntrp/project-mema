package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"media-manager/internal/events"
	"media-manager/internal/storage"
	mediatools "media-manager/internal/tools"
)

type MediaComponentExtractionArgs struct {
	ArtifactID string `json:"artifact_id" river:"unique"`
}

func (MediaComponentExtractionArgs) Kind() string {
	return "media.component_extract"
}

type ComponentExtractionWorker struct {
	river.WorkerDefaults[MediaComponentExtractionArgs]

	settings *storage.SettingsStore
	events   *events.Broker
	runner   componentExtractionRunner
}

type componentExtractionRunner interface {
	Extract(ctx context.Context, args []string) (string, error)
}

type mkvExtractRunner struct{}

func (mkvExtractRunner) Extract(ctx context.Context, args []string) (string, error) {
	output, err := mediatools.RunOutput(ctx, mediatools.CommandSpec{
		Name:           "mkvextract",
		Args:           args,
		Timeout:        30 * time.Minute,
		MaxOutputBytes: 64 * 1024,
		MaxStderrBytes: 64 * 1024,
	})
	return string(output), err
}

func (w *ComponentExtractionWorker) Work(
	ctx context.Context,
	job *river.Job[MediaComponentExtractionArgs],
) (err error) {
	ctx = withJobExecution(ctx, job.JobRow.ID)
	recordJobUpdated(ctx, w.settings, w.events, job.JobRow, "running")
	defer func() { recordJobFinished(ctx, w.settings, w.events, job.JobRow, err) }()
	recordJobProgress(ctx, w.settings, w.events, nil, "Extracting media components")

	artifactID, err := uuid.Parse(job.Args.ArtifactID)
	if err != nil {
		return fmt.Errorf("parse artifact id: %w", err)
	}
	artifact, err := w.settings.StartMediaComponentArtifact(ctx, artifactID)
	if err != nil {
		return fmt.Errorf("start component extraction: %w", err)
	}
	source, err := w.settings.GetMediaComponentSource(ctx, artifact.MediaItemID, artifact.SourceID)
	if err != nil {
		return w.failArtifact(ctx, artifactID, "", fmt.Errorf("load component source: %w", err))
	}
	args, err := MkvExtractTrackArgs(source.RetainedPath, artifact.StreamID, artifact.OutputPath)
	if err != nil {
		return w.failArtifact(ctx, artifactID, "", err)
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "media", "Component extraction started", map[string]any{
		"artifactId":  artifact.ID.String(),
		"mediaItemId": artifact.MediaItemID.String(),
		"sourceId":    source.ID.String(),
	})
	summary, err := w.runnerOrDefault().Extract(ctx, args)
	if err != nil {
		slog.Error("component extraction failed", "artifactId", artifactID, "error", err)
		return w.failArtifact(ctx, artifactID, summary, err)
	}
	artifact, err = w.settings.CompleteMediaComponentArtifact(ctx, artifactID, summary)
	if err != nil {
		return fmt.Errorf("complete component extraction: %w", err)
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventInfo, "media", "Component extraction finished", map[string]any{
		"artifactId":  artifact.ID.String(),
		"mediaItemId": artifact.MediaItemID.String(),
		"sourceId":    artifact.SourceID.String(),
		"outputPath":  artifact.OutputPath,
	})
	return nil
}

func (w *ComponentExtractionWorker) runnerOrDefault() componentExtractionRunner {
	if w.runner != nil {
		return w.runner
	}
	return mkvExtractRunner{}
}

func (w *ComponentExtractionWorker) failArtifact(
	ctx context.Context,
	artifactID uuid.UUID,
	summary string,
	err error,
) error {
	message := strings.TrimSpace(err.Error())
	artifact, markErr := w.settings.FailMediaComponentArtifact(ctx, artifactID, summary, message)
	if markErr != nil {
		return fmt.Errorf("%w; mark failed: %v", err, markErr)
	}
	publishSystemEvent(ctx, w.settings, w.events, jobEventError, "media", "Component extraction failed", map[string]any{
		"artifactId":  artifact.ID.String(),
		"mediaItemId": artifact.MediaItemID.String(),
		"sourceId":    artifact.SourceID.String(),
		"error":       message,
	})
	return err
}

func MkvExtractTrackArgs(sourcePath string, streamID int32, outputPath string) ([]string, error) {
	if streamID < 0 {
		return nil, fmt.Errorf("stream id must be non-negative")
	}
	if err := mediatools.SafePathArg(sourcePath); err != nil {
		return nil, err
	}
	if err := mediatools.SafePathArg(outputPath); err != nil {
		return nil, err
	}
	return []string{"tracks", sourcePath, strconv.Itoa(int(streamID)) + ":" + outputPath}, nil
}
