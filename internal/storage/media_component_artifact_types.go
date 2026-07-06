package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaComponentArtifact struct {
	ID           uuid.UUID
	MediaItemID  uuid.UUID
	SourceID     uuid.UUID
	StreamID     int32
	StreamType   string
	Language     *string
	OutputPath   string
	Status       string
	ToolName     string
	ToolSummary  string
	ErrorMessage *string
	JobID        *string
	SizeBytes    *int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CompletedAt  *time.Time
}

type MediaComponentArtifactInput struct {
	StreamID   int32
	StreamType string
	Language   *string
	JobID      *string
}
