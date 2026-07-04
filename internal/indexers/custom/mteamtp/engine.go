package mteamtp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://kp.m-team.cc"

type Engine struct {
	client engine.HTTPDoer
}

func New(clients ...engine.HTTPDoer) *Engine {
	var client engine.HTTPDoer
	if len(clients) > 0 {
		client = clients[0]
	}
	return &Engine{client: client}
}

func (e *Engine) Test(ctx context.Context, config engine.Config) engine.TestResult {
	_, err := e.Search(ctx, config, "test", "")
	if err != nil {
		return engine.FailedResult("Invalid M-Team request", "error", err.Error())
	}
	return engine.SuccessResult("M-Team indexer reachable")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := apiEndpoint(config, "/api/torrent/search")
	if err != nil {
		return nil, err
	}
	payload := searchPayload(config, query)
	var decoded response
	if err := postJSON(ctx, e.client, endpoint, apiKey(config), payload, &decoded); err != nil {
		return nil, err
	}
	if decoded.Message != "" && !strings.EqualFold(decoded.Message, "SUCCESS") && decoded.Data.Torrents == nil {
		return nil, fmt.Errorf("m-team api response: %s", decoded.Message)
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	apiBase, _ := apiBaseURL(config)
	releases := make([]engine.Release, 0, len(decoded.Data.Torrents))
	for _, item := range decoded.Data.Torrents {
		release := item.release(config, baseURL, apiBase)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchPayload(config engine.Config, query string) map[string]any {
	payload := map[string]any{
		"mode":       "Normal",
		"categories": categoryStrings(config),
		"pageNumber": 1,
		"pageSize":   100,
	}
	if search := strings.TrimSpace(query); search != "" {
		payload["keyword"] = search
	}
	if common.FieldBool(config, "freeleechOnly") || common.FieldBool(config, "freeleech") {
		payload["discount"] = "FREE"
	}
	return payload
}

func categoryStrings(config engine.Config) []string {
	values := make([]string, 0, len(config.Categories))
	for _, category := range config.Categories {
		values = append(values, strconv.FormatInt(int64(category), 10))
	}
	return values
}

func apiKey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

func apiEndpoint(config engine.Config, path string) (string, error) {
	base, err := apiBaseURL(config)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(path, "/"), nil
}

func apiBaseURL(config engine.Config) (string, error) {
	baseURL := common.BaseURL(config, defaultBaseURL)
	parsed, err := engine.ParseBaseURL(baseURL)
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(parsed.Host, ".", 2)
	if len(parts) == 2 {
		parsed.Host = "api." + parts[1]
	}
	parsed.Path = ""
	parsed.RawQuery = ""
	return parsed.String(), nil
}

func postJSON(ctx context.Context, client engine.HTTPDoer, endpoint string, key string, payload any, target any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", key)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.35")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.HTTPStatusError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

type response struct {
	Data    data   `json:"data"`
	Message string `json:"message"`
}

type data struct {
	Torrents []torrent `json:"data"`
}

type torrent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"smallDescr"`
	Category    string `json:"category"`
	NumFiles    string `json:"numfiles"`
	Size        string `json:"size"`
	Status      status `json:"status"`
}

type status struct {
	CreatedDate    string `json:"createdDate"`
	Discount       string `json:"discount"`
	TimesCompleted string `json:"timesCompleted"`
	Seeders        string `json:"seeders"`
	Leechers       string `json:"leechers"`
}

func (t torrent) release(config engine.Config, baseURL string, apiBase string) engine.Release {
	infoURL := strings.TrimRight(baseURL, "/") + "/detail/" + t.ID
	downloadURL := strings.TrimRight(apiBase, "/") + "/api/torrent/genDlToken?id=" + url.QueryEscape(t.ID)
	seeders := atoi32(t.Status.Seeders)
	leechers := atoi32(t.Status.Leechers)
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.Join(strings.Fields(t.Name), " "),
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       atoi64(t.Size),
		Seeders:         common.Int32Ptr(int(seeders)),
		Peers:           common.Int32Ptr(int(seeders + leechers)),
		PublishedAt:     common.ParseFlexibleTime(t.Status.CreatedDate + " +0800"),
	}
}

func atoi32(value string) int32 {
	parsed, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 32)
	return int32(parsed)
}

func atoi64(value string) int64 {
	parsed, _ := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	return parsed
}
