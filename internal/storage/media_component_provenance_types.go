package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaComponentProvenance struct {
	ID                  uuid.UUID
	MediaItemID         uuid.UUID
	ComponentType       string
	ComponentKey        string
	ReleaseGroup        string
	ReleaseName         string
	ReleaseID           *string
	SourceProvider      *string
	SourceFilePath      *string
	RetainedSourceID    *uuid.UUID
	SourceStreamID      *int32
	TransformationChain []map[string]any
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type MediaComponentProvenanceInput struct {
	MediaItemID         uuid.UUID
	ComponentType       string
	ComponentKey        string
	ReleaseGroup        string
	ReleaseName         string
	ReleaseID           *string
	SourceProvider      *string
	SourceFilePath      *string
	RetainedSourceID    *uuid.UUID
	SourceStreamID      *int32
	TransformationChain []map[string]any
}
