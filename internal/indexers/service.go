package indexers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	client HTTPDoer
	loader *cardigannLoader
}

func NewService(client HTTPDoer) *Service {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Service{client: client, loader: newCardigannLoader(client)}
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
	if config.usesCardigannDefinition() {
		return s.testCardigann(ctx, config)
	}
	switch config.Protocol {
	case "torrent", "usenet":
		return s.testCaps(ctx, config)
	default:
		return failedResult("Unsupported indexer protocol", "protocol", config.Protocol)
	}
}

func (s *Service) Search(ctx context.Context, config Config, query string, mediaType string) ([]Release, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}
	if config.usesCardigannDefinition() {
		return s.searchCardigann(ctx, config, query, mediaType)
	}
	switch config.Protocol {
	case "torrent", "usenet":
		return s.searchTorznab(ctx, config, query, mediaType)
	default:
		return nil, fmt.Errorf("unsupported indexer protocol %q", config.Protocol)
	}
}

func (config Config) usesCardigannDefinition() bool {
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
