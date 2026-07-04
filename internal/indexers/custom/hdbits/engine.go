package hdbits

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://hdbits.org"

var nonWord = regexp.MustCompile(`[\W]+`)

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
		return engine.FailedResult("Invalid HDBits request", "error", err.Error())
	}
	return engine.SuccessResult("HDBits indexer reachable")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := engine.ParseBaseURL(common.BaseURL(config, defaultBaseURL) + "/api/torrents")
	if err != nil {
		return nil, err
	}
	payload := searchPayload(config, query)
	var decoded apiResponse
	if err := postJSON(ctx, e.client, endpoint.String(), payload, &decoded); err != nil {
		return nil, err
	}
	if decoded.Status != 0 {
		return nil, fmt.Errorf("hdbits api status %d: %s", decoded.Status, decoded.Message)
	}
	var torrents []torrent
	if len(decoded.Data) > 0 {
		if err := json.Unmarshal(decoded.Data, &torrents); err != nil {
			return nil, err
		}
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := make([]engine.Release, 0, len(torrents))
	for _, item := range torrents {
		release := item.release(config, baseURL, passkey(config))
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchPayload(config engine.Config, query string) map[string]any {
	payload := map[string]any{
		"username": username(config),
		"passkey":  passkey(config),
		"limit":    100,
	}
	if search := strings.TrimSpace(query); search != "" {
		payload["search"] = strings.TrimSpace(nonWord.ReplaceAllString(search, " "))
	}
	if len(config.Categories) > 0 {
		categories := make([]int, 0, len(config.Categories))
		for _, category := range config.Categories {
			categories = append(categories, int(category))
		}
		payload["category"] = categories
	}
	return payload
}

func username(config engine.Config) string {
	return common.FieldString(config, "username")
}

func passkey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "passkey", "apiKey", "apikey"), engine.StringValue(config.APIKey))
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

type apiResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type torrent struct {
	ID             string `json:"id"`
	Hash           string `json:"hash"`
	Leechers       int32  `json:"leechers"`
	Seeders        int32  `json:"seeders"`
	Name           string `json:"name"`
	TimesCompleted int32  `json:"times_completed"`
	Size           int64  `json:"size"`
	UtAdded        int64  `json:"utadded"`
	Added          string `json:"added"`
	NumFiles       int32  `json:"numfiles"`
	FileName       string `json:"filename"`
	TypeCategory   int    `json:"type_category"`
	TypeMedium     int    `json:"type_medium"`
}

func (t torrent) release(config engine.Config, baseURL string, key string) engine.Release {
	infoURL, _ := common.URLWithQuery(baseURL, "/details.php", map[string]string{"id": t.ID})
	downloadURL, _ := common.URLWithQuery(baseURL, "/download.php", map[string]string{
		"id":      t.ID,
		"passkey": key,
	})
	title := strings.TrimSpace(t.Name)
	if useFilenames(config) && t.TypeCategory != 7 && t.TypeMedium != 1 && strings.TrimSpace(t.FileName) != "" {
		title = strings.TrimSuffix(strings.TrimSpace(t.FileName), ".torrent")
	}
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           title,
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            engine.FirstNonEmpty("HDBits-"+t.ID, infoURL),
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(int(t.Seeders)),
		Peers:           common.Int32Ptr(int(t.Seeders + t.Leechers)),
		PublishedAt:     hdbitsTime(t.Added, t.UtAdded),
	}
}

func hdbitsTime(value string, unix int64) *time.Time {
	if parsed := common.ParseFlexibleTime(value); parsed != nil {
		return parsed
	}
	return common.UnixTime(unix)
}

func useFilenames(config engine.Config) bool {
	if len(config.Fields) == 0 {
		return true
	}
	var fields []struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}
	if err := json.Unmarshal(config.Fields, &fields); err != nil {
		return true
	}
	for _, field := range fields {
		if !strings.EqualFold(field.Name, "useFilenames") {
			continue
		}
		switch value := field.Value.(type) {
		case bool:
			return value
		case string:
			return strings.EqualFold(strings.TrimSpace(value), "true")
		default:
			return true
		}
	}
	return true
}
