package torrentrss

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/indexers/engine"
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
	req, err := engine.Get(ctx, config.BaseURL)
	if err != nil {
		return engine.FailedResult("Invalid RSS URL", "error", err.Error())
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return engine.RequestFailedResult(err)
	}
	defer engine.CloseBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return engine.StatusFailedResult(resp.StatusCode)
	}
	return engine.SuccessResult("Torrent RSS feed reachable", "endpoint", config.BaseURL)
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	req, err := engine.Get(ctx, config.BaseURL)
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
	query = strings.ToLower(strings.TrimSpace(query))
	releases := make([]engine.Release, 0, len(items))
	for _, item := range items {
		release := item.toRelease(config)
		if release.DownloadURL == "" || release.Title == "" {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(release.Title), query) {
			continue
		}
		releases = append(releases, release)
	}
	return releases, nil
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
	PubDate   string `xml:"pubDate"`
	Published string `xml:"published"`
	Size      int64  `xml:"size"`
	Enclosure struct {
		URL    string `xml:"url,attr"`
		Length int64  `xml:"length,attr"`
	} `xml:"enclosure"`
}

func (i item) toRelease(config engine.Config) engine.Release {
	size := i.Size
	if size == 0 {
		size = i.Enclosure.Length
	}
	published := parseRSSDate(engine.FirstNonEmpty(i.PubDate, i.Published))
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(i.Title),
		DownloadURL:     engine.FirstNonEmpty(i.Enclosure.URL, i.Link, httpGUID(i.GUID)),
		InfoURL:         engine.FirstNonEmpty(i.Link, httpGUID(i.GUID)),
		GUID:            strings.TrimSpace(engine.FirstNonEmpty(i.GUID, i.Link, i.Enclosure.URL)),
		SizeBytes:       size,
		PublishedAt:     published,
	}
}

func parseRSSDate(value string) *time.Time {
	value = strings.TrimSpace(value)
	for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC3339, time.RFC3339Nano} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			converted := parsed.UTC()
			return &converted
		}
	}
	if unix, err := strconv.ParseInt(value, 10, 64); err == nil && unix > 0 {
		parsed := time.Unix(unix, 0).UTC()
		return &parsed
	}
	return nil
}

func httpGUID(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "magnet:") {
		return value
	}
	return ""
}
