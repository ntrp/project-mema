package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaItem struct {
	ID               uuid.UUID
	Type             string
	Title            string
	Year             *int32
	Monitored        bool
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
	MediaMetadataSnapshot
	MonitorMode         string
	SeriesType          *string
	MinimumAvailability string
	QualityProfileID    *string
	QualityProfileName  *string
	Status              string
	LibraryFolderID     *uuid.UUID
	LibraryFolderPath   *string
	MediaFolderPath     *string
	FilePaths           []string
	MetadataFilePaths   []string
	Tags                []string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type MediaItemInput struct {
	Type             string
	Title            string
	Year             *int32
	Monitored        bool
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
	MediaMetadataSnapshot
	MonitorMode         string
	SeriesType          *string
	MinimumAvailability string
	QualityProfileID    *string
	LibraryFolderID     *uuid.UUID
	Tags                []string
}

type DownloadActivity struct {
	ID                 uuid.UUID
	MediaItemID        uuid.UUID
	MediaTitle         string
	MediaType          string
	MediaYear          *int32
	ReleaseTitle       string
	IndexerName        string
	DownloadClientName string
	DownloadID         *string
	DownloadURL        string
	Status             string
	ProgressPercent    *int
	Error              *string
	FailureType        *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type DownloadActivityInput struct {
	MediaItemID        uuid.UUID
	ReleaseTitle       string
	IndexerName        string
	DownloadClientName string
	DownloadID         *string
	DownloadURL        string
	Status             string
	Error              *string
	FailureType        *string
}

type ReleaseCandidate struct {
	ID          uuid.UUID
	MediaItemID uuid.UUID
	IndexerID   *uuid.UUID
	IndexerName string
	IndexerType string
	Title       string
	DownloadURL string
	InfoURL     *string
	GUID        *string
	SizeBytes   int64
	Seeders     *int32
	Peers       *int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReleaseCandidateInput struct {
	MediaItemID uuid.UUID
	IndexerID   *uuid.UUID
	IndexerName string
	IndexerType string
	Title       string
	DownloadURL string
	InfoURL     *string
	GUID        *string
	SizeBytes   int64
	Seeders     *int32
	Peers       *int32
}

type ReleaseSearchSnapshot struct {
	Releases []ReleaseCandidate
	Errors   []string
}
