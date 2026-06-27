package downloadclients

import (
	"context"
	"net/http"
	"strings"
)

const transmissionSessionHeader = "X-Transmission-Session-Id"

func (s *Service) testTransmission(ctx context.Context, config Config) TestResult {
	endpoint, err := endpointWithPath(config.BaseURL, "/transmission/rpc", func(value string) bool {
		return strings.HasSuffix(strings.TrimRight(value, "/"), "/transmission/rpc")
	})
	if err != nil {
		return failedResult("Invalid Transmission URL", "error", err.Error())
	}

	resp, err := s.doTransmissionRequest(ctx, endpoint, config, "")
	if err != nil {
		return requestFailedResult(err)
	}

	if resp.StatusCode == http.StatusConflict {
		sessionID := resp.Header.Get(transmissionSessionHeader)
		closeBody(resp.Body)
		if sessionID == "" {
			return failedResult("Transmission session id was not returned", "statusCode", resp.StatusCode)
		}
		resp, err = s.doTransmissionRequest(ctx, endpoint, config, sessionID)
		if err != nil {
			return requestFailedResult(err)
		}
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return statusFailedResult(resp.StatusCode)
	}

	var payload struct {
		Result    string `json:"result"`
		Arguments struct {
			Version    string `json:"version"`
			RPCVersion int    `json:"rpc-version"`
		} `json:"arguments"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return failedResult("Transmission response could not be parsed", "error", err.Error())
	}
	if payload.Result != "success" {
		return formatResultFailure("Transmission", payload.Result)
	}

	return successResult(
		"Transmission connection OK",
		"version", payload.Arguments.Version,
		"rpcVersion", payload.Arguments.RPCVersion,
	)
}

func (s *Service) doTransmissionRequest(ctx context.Context, endpoint string, config Config, sessionID string) (*http.Response, error) {
	req, err := jsonPost(ctx, endpoint, map[string]string{"method": "session-get"})
	if err != nil {
		return nil, err
	}
	if sessionID != "" {
		req.Header.Set(transmissionSessionHeader, sessionID)
	}
	addBasicAuth(req, config.Username, config.Password)
	return s.client.Do(req)
}
