package downloadclients

import (
	"context"
	"net/http"
	"strings"
)

func (s *Service) testSABnzbd(ctx context.Context, config Config) TestResult {
	endpoint, err := endpointWithPath(config.BaseURL, "/api", func(value string) bool {
		return strings.HasSuffix(strings.TrimRight(value, "/"), "/api")
	})
	if err != nil {
		return failedResult("Invalid SABnzbd URL", "error", err.Error())
	}
	endpoint, err = endpointWithQuery(endpoint, map[string]string{
		"mode":   "version",
		"output": "json",
		"apikey": stringValue(config.APIKey),
	})
	if err != nil {
		return failedResult("Invalid SABnzbd URL", "error", err.Error())
	}

	req, err := get(ctx, endpoint)
	if err != nil {
		return failedResult("Invalid SABnzbd request", "error", err.Error())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return requestFailedResult(err)
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return statusFailedResult(resp.StatusCode)
	}

	var payload struct {
		Version string `json:"version"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return failedResult("SABnzbd response could not be parsed", "error", err.Error())
	}
	if payload.Version == "" {
		return failedResult("SABnzbd version was not returned")
	}

	return successResult("SABnzbd connection OK", "version", payload.Version)
}
