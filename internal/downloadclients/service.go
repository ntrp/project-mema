package downloadclients

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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

func (s *Service) Add(ctx context.Context, config Config, request AddRequest) AddResult {
	if strings.TrimSpace(request.URL) == "" {
		return addFailedResult("Download URL is required")
	}
	switch config.Type {
	case "transmission":
		return s.addTransmission(ctx, config, request)
	case "sabnzbd":
		return s.addSABnzbd(ctx, config, request)
	default:
		return addFailedResult("Unsupported download client type", "type", config.Type)
	}
}

func (s *Service) Status(ctx context.Context, config Config, request StatusRequest) StatusResult {
	if strings.TrimSpace(request.DownloadID) == "" {
		return statusLookupFailedResult("Download ID is required")
	}
	switch config.Type {
	case "transmission":
		return s.statusTransmission(ctx, config, request)
	case "sabnzbd":
		return s.statusSABnzbd(ctx, config, request)
	default:
		return statusLookupFailedResult("Unsupported download client type", "type", config.Type)
	}
}

func (s *Service) Cancel(ctx context.Context, config Config, request CancelRequest) CancelResult {
	if strings.TrimSpace(request.DownloadID) == "" {
		return cancelFailedResult("Download ID is required")
	}
	switch config.Type {
	case "transmission":
		return s.cancelTransmission(ctx, config, request)
	case "sabnzbd":
		return s.cancelSABnzbd(ctx, config, request)
	default:
		return cancelFailedResult("Unsupported download client type", "type", config.Type)
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

func addFailedResult(message string, pairs ...interface{}) AddResult {
	return AddResult{
		Success: false,
		Message: message,
		Details: details(pairs...),
	}
}

func addSuccessResult(message string, downloadID string, pairs ...interface{}) AddResult {
	return AddResult{
		Success:    true,
		Message:    message,
		DownloadID: downloadID,
		Details:    details(pairs...),
	}
}

func cancelFailedResult(message string, pairs ...interface{}) CancelResult {
	return CancelResult{
		Success: false,
		Message: message,
		Details: details(pairs...),
	}
}

func cancelSuccessResult(message string, pairs ...interface{}) CancelResult {
	return CancelResult{
		Success: true,
		Message: message,
		Details: details(pairs...),
	}
}

func statusLookupFailedResult(message string, pairs ...interface{}) StatusResult {
	return StatusResult{
		Success: false,
		Message: message,
		Details: details(pairs...),
	}
}

func statusLookupNotFoundResult(message string, pairs ...interface{}) StatusResult {
	return StatusResult{
		Success: true,
		Found:   false,
		Message: message,
		Details: details(pairs...),
	}
}

func statusLookupResult(status string, message string, pairs ...interface{}) StatusResult {
	return StatusResult{
		Success: true,
		Found:   true,
		Status:  status,
		Message: message,
		Details: details(pairs...),
	}
}

func statusLookupResultWithFiles(status string, message string, files []StatusFile, pairs ...interface{}) StatusResult {
	result := statusLookupResult(status, message, pairs...)
	result.Files = files
	return result
}

func statusLookupResultWithProgressAndFiles(status string, progressPercent *int, message string, files []StatusFile, pairs ...interface{}) StatusResult {
	result := statusLookupResultWithFiles(status, message, files, pairs...)
	result.ProgressPercent = progressPercent
	return result
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
