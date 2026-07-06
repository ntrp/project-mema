package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaProviderMapping struct {
	ID                 uuid.UUID
	MediaItemID        uuid.UUID
	SeasonID           *uuid.UUID
	EpisodeID          *uuid.UUID
	EntityType         string
	ProviderName       string
	ProviderEntityType string
	ExternalID         string
	Canonical          bool
	Confidence         *float64
	Source             map[string]any
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type MediaProviderMappingInput struct {
	SeasonID           *uuid.UUID
	EpisodeID          *uuid.UUID
	EntityType         string
	ProviderName       string
	ProviderEntityType string
	ExternalID         string
	Canonical          bool
	Confidence         *float64
	Source             map[string]any
}

type MediaItemAlias struct {
	ID                uuid.UUID
	MediaItemID       uuid.UUID
	Alias             string
	NormalizedAlias   string
	Language          *string
	Kind              string
	ProviderName      *string
	ProviderMappingID *uuid.UUID
	Source            map[string]any
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type MediaItemAliasInput struct {
	Alias             string
	Language          *string
	Kind              string
	ProviderName      *string
	ProviderMappingID *uuid.UUID
	Source            map[string]any
}

type MediaEpisodeNumbering struct {
	ID              uuid.UUID
	MediaItemID     uuid.UUID
	SeasonID        *uuid.UUID
	EpisodeID       uuid.UUID
	ProviderName    string
	NumberingScheme string
	SeasonNumber    *int32
	EpisodeNumber   *int32
	AbsoluteNumber  *int32
	Source          map[string]any
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type MediaEpisodeNumberingInput struct {
	SeasonID        *uuid.UUID
	EpisodeID       uuid.UUID
	ProviderName    string
	NumberingScheme string
	SeasonNumber    *int32
	EpisodeNumber   *int32
	AbsoluteNumber  *int32
	Source          map[string]any
}
