package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

func mediaComponentAssemblyRunWithInputs(
	ctx context.Context,
	q storagegen.DBTX,
	row storagegen.AppMediaComponentAssemblyRun,
	err error,
) (MediaComponentAssemblyRun, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaComponentAssemblyRun{}, ErrNotFound
	}
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	run := mediaComponentAssemblyRunFromRow(row)
	inputs, err := storagegen.New(q).ListMediaComponentAssemblyInputs(ctx, row.ID)
	if err != nil {
		return MediaComponentAssemblyRun{}, err
	}
	for _, input := range inputs {
		run.Inputs = append(run.Inputs, mediaComponentAssemblyInputFromRow(input))
	}
	return run, nil
}

func mediaComponentAssemblyRunFromRow(row storagegen.AppMediaComponentAssemblyRun) MediaComponentAssemblyRun {
	return MediaComponentAssemblyRun{
		ID:           row.ID,
		MediaItemID:  row.MediaItemID,
		BaseSourceID: row.BaseSourceID,
		OutputPath:   row.OutputPath,
		Status:       row.Status,
		ToolName:     row.ToolName,
		ToolSummary:  row.ToolSummary,
		ErrorMessage: textPtr(row.ErrorMessage),
		JobID:        textPtr(row.JobID),
		SizeBytes:    int8Ptr(row.SizeBytes),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		CompletedAt:  row.CompletedAt,
	}
}

func mediaComponentAssemblyInputFromRow(row storagegen.AppMediaComponentAssemblyInput) MediaComponentAssemblyInput {
	return MediaComponentAssemblyInput{
		ID:         row.ID,
		RunID:      row.RunID,
		SourceID:   row.SourceID,
		ArtifactID: row.ArtifactID,
		StreamType: row.StreamType,
		InputPath:  row.InputPath,
		Provenance: jsonMap(row.Provenance),
		CreatedAt:  row.CreatedAt,
	}
}
