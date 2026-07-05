package storage

import (
	"time"

	"github.com/google/uuid"
)

type ImportAttempt struct {
	ID                     uuid.UUID
	ActivityID             uuid.UUID
	MediaItemID            uuid.UUID
	SourcePath             *string
	TargetPath             *string
	ImportMode             string
	Status                 string
	FailureStage           *string
	ErrorMessage           *string
	CreatedTargets         []string
	InsertedMediaFilePaths []string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type ImportAttemptInput struct {
	ActivityID             uuid.UUID
	MediaItemID            uuid.UUID
	SourcePath             *string
	TargetPath             *string
	ImportMode             string
	Status                 string
	FailureStage           *string
	ErrorMessage           *string
	CreatedTargets         []string
	InsertedMediaFilePaths []string
}
