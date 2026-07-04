package rarbg

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://torrentapi.org/"

var magnetGUID = regexp.MustCompile(`(?i)^magnet:\?xt=urn:btih:([a-f0-9]+)`)

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
		return engine.FailedResult("Invalid Rarbg request", "error", err.Error())
	}
	return engine.SuccessResult("Rarbg indexer reachable")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	token, err := e.token(ctx, config)
	if err != nil {
		return nil, err
	}
	endpoint, err := searchURL(config, token, query)
	if err != nil {
		return nil, err
	}
	var decoded response
	if err := common.GetJSON(ctx, e.client, endpoint, &decoded); err != nil {
		return nil, err
	}
	if decoded.ErrorCode != nil {
		switch *decoded.ErrorCode {
		case 5, 8, 9, 10, 13, 14, 20:
			return []engine.Release{}, nil
		default:
			return nil, fmt.Errorf("rarbg api error %d: %s", *decoded.ErrorCode, decoded.Error)
		}
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

func (e *Engine) token(ctx context.Context, config engine.Config) (string, error) {
	if token := common.FieldString(config, "token"); token != "" {
		return token, nil
	}
	endpoint, err := common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/pubapi_v2.php", map[string]string{
		"get_token": "get_token",
		"app_id":    "rralworP_project-mema",
	})
	if err != nil {
		return "", err
	}
	var decoded struct {
		Token string `json:"token"`
	}
	if err := common.GetJSON(ctx, e.client, endpoint, &decoded); err != nil {
		return "", err
	}
	if strings.TrimSpace(decoded.Token) == "" {
		return "", fmt.Errorf("rarbg token response did not include token")
	}
	return strings.TrimSpace(decoded.Token), nil
}

func searchURL(config engine.Config, token string, query string) (string, error) {
	values := map[string]string{
		"limit":  "100",
		"token":  token,
		"format": "json_extended",
		"app_id": "rralworP_project-mema",
		"mode":   "search",
	}
	if strings.TrimSpace(query) == "" {
		values["mode"] = "list"
	} else {
		values["search_string"] = strings.TrimSpace(query)
	}
	if !common.FieldBool(config, "rankedOnly") {
		values["ranked"] = "0"
	}
	if len(config.Categories) > 0 {
		categories := make([]string, 0, len(config.Categories))
		for _, category := range config.Categories {
			categories = append(categories, strconv.FormatInt(int64(category), 10))
		}
		values["category"] = strings.Join(categories, ";")
	}
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/pubapi_v2.php", values)
}

type response struct {
	Error     string    `json:"error"`
	ErrorCode *int      `json:"error_code"`
	RateLimit int       `json:"rate_limit"`
	Results   []torrent `json:"torrent_results"`
}

type torrent struct {
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Download string    `json:"download"`
	Seeders  int       `json:"seeders"`
	Leechers int       `json:"leechers"`
	Size     int64     `json:"size"`
	PubDate  time.Time `json:"pubdate"`
	InfoPage string    `json:"info_page"`
}

func (t torrent) release(config engine.Config) engine.Release {
	infoURL := t.InfoPage
	if infoURL != "" {
		sep := "?"
		if strings.Contains(infoURL, "?") {
			sep = "&"
		}
		infoURL += sep + "app_id=" + url.QueryEscape("rralworP_project-mema")
	}
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(t.Title),
		DownloadURL:     strings.TrimSpace(t.Download),
		InfoURL:         infoURL,
		GUID:            rarbgGUID(t.Download),
		SizeBytes:       t.Size,
		Seeders:         common.Int32Ptr(t.Seeders),
		Peers:           common.Int32Ptr(t.Seeders + t.Leechers),
		PublishedAt:     &t.PubDate,
	}
}

func rarbgGUID(download string) string {
	if match := magnetGUID.FindStringSubmatch(download); len(match) == 2 {
		return "rarbg-" + match[1]
	}
	return "rarbg-" + download
}
