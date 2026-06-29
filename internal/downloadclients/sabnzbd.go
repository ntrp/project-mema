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

func (s *Service) statusSABnzbd(ctx context.Context, config Config, request StatusRequest) StatusResult {
	queueResult := s.statusSABnzbdQueue(ctx, config, request)
	if queueResult.Success && queueResult.Found {
		return queueResult
	}
	historyResult := s.statusSABnzbdHistory(ctx, config, request)
	if historyResult.Success && historyResult.Found {
		return historyResult
	}
	if !queueResult.Success {
		return queueResult
	}
	if !historyResult.Success {
		return historyResult
	}
	return statusLookupNotFoundResult("SABnzbd download was not found")
}

func (s *Service) cancelSABnzbd(ctx context.Context, config Config, request CancelRequest) CancelResult {
	endpoint, err := sabnzbdAPIEndpoint(config, map[string]string{
		"mode":  "queue",
		"name":  "delete",
		"value": request.DownloadID,
	})
	if err != nil {
		return cancelFailedResult("Invalid SABnzbd URL", "error", err.Error())
	}
	req, err := get(ctx, endpoint)
	if err != nil {
		return cancelFailedResult("Invalid SABnzbd request", "error", err.Error())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return cancelFailedResult("Connection failed", "error", err.Error())
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return cancelFailedResult("Unexpected response status", "statusCode", resp.StatusCode)
	}
	var payload struct {
		Status bool   `json:"status"`
		Error  string `json:"error"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return cancelFailedResult("SABnzbd response could not be parsed", "error", err.Error())
	}
	if !payload.Status {
		if payload.Error == "" {
			payload.Error = "SABnzbd rejected the cancel request"
		}
		return cancelFailedResult(payload.Error)
	}
	return cancelSuccessResult("SABnzbd download cancelled")
}

func (s *Service) statusSABnzbdQueue(ctx context.Context, config Config, request StatusRequest) StatusResult {
	endpoint, err := sabnzbdAPIEndpoint(config, map[string]string{
		"mode":   "queue",
		"search": request.DownloadID,
	})
	if err != nil {
		return statusLookupFailedResult("Invalid SABnzbd URL", "error", err.Error())
	}
	req, err := get(ctx, endpoint)
	if err != nil {
		return statusLookupFailedResult("Invalid SABnzbd request", "error", err.Error())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return statusLookupFailedResult("Connection failed", "error", err.Error())
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return statusLookupFailedResult("Unexpected response status", "statusCode", resp.StatusCode)
	}

	var payload struct {
		Queue struct {
			Slots []struct {
				NZOID  string `json:"nzo_id"`
				Status string `json:"status"`
			} `json:"slots"`
		} `json:"queue"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return statusLookupFailedResult("SABnzbd response could not be parsed", "error", err.Error())
	}
	for _, slot := range payload.Queue.Slots {
		if slot.NZOID != request.DownloadID {
			continue
		}
		status := strings.ToLower(slot.Status)
		if strings.Contains(status, "pause") {
			return statusLookupResult("grabbed", "SABnzbd download is paused", "nzoId", slot.NZOID, "clientStatus", slot.Status)
		}
		return statusLookupResult("downloading", "SABnzbd download is active", "nzoId", slot.NZOID, "clientStatus", slot.Status)
	}
	return statusLookupNotFoundResult("SABnzbd download was not found in queue")
}

func (s *Service) statusSABnzbdHistory(ctx context.Context, config Config, request StatusRequest) StatusResult {
	endpoint, err := sabnzbdAPIEndpoint(config, map[string]string{
		"mode":   "history",
		"search": request.DownloadID,
	})
	if err != nil {
		return statusLookupFailedResult("Invalid SABnzbd URL", "error", err.Error())
	}
	req, err := get(ctx, endpoint)
	if err != nil {
		return statusLookupFailedResult("Invalid SABnzbd request", "error", err.Error())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return statusLookupFailedResult("Connection failed", "error", err.Error())
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return statusLookupFailedResult("Unexpected response status", "statusCode", resp.StatusCode)
	}

	var payload struct {
		History struct {
			Slots []struct {
				NZOID   string `json:"nzo_id"`
				Status  string `json:"status"`
				FailMsg string `json:"fail_message"`
				Storage string `json:"storage"`
			} `json:"slots"`
		} `json:"history"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return statusLookupFailedResult("SABnzbd response could not be parsed", "error", err.Error())
	}
	for _, slot := range payload.History.Slots {
		if slot.NZOID != request.DownloadID {
			continue
		}
		status := strings.ToLower(slot.Status)
		if strings.Contains(status, "fail") {
			message := strings.TrimSpace(slot.FailMsg)
			if message == "" {
				message = "SABnzbd download failed"
			}
			return statusLookupResult("failed", message, "nzoId", slot.NZOID, "clientStatus", slot.Status)
		}
		return statusLookupResultWithFiles(
			"completed",
			"SABnzbd download completed",
			sabnzbdStatusFiles(slot.Storage),
			"nzoId", slot.NZOID,
			"clientStatus", slot.Status,
		)
	}
	return statusLookupNotFoundResult("SABnzbd download was not found in history")
}

func sabnzbdStatusFiles(storage string) []StatusFile {
	storage = strings.TrimSpace(storage)
	if storage == "" {
		return nil
	}
	return []StatusFile{{Path: storage, Complete: true}}
}

func sabnzbdAPIEndpoint(config Config, values map[string]string) (string, error) {
	endpoint, err := endpointWithPath(config.BaseURL, "/api", func(value string) bool {
		return strings.HasSuffix(strings.TrimRight(value, "/"), "/api")
	})
	if err != nil {
		return "", err
	}
	values["output"] = "json"
	values["apikey"] = stringValue(config.APIKey)
	return endpointWithQuery(endpoint, values)
}
