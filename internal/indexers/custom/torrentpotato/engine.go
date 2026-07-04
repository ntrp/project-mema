package torrentpotato

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

var magnetHashRegex = regexp.MustCompile(`(?i)^magnet:\?xt=urn:btih:([a-f0-9]+)`)

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
		return engine.FailedResult("Invalid TorrentPotato request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "TorrentPotato")
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
	if payload.Error != "" {
		return nil, fmt.Errorf("torrentpotato error: %s", payload.Error)
	}
	releases := make([]engine.Release, 0, len(payload.Results))
	for _, item := range payload.Results {
		if item.Name == "" || item.DownloadURL == "" {
			continue
		}
		seeders := item.Seeders
		releases = append(releases, engine.Release{
			IndexerID:       config.ID,
			IndexerName:     config.Name,
			IndexerProtocol: config.Protocol,
			Title:           html.UnescapeString(item.Name),
			DownloadURL:     item.DownloadURL,
			InfoURL:         item.DetailsURL,
			GUID:            potatoGUID(item.DownloadURL),
			SizeBytes:       item.Size * 1000 * 1000,
			Seeders:         common.Int32Ptr(seeders),
			Peers:           common.Int32Ptr(seeders + item.Leechers),
			PublishedAt:     item.publishTime(),
		})
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	baseURL := common.BaseURL(config, "http://127.0.0.1")
	user := common.FieldString(config, "user", "username")
	passkey := common.FieldString(config, "passkey", "apiKey", "apikey")
	if passkey == "" {
		passkey = engine.StringValue(config.APIKey)
	}
	return common.URLWithQuery(baseURL, "", map[string]string{
		"passkey": passkey,
		"user":    user,
		"search":  strings.TrimSpace(query),
	})
}

func potatoGUID(downloadURL string) string {
	if match := magnetHashRegex.FindStringSubmatch(downloadURL); len(match) == 2 {
		return "potato-" + match[1]
	}
	return "potato-" + downloadURL
}

type response struct {
	Results []result `json:"results"`
	Error   string   `json:"error"`
}

type result struct {
	Name        string    `json:"release_name"`
	DetailsURL  string    `json:"details_url"`
	DownloadURL string    `json:"download_url"`
	Size        int64     `json:"size"`
	Leechers    int       `json:"leechers"`
	Seeders     int       `json:"seeders"`
	PublishDate time.Time `json:"publish_date"`
}

func (r result) publishTime() *time.Time {
	value := r.PublishDate.UTC()
	if value.IsZero() {
		return nil
	}
	return &value
}
