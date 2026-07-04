package indexers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	cardigannengine "media-manager/internal/indexers/cardigann"
	"media-manager/internal/indexers/custom"
)

type Service struct {
	client    HTTPDoer
	cardigann *cardigannengine.Engine
	custom    *custom.Registry
}

func NewService(client HTTPDoer) *Service {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Service{
		client:    client,
		cardigann: cardigannengine.New(client),
		custom:    custom.NewRegistry(client),
	}
}

func (s *Service) Test(ctx context.Context, config Config) TestResult {
	startedAt := time.Now()
	result := s.test(ctx, config)
	result.Latency = time.Since(startedAt)
	if result.Details == nil {
		result.Details = map[string]interface{}{}
	}
	return result
}

func (s *Service) test(ctx context.Context, config Config) TestResult {
	if usesCardigannDefinition(config) {
		return s.cardigann.Test(ctx, config)
	}
	if indexer, ok := s.custom.EngineFor(config); ok {
		return indexer.Test(ctx, config)
	}
	return failedResult("Unsupported indexer protocol", "protocol", config.Protocol)
}

func (s *Service) Search(ctx context.Context, config Config, query string, mediaType string) ([]Release, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}
	if usesCardigannDefinition(config) {
		return s.cardigann.Search(ctx, config, query, mediaType)
	}
	if indexer, ok := s.custom.EngineFor(config); ok {
		return indexer.Search(ctx, config, query, mediaType)
	}
	return nil, fmt.Errorf("unsupported indexer protocol %q", config.Protocol)
}

func usesCardigannDefinition(config Config) bool {
	if config.DefinitionID == "" || strings.HasPrefix(config.DefinitionID, "generic-") {
		return false
	}
	return strings.EqualFold(config.Implementation, "Cardigann") || config.Implementation == ""
}

func failedResult(message string, pairs ...interface{}) TestResult {
	return TestResult{
		Success: false,
		Message: message,
		Details: details(pairs...),
	}
}

func successResult(message string, pairs ...interface{}) TestResult {
	return TestResult{
		Success: true,
		Message: message,
		Details: details(pairs...),
	}
}

func requestFailedResult(err error) TestResult {
	return failedResult("Connection failed", "error", err.Error())
}

func statusFailedResult(statusCode int) TestResult {
	return failedResult("Unexpected response status", "statusCode", statusCode)
}

func details(pairs ...interface{}) map[string]interface{} {
	values := map[string]interface{}{}
	for i := 0; i+1 < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok || key == "" {
			continue
		}
		values[key] = pairs[i+1]
	}
	return values
}
