package storage

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) ListMediaComponentAssemblyRuns(
	ctx context.Context,
	mediaItemID uuid.UUID,
) ([]MediaComponentAssemblyRun, error) {
	return listMediaComponentAssemblyRuns(ctx, s.pool, mediaItemID)
}

func (s *SettingsStore) GetMediaComponentAssemblyRun(
	ctx context.Context,
	runID uuid.UUID,
) (MediaComponentAssemblyRun, error) {
	row, err := storagegen.New(s.pool).GetMediaComponentAssemblyRun(ctx, runID)
	return mediaComponentAssemblyRunWithInputs(ctx, s.pool, row, err)
}

func (s *SettingsStore) CreateMediaComponentAssemblyRun(
	ctx context.Context,
	mediaItemID uuid.UUID,
	input MediaComponentAssemblyRunInput,
) (MediaComponentAssemblyRun, error) {
	item, err := s.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	base, err := s.GetMediaComponentSource(ctx, mediaItemID, input.BaseSourceID)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	finalContainer, err := s.mediaComponentAssemblyContainer(ctx, item)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	runID, outputPath, err := mediaComponentAssemblyTarget(item, base, input.ArtifactIDs, finalContainer)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	artifactInputs, err := s.componentAssemblyInputs(ctx, mediaItemID, base, input.ArtifactIDs)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	defer tx.Rollback(ctx)

	q := storagegen.New(tx)
	row, err := q.CreateMediaComponentAssemblyRun(ctx, storagegen.CreateMediaComponentAssemblyRunParams{
		ID:           runID,
		MediaItemID:  mediaItemID,
		BaseSourceID: base.ID,
		OutputPath:   outputPath,
		JobID:        textValue(input.JobID),
	})
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	if _, err := createMediaComponentAssemblyInput(ctx, q, runID, assemblyBaseInput(base)); err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	for _, artifactInput := range artifactInputs {
		if _, err := createMediaComponentAssemblyInput(ctx, q, runID, artifactInput); err != nil {
			return MediaComponentAssemblyRun{}, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	run, err := mediaComponentAssemblyRunWithInputs(ctx, s.pool, row, nil)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	if err := s.WriteMediaComponentAssemblyProvenance(ctx, run.ID); err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	return run, nil
}

func (s *SettingsStore) AssignMediaComponentAssemblyJob(
	ctx context.Context,
	runID uuid.UUID,
	jobID string,
) (MediaComponentAssemblyRun, error) {
	row, err := storagegen.New(s.pool).AssignMediaComponentAssemblyJob(ctx, storagegen.AssignMediaComponentAssemblyJobParams{
		ID:    runID,
		JobID: textValue(&jobID),
	})
	return mediaComponentAssemblyRunWithInputs(ctx, s.pool, row, err)
}

func (s *SettingsStore) StartMediaComponentAssemblyRun(
	ctx context.Context,
	runID uuid.UUID,
) (MediaComponentAssemblyRun, error) {
	row, err := storagegen.New(s.pool).StartMediaComponentAssemblyRun(ctx, runID)
	return mediaComponentAssemblyRunWithInputs(ctx, s.pool, row, err)
}

func (s *SettingsStore) CompleteMediaComponentAssemblyRun(
	ctx context.Context,
	runID uuid.UUID,
	toolSummary string,
) (MediaComponentAssemblyRun, error) {
	run, err := s.GetMediaComponentAssemblyRun(ctx, runID)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	info, err := os.Stat(run.OutputPath)
	if err != nil || info.IsDir() {
		return MediaComponentAssemblyRun{}, ErrInvalidInput
	}
	size := info.Size()
	row, err := storagegen.New(s.pool).CompleteMediaComponentAssemblyRun(ctx, storagegen.CompleteMediaComponentAssemblyRunParams{
		ID:          runID,
		ToolSummary: strings.TrimSpace(toolSummary),
		SizeBytes:   int8Value(&size),
	})
	return mediaComponentAssemblyRunWithInputs(ctx, s.pool, row, err)
}

func (s *SettingsStore) FailMediaComponentAssemblyRun(
	ctx context.Context,
	runID uuid.UUID,
	toolSummary string,
	errMessage string,
) (MediaComponentAssemblyRun, error) {
	errMessage = strings.TrimSpace(errMessage)
	row, err := storagegen.New(s.pool).FailMediaComponentAssemblyRun(ctx, storagegen.FailMediaComponentAssemblyRunParams{
		ID:           runID,
		ToolSummary:  strings.TrimSpace(toolSummary),
		ErrorMessage: textValue(&errMessage),
	})
	return mediaComponentAssemblyRunWithInputs(ctx, s.pool, row, err)
}

func listMediaComponentAssemblyRuns(
	ctx context.Context,
	q storagegen.DBTX,
	mediaItemID uuid.UUID,
) ([]MediaComponentAssemblyRun, error) {
	rows, err := storagegen.New(q).ListMediaComponentAssemblyRuns(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}
	runs := make([]MediaComponentAssemblyRun, 0, len(rows))
	for _, row := range rows {
		run, err := mediaComponentAssemblyRunWithInputs(ctx, q, row, nil)
		if err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, nil
}

func (s *SettingsStore) componentAssemblyInputs(
	ctx context.Context,
	mediaItemID uuid.UUID,
	base MediaComponentSource,
	artifactIDs []uuid.UUID,
) ([]MediaComponentAssemblyInput, error) {
	if len(artifactIDs) == 0 {
		return nil, ErrInvalidInput
	}
	inputs := make([]MediaComponentAssemblyInput, 0, len(artifactIDs))
	for _, artifactID := range artifactIDs {
		artifact, err := s.GetMediaComponentArtifact(ctx, artifactID)
		if err != nil {
			return nil, err
		}
		if artifact.MediaItemID != mediaItemID || artifact.Status != "succeeded" {
			return nil, ErrInvalidInput
		}
		if !s.componentAssemblyArtifactAllowed(ctx, mediaItemID, base.ID, artifact.SourceID) {
			return nil, ErrInvalidInput
		}
		inputs = append(inputs, assemblyArtifactInput(artifact))
	}
	return inputs, nil
}

func (s *SettingsStore) componentAssemblyArtifactAllowed(
	ctx context.Context,
	mediaItemID uuid.UUID,
	baseSourceID uuid.UUID,
	componentSourceID uuid.UUID,
) bool {
	decisions, err := listMediaComponentCompatibilityForSource(ctx, s.pool, componentSourceID)
	if err != nil {
		return false
	}
	for _, decision := range decisions {
		if decision.MediaItemID == mediaItemID && decision.BaseSourceID == baseSourceID &&
			decision.AutomationState == "allowed" {
			return true
		}
	}
	return false
}

func mediaComponentAssemblyTarget(
	item MediaItem,
	base MediaComponentSource,
	artifactIDs []uuid.UUID,
	finalContainer string,
) (uuid.UUID, string, error) {
	if base.SourceRole != "baseVideo" || base.RetentionState != "retained" || len(artifactIDs) == 0 {
		return uuid.Nil, "", ErrInvalidInput
	}
	if _, err := mediaComponentSourceTarget(item, base.RetainedPath); err != nil {
		return uuid.Nil, "", err
	}
	id := uuid.New()
	extension := ".mkv"
	if finalContainer == "mp4" {
		extension = ".mp4"
	}
	target, err := mediaComponentSourceTarget(item, filepath.Join(".mema", "assemblies", id.String(), "assembled"+extension))
	return id, target, err
}

func (s *SettingsStore) mediaComponentAssemblyContainer(ctx context.Context, item MediaItem) (string, error) {
	if item.QualityProfileID == nil || strings.TrimSpace(*item.QualityProfileID) == "" {
		return "mkv", nil
	}
	profile, err := s.GetMediaProfile(ctx, *item.QualityProfileID)
	if err != nil {
		return "", err
	}
	if profile.FinalContainer == "mp4" {
		return "mp4", nil
	}
	return "mkv", nil
}

func assemblyBaseInput(base MediaComponentSource) MediaComponentAssemblyInput {
	return MediaComponentAssemblyInput{
		SourceID:   &base.ID,
		StreamType: "video",
		InputPath:  base.RetainedPath,
		Provenance: map[string]any{
			"kind":     "baseSource",
			"sourceId": base.ID.String(),
			"role":     base.SourceRole,
		},
	}
}

func assemblyArtifactInput(artifact MediaComponentArtifact) MediaComponentAssemblyInput {
	return MediaComponentAssemblyInput{
		SourceID:   &artifact.SourceID,
		ArtifactID: &artifact.ID,
		StreamType: artifact.StreamType,
		InputPath:  artifact.OutputPath,
		Provenance: map[string]any{
			"kind":       "extractedArtifact",
			"sourceId":   artifact.SourceID.String(),
			"artifactId": artifact.ID.String(),
			"streamId":   artifact.StreamID,
			"streamType": artifact.StreamType,
		},
	}
}

func createMediaComponentAssemblyInput(
	ctx context.Context,
	q *storagegen.Queries,
	runID uuid.UUID,
	input MediaComponentAssemblyInput,
) (MediaComponentAssemblyInput, error) {
	row, err := q.CreateMediaComponentAssemblyInput(ctx, storagegen.CreateMediaComponentAssemblyInputParams{
		ID:         uuid.New(),
		RunID:      runID,
		SourceID:   input.SourceID,
		ArtifactID: input.ArtifactID,
		StreamType: input.StreamType,
		InputPath:  input.InputPath,
		Provenance: jsonObject(input.Provenance),
	})
	if err != nil {
		return MediaComponentAssemblyInput{}, err
	}
	return mediaComponentAssemblyInputFromRow(row), nil
}
