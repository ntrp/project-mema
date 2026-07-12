package subtitriid

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/nativeutil"
)

const baseURL = "https://subtitri.do.am"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return nativeutil.Test(ctx, svc, cfg, "subtitriid", baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && sr.MediaType != "movie" {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	form := url.Values{"q": {sr.Title}}
	data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subtitriid", BaseURL: baseURL, Path: "/search/", Form: form})
	if err != nil {
		return nil, err
	}
	return parseSearch(ctx, svc, cfg, data), nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return nativeutil.DownloadSubtitle(ctx, svc, cfg, "subtitriid", baseURL, cand)
}

func parseSearch(ctx context.Context, svc providercore.Service, cfg providercore.Config, data []byte) []providercore.Candidate {
	doc, err := nativeutil.Document(data)
	if err != nil {
		return nil
	}
	out := []providercore.Candidate{}
	doc.Find(".eBlock").Each(func(_ int, row *goquery.Selection) {
		page := nativeutil.Attr(row, ".eTitle > a", "href")
		if page == "" {
			return
		}
		candidate, ok := detail(ctx, svc, cfg, page)
		if ok {
			out = append(out, candidate)
		}
	})
	return out
}

func detail(ctx context.Context, svc providercore.Service, cfg providercore.Config, page string) (providercore.Candidate, bool) {
	data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subtitriid", BaseURL: baseURL, Path: page})
	if err != nil {
		return providercore.Candidate{}, false
	}
	doc, err := nativeutil.Document(data)
	if err != nil {
		return providercore.Candidate{}, false
	}
	down := nativeutil.Attr(doc.Selection, ".hvr", "href")
	if down == "" {
		return providercore.Candidate{}, false
	}
	year, _ := strconv.Atoi(nativeutil.FirstText(doc.Selection, "#film-page-year"))
	name := strings.TrimSpace(nativeutil.FirstText(doc.Selection, ".main-header"))
	if parts := strings.Split(name, " / "); len(parts) > 0 {
		name = parts[len(parts)-1]
	}
	return providercore.Candidate{ProviderName: "subtitriid", LanguageID: "lav", Format: nativeutil.Format(down), ReleaseName: strings.TrimSpace(name + " " + strconv.Itoa(year)), SourceURL: nativeutil.Absolute(cfg, baseURL, down)}, true
}
