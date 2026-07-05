package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaFileHistoryEntry struct {
	ID              uuid.UUID
	MediaItemID     *uuid.UUID
	FilePath        string
	SourcePath      *string
	DestinationPath *string
	Operation       string
	Status          string
	ActorType       string
	ActorID         *string
	JobID           *string
	Details         map[string]any
	FailureDetails  *string
	CreatedAt       time.Time
}

type MediaFileHistoryInput struct {
	MediaItemID     *uuid.UUID
	FilePath        string
	SourcePath      *string
	DestinationPath *string
	Operation       string
	Status          string
	ActorType       string
	ActorID         *string
	JobID           *string
	Details         map[string]any
	FailureDetails  *string
}
