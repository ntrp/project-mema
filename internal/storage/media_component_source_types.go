package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaComponentSource struct {
	ID              uuid.UUID
	MediaItemID     uuid.UUID
	SourceRole      string
	SourceFilePath  string
	RetainedPath    string
	ReleaseTitle    *string
	ReleaseGroup    *string
	ReleaseName     *string
	ReleaseID       *string
	SourceMetadata  *string
	StreamInventory string
	Checksum        *string
	SizeBytes       *int64
	RetentionState  string
	RetainedAt      time.Time
	ReleasedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Artifacts       []MediaComponentArtifact
	Compatibility   []MediaComponentCompatibilityDecision
}

type MediaComponentSourceInput struct {
	SourceRole      string
	SourceFilePath  string
	ReleaseTitle    *string
	ReleaseGroup    *string
	ReleaseName     *string
	ReleaseID       *string
	SourceMetadata  *string
	StreamInventory string
	Checksum        *string
}
