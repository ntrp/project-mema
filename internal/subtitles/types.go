package subtitles

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	Name     string
	Type     string
	BaseURL  string
	Username *string
	Password *string
	APIKey   *string
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
	err := s.testOpenSubtitles(ctx, config)
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
	return s.searchOpenSubtitles(ctx, config, request)
}

func (s *Service) Download(ctx context.Context, config Config, candidate Candidate) (Download, error) {
	return s.downloadOpenSubtitles(ctx, config, candidate)
}
