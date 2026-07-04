package beyondhd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://beyond-hd.me"

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
		return engine.FailedResult("Invalid BeyondHD request", "error", err.Error())
	}
	return engine.SuccessResult("BeyondHD indexer reachable")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config)
	if err != nil {
		return nil, err
	}
	payload := searchPayload(config, query)
	var decoded response
	if err := postJSON(ctx, e.client, endpoint, payload, &decoded); err != nil {
		return nil, err
	}
	if decoded.StatusCode == 0 {
		return nil, fmt.Errorf("beyondhd api error: %s", decoded.StatusMessage)
	}
	releases := make([]engine.Release, 0, len(decoded.Results))
	for _, item := range decoded.Results {
		release := item.release(config)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchURL(config engine.Config) (string, error) {
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/api/torrents/"+apiKey(config), nil)
}

func searchPayload(config engine.Config, query string) map[string]any {
	payload := map[string]any{
		"action": "search",
		"rsskey": rssKey(config),
	}
	if search := strings.TrimSpace(query); search != "" {
		payload["search"] = search
	}
	if common.FieldBool(config, "freeleechOnly") || common.FieldBool(config, "freeleech") {
		payload["freeleech"] = 1
	}
	if common.FieldBool(config, "limitedOnly") {
		payload["limited"] = 1
	}
	return payload
}

func apiKey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

func rssKey(config engine.Config) string {
	return common.FieldString(config, "rssKey", "rsskey")
}

func postJSON(ctx context.Context, client engine.HTTPDoer, endpoint string, payload any, target any) error {
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
	StatusCode    int       `json:"status_code"`
	StatusMessage string    `json:"status_message"`
	Results       []torrent `json:"results"`
}

type torrent struct {
	Name         string `json:"name"`
	InfoHash     string `json:"info_hash"`
	Category     string `json:"category"`
	Size         int64  `json:"size"`
	Grabs        int32  `json:"times_completed"`
	Seeders      int32  `json:"seeders"`
	Leechers     int32  `json:"leechers"`
	CreatedAt    string `json:"created_at"`
	DownloadLink string `json:"download_url"`
	InfoURL      string `json:"url"`
	IMDBID       string `json:"imdb_id"`
	Freeleech    bool   `json:"freeleech"`
	Limited      bool   `json:"limited"`
}

func (t torrent) release(config engine.Config) engine.Release {
	seeders := t.Seeders
	leechers := t.Leechers
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(t.Name),
		DownloadURL:     strings.TrimSpace(t.DownloadLink),
		InfoURL:         strings.TrimSpace(t.InfoURL),
		GUID:            strings.TrimSpace(t.InfoURL),
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(int(seeders)),
		Peers:           common.Int32Ptr(int(seeders + leechers)),
		PublishedAt:     common.ParseFlexibleTime(t.CreatedAt),
	}
}
