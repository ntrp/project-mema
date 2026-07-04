package binsearch

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://binsearch.info/"

var titleRegex = regexp.MustCompile(`"(?P<title>.*)(?:\.(rar|nfo|mkv|par2|001|nzb|url|zip|r[0-9]{2}))"`)

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
		return engine.FailedResult("Invalid BinSearch request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "BinSearch")
}

func (e *Engine) Search(ctx context.Context, config engine.Config, query string, mediaType string) ([]engine.Release, error) {
	endpoint, err := searchURL(config, query)
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
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	baseURL := common.BaseURL(config, defaultBaseURL)
	releases := []engine.Release{}
	doc.Find("table.xMenuT > tbody > tr").Each(func(index int, row *goquery.Selection) {
		if index == 0 {
			return
		}
		release := releaseFromRow(config, baseURL, row)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	})
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "", map[string]string{
		"adv_col":  "on",
		"postdate": "date",
		"adv_sort": "date",
		"q":        strings.TrimSpace(query),
		"m":        "0",
		"max":      "100",
	})
}

func releaseFromRow(config engine.Config, baseURL string, row *goquery.Selection) engine.Release {
	rawTitle := row.Find("td > span.s").First().Text()
	title := parsedTitle(rawTitle)
	guid, _ := row.Find(`input[type="checkbox"]`).First().Attr("name")
	link, _ := row.Find("a").First().Attr("href")
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           title,
		DownloadURL:     common.ResolveURL(baseURL, "/?action=nzb&"+guid+"=1"),
		InfoURL:         common.ResolveURL(baseURL, link),
		GUID:            strings.TrimSpace(guid),
		SizeBytes:       common.ParseSizeBytes(row.Find("td > span.d").First().Text()),
		PublishedAt:     common.ParseFlexibleTime(row.Find("td:nth-child(6)").Text()),
	}
}

func parsedTitle(value string) string {
	match := titleRegex.FindStringSubmatch(value)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return strings.Trim(value, `" `)
}
