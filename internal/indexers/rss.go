package indexers

import (
	"context"
	"net/http"
	"strings"
)

type rssDocument struct {
	Channel struct {
		Title string    `xml:"title"`
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
}

type atomDocument struct {
	Title   string      `xml:"title"`
	Entries []atomEntry `xml:"entry"`
}

type rssItem struct{}

type atomEntry struct{}

func (s *Service) searchRSS(ctx context.Context, config Config, query string) ([]Release, error) {
	endpoint, err := endpointWithQuery(config.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	req, err := get(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, httpStatusError(resp.StatusCode)
	}

	body, err := readLimitedBody(resp.Body)
	if err != nil {
		return nil, err
	}
	items, err := decodeReleaseFeed(body)
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	releases := []Release{}
	for _, item := range items {
		if !strings.Contains(strings.ToLower(item.Title), query) {
			continue
		}
		release := item.toRelease(config)
		if release.DownloadURL == "" {
			continue
		}
		releases = append(releases, release)
	}
	return releases, nil
}

func (s *Service) testRSS(ctx context.Context, config Config) TestResult {
	endpoint, err := endpointWithQuery(config.BaseURL, nil)
	if err != nil {
		return failedResult("Invalid RSS URL", "error", err.Error())
	}

	req, err := get(ctx, endpoint)
	if err != nil {
		return failedResult("Invalid RSS request", "error", err.Error())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return requestFailedResult(err)
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return statusFailedResult(resp.StatusCode)
	}

	var rss rssDocument
	if err := decodeLimitedXML(resp.Body, &rss); err != nil {
		return failedResult("RSS response could not be parsed", "error", err.Error())
	}
	if rss.Channel.Title != "" || len(rss.Channel.Items) > 0 {
		return successResult("RSS feed OK", "title", rss.Channel.Title, "itemCount", len(rss.Channel.Items))
	}

	return failedResult("RSS feed did not include a channel")
}
