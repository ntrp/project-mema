package indexers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type StatusError struct {
	StatusCode int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("unexpected response status %d", e.StatusCode)
}

func get(ctx context.Context, endpoint string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
}

func closeBody(body io.ReadCloser) {
	_, _ = io.Copy(io.Discard, io.LimitReader(body, 1024))
	_ = body.Close()
}

func readLimitedBody(body io.Reader) ([]byte, error) {
	return io.ReadAll(io.LimitReader(body, 1024*1024))
}

func decodeLimitedXML(body io.Reader, target interface{}) error {
	return xml.NewDecoder(io.LimitReader(body, 1024*1024)).Decode(target)
}

func endpointWithQuery(baseURL string, values map[string]string) (string, error) {
	parsed, err := parseBaseURL(baseURL)
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

func parseBaseURL(baseURL string) (*url.URL, error) {
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

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
