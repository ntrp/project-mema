package torrentday

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://www.torrentday.com/"

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
		return engine.FailedResult("Invalid TorrentDay request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "TorrentDay")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
	if err != nil {
		return nil, err
	}
	var rows []row
	if err := common.GetJSON(ctx, e.client, endpoint, &rows); err != nil {
		return nil, err
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := make([]engine.Release, 0, len(rows))
	for _, item := range rows {
		release := item.release(config, baseURL)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	baseURL := common.BaseURL(config, defaultBaseURL)
	path := "t.json?"
	categories := []string{}
	for _, category := range config.Categories {
		categories = append(categories, strconv.FormatInt(int64(category), 10))
	}
	if len(categories) > 0 {
		path += strings.Join(categories, ";")
	}
	if common.FieldBool(config, "freeleechOnly") || common.FieldBool(config, "freeleech") {
		if !strings.HasSuffix(path, "?") {
			path += ";"
		}
		path += "free"
	}
	path += ";q=" + url.QueryEscape(strings.TrimSpace(query))
	return common.ResolveURL(baseURL, path), nil
}

type row struct {
	ID        int    `json:"t"`
	Category  int    `json:"c"`
	Name      string `json:"name"`
	Seeders   int    `json:"seeders"`
	Leechers  int    `json:"leechers"`
	Size      int64  `json:"size"`
	Created   int64  `json:"ctime"`
	Files     int32  `json:"files"`
	Completed int32  `json:"completed"`
}

func (r row) release(config engine.Config, baseURL string) engine.Release {
	id := strconv.Itoa(r.ID)
	infoURL, _ := common.URLWithQuery(baseURL, "/details.php", map[string]string{"id": id})
	downloadURL := common.ResolveURL(baseURL, "download.php/"+id+"/"+id+".torrent")
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(r.Name),
		DownloadURL:     downloadURL,
		InfoURL:         infoURL,
		GUID:            infoURL,
		SizeBytes:       r.Size,
		Seeders:         common.Int32Ptr(r.Seeders),
		Peers:           common.Int32Ptr(r.Seeders + r.Leechers),
		PublishedAt:     common.UnixTime(r.Created),
	}
}
