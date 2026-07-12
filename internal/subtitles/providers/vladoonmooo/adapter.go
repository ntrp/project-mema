package vladoonmooo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/sitehtml"
)

const key = "vladoonmooo"
const defaultBaseURL = "https://vladoon.mooo.com"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sitehtml.BaseURL(cfg, defaultBaseURL), nil)
	if err != nil {
		return err
	}
	return sitehtml.Test(req, svc, key)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if !sitehtml.Supports(sr.MediaType, "movie", "serie") {
		return nil, sitehtml.Unsupported(key, sr.MediaType)
	}
	values := url.Values{"q": {sr.Title}}
	if sr.Year != nil {
		values.Set("year", fmt.Sprint(*sr.Year))
	}
	if sr.SeasonNumber != nil {
		values.Set("season", fmt.Sprint(*sr.SeasonNumber))
	}
	if sr.EpisodeNumber != nil {
		values.Set("episode", fmt.Sprint(*sr.EpisodeNumber))
	}
	searchURL := sitehtml.BaseURL(cfg, defaultBaseURL) + "/search.php?" + values.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return nil, err
	}
	return parseVladoonmooo(doc.Selection, searchURL, sr.LanguageID)
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(cand.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: candidate has no source URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cand.SourceURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	return sitehtml.Download(req, svc, key, true, cand.ReleaseName)
}

func parseVladoonmooo(root *goquery.Selection, pageURL, languageID string) ([]providercore.Candidate, error) {
	var candidates []providercore.Candidate
	root.Find("[data-subtitle], [data-subtitle], .subtitle, .result, table tr, li").Each(func(_ int, row *goquery.Selection) {
		href := ""
		row.Find("a[href]").EachWithBreak(func(_ int, link *goquery.Selection) bool {
			candidate, _ := link.Attr("href")
			lower := strings.ToLower(candidate)
			if strings.Contains(lower, "download") || strings.Contains(lower, "down") || strings.HasSuffix(lower, ".zip") || strings.HasSuffix(lower, ".srt") {
				href = candidate
				return false
			}
			return true
		})
		if href == "" {
			return
		}
		lang := sitehtml.Attr(row, "data-lang", "data-language")
		if lang == "" {
			lang = languageID
		}
		release := sitehtml.Attr(row, "data-release", "data-title")
		if release == "" {
			release = sitehtml.Text(row, ".release", ".title", "a")
		}
		candidates = append(candidates, providercore.Candidate{ProviderName: key, LanguageID: lang, Format: "srt", ReleaseName: release, SourceURL: sitehtml.Resolve(pageURL, href), SourceRef: pageURL})
	})
	if len(candidates) == 0 {
		return nil, fmt.Errorf("%w: no subtitle links found", providercore.ErrProviderBrokenUpstream)
	}
	return candidates, nil
}
