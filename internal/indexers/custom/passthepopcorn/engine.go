package passthepopcorn

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

const defaultBaseURL = "https://passthepopcorn.me"

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
		return engine.FailedResult("Invalid PassThePopcorn request", "error", err.Error())
	}
	req, err := authenticatedRequest(ctx, config, endpoint)
	if err != nil {
		return engine.FailedResult("Invalid PassThePopcorn request", "error", err.Error())
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return engine.RequestFailedResult(err)
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.StatusFailedResult(resp.StatusCode)
	}
	return engine.SuccessResult("PassThePopcorn indexer reachable", "endpoint", endpoint)
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
	var decoded response
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := []engine.Release{}
	for _, movie := range decoded.Movies {
		for _, item := range movie.Torrents {
			release := item.release(config, baseURL, movie)
			if release.Title != "" && release.DownloadURL != "" {
				releases = append(releases, release)
			}
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	parsed, err := engine.ParseBaseURL(common.BaseURL(config, defaultBaseURL) + "/torrents.php")
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("action", "advanced")
	values.Set("json", "noredirect")
	values.Set("grouping", "0")
	values.Set("order_by", "time")
	values.Set("order_way", "desc")
	values.Set("searchstr", strings.TrimSpace(query))
	if common.FieldBool(config, "freeleechOnly") || common.FieldBool(config, "freeleech") {
		values.Set("freetorrent", "1")
	}
	if common.FieldBool(config, "goldenPopcornOnly") {
		values.Set("scene", "2")
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
	req.Header.Set("ApiUser", apiUser(config))
	req.Header.Set("ApiKey", apiKey(config))
	return req, nil
}

func apiUser(config engine.Config) string {
	return common.FieldString(config, "apiUser", "apiuser")
}

func apiKey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

type response struct {
	TotalResults string  `json:"TotalResults"`
	Movies       []movie `json:"Movies"`
}

type movie struct {
	GroupID    string    `json:"GroupId"`
	CategoryID string    `json:"CategoryId"`
	Title      string    `json:"Title"`
	Year       string    `json:"Year"`
	IMDBID     string    `json:"ImdbId"`
	Torrents   []torrent `json:"Torrents"`
}

type torrent struct {
	ID            int    `json:"Id"`
	Size          string `json:"Size"`
	UploadTime    string `json:"UploadTime"`
	Snatched      string `json:"Snatched"`
	Seeders       string `json:"Seeders"`
	Leechers      string `json:"Leechers"`
	ReleaseName   string `json:"ReleaseName"`
	FreeleechType string `json:"FreeleechType"`
}

func (t torrent) release(config engine.Config, baseURL string, parent movie) engine.Release {
	id := strconv.Itoa(t.ID)
	infoURL, _ := common.URLWithQuery(baseURL, "/torrents.php", map[string]string{
		"id":        parent.GroupID,
		"torrentid": id,
	})
	downloadURL, _ := common.URLWithQuery(baseURL, "/torrents.php", map[string]string{
		"action": "download",
		"id":     id,
	})
	seeders := atoi32(t.Seeders)
	leechers := atoi32(t.Leechers)
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(t.ReleaseName),
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            "PassThePopcorn-" + id,
		SizeBytes:       atoi64(t.Size),
		Seeders:         common.Int32Ptr(int(seeders)),
		Peers:           common.Int32Ptr(int(seeders + leechers)),
		PublishedAt:     common.ParseFlexibleTime(t.UploadTime + " +0000"),
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
