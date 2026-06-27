package downloadclients

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	client HTTPDoer
}

func NewService(client HTTPDoer) *Service {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &Service{client: client}
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
	switch config.Type {
	case "transmission":
		return s.testTransmission(ctx, config)
	case "sabnzbd":
		return s.testSABnzbd(ctx, config)
	default:
		return failedResult("Unsupported download client type", "type", config.Type)
	}
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

func formatResultFailure(clientName, result string) TestResult {
	if result == "" {
		result = "empty"
	}
	return failedResult(fmt.Sprintf("%s rejected the test request", clientName), "result", result)
}
