package storage

import (
	"time"

	"github.com/google/uuid"
)

type MediaSeriesSeason struct {
	ID               uuid.UUID
	MediaItemID      uuid.UUID
	ExternalProvider *string
	ExternalID       *string
	SeasonNumber     int32
	Name             string
	Overview         *string
	AirDate          *string
	PosterPath       *string
	EpisodeCount     *int32
	Monitored        bool
	Source           map[string]any
	Episodes         []MediaSeriesEpisode
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type MediaSeriesEpisode struct {
	ID               uuid.UUID
	SeasonID         uuid.UUID
	MediaItemID      uuid.UUID
	ExternalProvider *string
	ExternalID       *string
	SeasonNumber     int32
	EpisodeNumber    int32
	Name             string
	Overview         *string
	AirDate          *string
	StillPath        *string
	RuntimeMinutes   *int32
	Monitored        bool
	Source           map[string]any
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type MediaSeriesSeasonInput struct {
	ExternalProvider *string
	ExternalID       *string
	SeasonNumber     int32
	Name             string
	Overview         *string
	AirDate          *string
	PosterPath       *string
	EpisodeCount     *int32
	Monitored        bool
	Source           map[string]any
	Episodes         []MediaSeriesEpisodeInput
}

type MediaSeriesEpisodeInput struct {
	ExternalProvider *string
	ExternalID       *string
	EpisodeNumber    int32
	Name             string
	Overview         *string
	AirDate          *string
	StillPath        *string
	RuntimeMinutes   *int32
	Monitored        bool
	Source           map[string]any
}
