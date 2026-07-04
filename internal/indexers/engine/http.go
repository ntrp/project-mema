package engine

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"media-manager/internal/ratelimit"
)

type StatusError struct {
	StatusCode int
	RetryAfter time.Duration
}

func (e StatusError) Error() string {
	return fmt.Sprintf("unexpected response status %d", e.StatusCode)
}

func Get(ctx context.Context, endpoint string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
}

func CloseBody(body io.ReadCloser) {
	_, _ = io.Copy(io.Discard, io.LimitReader(body, 1024))
	_ = body.Close()
}

func ReadLimitedBody(body io.Reader) ([]byte, error) {
	return io.ReadAll(io.LimitReader(body, 1024*1024))
}

func DecodeLimitedXML(body io.Reader, target interface{}) error {
	return xml.NewDecoder(io.LimitReader(body, 1024*1024)).Decode(target)
}

func EndpointWithQuery(baseURL string, values map[string]string) (string, error) {
	parsed, err := ParseBaseURL(baseURL)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	for key, value := range values {
		if value == "" {
			continue
		}
		query.Set(key, value)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func ParseBaseURL(baseURL string) (*url.URL, error) {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("base URL must use http or https")
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("base URL must include a host")
	}
	return parsed, nil
}

func StringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func HTTPStatusError(resp *http.Response) error {
	return StatusError{StatusCode: resp.StatusCode, RetryAfter: ratelimit.DelayFromHeaders(resp.Header)}
}

func FailedResult(message string, pairs ...interface{}) TestResult {
	return TestResult{
		Success: false,
		Message: message,
		Details: Details(pairs...),
	}
}

func SuccessResult(message string, pairs ...interface{}) TestResult {
	return TestResult{
		Success: true,
		Message: message,
		Details: Details(pairs...),
	}
}

func RequestFailedResult(err error) TestResult {
	return FailedResult("Connection failed", "error", err.Error())
}

func StatusFailedResult(statusCode int) TestResult {
	return FailedResult("Unexpected response status", "statusCode", statusCode)
}

func Details(pairs ...interface{}) map[string]interface{} {
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
