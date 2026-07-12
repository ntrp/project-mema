package subtitles

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"media-manager/internal/subtitles/catalog"
)

var (
	ErrCredentialsRequired = errors.New("subtitle provider API key is required")
	ErrUnsupportedProvider = errors.New("unsupported subtitle provider")
	ErrCatalogOnlyProvider = errors.New("subtitle provider is catalog-only; runtime support is not implemented yet")
)

type providerAdapter interface {
	Test(context.Context, *Service, Config) error
	Search(context.Context, *Service, Config, SearchRequest) ([]Candidate, error)
	Download(context.Context, *Service, Config, Candidate) (Download, error)
}

var providerRegistry = map[string]providerAdapter{
	"mock":             mockAdapter{},
	"opensubtitles":    openSubtitlesAdapter{providerKey: "opensubtitlescom"},
	"opensubtitlescom": openSubtitlesAdapter{providerKey: "opensubtitlescom"},
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
		Details: map[string]any{"provider": canonicalProviderKey(config.Type)},
	}
}

func (s *Service) Search(
	ctx context.Context,
	config Config,
	request SearchRequest,
) ([]Candidate, error) {
	adapter, err := adapterFor(config.Type)
	if err != nil {
		return nil, err
	}
	return adapter.Search(ctx, s, config, request)
}

func (s *Service) Download(ctx context.Context, config Config, candidate Candidate) (Download, error) {
	adapter, err := adapterFor(config.Type)
	if err != nil {
		return Download{}, err
	}
	return adapter.Download(ctx, s, config, candidate)
}

func (s *Service) testProvider(ctx context.Context, config Config) error {
	adapter, err := adapterFor(config.Type)
	if err != nil {
		return err
	}
	return adapter.Test(ctx, s, config)
}

func RuntimeSupported(providerType string) bool {
	_, err := adapterFor(providerType)
	return err == nil
}

func UnsupportedRuntimeError(providerType string) error {
	_, err := adapterFor(providerType)
	if err == nil {
		return nil
	}
	return err
}

func adapterFor(providerType string) (providerAdapter, error) {
	key := canonicalProviderKey(providerType)
	if adapter, ok := providerRegistry[key]; ok {
		return adapter, nil
	}
	if entry, ok := catalog.Lookup(key); ok {
		if entry.RuntimeStatus != catalog.RuntimeSupported {
			return nil, fmt.Errorf("%w: %s", ErrCatalogOnlyProvider, entry.RuntimeMessage)
		}
	}
	return nil, ErrUnsupportedProvider
}

func canonicalProviderKey(providerType string) string {
	key := strings.ToLower(strings.TrimSpace(providerType))
	if key == "opensubtitles" {
		return "opensubtitlescom"
	}
	return key
}
