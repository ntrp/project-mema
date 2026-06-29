package downloadclients

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

func (s *Service) cancelTransmission(ctx context.Context, config Config, request CancelRequest) CancelResult {
	endpoint, err := endpointWithPath(config.BaseURL, "/transmission/rpc", func(value string) bool {
		return strings.HasSuffix(strings.TrimRight(value, "/"), "/transmission/rpc")
	})
	if err != nil {
		return cancelFailedResult("Invalid Transmission URL", "error", err.Error())
	}

	resp, err := s.doTransmissionCancelRequest(ctx, endpoint, config, request, "")
	if err != nil {
		return cancelFailedResult("Connection failed", "error", err.Error())
	}
	if resp.StatusCode == http.StatusConflict {
		sessionID := resp.Header.Get(transmissionSessionHeader)
		closeBody(resp.Body)
		if sessionID == "" {
			return cancelFailedResult("Transmission session id was not returned", "statusCode", resp.StatusCode)
		}
		resp, err = s.doTransmissionCancelRequest(ctx, endpoint, config, request, sessionID)
		if err != nil {
			return cancelFailedResult("Connection failed", "error", err.Error())
		}
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return cancelFailedResult("Unexpected response status", "statusCode", resp.StatusCode)
	}

	var payload struct {
		Result string `json:"result"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return cancelFailedResult("Transmission response could not be parsed", "error", err.Error())
	}
	if payload.Result != "success" {
		return cancelFailedResult("Transmission rejected the cancel request", "result", payload.Result)
	}
	return cancelSuccessResult("Transmission download cancelled")
}

func (s *Service) doTransmissionCancelRequest(ctx context.Context, endpoint string, config Config, request CancelRequest, sessionID string) (*http.Response, error) {
	id, err := parseTransmissionID(request.DownloadID)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"method": "torrent-remove",
		"arguments": map[string]interface{}{
			"ids":               []int{id},
			"delete-local-data": false,
		},
	}
	req, err := jsonPost(ctx, endpoint, body)
	if err != nil {
		return nil, err
	}
	if sessionID != "" {
		req.Header.Set(transmissionSessionHeader, sessionID)
	}
	addBasicAuth(req, config.Username, config.Password)
	return s.client.Do(req)
}

func parseTransmissionID(value string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(value))
}
