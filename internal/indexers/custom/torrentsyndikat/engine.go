package torrentsyndikat

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://torrent-syndikat.org"

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
		return engine.FailedResult("Invalid TorrentSyndikat request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "TorrentSyndikat")
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
	baseURL := common.BaseURL(config, defaultBaseURL)
	key := apiKey(config)
	releases := make([]engine.Release, 0, len(decoded.Rows))
	for _, item := range decoded.Rows {
		release := item.release(config, baseURL, key)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	parsed, err := engine.ParseBaseURL(common.BaseURL(config, defaultBaseURL) + "/api_9djWe8Tb2NE3p6opyqnh/v1/browse.php")
	if err != nil {
		return "", err
	}
	values := url.Values{}
	values.Set("apikey", apiKey(config))
	values.Set("limit", "50")
	values.Set("ponly", strconv.FormatBool(common.FieldBool(config, "productsOnly")))
	if search := strings.TrimSpace(query); search != "" {
		values.Set("searchstring", seasonPackWildcard(search))
	}
	if len(config.Categories) > 0 {
		categories := make([]string, 0, len(config.Categories))
		for _, category := range config.Categories {
			categories = append(categories, strconv.FormatInt(int64(category), 10))
		}
		values.Set("cats", strings.Join(categories, ","))
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func apiKey(config engine.Config) string {
	return engine.FirstNonEmpty(common.FieldString(config, "apiKey", "apikey"), engine.StringValue(config.APIKey))
}

func seasonPackWildcard(query string) string {
	parts := strings.Fields(strings.TrimSpace(query))
	for i, part := range parts {
		if len(part) == 3 && (part[0] == 's' || part[0] == 'S') && part[1] >= '0' && part[1] <= '9' && part[2] >= '0' && part[2] <= '9' {
			parts[i] = part + "*"
			break
		}
	}
	return strings.Join(parts, " ")
}

type response struct {
	Rows []torrent `json:"rows"`
}

type torrent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Added    int64  `json:"added"`
	Size     int64  `json:"size"`
	NumFiles int32  `json:"numfiles"`
	Seeders  int32  `json:"seeders"`
	Leechers int32  `json:"leechers"`
	Snatched int32  `json:"snatched"`
	Category int    `json:"category"`
	Poster   string `json:"poster"`
	IMDBID   int    `json:"imdbId"`
	TVDBID   int    `json:"tvdbId"`
	TMDBID   int    `json:"tmdbId"`
}

func (t torrent) release(config engine.Config, baseURL string, key string) engine.Release {
	infoURL, _ := common.URLWithQuery(baseURL, "/details.php", map[string]string{"id": t.ID})
	downloadURL, _ := common.URLWithQuery(baseURL, "/download.php", map[string]string{
		"id":     t.ID,
		"apikey": key,
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
		PublishedAt:     common.UnixTime(t.Added),
	}
}
