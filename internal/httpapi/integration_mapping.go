package httpapi

import (
	"time"

	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/metadata"
)

func downloadClientTestResponse(checkedAt time.Time, result downloadclients.TestResult) IntegrationTestResponse {
	return integrationTestResponse(checkedAt, result.Success, result.Message, result.Latency, result.Details)
}

func indexerTestResponse(checkedAt time.Time, result indexers.TestResult) IntegrationTestResponse {
	return integrationTestResponse(checkedAt, result.Success, result.Message, result.Latency, result.Details)
}

func metadataProviderTestResponse(checkedAt time.Time, result metadata.TestResult) IntegrationTestResponse {
	return integrationTestResponse(checkedAt, result.Success, result.Message, result.Latency, result.Details)
}

func integrationTestResponse(checkedAt time.Time, success bool, message string, latency time.Duration, details map[string]interface{}) IntegrationTestResponse {
	if details == nil {
		details = map[string]interface{}{}
	}
	return IntegrationTestResponse{
		Success:   success,
		Message:   message,
		LatencyMs: durationMilliseconds(latency),
		CheckedAt: checkedAt,
		Details:   details,
	}
}

func durationMilliseconds(duration time.Duration) int32 {
	const maxInt32 = int64(2147483647)
	milliseconds := duration.Milliseconds()
	if milliseconds > maxInt32 {
		return int32(maxInt32)
	}
	if milliseconds < 0 {
		return 0
	}
	return int32(milliseconds)
}
