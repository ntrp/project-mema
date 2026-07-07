package subtitles

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	Name          string
	Type          string
	BaseURL       string
	Username      *string
	Password      *string
	APIKey        *string
	MockSubtitles []MockSubtitle
}

type MockSubtitle struct {
	Title      string
	LanguageID string
	Format     string
}

type TestResult struct {
	Success bool
	Message string
	Latency time.Duration
	Details map[string]any
}

type SearchRequest struct {
	MediaType     string
	Title         string
	LanguageID    string
	Year          *int32
	SeasonNumber  *int32
	EpisodeNumber *int32
	FilePath      string
}

type Candidate struct {
	ProviderName  string
	LanguageID    string
	FileID        int64
	Format        string
	ReleaseName   string
	DownloadCount int
	SourceURL     string
	SourceRef     string
}

type Download struct {
	Content []byte
	URL     string
}

type Service struct {
	client *http.Client
}

func NewService(client *http.Client) *Service {
	if client == nil {
		client = http.DefaultClient
	}
	return &Service{client: client}
}

func (s *Service) Test(ctx context.Context, config Config) TestResult {
	start := time.Now()
	err := s.testProvider(ctx, config)
	latency := time.Since(start)
	if err != nil {
		return TestResult{
			Success: false,
			Message: err.Error(),
			Latency: latency,
			Details: map[string]any{"provider": config.Type},
		}
	}
	return TestResult{
		Success: true,
		Message: "Subtitle provider connection OK",
		Latency: latency,
		Details: map[string]any{"provider": config.Type},
	}
}

func (s *Service) Search(
	ctx context.Context,
	config Config,
	request SearchRequest,
) ([]Candidate, error) {
	switch config.Type {
	case "opensubtitles":
		return s.searchOpenSubtitles(ctx, config, request)
	case "mock":
		return s.searchMock(config, request), nil
	default:
		return nil, ErrUnsupportedProvider
	}
}

func (s *Service) Download(ctx context.Context, config Config, candidate Candidate) (Download, error) {
	switch config.Type {
	case "opensubtitles":
		return s.downloadOpenSubtitles(ctx, config, candidate)
	case "mock":
		return s.downloadMock(candidate), nil
	default:
		return Download{}, ErrUnsupportedProvider
	}
}

func (s *Service) testProvider(ctx context.Context, config Config) error {
	switch config.Type {
	case "opensubtitles":
		return s.testOpenSubtitles(ctx, config)
	case "mock":
		return nil
	default:
		return ErrUnsupportedProvider
	}
}
