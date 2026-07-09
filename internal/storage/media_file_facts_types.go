package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaFileFact struct {
	ID                  uuid.UUID
	MediaItemID         uuid.UUID
	SeasonID            *uuid.UUID
	EpisodeID           *uuid.UUID
	FilePath            string
	QualityID           *string
	ContainerFormat     *string
	ContainerFormatName *string
	ContainerBitrate    *int64
	DurationMs          *int64
	SizeBytes           *int64
	SourceKind          string
	ProbedAt            time.Time
	Tracks              []MediaFileTrackFact
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type MediaFileTrackFact struct {
	ID              uuid.UUID
	MediaFileFactID uuid.UUID
	MediaItemID     uuid.UUID
	FilePath        string
	StreamIndex     int32
	TrackType       string
	LanguageID      *string
	Codec           *string
	Channels        *string
	DurationMs      *int64
	BitrateKbps     *int32
	Width           *int32
	Height          *int32
	HDRFormat       *string
	PixelFormat     *string
	BitDepth        *int32
	Format          *string
	Title           *string
	Disposition     map[string]any
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type MediaFileFactInput struct {
	MediaItemID         uuid.UUID
	SeasonID            *uuid.UUID
	EpisodeID           *uuid.UUID
	FilePath            string
	QualityID           *string
	ContainerFormat     *string
	ContainerFormatName *string
	ContainerBitrate    *int64
	DurationMs          *int64
	SizeBytes           *int64
	SourceKind          string
	ProbedAt            time.Time
	Tracks              []MediaFileTrackFactInput
}

type MediaFileTrackFactInput struct {
	StreamIndex int32
	TrackType   string
	LanguageID  *string
	Codec       *string
	Channels    *string
	DurationMs  *int64
	BitrateKbps *int32
	Width       *int32
	Height      *int32
	HDRFormat   *string
	PixelFormat *string
	BitDepth    *int32
	Format      *string
	Title       *string
	Disposition map[string]any
}
