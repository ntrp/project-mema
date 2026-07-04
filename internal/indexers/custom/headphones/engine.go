package headphones

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://indexer.codeshy.com"

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
		return engine.FailedResult("Invalid Headphones request", "error", err.Error())
	}
	req, err := request(ctx, config, endpoint)
	if err != nil {
		return engine.FailedResult("Invalid Headphones request", "error", err.Error())
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return engine.RequestFailedResult(err)
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.StatusFailedResult(resp.StatusCode)
	}
	return engine.SuccessResult("Headphones indexer reachable", "endpoint", endpoint)
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
	if err != nil {
		return nil, err
	}
	req, err := request(ctx, config, endpoint)
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
	body, err := engine.ReadLimitedBody(resp.Body)
	if err != nil {
		return nil, err
	}
	items, err := decodeFeed(body)
	if err != nil {
		return nil, err
	}
	releases := make([]engine.Release, 0, len(items))
	for _, item := range items {
		release := item.toRelease(config)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	apiPath := common.FieldString(config, "apiPath")
	if apiPath == "" {
		apiPath = "/api"
	}
	apiKey := common.FieldString(config, "apiKey", "apikey")
	if apiKey == "" {
		apiKey = engine.StringValue(config.APIKey)
	}
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), apiPath, map[string]string{
		"t":        "search",
		"extended": "1",
		"apikey":   apiKey,
		"q":        strings.ReplaceAll(strings.TrimSpace(query), "+", "%20"),
	})
}

func request(ctx context.Context, config engine.Config, endpoint string) (*http.Request, error) {
	req, err := engine.Get(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	username := common.FieldString(config, "username", "user")
	password := common.FieldString(config, "password")
	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}
	req.Header.Set("Accept", "application/rss+xml, application/xml")
	return req, nil
}

func decodeFeed(body []byte) ([]item, error) {
	var feed rssFeed
	if err := xml.NewDecoder(bytes.NewReader(body)).Decode(&feed); err != nil {
		return nil, err
	}
	return feed.Channel.Items, nil
}

type rssFeed struct {
	Channel struct {
		Items []item `xml:"item"`
	} `xml:"channel"`
}

type item struct {
	Title     string `xml:"title"`
	Link      string `xml:"link"`
	GUID      string `xml:"guid"`
	Comments  string `xml:"comments"`
	PubDate   string `xml:"pubDate"`
	Attrs     []attr `xml:"attr"`
	Enclosure struct {
		URL    string `xml:"url,attr"`
		Length int64  `xml:"length,attr"`
	} `xml:"enclosure"`
}

type attr struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (i item) toRelease(config engine.Config) engine.Release {
	size := i.int64Attr("size")
	if size == 0 {
		size = i.Enclosure.Length
	}
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(i.Title),
		DownloadURL:     engine.FirstNonEmpty(i.Enclosure.URL, i.Link, httpGUID(i.GUID)),
		InfoURL:         strings.TrimSuffix(engine.FirstNonEmpty(i.Comments, i.Link, httpGUID(i.GUID)), "#comments"),
		GUID:            strings.TrimSpace(engine.FirstNonEmpty(i.GUID, i.Link, i.Enclosure.URL)),
		SizeBytes:       size,
		PublishedAt:     common.ParseFlexibleTime(engine.FirstNonEmpty(i.attr("usenetdate"), i.PubDate)),
	}
}

func (i item) attr(name string) string {
	for _, attr := range i.Attrs {
		if strings.EqualFold(attr.Name, name) {
			return attr.Value
		}
	}
	return ""
}

func (i item) int64Attr(name string) int64 {
	var value int64
	if _, err := fmt.Sscanf(i.attr(name), "%d", &value); err == nil {
		return value
	}
	return 0
}

func httpGUID(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "magnet:") {
		return value
	}
	return ""
}
