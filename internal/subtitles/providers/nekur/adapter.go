package nekur

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/htmlutil"
)

const (
	provider = "nekur"
	baseURL  = "http://subtitri.nekur.net"
)

type adapter struct{}

var Adapter providercore.Adapter = adapter{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return htmlutil.Test(ctx, svc, cfg, provider, baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && !strings.EqualFold(sr.MediaType, "movie") {
		return nil, fmt.Errorf("%w: nekur does not support %s", providercore.ErrProviderPrerequisiteMissing, sr.MediaType)
	}
	endpoint, err := url.Parse(htmlutil.BaseURL(cfg, baseURL) + "/modules/Subtitles.php")
	if err != nil {
		return nil, err
	}
	form := url.Values{"ajax": {"1"}, "sSearch": {sr.Title}}
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodPost, endpoint.String(), strings.NewReader(form.Encode()), false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "search")
	if err != nil {
		return nil, err
	}
	return parse(endpoint.String(), body)
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return htmlutil.Download(ctx, svc, provider, cand.SourceURL, cand.ReleaseName, true)
}

func parse(pageURL string, body []byte) ([]providercore.Candidate, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var out []providercore.Candidate
	doc.Find("tbody > tr").Each(func(_ int, row *goquery.Selection) {
		link := row.Find(".title > a[href]").First()
		href, ok := link.Attr("href")
		if !ok {
			return
		}
		title := strings.TrimSpace(link.Contents().Not("span").Text())
		year := strings.Trim(row.Find(".year").First().Text(), "() ")
		notes := strings.Join(strings.Fields(row.Find(".notes").First().Text()), " ")
		fps := strings.TrimSpace(row.Find(".fps").First().Text())
		release := strings.TrimSpace(strings.Join([]string{title, year, fps, notes}, " "))
		out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: "lav", Format: "srt", ReleaseName: release, SourceURL: htmlutil.Resolve(pageURL, href), SourceRef: pageURL})
	})
	if len(out) == 0 {
		return nil, fmt.Errorf("%w: no nekur rows found", providercore.ErrProviderBrokenUpstream)
	}
	return out, nil
}
