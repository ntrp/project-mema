package storage

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) EvaluateMediaComponentCompatibility(
	ctx context.Context,
	mediaItemID uuid.UUID,
	componentSourceID uuid.UUID,
	baseSourceID uuid.UUID,
) (MediaComponentCompatibilityDecision, error) {
	base, component, err := s.compatibilitySources(ctx, mediaItemID, componentSourceID, baseSourceID)
	if err != nil {
		return MediaComponentCompatibilityDecision{}, err
	}
	assessment := assessComponentCompatibility(base, component)
	row, err := storagegen.New(s.pool).UpsertMediaComponentCompatibility(
		ctx,
		storagegen.UpsertMediaComponentCompatibilityParams{
			ID:                uuid.New(),
			MediaItemID:       mediaItemID,
			BaseSourceID:      baseSourceID,
			ComponentSourceID: componentSourceID,
			ConfidenceState:   assessment.confidence,
			AutomationState:   assessment.automation,
			ReviewState:       assessment.review,
			Reason:            assessment.reason,
			RuntimeDeltaMs:    int4Value(assessment.runtimeDeltaMs),
			Evidence:          jsonObject(assessment.evidence),
		},
	)
	return mediaComponentCompatibilityRow(row, err)
}

func (s *SettingsStore) ReviewMediaComponentCompatibility(
	ctx context.Context,
	mediaItemID uuid.UUID,
	componentSourceID uuid.UUID,
	decisionID uuid.UUID,
	input MediaComponentCompatibilityReviewInput,
) (MediaComponentCompatibilityDecision, error) {
	reviewState := strings.TrimSpace(input.ReviewState)
	automationState := "blocked"
	switch reviewState {
	case "approved":
		automationState = "allowed"
	case "rejected":
	default:
		return MediaComponentCompatibilityDecision{}, ErrInvalidInput
	}
	row, err := storagegen.New(s.pool).ReviewMediaComponentCompatibility(
		ctx,
		storagegen.ReviewMediaComponentCompatibilityParams{
			ID:                decisionID,
			MediaItemID:       mediaItemID,
			ComponentSourceID: componentSourceID,
			ReviewState:       reviewState,
			AutomationState:   automationState,
			ReviewReason:      textValue(input.Reason),
		},
	)
	return mediaComponentCompatibilityRow(row, err)
}

func listMediaComponentCompatibilityForSource(
	ctx context.Context,
	q storagegen.DBTX,
	sourceID uuid.UUID,
) ([]MediaComponentCompatibilityDecision, error) {
	rows, err := storagegen.New(q).ListMediaComponentCompatibilityForSource(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	decisions := make([]MediaComponentCompatibilityDecision, 0, len(rows))
	for _, row := range rows {
		decisions = append(decisions, mediaComponentCompatibilityFromRow(row))
	}
	return decisions, nil
}

func (s *SettingsStore) compatibilitySources(
	ctx context.Context,
	mediaItemID uuid.UUID,
	componentSourceID uuid.UUID,
	baseSourceID uuid.UUID,
) (MediaComponentSource, MediaComponentSource, error) {
	base, err := s.GetMediaComponentSource(ctx, mediaItemID, baseSourceID)
	if err != nil {
		return MediaComponentSource{}, MediaComponentSource{}, err
	}
	component, err := s.GetMediaComponentSource(ctx, mediaItemID, componentSourceID)
	if err != nil {
		return MediaComponentSource{}, MediaComponentSource{}, err
	}
	if base.ID == component.ID || base.SourceRole != "baseVideo" {
		return MediaComponentSource{}, MediaComponentSource{}, ErrInvalidInput
	}
	return base, component, nil
}

type componentCompatibilityAssessment struct {
	confidence     string
	automation     string
	review         string
	reason         string
	runtimeDeltaMs *int32
	evidence       map[string]any
}

func assessComponentCompatibility(base MediaComponentSource, component MediaComponentSource) componentCompatibilityAssessment {
	baseMeta := componentCompatibilityMetadata(base)
	componentMeta := componentCompatibilityMetadata(component)
	delta := runtimeDeltaMs(baseMeta.runtimeMs, componentMeta.runtimeMs)
	cutMismatch := baseMeta.cutHint != "" && componentMeta.cutHint != "" && baseMeta.cutHint != componentMeta.cutHint
	assessment := componentCompatibilityAssessment{
		confidence:     "uncertain",
		automation:     "blocked",
		review:         "pending",
		reason:         "Runtime metadata is incomplete or requires review",
		runtimeDeltaMs: delta,
		evidence: map[string]any{
			"baseRuntimeMs":      baseMeta.runtimeMs,
			"componentRuntimeMs": componentMeta.runtimeMs,
			"baseCutHint":        baseMeta.cutHint,
			"componentCutHint":   componentMeta.cutHint,
		},
	}
	switch {
	case cutMismatch:
		assessment.confidence = "incompatible"
		assessment.reason = "Edition or cut hints do not match"
	case delta == nil:
	case *delta <= 1000:
		assessment.confidence = "exact"
		assessment.automation = "allowed"
		assessment.review = "notRequired"
		assessment.reason = "Runtime matches within 1 second"
	case *delta <= 3000:
		assessment.confidence = "likely"
		assessment.automation = "allowed"
		assessment.review = "notRequired"
		assessment.reason = "Runtime matches within 3 seconds"
	case *delta > 10000:
		assessment.confidence = "incompatible"
		assessment.reason = "Runtime differs by more than 10 seconds"
	default:
		assessment.reason = "Runtime differs enough to require review"
	}
	return assessment
}

type componentCompatibilityInfo struct {
	runtimeMs *int32
	cutHint   string
}

func componentCompatibilityMetadata(source MediaComponentSource) componentCompatibilityInfo {
	metadata := componentCompatibilityInfo{}
	for _, payload := range []string{source.StreamInventory, stringPtrValue(source.SourceMetadata)} {
		if metadata.runtimeMs == nil {
			metadata.runtimeMs = durationMsFromJSON(payload)
		}
		if metadata.cutHint == "" {
			metadata.cutHint = cutHint(payload)
		}
	}
	if metadata.cutHint == "" {
		metadata.cutHint = cutHint(stringPtrValue(source.ReleaseTitle))
	}
	return metadata
}

func durationMsFromJSON(payload string) *int32 {
	payload = strings.TrimSpace(payload)
	if payload == "" || (!strings.HasPrefix(payload, "{") && !strings.HasPrefix(payload, "[")) {
		return nil
	}
	var value any
	if err := json.Unmarshal([]byte(payload), &value); err != nil {
		return nil
	}
	return durationMsFromValue(value)
}

func durationMsFromValue(value any) *int32 {
	switch typed := value.(type) {
	case map[string]any:
		for _, key := range []string{"durationMs", "runtimeMs", "duration_ms"} {
			if duration := durationNumberMs(typed[key], 1); duration != nil {
				return duration
			}
		}
		if format, ok := typed["format"].(map[string]any); ok {
			if duration := durationNumberMs(format["duration"], 1000); duration != nil {
				return duration
			}
		}
		if streams, ok := typed["streams"].([]any); ok {
			return durationMsFromValue(streams)
		}
	case []any:
		for _, item := range typed {
			if duration := durationMsFromValue(item); duration != nil {
				return duration
			}
		}
	}
	return nil
}

func durationNumberMs(value any, multiplier float64) *int32 {
	var number float64
	switch typed := value.(type) {
	case float64:
		number = typed
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return nil
		}
		number = parsed
	default:
		return nil
	}
	ms := int32(math.Round(number * multiplier))
	return &ms
}

func runtimeDeltaMs(left *int32, right *int32) *int32 {
	if left == nil || right == nil {
		return nil
	}
	delta := *left - *right
	if delta < 0 {
		delta = -delta
	}
	return &delta
}

func cutHint(value string) string {
	normalized := strings.ToLower(value)
	for _, hint := range []string{"director", "extended", "uncut", "theatrical", "remaster"} {
		if strings.Contains(normalized, hint) {
			return hint
		}
	}
	return ""
}

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func mediaComponentCompatibilityRow(
	row storagegen.AppMediaComponentCompatibilityDecision,
	err error,
) (MediaComponentCompatibilityDecision, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaComponentCompatibilityDecision{}, ErrNotFound
	}
	if err != nil {
		return MediaComponentCompatibilityDecision{}, err
	}
	return mediaComponentCompatibilityFromRow(row), nil
}

func mediaComponentCompatibilityFromRow(
	row storagegen.AppMediaComponentCompatibilityDecision,
) MediaComponentCompatibilityDecision {
	return MediaComponentCompatibilityDecision{
		ID:                row.ID,
		MediaItemID:       row.MediaItemID,
		BaseSourceID:      row.BaseSourceID,
		ComponentSourceID: row.ComponentSourceID,
		ConfidenceState:   row.ConfidenceState,
		AutomationState:   row.AutomationState,
		ReviewState:       row.ReviewState,
		Reason:            row.Reason,
		RuntimeDeltaMs:    int4Ptr(row.RuntimeDeltaMs),
		Evidence:          jsonMap(row.Evidence),
		ReviewReason:      textPtr(row.ReviewReason),
		ReviewedAt:        row.ReviewedAt,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}
