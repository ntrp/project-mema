package scenehd

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://scenehd.org"

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
		return engine.FailedResult("Invalid SceneHD request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "SceneHD")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
	if err != nil {
		return nil, err
	}
	var decoded []torrent
	if err := common.GetJSON(ctx, e.client, endpoint, &decoded); err != nil {
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
	parsed, err := engine.ParseBaseURL(common.BaseURL(config, defaultBaseURL) + "/browse.php")
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("api", "")
	values.Set("passkey", passkey(config))
	if search := strings.TrimSpace(query); search != "" {
		values.Set("search", search)
	}
	if len(config.Categories) > 0 {
		categories := make([]string, 0, len(config.Categories))
		for _, category := range config.Categories {
			categories = append(categories, strconv.FormatInt(int64(category), 10))
		}
		values.Set("cat", strings.Join(categories, ","))
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func passkey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "passkey"), engine.StringValue(config.APIKey))
}

type torrent struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Added          string `json:"added"`
	Size           int64  `json:"size"`
	TimesCompleted int32  `json:"times_completed"`
	NumFiles       int32  `json:"numfiles"`
	Seeders        int32  `json:"seeders"`
	Leechers       int32  `json:"leechers"`
	Freeleech      int    `json:"is_freeleech"`
}

func (t torrent) release(config engine.Config, baseURL string, passkey string) engine.Release {
	id := strconv.FormatInt(t.ID, 10)
	infoURL, _ := common.URLWithQuery(baseURL, "/details.php", map[string]string{"id": id})
	downloadURL, _ := common.URLWithQuery(baseURL, "/download.php", map[string]string{
		"id":      id,
		"passkey": passkey,
	})
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(t.Name),
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(int(t.Seeders)),
		Peers:           common.Int32Ptr(int(t.Seeders + t.Leechers)),
		PublishedAt:     common.ParseFlexibleTime(t.Added),
	}
}
