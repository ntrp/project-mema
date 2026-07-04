package torrentscsv

import (
	"context"
	"net/url"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://torrents-csv.com/"

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
	baseURL := common.BaseURL(config, defaultBaseURL)
	return common.TestURL(ctx, e.client, baseURL+"/", "TorrentsCSV")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	query = strings.TrimSpace(query)
	if len(query) < 3 {
		return []engine.Release{}, nil
	}
	endpoint, err := common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/service/search", map[string]string{
		"size": "100",
		"q":    query,
	})
	if err != nil {
		return nil, err
	}
	var payload response
	if err := common.GetJSON(ctx, e.client, endpoint, &payload); err != nil {
		return nil, err
	}
	releases := make([]engine.Release, 0, len(payload.Torrents))
	for _, torrent := range payload.Torrents {
		if strings.TrimSpace(torrent.InfoHash) == "" || strings.TrimSpace(torrent.Name) == "" {
			continue
		}
		magnet := common.Magnet(torrent.InfoHash)
		seeders := torrent.Seeders
		leechers := torrent.Leechers
		releases = append(releases, engine.Release{
			IndexerID:       config.ID,
			IndexerName:     config.Name,
			IndexerProtocol: config.Protocol,
			Title:           strings.TrimSpace(torrent.Name),
			DownloadURL:     magnet,
			InfoURL:         common.BaseURL(config, defaultBaseURL) + "/search?q=" + url.QueryEscape(strings.TrimSpace(torrent.Name)),
			GUID:            magnet,
			SizeBytes:       torrent.Size,
			Seeders:         common.Int32Ptr(seeders),
			Peers:           common.Int32Ptr(seeders + leechers),
			PublishedAt:     common.UnixTime(torrent.Created),
		})
	}
	return releases, nil
}

type response struct {
	Torrents []torrent `json:"torrents"`
}

type torrent struct {
	InfoHash  string `json:"infohash"`
	Name      string `json:"name"`
	Size      int64  `json:"size_bytes"`
	Created   int64  `json:"created_unix"`
	Leechers  int    `json:"leechers"`
	Seeders   int    `json:"seeders"`
	Completed int    `json:"completed"`
}
