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
	provider, err := s.provider(config)
	if err != nil {
		return nil, err
	}
	return provider.Search(ctx, request)
}

func (s *Service) Discover(ctx context.Context, config Config, request DiscoverRequest) ([]SearchResult, error) {
	provider, err := s.provider(config)
	if err != nil {
		return nil, err
	}
	discoverProvider, ok := provider.(DiscoverProvider)
	if !ok {
		return nil, ErrUnsupportedProvider
	}
	return discoverProvider.Discover(ctx, request)
}

func (s *Service) Details(ctx context.Context, config Config, request DetailsRequest) (Details, error) {
	provider, err := s.provider(config)
	if err != nil {
		return Details{}, err
	}
	return provider.Details(ctx, request)
}

func (s *Service) provider(config Config) (Provider, error) {
	switch config.Type {
	case "tmdb":
		return tmdbProvider{service: s, config: config}, nil
	case "tvdb":
		return tvdbProvider{service: s, config: config}, nil
	default:
		return nil, ErrUnsupportedProvider
	}
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
