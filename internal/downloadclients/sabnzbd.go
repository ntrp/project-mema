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

func (s *Service) addSABnzbd(ctx context.Context, config Config, request AddRequest) AddResult {
	endpoint, err := endpointWithPath(config.BaseURL, "/api", func(value string) bool {
		return strings.HasSuffix(strings.TrimRight(value, "/"), "/api")
	})
	if err != nil {
		return addFailedResult("Invalid SABnzbd URL", "error", err.Error())
	}
	category := stringValue(config.Category)
	if request.Category != nil && *request.Category != "" {
		category = *request.Category
	}
	endpoint, err = endpointWithQuery(endpoint, map[string]string{
		"mode":   "addurl",
		"name":   request.URL,
		"cat":    category,
		"output": "json",
		"apikey": stringValue(config.APIKey),
	})
	if err != nil {
		return addFailedResult("Invalid SABnzbd URL", "error", err.Error())
	}

	req, err := get(ctx, endpoint)
	if err != nil {
		return addFailedResult("Invalid SABnzbd request", "error", err.Error())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return addFailedResult("Connection failed", "error", err.Error())
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return addFailedResult("Unexpected response status", "statusCode", resp.StatusCode)
	}

	var payload struct {
		Status bool     `json:"status"`
		NZOIDs []string `json:"nzo_ids"`
		Error  string   `json:"error"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return addFailedResult("SABnzbd response could not be parsed", "error", err.Error())
	}
	if !payload.Status {
		if payload.Error == "" {
			payload.Error = "SABnzbd rejected the download"
		}
		return addFailedResult(payload.Error)
	}
	downloadID := ""
	if len(payload.NZOIDs) > 0 {
		downloadID = payload.NZOIDs[0]
	}
	return addSuccessResult("SABnzbd download queued", downloadID, "nzoIds", payload.NZOIDs)
}
