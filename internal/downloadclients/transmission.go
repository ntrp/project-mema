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

func (s *Service) addTransmission(ctx context.Context, config Config, request AddRequest) AddResult {
	endpoint, err := endpointWithPath(config.BaseURL, "/transmission/rpc", func(value string) bool {
		return strings.HasSuffix(strings.TrimRight(value, "/"), "/transmission/rpc")
	})
	if err != nil {
		return addFailedResult("Invalid Transmission URL", "error", err.Error())
	}

	resp, err := s.doTransmissionAddRequest(ctx, endpoint, config, request, "")
	if err != nil {
		return addFailedResult("Connection failed", "error", err.Error())
	}

	if resp.StatusCode == http.StatusConflict {
		sessionID := resp.Header.Get(transmissionSessionHeader)
		closeBody(resp.Body)
		if sessionID == "" {
			return addFailedResult("Transmission session id was not returned", "statusCode", resp.StatusCode)
		}
		resp, err = s.doTransmissionAddRequest(ctx, endpoint, config, request, sessionID)
		if err != nil {
			return addFailedResult("Connection failed", "error", err.Error())
		}
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return addFailedResult("Unexpected response status", "statusCode", resp.StatusCode)
	}

	var payload struct {
		Result    string `json:"result"`
		Arguments struct {
			TorrentAdded struct {
				ID int `json:"id"`
			} `json:"torrent-added"`
			TorrentDuplicate struct {
				ID int `json:"id"`
			} `json:"torrent-duplicate"`
		} `json:"arguments"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return addFailedResult("Transmission response could not be parsed", "error", err.Error())
	}
	if payload.Result != "success" {
		return addFailedResult("Transmission rejected the download", "result", payload.Result)
	}

	id := payload.Arguments.TorrentAdded.ID
	if id == 0 {
		id = payload.Arguments.TorrentDuplicate.ID
	}
	return addSuccessResult("Transmission download queued", intString(id), "torrentId", id)
}

func (s *Service) statusTransmission(ctx context.Context, config Config, request StatusRequest) StatusResult {
	endpoint, err := endpointWithPath(config.BaseURL, "/transmission/rpc", func(value string) bool {
		return strings.HasSuffix(strings.TrimRight(value, "/"), "/transmission/rpc")
	})
	if err != nil {
		return statusLookupFailedResult("Invalid Transmission URL", "error", err.Error())
	}

	resp, err := s.doTransmissionStatusRequest(ctx, endpoint, config, request, "")
	if err != nil {
		return statusLookupFailedResult("Connection failed", "error", err.Error())
	}
	if resp.StatusCode == http.StatusConflict {
		sessionID := resp.Header.Get(transmissionSessionHeader)
		closeBody(resp.Body)
		if sessionID == "" {
			return statusLookupFailedResult("Transmission session id was not returned", "statusCode", resp.StatusCode)
		}
		resp, err = s.doTransmissionStatusRequest(ctx, endpoint, config, request, sessionID)
		if err != nil {
			return statusLookupFailedResult("Connection failed", "error", err.Error())
		}
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return statusLookupFailedResult("Unexpected response status", "statusCode", resp.StatusCode)
	}

	var payload struct {
		Result    string `json:"result"`
		Arguments struct {
			Torrents []struct {
				ID            int     `json:"id"`
				Status        int     `json:"status"`
				ErrorString   string  `json:"errorString"`
				IsFinished    bool    `json:"isFinished"`
				PercentDone   float64 `json:"percentDone"`
				LeftUntilDone int64   `json:"leftUntilDone"`
				DownloadDir   string  `json:"downloadDir"`
				Name          string  `json:"name"`
				Files         []struct {
					Name           string `json:"name"`
					Length         int64  `json:"length"`
					BytesCompleted int64  `json:"bytesCompleted"`
				} `json:"files"`
			} `json:"torrents"`
		} `json:"arguments"`
	}
	if err := decodeLimitedJSON(resp.Body, &payload); err != nil {
		return statusLookupFailedResult("Transmission response could not be parsed", "error", err.Error())
	}
	if payload.Result != "success" {
		return statusLookupFailedResult("Transmission rejected the status request", "result", payload.Result)
	}
	if len(payload.Arguments.Torrents) == 0 {
		return statusLookupNotFoundResult("Transmission download was not found")
	}
	torrent := payload.Arguments.Torrents[0]
	if torrent.ErrorString != "" {
		return statusLookupResult("failed", torrent.ErrorString, "torrentId", torrent.ID, "statusCode", torrent.Status)
	}
	if torrent.IsFinished || torrent.PercentDone >= 1 || torrent.LeftUntilDone == 0 && torrent.Status >= 5 {
		return statusLookupResultWithFiles(
			"completed",
			"Transmission download completed",
			transmissionStatusFiles(torrent.DownloadDir, torrent.Name, torrent.Files),
			"torrentId", torrent.ID,
			"statusCode", torrent.Status,
		)
	}
	if torrent.Status == 0 {
		return statusLookupResult("grabbed", "Transmission download is stopped", "torrentId", torrent.ID, "statusCode", torrent.Status)
	}
	return statusLookupResult("downloading", "Transmission download is active", "torrentId", torrent.ID, "statusCode", torrent.Status)
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

func (s *Service) doTransmissionStatusRequest(ctx context.Context, endpoint string, config Config, request StatusRequest, sessionID string) (*http.Response, error) {
	id, err := parseTransmissionID(request.DownloadID)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"method": "torrent-get",
		"arguments": map[string]interface{}{
			"ids": []int{id},
			"fields": []string{
				"id",
				"status",
				"errorString",
				"isFinished",
				"percentDone",
				"leftUntilDone",
				"downloadDir",
				"name",
				"files",
			},
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

func transmissionStatusFiles(downloadDir string, name string, files []struct {
	Name           string `json:"name"`
	Length         int64  `json:"length"`
	BytesCompleted int64  `json:"bytesCompleted"`
}) []StatusFile {
	if len(files) == 0 && name != "" {
		return []StatusFile{{Path: pathJoin(downloadDir, name), Complete: true}}
	}
	results := make([]StatusFile, 0, len(files))
	for _, file := range files {
		results = append(results, StatusFile{
			Path:      pathJoin(downloadDir, file.Name),
			SizeBytes: file.Length,
			Complete:  file.Length == 0 || file.BytesCompleted >= file.Length,
		})
	}
	return results
}

func (s *Service) doTransmissionAddRequest(ctx context.Context, endpoint string, config Config, request AddRequest, sessionID string) (*http.Response, error) {
	body := map[string]interface{}{
		"method": "torrent-add",
		"arguments": map[string]interface{}{
			"filename": request.URL,
			"paused":   false,
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
