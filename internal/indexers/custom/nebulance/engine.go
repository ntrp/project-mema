package nebulance

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://nebulance.io"

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
		return engine.FailedResult("Invalid Nebulance request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "Nebulance")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
	if err != nil {
		return nil, err
	}
	var decoded response
	if err := common.GetJSON(ctx, e.client, endpoint, &decoded); err != nil {
		return nil, err
	}
	if decoded.Error != nil {
		return nil, fmt.Errorf("nebulance api error: %s", decoded.Error.Message)
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := make([]engine.Release, 0, len(decoded.Items))
	for _, item := range decoded.Items {
		release := item.release(config, baseURL)
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
	values.Set("action", "search")
	values.Set("age", ">0")
	values.Set("api_key", apiKey(config))
	values.Set("per_page", "100")
	if search := strings.TrimSpace(query); search != "" {
		values.Set("release", search)
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func apiKey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

type response struct {
	TotalResults int       `json:"total_results"`
	Items        []torrent `json:"items"`
	Error        *apiError `json:"error"`
}

type apiError struct {
	Message string `json:"message"`
}

type torrent struct {
	ReleaseTitle   string   `json:"rls_name"`
	Size           int64    `json:"size"`
	Seeders        int32    `json:"seed"`
	Leechers       int32    `json:"leech"`
	Snatches       int32    `json:"snatch"`
	DownloadLink   string   `json:"download"`
	FileList       []string `json:"file_list"`
	GroupName      string   `json:"group_name"`
	TorrentID      int      `json:"group_id"`
	PublishDateUTC string   `json:"rls_utc"`
}

func (t torrent) release(config engine.Config, baseURL string) engine.Release {
	id := strconv.Itoa(t.TorrentID)
	infoURL, _ := common.URLWithQuery(baseURL, "/torrents.php", map[string]string{"id": id})
	title := engine.FirstNonEmpty(t.ReleaseTitle, t.GroupName)
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(title),
		DownloadURL:     strings.TrimSpace(t.DownloadLink),
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(int(t.Seeders)),
		Peers:           common.Int32Ptr(int(t.Seeders + t.Leechers)),
		PublishedAt:     common.ParseFlexibleTime(t.PublishDateUTC),
	}
}
