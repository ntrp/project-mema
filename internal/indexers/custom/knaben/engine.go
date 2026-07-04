package knaben

import (
	"context"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const (
	defaultBaseURL = "https://knaben.org/"
	apiEndpoint    = "https://api.knaben.org/v1"
)

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
	return common.TestURL(ctx, e.client, common.BaseURL(config, defaultBaseURL)+"/", "Knaben")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	body := map[string]any{
		"order_by":        "date",
		"order_direction": "desc",
		"from":            0,
		"size":            100,
		"hide_unsafe":     true,
	}
	if strings.TrimSpace(query) != "" {
		body["search_type"] = "100%"
		body["search_field"] = "title"
		body["query"] = strings.TrimSpace(query)
	}
	var payload response
	if err := common.PostJSON(ctx, e.client, apiEndpoint, body, &payload); err != nil {
		return nil, err
	}
	releases := make([]engine.Release, 0, len(payload.Hits))
	for _, hit := range payload.Hits {
		if hit.Seeders <= 0 || strings.TrimSpace(hit.Title) == "" {
			continue
		}
		downloadURL := strings.TrimSpace(hit.DownloadURL)
		if downloadURL == "" {
			downloadURL = strings.TrimSpace(hit.MagnetURL)
		}
		if downloadURL == "" && hit.InfoHash != "" {
			downloadURL = common.Magnet(hit.InfoHash)
		}
		if downloadURL == "" {
			continue
		}
		releases = append(releases, engine.Release{
			IndexerID:       config.ID,
			IndexerName:     config.Name,
			IndexerProtocol: config.Protocol,
			Title:           strings.TrimSpace(hit.Title),
			DownloadURL:     downloadURL,
			InfoURL:         strings.TrimSpace(hit.InfoURL),
			GUID:            strings.TrimSpace(hit.InfoURL),
			SizeBytes:       hit.Size,
			Seeders:         common.Int32Ptr(hit.Seeders),
			Peers:           common.Int32Ptr(hit.Seeders + hit.Leechers),
			PublishedAt:     common.ParseTime(hit.Date),
		})
	}
	return releases, nil
}

type response struct {
	Hits []hit `json:"hits"`
}

type hit struct {
	Title       string `json:"title"`
	InfoHash    string `json:"hash"`
	InfoURL     string `json:"details"`
	DownloadURL string `json:"link"`
	MagnetURL   string `json:"magnetUrl"`
	Size        int64  `json:"bytes"`
	Seeders     int    `json:"seeders"`
	Leechers    int    `json:"peers"`
	Date        string `json:"date"`
}
