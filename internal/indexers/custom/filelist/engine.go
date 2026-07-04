package filelist

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://filelist.io"

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
	endpoint, err := searchURL(config, "test")
	if err != nil {
		return engine.FailedResult("Invalid FileList request", "error", err.Error())
	}
	req, err := authenticatedRequest(ctx, config, endpoint)
	if err != nil {
		return engine.FailedResult("Invalid FileList request", "error", err.Error())
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return engine.RequestFailedResult(err)
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.StatusFailedResult(resp.StatusCode)
	}
	return engine.SuccessResult("FileList indexer reachable", "endpoint", endpoint)
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
	if err != nil {
		return nil, err
	}
	req, err := authenticatedRequest(ctx, config, endpoint)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, engine.HTTPStatusError(resp)
	}
	var decoded []torrent
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := make([]engine.Release, 0, len(decoded))
	for _, item := range decoded {
		release := item.release(config, baseURL, passkey(config))
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	parsed, err := engine.ParseBaseURL(common.BaseURL(config, defaultBaseURL) + "/api.php")
	if err != nil {
		return "", err
	}
	values := url.Values{}
	if search := strings.TrimSpace(query); search != "" {
		values.Set("action", "search-torrents")
		values.Set("type", "name")
		values.Set("query", search)
	} else {
		values.Set("action", "latest-torrents")
	}
	if len(config.Categories) > 0 {
		categories := make([]string, 0, len(config.Categories))
		for _, category := range config.Categories {
			categories = append(categories, strconv.FormatInt(int64(category), 10))
		}
		values.Set("category", strings.Join(categories, ","))
	}
	if common.FieldBool(config, "freeleechOnly") || common.FieldBool(config, "freeleech") {
		values.Set("freeleech", "1")
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func authenticatedRequest(ctx context.Context, config engine.Config, endpoint string) (*http.Request, error) {
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(username(config), passkey(config))
	return req, nil
}

func username(config engine.Config) string {
	return common.FieldString(config, "username")
}

func passkey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "passkey", "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

type torrent struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Size           int64  `json:"size"`
	Leechers       int32  `json:"leechers"`
	Seeders        int32  `json:"seeders"`
	TimesCompleted int32  `json:"times_completed"`
	Files          int32  `json:"files"`
	UploadDate     string `json:"upload_date"`
	Category       string `json:"category"`
	Freeleech      bool   `json:"freeleech"`
	DoubleUp       bool   `json:"doubleup"`
}

func (t torrent) release(config engine.Config, baseURL string, key string) engine.Release {
	id := strconv.FormatUint(uint64(t.ID), 10)
	infoURL, _ := common.URLWithQuery(baseURL, "/details.php", map[string]string{"id": id})
	downloadURL, _ := common.URLWithQuery(baseURL, "/download.php", map[string]string{
		"id":      id,
		"passkey": strings.TrimSpace(key),
	})
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(t.Name),
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            "FileList-" + id,
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(int(t.Seeders)),
		Peers:           common.Int32Ptr(int(t.Seeders + t.Leechers)),
		PublishedAt:     common.ParseFlexibleTime(t.UploadDate + " +0300"),
	}
}
