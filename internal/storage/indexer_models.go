package storage

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Indexer struct {
	ID                   uuid.UUID
	DefinitionID         string
	Name                 string
	Implementation       string
	ImplementationName   string
	Protocol             string
	Privacy              string
	Language             string
	Encoding             *string
	Description          *string
	IndexerURLs          []string
	LegacyURLs           []string
	BaseURL              string
	APIKey               *string
	Categories           []int32
	Fields               json.RawMessage
	Capabilities         json.RawMessage
	Redirect             bool
	AppProfileID         string
	MinimumSeeders       *int32
	SeedRatio            *float64
	SeedTime             *int32
	PackSeedTime         *int32
	PreferMagnetURL      bool
	SupportsRSS          bool
	SupportsSearch       bool
	SupportsRedirect     bool
	SupportsPagination   bool
	Enabled              bool
	Priority             int32
	HealthStatus         string
	LastQueryAt          *time.Time
	LastSuccessAt        *time.Time
	LastFailureAt        *time.Time
	NextCheckAt          *time.Time
	LastStatusCode       *int32
	LastError            *string
	FailureCount         int32
	RSSMarkerPublishedAt *time.Time
	RSSMarkerGUID        *string
	RSSMarkerDownloadURL *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type RSSMarkerInput struct {
	PublishedAt *time.Time
	GUID        *string
	DownloadURL *string
}

type IndexerInput struct {
	DefinitionID       string
	Name               string
	Implementation     string
	ImplementationName string
	Protocol           string
	Privacy            string
	Language           string
	Encoding           *string
	Description        *string
	IndexerURLs        []string
	LegacyURLs         []string
	BaseURL            string
	APIKey             *string
	Categories         []int32
	Fields             json.RawMessage
	Capabilities       json.RawMessage
	Redirect           bool
	AppProfileID       string
	MinimumSeeders     *int32
	SeedRatio          *float64
	SeedTime           *int32
	PackSeedTime       *int32
	PreferMagnetURL    bool
	SupportsRSS        bool
	SupportsSearch     bool
	SupportsRedirect   bool
	SupportsPagination bool
	Enabled            bool
	Priority           int32
}

type IndexerBulkUpdateInput struct {
	IDs             []uuid.UUID
	Enabled         *bool
	AppProfileID    *string
	Priority        *int32
	MinimumSeeders  *int32
	SeedRatio       *float64
	SeedTime        *int32
	PackSeedTime    *int32
	PreferMagnetURL *bool
}
