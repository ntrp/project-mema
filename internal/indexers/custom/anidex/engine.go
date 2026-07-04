package anidex

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/indexers/custom/common"
	"media-manager/internal/indexers/engine"
)

const defaultBaseURL = "https://anidex.info/"

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
		return engine.FailedResult("Invalid Anidex request", "error", err.Error())
	}
	return common.TestURL(ctx, e.client, endpoint, "Anidex")
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
	doc.Find("div#content table > tbody > tr").Each(func(index int, row *goquery.Selection) {
		release := releaseFromRow(config, baseURL, row)
		if release.Title != "" && release.DownloadURL != "" {
			releases = append(releases, release)
		}
	})
	return releases, nil
}

func searchURL(config engine.Config, query string) (string, error) {
	values := map[string]string{
		"page":     "search",
		"s":        "upload_timestamp",
		"o":        "desc",
		"group_id": "0",
		"q":        strings.TrimSpace(query),
	}
	if common.FieldBool(config, "authorisedOnly") {
		values["a"] = "1"
	}
	return common.URLWithQuery(common.BaseURL(config, defaultBaseURL), "", values)
}

func releaseFromRow(config engine.Config, baseURL string, row *goquery.Selection) engine.Release {
	downloadPath, _ := row.Find(`a[href^="/dl/"]`).First().Attr("href")
	infoPath, _ := row.Find("td:nth-child(3) a").First().Attr("href")
	magnet, _ := row.Find(`a[href^="magnet:?"]`).First().Attr("href")
	title, _ := row.Find("td:nth-child(3) span").First().Attr("title")
	language, _ := row.Find("td:nth-child(1) img").First().Attr("title")
	if strings.TrimSpace(language) != "" {
		title = strings.TrimSpace(title) + " [" + strings.TrimSpace(language) + "]"
	}
	seeders, _ := strconv.Atoi(strings.TrimSpace(row.Find("td:nth-child(9)").Text()))
	leechers, _ := strconv.Atoi(strings.TrimSpace(row.Find("td:nth-child(10)").Text()))
	added, _ := row.Find("td:nth-child(8)").First().Attr("title")
	return engine.Release{
		IndexerID:       config.ID,
		IndexerName:     config.Name,
		IndexerProtocol: config.Protocol,
		Title:           strings.TrimSpace(title),
		DownloadURL:     engine.FirstNonEmpty(common.ResolveURL(baseURL, downloadPath), magnet),
		InfoURL:         common.ResolveURL(baseURL, infoPath),
		GUID:            common.ResolveURL(baseURL, infoPath),
		SizeBytes:       common.ParseSizeBytes(row.Find("td:nth-child(7)").Text()),
		Seeders:         common.Int32Ptr(seeders),
		Peers:           common.Int32Ptr(seeders + leechers),
		PublishedAt:     common.ParseFlexibleTime(strings.TrimSpace(added)),
	}
}
