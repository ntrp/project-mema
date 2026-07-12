package tvsubtitles

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

const key = "tvsubtitles"
const defaultBaseURL = "https://tvsubtitles.net"

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
	if !sitehtml.Supports(sr.MediaType, "serie") {
		return nil, sitehtml.Unsupported(key, sr.MediaType)
	}
	if sr.SeasonNumber == nil || sr.EpisodeNumber == nil {
		return nil, fmt.Errorf("%w: season and episode required", providercore.ErrProviderPrerequisiteMissing)
	}
	base := sitehtml.BaseURL(cfg, defaultBaseURL)
	showURL, err := tvSearchShow(ctx, svc, base, sr)
	if err != nil {
		return nil, err
	}
	seasonURL := sitehtml.Resolve(showURL, fmt.Sprintf("season-%d.html", *sr.SeasonNumber))
	episodeURL, err := tvFindEpisode(ctx, svc, seasonURL, sr)
	if err != nil {
		return nil, err
	}
	return tvFindSubtitles(ctx, svc, episodeURL, sr.LanguageID)
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(cand.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: candidate has no source URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cand.SourceURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	return sitehtml.Download(req, svc, key, false, cand.ReleaseName)
}

func tvSearchShow(ctx context.Context, svc providercore.Service, base string, sr providercore.SearchRequest) (string, error) {
	form := url.Values{"q": {sr.Title}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/search.php", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return "", err
	}
	showURL := ""
	doc.Find("a[href]").EachWithBreak(func(_ int, link *goquery.Selection) bool {
		href, _ := link.Attr("href")
		lower := strings.ToLower(href + " " + link.Text())
		if strings.Contains(lower, "tvshow") || strings.Contains(lower, strings.ToLower(sr.Title)) {
			showURL = sitehtml.Resolve(base+"/search.php", href)
			return false
		}
		return true
	})
	if showURL == "" {
		return "", fmt.Errorf("%w: show not found", providercore.ErrProviderBrokenUpstream)
	}
	return showURL, nil
}

func tvFindEpisode(ctx context.Context, svc providercore.Service, seasonURL string, sr providercore.SearchRequest) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, seasonURL, nil)
	if err != nil {
		return "", err
	}
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return "", err
	}
	episodeURL := ""
	episodeText := fmt.Sprintf("x%02d", *sr.EpisodeNumber)
	doc.Find("a[href]").EachWithBreak(func(_ int, link *goquery.Selection) bool {
		href, _ := link.Attr("href")
		text := strings.ToLower(link.Text() + " " + href)
		if strings.Contains(text, episodeText) || strings.Contains(text, fmt.Sprintf("episode-%d", *sr.EpisodeNumber)) {
			episodeURL = sitehtml.Resolve(seasonURL, href)
			return false
		}
		return true
	})
	if episodeURL == "" {
		return "", fmt.Errorf("%w: episode not found", providercore.ErrProviderBrokenUpstream)
	}
	return episodeURL, nil
}

func tvFindSubtitles(ctx context.Context, svc providercore.Service, episodeURL, languageID string) ([]providercore.Candidate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, episodeURL, nil)
	if err != nil {
		return nil, err
	}
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return nil, err
	}
	var candidates []providercore.Candidate
	doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, _ := link.Attr("href")
		lower := strings.ToLower(href)
		if !strings.Contains(lower, "download") && !strings.Contains(lower, "subtitle") {
			return
		}
		lang := languageID
		if data := sitehtml.Attr(link, "data-lang", "data-language"); data != "" {
			lang = data
		}
		candidates = append(candidates, providercore.Candidate{ProviderName: key, LanguageID: lang, Format: "srt", ReleaseName: sitehtml.Text(link), SourceURL: sitehtml.Resolve(episodeURL, href), SourceRef: episodeURL})
	})
	if len(candidates) == 0 {
		return nil, fmt.Errorf("%w: no subtitle links found", providercore.ErrProviderBrokenUpstream)
	}
	return candidates, nil
}
