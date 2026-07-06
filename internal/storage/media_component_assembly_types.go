package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaComponentAssemblyRun struct {
	ID           uuid.UUID
	MediaItemID  uuid.UUID
	BaseSourceID uuid.UUID
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
	Inputs       []MediaComponentAssemblyInput
}

type MediaComponentAssemblyInput struct {
	ID         uuid.UUID
	RunID      uuid.UUID
	SourceID   *uuid.UUID
	ArtifactID *uuid.UUID
	StreamType string
	InputPath  string
	Provenance map[string]any
	CreatedAt  time.Time
}

type MediaComponentAssemblyRunInput struct {
	BaseSourceID uuid.UUID
	ArtifactIDs  []uuid.UUID
	JobID        *string
}
