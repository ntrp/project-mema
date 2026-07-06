package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaComponentCompatibilityDecision struct {
	ID                uuid.UUID
	MediaItemID       uuid.UUID
	BaseSourceID      uuid.UUID
	ComponentSourceID uuid.UUID
	ConfidenceState   string
	AutomationState   string
	ReviewState       string
	Reason            string
	RuntimeDeltaMs    *int32
	Evidence          map[string]any
	ReviewReason      *string
	ReviewedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type MediaComponentCompatibilityReviewInput struct {
	ReviewState string
	Reason      *string
}
