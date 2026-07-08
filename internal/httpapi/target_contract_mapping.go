package httpapi

import (
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/satisfaction"
	"media-manager/internal/targets"
)

func targetSatisfactionSummaryResponse(
	targetRows []targets.Target,
	candidateRows []targets.Candidate,
) TargetSatisfactionSummary {
	response := TargetSatisfactionSummary{
		Targets:    make([]TargetSatisfactionTarget, 0, len(targetRows)),
		Candidates: make([]TargetSatisfactionCandidate, 0, len(candidateRows)),
	}
	for _, target := range targetRows {
		response.Targets = append(response.Targets, targetSatisfactionTargetResponse(target))
	}
	for _, candidate := range candidateRows {
		response.Candidates = append(response.Candidates, targetSatisfactionCandidateResponse(candidate))
	}
	return response
}

func targetSatisfactionTargetResponse(target targets.Target) TargetSatisfactionTarget {
	mediaItemID, _ := parseOpenAPIUUID(target.MediaItemID)
	mediaFileID, hasMediaFileID := parseOpenAPIUUID(target.MediaFileID)
	languageID := optionalNonEmptyString(target.LanguageID)
	response := TargetSatisfactionTarget{
		Id:          target.ID,
		Type:        TargetSatisfactionType(target.Type),
		State:       TargetSatisfactionState(target.State),
		MediaItemId: mediaItemID,
		LanguageId:  languageID,
		Reasons:     append([]string(nil), target.Reasons...),
	}
	if hasMediaFileID {
		response.MediaFileId = &mediaFileID
	}
	response.RequiredOperation = targetOperationMetadataResponse(target.RequiredOperation)
	return response
}

func targetSatisfactionCandidateResponse(candidate targets.Candidate) TargetSatisfactionCandidate {
	languageID := optionalNonEmptyString(candidate.LanguageID)
	response := TargetSatisfactionCandidate{
		Id:          candidate.ID,
		Type:        TargetCandidateType(candidate.Type),
		VisualState: TargetCandidateVisualState(candidate.VisualState),
		TargetIds:   append([]string(nil), candidate.TargetIDs...),
		LanguageId:  languageID,
		Operation:   targetOperationMetadataResponse(candidate.Operation),
	}
	if len(candidate.UnwantedRules) > 0 {
		rules := append([]string(nil), candidate.UnwantedRules...)
		response.UnwantedRules = &rules
	}
	return response
}

func targetOperationMetadataResponse(operation *targets.Operation) *TargetOperationMetadata {
	if operation == nil {
		return nil
	}
	response := TargetOperationMetadata{
		Type:      TargetOperationType(operation.Type),
		Manual:    operation.Manual,
		Automatic: operation.Automatic,
		Reason:    operation.Reason,
	}
	if jobID, ok := parseOpenAPIUUID(operation.JobID); ok {
		response.JobId = &jobID
	}
	return &response
}

func mediaRollupSummaryResponse(state MediaRollupState, targetStates []targets.State, reasons []string) MediaRollupSummary {
	return MediaRollupSummary{
		State:        state,
		TargetCounts: targetStateCountsResponse(targetStates),
		Reasons:      append([]string(nil), reasons...),
	}
}

func mediaItemRollupSummary(status string) MediaRollupSummary {
	switch strings.TrimSpace(status) {
	case "downloading":
		return mediaRollupSummaryResponse(MediaRollupStateDownloading, nil, []string{"Media item has active download work."})
	case "downloaded":
		return mediaRollupSummaryResponse(MediaRollupStateDownloaded, nil, []string{"Media item has imported media."})
	default:
		return mediaRollupSummaryResponse(MediaRollupStateMissing, nil, []string{"Media item has no usable file."})
	}
}

func mediaFileRollupSummary(status MediaFileInfoStatus) MediaRollupSummary {
	if status == MediaFileInfoStatusAvailable {
		return mediaRollupSummaryResponse(MediaRollupStateDownloaded, nil, []string{"Media file is available."})
	}
	return mediaRollupSummaryResponse(MediaRollupStateMissing, nil, []string{"Media file is missing."})
}

func wantedRowResponse(row satisfaction.WantedRow) WantedRow {
	mediaItemID, _ := parseOpenAPIUUID(row.MediaItemID)
	response := WantedRow{
		Id:                row.ID,
		Kind:              WantedRowKind(row.Kind),
		MediaItemId:       mediaItemID,
		MediaTitle:        row.MediaTitle,
		MediaType:         MediaType(row.MediaType),
		SeasonNumber:      row.SeasonNumber,
		EpisodeNumber:     row.EpisodeNumber,
		FileLabel:         optionalNonEmptyString(row.FileLabel),
		FilePath:          optionalNonEmptyString(row.FilePath),
		LanguageId:        optionalNonEmptyString(row.LanguageID),
		RequiredOperation: targetOperationMetadataResponse(row.RequiredOperation),
		CurrentScore:      row.CurrentScore,
		TargetScore:       row.TargetScore,
	}
	if row.TargetType != "" {
		targetType := TargetSatisfactionType(row.TargetType)
		response.TargetType = &targetType
	}
	if row.TargetState != "" {
		targetState := TargetSatisfactionState(row.TargetState)
		response.TargetState = &targetState
	}
	return response
}

func targetStateCountsResponse(states []targets.State) TargetStateCounts {
	counts := TargetStateCounts{}
	for _, state := range states {
		switch state {
		case targets.StateMissing:
			counts.Missing++
		case targets.StatePartial:
			counts.Partial++
		case targets.StatePending:
			counts.Pending++
		case targets.StateSatisfied:
			counts.Satisfied++
		case targets.StateUpgradeable:
			counts.Upgradeable++
		case targets.StateBlocked:
			counts.Blocked++
		case targets.StateFailed:
			counts.Failed++
		}
	}
	return counts
}

func parseOpenAPIUUID(value string) (openapi_types.UUID, bool) {
	parsed, err := uuid.Parse(strings.TrimSpace(value))
	return openapi_types.UUID(parsed), err == nil
}

func optionalNonEmptyString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
