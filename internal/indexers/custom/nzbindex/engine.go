package nzbindex

import (
	"context"
	"regexp"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://nzbindex.com/"

var titleRegex = regexp.MustCompile(`"(?P<title>[^:/]*?)(?:\.(rar|nfo|mkv|par2|001|nzb|url|zip|r[0-9]{2}))?"`)

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
		return engine.FailedResult("Invalid NZBIndex request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "NZBIndex")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
	if err != nil {
		return nil, err
	}
	var payload response
	if err := common.GetJSON(ctx, e.client, endpoint, &payload); err != nil {
		return nil, err
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := make([]engine.Release, 0, len(payload.Results))
	for _, item := range payload.Results {
		title := parsedTitle(item.Name)
		if title == "" || item.ID == "" {
			continue
		}
		published := time.UnixMilli(item.Posted).UTC()
		releases = append(releases, engine.Release{
			IndexerID:       config.ID,
			IndexerName:     config.Name,
			IndexerProtocol: config.Protocol,
			Title:           title,
			DownloadURL:     baseURL + "/download/" + item.ID,
			InfoURL:         baseURL + "/collection/" + item.ID,
			GUID:            baseURL + "/collection/" + item.ID,
			SizeBytes:       item.Size,
			PublishedAt:     &published,
		})
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	apiKey := common.FieldString(config, "apiKey", "apikey", "key")
	if apiKey == "" {
		apiKey = engine.StringValue(config.APIKey)
	}
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "/api/v3/search/", map[string]string{
		"key": apiKey,
		"max": "100",
		"q":   strings.TrimSpace(query),
		"p":   "0",
	})
}

func parsedTitle(value string) string {
	match := titleRegex.FindStringSubmatch(value)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return strings.TrimSpace(value)
}

type response struct {
	Results []result `json:"results"`
}

type result struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Posted int64  `json:"posted"`
	Size   int64  `json:"size"`
}
