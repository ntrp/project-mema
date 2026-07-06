package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *SettingsStore) WriteMediaComponentAssemblyProvenance(
	ctx context.Context,
	runID uuid.UUID,
) error {
	run, err := s.GetMediaComponentAssemblyRun(ctx, runID)
	if err != nil {
		return err
	}
	for _, input := range run.Inputs {
		if input.SourceID == nil {
			continue
		}
		source, err := s.GetMediaComponentSource(ctx, run.MediaItemID, *input.SourceID)
		if err != nil {
			return err
		}
		provenance := provenanceInputForAssembly(run, input, source)
		if _, err := s.UpsertMediaComponentProvenance(ctx, provenance); err != nil {
			return err
		}
	}
	container := containerProvenanceInput(run)
	_, err = s.UpsertMediaComponentProvenance(ctx, container)
	return err
}

func provenanceInputForAssembly(
	run MediaComponentAssemblyRun,
	input MediaComponentAssemblyInput,
	source MediaComponentSource,
) MediaComponentProvenanceInput {
	streamID := assemblyInputStreamID(input)
	return MediaComponentProvenanceInput{
		MediaItemID:         run.MediaItemID,
		ComponentType:       input.StreamType,
		ComponentKey:        assemblyComponentKey(input),
		ReleaseGroup:        stringValue(source.ReleaseGroup),
		ReleaseName:         firstStringValue(source.ReleaseName, source.ReleaseTitle, &source.SourceFilePath),
		ReleaseID:           source.ReleaseID,
		SourceProvider:      source.ReleaseName,
		SourceFilePath:      &source.SourceFilePath,
		RetainedSourceID:    &source.ID,
		SourceStreamID:      streamID,
		TransformationChain: assemblyTransformationChain(run, input, source),
	}
}

func containerProvenanceInput(run MediaComponentAssemblyRun) MediaComponentProvenanceInput {
	return MediaComponentProvenanceInput{
		MediaItemID:    run.MediaItemID,
		ComponentType:  "container",
		ComponentKey:   run.ID.String(),
		ReleaseName:    run.OutputPath,
		SourceFilePath: &run.OutputPath,
		TransformationChain: []map[string]any{{
			"kind":       "componentAssembly",
			"runId":      run.ID.String(),
			"outputPath": run.OutputPath,
		}},
	}
}

func assemblyComponentKey(input MediaComponentAssemblyInput) string {
	if input.ArtifactID != nil {
		return input.ArtifactID.String()
	}
	if input.SourceID != nil {
		return input.SourceID.String()
	}
	return fmt.Sprintf("%s:%s", input.StreamType, input.InputPath)
}

func assemblyInputStreamID(input MediaComponentAssemblyInput) *int32 {
	value, ok := input.Provenance["streamId"].(int32)
	if ok {
		return &value
	}
	floatValue, ok := input.Provenance["streamId"].(float64)
	if ok {
		value := int32(floatValue)
		return &value
	}
	return nil
}

func assemblyTransformationChain(
	run MediaComponentAssemblyRun,
	input MediaComponentAssemblyInput,
	source MediaComponentSource,
) []map[string]any {
	return []map[string]any{
		input.Provenance,
		{
			"kind":       "componentAssembly",
			"runId":      run.ID.String(),
			"sourceId":   source.ID.String(),
			"inputPath":  input.InputPath,
			"outputPath": run.OutputPath,
		},
	}
}

func firstStringValue(values ...*string) string {
	for _, value := range values {
		if value != nil && *value != "" {
			return *value
		}
	}
	return ""
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
