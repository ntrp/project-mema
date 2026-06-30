package metadata

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCredentialsRequired = errors.New("metadata provider credentials are required")
	ErrUnsupportedProvider = errors.New("unsupported metadata provider")
)

type ProviderHTTPError struct {
	StatusCode int
}

func (e ProviderHTTPError) Error() string {
	return fmt.Sprintf("metadata provider returned HTTP %d", e.StatusCode)
}

func (s *Service) Search(ctx context.Context, config Config, request SearchRequest) ([]SearchResult, error) {
	switch config.Type {
	case "tmdb":
		if request.MediaType != "movie" && request.MediaType != "series" {
			return nil, nil
		}
		return s.searchTMDB(ctx, config, request)
	case "tvdb":
		return s.searchTVDB(ctx, config, request)
	default:
		return nil, ErrUnsupportedProvider
	}
}

func (s *Service) Discover(ctx context.Context, config Config, request DiscoverRequest) ([]SearchResult, error) {
	if config.Type != "tmdb" {
		return nil, ErrUnsupportedProvider
	}
	return s.discoverTMDB(ctx, config, request)
}

func (s *Service) Details(ctx context.Context, config Config, request DetailsRequest) (Details, error) {
	if config.Type != "tmdb" {
		return Details{}, ErrUnsupportedProvider
	}
	return s.detailsTMDB(ctx, config, request)
}

func (s *Service) Test(ctx context.Context, config Config) TestResult {
	start := time.Now()
	results, err := s.Search(ctx, config, SearchRequest{Query: "test", MediaType: "movie"})
	latency := time.Since(start)
	if err != nil {
		return TestResult{
			Success: false,
			Message: err.Error(),
			Latency: latency,
			Details: map[string]interface{}{"provider": config.Type},
		}
	}
	return TestResult{
		Success: true,
		Message: "Metadata provider connection OK",
		Latency: latency,
		Details: map[string]interface{}{
			"provider": config.Type,
			"results":  len(results),
		},
	}
}
