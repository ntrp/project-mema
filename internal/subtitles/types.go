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
