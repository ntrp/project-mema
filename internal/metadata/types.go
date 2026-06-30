package metadata

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TokenStore interface {
	UpdateMetadataProviderSessionToken(ctx context.Context, id uuid.UUID, token string, expiresAt time.Time) error
}

type Config struct {
	ID                    uuid.UUID
	Name                  string
	Type                  string
	BaseURL               string
	APIKey                *string
	PIN                   *string
	AccessToken           *string
	SessionToken          *string
	SessionTokenExpiresAt *time.Time
}

type SearchRequest struct {
	Query     string
	MediaType string
	Year      *int32
}

type DiscoverRequest struct {
	MediaType string
	Section   string
	Limit     int
}

type DetailsRequest struct {
	MediaType  string
	ExternalID string
}

type SearchResult struct {
	Title            string  `json:"title"`
	Type             string  `json:"type"`
	Year             *int32  `json:"year,omitempty"`
	ExternalProvider string  `json:"externalProvider"`
	ExternalID       string  `json:"externalId"`
	Overview         *string `json:"overview,omitempty"`
	PosterPath       *string `json:"posterPath,omitempty"`
}

type Details struct {
	Title            string
	Type             string
	Year             *int32
	ExternalProvider string
	ExternalID       string
	Overview         *string
	PosterPath       *string
	CollectionID     *string
	CollectionName   *string
	BackdropPath     *string
	Status           *string
	OriginalLanguage *string
	ReleaseDate      *string
	FirstAirDate     *string
	RuntimeMinutes   *int32
	SeasonCount      *int32
	EpisodeCount     *int32
	VoteAverage      *float64
	Genres           []string
	Facts            []Fact
	Seasons          []Season
	Cast             []Person
}

type Collection struct {
	ID           string
	Name         string
	Overview     *string
	PosterPath   *string
	BackdropPath *string
	Parts        []SearchResult
}

type Fact struct {
	Label string
	Value string
}

type Season struct {
	Name         string
	EpisodeCount *int32
	AirDate      *string
	PosterPath   *string
	Episodes     []Episode
}

type Episode struct {
	Name          string
	EpisodeNumber int32
	Overview      *string
	AirDate       *string
	StillPath     *string
}

type Person struct {
	Name        string
	Role        *string
	ProfilePath *string
}

type TestResult struct {
	Success bool
	Message string
	Details map[string]interface{}
	Latency time.Duration
}

type Service struct {
	httpClient *http.Client
	tokenStore TokenStore
	mu         sync.Mutex
	lastByID   map[string]time.Time
}

func NewService(httpClient *http.Client, tokenStore TokenStore) *Service {
	return &Service{
		httpClient: httpClient,
		tokenStore: tokenStore,
		lastByID:   map[string]time.Time{},
	}
}
