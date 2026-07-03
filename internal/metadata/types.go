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
	Page      int
}

type DetailsRequest struct {
	MediaType  string
	ExternalID string
}

type SearchResult struct {
	Title            string   `json:"title"`
	Type             string   `json:"type"`
	Year             *int32   `json:"year,omitempty"`
	ExternalProvider string   `json:"externalProvider"`
	ExternalID       string   `json:"externalId"`
	Overview         *string  `json:"overview,omitempty"`
	PosterPath       *string  `json:"posterPath,omitempty"`
	Popularity       *float64 `json:"popularity,omitempty"`
	ReleaseDate      *string  `json:"releaseDate,omitempty"`
	RuntimeMinutes   *int32   `json:"runtimeMinutes,omitempty"`
	VoteAverage      *float64 `json:"voteAverage,omitempty"`
	VoteCount        *int32   `json:"voteCount,omitempty"`
	OriginalLanguage *string  `json:"originalLanguage,omitempty"`
	ContentRating    *string  `json:"contentRating,omitempty"`
	Genres           []string `json:"genres,omitempty"`
	Keywords         []string `json:"keywords,omitempty"`
	Studios          []string `json:"studios,omitempty"`
	BackdropPath     *string  `json:"backdropPath,omitempty"`
}

type PersonSearchResult struct {
	Name             string   `json:"name"`
	ExternalProvider string   `json:"externalProvider"`
	ExternalID       string   `json:"externalId"`
	ProfilePath      *string  `json:"profilePath,omitempty"`
	Popularity       *float64 `json:"popularity,omitempty"`
	KnownFor         []string `json:"knownFor,omitempty"`
}

type DiscoverMovieRequest struct {
	Sort              string
	Page              int
	ReleaseDateFrom   *string
	ReleaseDateTo     *string
	Studios           []string
	Genres            []string
	Keywords          []string
	WithoutGenres     []string
	WithoutKeywords   []string
	OriginalLanguages []string
	ContentRatings    []string
	RuntimeMin        *int32
	RuntimeMax        *int32
	ScoreMin          *float64
	ScoreMax          *float64
	MinVoteCount      *int32
}

type DiscoverSeriesRequest struct {
	Sort              string
	Page              int
	ReleaseDateFrom   *string
	ReleaseDateTo     *string
	Studios           []string
	Genres            []string
	Keywords          []string
	WithoutGenres     []string
	WithoutKeywords   []string
	OriginalLanguages []string
	ContentRatings    []string
	Status            []string
	RuntimeMin        *int32
	RuntimeMax        *int32
	ScoreMin          *float64
	ScoreMax          *float64
	MinVoteCount      *int32
}

type FacetOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
	TrailerURL       *string
	Status           *string
	OriginalLanguage *string
	ReleaseDate      *string
	FirstAirDate     *string
	RuntimeMinutes   *int32
	SeasonCount      *int32
	EpisodeCount     *int32
	VoteAverage      *float64
	Genres           []string
	Keywords         []string
	Facts            []Fact
	Seasons          []Season
	Cast             []Person
	Crew             []Person
	Recommendations  []SearchResult
	Similar          []SearchResult
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
	ExternalProvider *string
	ExternalID       *string
	Name             string
	Role             *string
	ProfilePath      *string
}

type PersonDetails struct {
	ID           string
	Name         string
	Biography    *string
	Birthday     *string
	Deathday     *string
	PlaceOfBirth *string
	ProfilePath  *string
	AlsoKnownAs  []string
	Appearances  []PersonAppearance
}

type PersonAppearance struct {
	Title            string
	Type             string
	Year             *int32
	ExternalProvider string
	ExternalID       string
	Overview         *string
	PosterPath       *string
	BackdropPath     *string
	Role             *string
	ReleaseDate      *string
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
