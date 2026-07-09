package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaItem struct {
	ID               uuid.UUID
	Type             string
	ContentKind      string
	Title            string
	Year             *int32
	Monitored        bool
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
	MediaMetadataSnapshot
	MonitorMode                  string
	SeriesType                   *string
	NumberingStrategy            *string
	MinimumAvailability          string
	QualityProfileID             *string
	QualityProfileName           *string
	FinalContainer               string
	Status                       string
	LibraryFolderID              *uuid.UUID
	LibraryFolderPath            *string
	MediaFolderPath              *string
	FilePaths                    []string
	MetadataFilePaths            []string
	VideoTarget                  MediaProfileVideoTarget
	AudioTargets                 []MediaProfileAudioTarget
	SubtitleTargets              []MediaProfileSubtitleTarget
	SubtitleMode                 string
	RemoveUnwantedAudio          bool
	RemoveUnwantedSubtitles      bool
	AllowSubtitleReleaseFallback bool
	ExternalSubtitles            []MediaItemSubtitle
	Sidecars                     []MediaItemSidecar
	FileFacts                    []MediaFileFact
	ComponentSources             []MediaComponentSource
	AssemblyRuns                 []MediaComponentAssemblyRun
	ComponentProvenance          []MediaComponentProvenance
	ProviderMappings             []MediaProviderMapping
	Aliases                      []MediaItemAlias
	EpisodeNumbering             []MediaEpisodeNumbering
	Tags                         []string
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
}

type MediaItemInput struct {
	Type             string
	ContentKind      string
	Title            string
	Year             *int32
	QualityFull      string
	Monitored        bool
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
	MediaMetadataSnapshot
	MonitorMode         string
	SeriesType          *string
	NumberingStrategy   *string
	MinimumAvailability string
	QualityProfileID    *string
	LibraryFolderID     *uuid.UUID
	ProviderMappings    []MediaProviderMappingInput
	Aliases             []MediaItemAliasInput
	EpisodeNumbering    []MediaEpisodeNumberingInput
	Tags                []string
}

type MediaItemOptionsInput struct {
	QualityProfileID     *string
	MinimumAvailability  *string
	LibraryFolderID      *uuid.UUID
	Monitored            *bool
	MonitorMode          *string
	Seasons              *[]MediaSeason
	MonitorSeasonName    *string
	MonitorEpisodeNumber *int32
	SeasonMonitored      *bool
	EpisodeMonitored     *bool
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
	ID                   uuid.UUID
	MediaItemID          uuid.UUID
	SeasonID             *uuid.UUID
	EpisodeID            *uuid.UUID
	IndexerID            *uuid.UUID
	IndexerName          string
	IndexerProtocol      string
	Title                string
	DownloadURL          string
	InfoURL              *string
	GUID                 *string
	SizeBytes            int64
	Seeders              *int32
	Peers                *int32
	PublishedAt          *time.Time
	SearchKind           string
	RequestedSeason      *int32
	RequestedEpisode     *int32
	Sources              []ReleaseCandidateSource
	CustomFormatScore    int32
	MatchedCustomFormats []ReleaseCandidateCustomFormatMatch
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ReleaseCandidateInput struct {
	MediaItemID          uuid.UUID
	SeasonID             *uuid.UUID
	EpisodeID            *uuid.UUID
	IndexerID            *uuid.UUID
	IndexerName          string
	IndexerProtocol      string
	Title                string
	DownloadURL          string
	InfoURL              *string
	GUID                 *string
	SizeBytes            int64
	Seeders              *int32
	Peers                *int32
	PublishedAt          *time.Time
	SearchKind           string
	RequestedSeason      *int32
	RequestedEpisode     *int32
	Sources              []ReleaseCandidateSource
	CustomFormatScore    int32
	MatchedCustomFormats []ReleaseCandidateCustomFormatMatch
}

type ReleaseCandidateCustomFormatMatch struct {
	Name  string `json:"name"`
	Score int32  `json:"score"`
}

type ReleaseCandidateSource struct {
	IndexerID       *uuid.UUID `json:"indexerId,omitempty"`
	IndexerName     string     `json:"indexerName"`
	IndexerProtocol string     `json:"indexerProtocol"`
	Title           string     `json:"title"`
	DownloadURL     string     `json:"downloadUrl"`
	InfoURL         *string    `json:"infoUrl,omitempty"`
	GUID            *string    `json:"guid,omitempty"`
}

type ReleaseBlocklistItem struct {
	ID                 uuid.UUID
	MediaItemID        uuid.UUID
	MediaTitle         string
	MediaType          string
	ReleaseTitle       string
	IndexerName        string
	IndexerProtocol    string
	DownloadClientName string
	DownloadURL        *string
	InfoURL            *string
	GUID               *string
	Reason             string
	Source             string
	Temporary          bool
	ExpiresAt          *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ReleaseBlocklistInput struct {
	MediaItemID        uuid.UUID
	ReleaseTitle       string
	IndexerName        string
	IndexerProtocol    string
	DownloadClientName string
	DownloadURL        string
	InfoURL            *string
	GUID               *string
	Reason             string
	Source             string
	Temporary          bool
	ExpiresAt          *time.Time
}

type ReleaseSearchSnapshot struct {
	Releases []ReleaseCandidate
	Errors   []string
}

type IndexerSearchSettings struct {
	CacheDurationMinutes         int32
	HistoryRetentionDays         int32
	AutomaticBlocklistExpiryDays int32
}

type IndexerSearchCacheStats struct {
	TotalEntries   int32
	ActiveEntries  int32
	ExpiredEntries int32
	IndexerCount   int32
}

type QueryHistoryStats struct {
	TotalEntries int32
	CacheHits    int32
	CacheMisses  int32
	Failures     int32
}

type IndexerSearchCacheEntry struct {
	IndexerID       uuid.UUID
	IndexerName     string
	IndexerProtocol string
	MediaType       string
	Query           string
	ResultCount     int32
	ExpiresAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Expired         bool
}

type IndexerSearchHistoryEntry struct {
	IndexerName     string
	IndexerProtocol string
	MediaType       string
	Query           string
	CacheHit        bool
	Success         bool
	ResultCount     int32
	Error           *string
	Response        string
	CreatedAt       time.Time
}

type MetadataSearchHistoryInput struct {
	ProviderID   uuid.UUID
	ProviderName string
	ProviderType string
	MediaType    string
	Query        string
	Year         *int32
	CacheHit     bool
	Success      bool
	ItemCount    int32
	Error        *string
	Response     any
}
