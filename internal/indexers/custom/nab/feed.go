package nab

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

type releaseFeed struct {
	Channel struct {
		Items []releaseItem `xml:"item"`
	} `xml:"channel"`
}

type releaseItem struct {
	Title     string        `xml:"title"`
	Link      string        `xml:"link"`
	GUID      string        `xml:"guid"`
	PubDate   string        `xml:"pubDate"`
	Published string        `xml:"published"`
	Updated   string        `xml:"updated"`
	Size      int64         `xml:"size"`
	Attrs     []torznabAttr `xml:"attr"`
	Enclosure struct {
		URL    string `xml:"url,attr"`
		Length int64  `xml:"length,attr"`
	} `xml:"enclosure"`
}

type torznabAttr struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (e *Engine) search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	searchType := "search"
	if mediaType == "movie" {
		searchType = "movie"
	}
	if mediaType == "serie" {
		searchType = "tvsearch"
	}

	values := map[string]string{
		"t":      searchType,
		"q":      query,
		"apikey": engine.StringValue(config.APIKey),
	}
	if len(config.Categories) > 0 {
		values["cat"] = categoryQuery(config.Categories)
	}

	endpoint, err := engine.EndpointWithQuery(config.BaseURL, values)
	if err != nil {
		return nil, err
	}

	req, err := engine.Get(ctx, endpoint)
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
	items, err := decodeReleaseFeed(body)
	if err != nil {
		return nil, err
	}

	releases := make([]engine.Release, 0, len(items))
	for _, item := range items {
		release := item.toRelease(config)
		if release.DownloadURL == "" {
			continue
		}
		releases = append(releases, release)
	}
	return releases, nil
}

func decodeReleaseFeed(body []byte) ([]releaseItem, error) {
	var feed releaseFeed
	if err := xml.NewDecoder(bytes.NewReader(body)).Decode(&feed); err != nil {
		return nil, err
	}
	return feed.Channel.Items, nil
}

func (item releaseItem) toRelease(config engine.Config) engine.Release {
	size := item.Size
	if size == 0 {
		size = item.Enclosure.Length
	}
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(item.Title),
		DownloadURL:     engine.FirstNonEmpty(item.Enclosure.URL, item.Link, httpGUID(item.GUID)),
		InfoURL:         engine.FirstNonEmpty(item.Link, httpGUID(item.GUID)),
		GUID:            strings.TrimSpace(item.GUID),
		SizeBytes:       size,
		Seeders:         item.int32Attr("seeders"),
		Peers:           item.int32Attr("peers"),
		PublishedAt:     item.publishedAt(),
	}
}

func (item releaseItem) publishedAt() *time.Time {
	for _, value := range []string{item.PubDate, item.Published, item.Updated, item.attr("publishdate")} {
		parsed, ok := parseFeedTime(value)
		if ok {
			return &parsed
		}
	}
	return nil
}

func parseFeedTime(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC3339, time.RFC3339Nano} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func (item releaseItem) attr(name string) string {
	for _, attr := range item.Attrs {
		if attr.Name == name {
			return attr.Value
		}
	}
	return ""
}

func (item releaseItem) int32Attr(name string) *int32 {
	for _, attr := range item.Attrs {
		if attr.Name != name {
			continue
		}
		parsed, err := strconv.ParseInt(attr.Value, 10, 32)
		if err != nil {
			return nil
		}
		value := int32(parsed)
		return &value
	}
	return nil
}

func categoryQuery(categories []int32) string {
	values := make([]string, 0, len(categories))
	for _, category := range categories {
		values = append(values, strconv.FormatInt(int64(category), 10))
	}
	return strings.Join(values, ",")
}

func httpGUID(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "magnet:") {
		return value
	}
	return ""
}
