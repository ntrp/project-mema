package storage

import (
	"time"

	"github.com/google/uuid"
)

type SubtitleRetentionMode string

const (
	SubtitleRetentionExternal SubtitleRetentionMode = "external"
	SubtitleRetentionMux      SubtitleRetentionMode = "mux"
	SubtitleRetentionIgnore   SubtitleRetentionMode = "ignore"
)

type MediaItemSubtitle struct {
	ID                 uuid.UUID
	MediaItemID        uuid.UUID
	SeasonID           *uuid.UUID
	EpisodeID          *uuid.UUID
	ProviderID         *uuid.UUID
	ProviderName       string
	LanguageID         string
	Format             string
	FilePath           string
	SourceURL          *string
	SourceRef          *string
	ReleaseName        *string
	ProviderSubtitleID *string
	Checksum           *string
	SizeBytes          *int64
	DownloadedAt       time.Time
	Selected           bool
	RetentionMode      SubtitleRetentionMode
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type MediaItemSubtitleInput struct {
	MediaItemID        uuid.UUID
	SeasonID           *uuid.UUID
	EpisodeID          *uuid.UUID
	ProviderID         *uuid.UUID
	ProviderName       string
	LanguageID         string
	Format             string
	FilePath           string
	SourceURL          *string
	SourceRef          *string
	ReleaseName        *string
	ProviderSubtitleID *string
	Checksum           *string
	SizeBytes          *int64
	DownloadedAt       time.Time
	Selected           *bool
	RetentionMode      SubtitleRetentionMode
}

type MediaItemSubtitleSelectionInput struct {
	Selected      bool
	RetentionMode SubtitleRetentionMode
}

type SubtitleAssemblyArtifact struct {
	ID                 uuid.UUID
	MediaItemID        uuid.UUID
	LanguageID         string
	Format             string
	FilePath           string
	RetentionMode      SubtitleRetentionMode
	ProviderName       string
	SourceURL          *string
	SourceRef          *string
	ProviderSubtitleID *string
	Checksum           *string
	SizeBytes          *int64
	DownloadedAt       time.Time
}
