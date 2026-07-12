package tvsubtitles

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/sitehtml"
)

const key = "tvsubtitles"
const defaultBaseURL = "https://tvsubtitles.net"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

var idPattern = regexp.MustCompile(`(tvshow|episode|subtitle)-(\d+)`)
var scriptPartPattern = regexp.MustCompile(`(?m)s\d+\s*=\s*'([^']*)';`)

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
	base := strings.TrimRight(sitehtml.BaseURL(cfg, defaultBaseURL), "/")
	showID, err := searchShowID(ctx, svc, base, sr.Title)
	if err != nil {
		return nil, err
	}
	episodeID, err := findEpisodeID(ctx, svc, base, showID, *sr.SeasonNumber, *sr.EpisodeNumber)
	if err != nil {
		return nil, err
	}
	return findSubtitles(ctx, svc, base, episodeID, sr.LanguageID)
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(cand.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: candidate has no source URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cand.SourceURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return providercore.Download{}, err
	}
	parts := scriptPartPattern.FindAllStringSubmatch(doc.Text(), -1)
	if len(parts) == 0 {
		html, htmlErr := doc.Html()
		if htmlErr != nil {
			return providercore.Download{}, htmlErr
		}
		parts = scriptPartPattern.FindAllStringSubmatch(html, -1)
	}
	var relative strings.Builder
	for _, part := range parts {
		relative.WriteString(part[1])
	}
	if relative.Len() == 0 {
		return providercore.Download{}, fmt.Errorf("%w: download link missing", providercore.ErrProviderBrokenUpstream)
	}
	downloadURL := sitehtml.Resolve(cand.SourceURL, relative.String())
	downloadReq, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	return sitehtml.Download(downloadReq, svc, key, true, "subtitle.zip")
}

func searchShowID(ctx context.Context, svc providercore.Service, base, title string) (int64, error) {
	form := url.Values{"qs": {title}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/search1.php", strings.NewReader(form.Encode()))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return 0, err
	}
	var id int64
	doc.Find(`div.left li div a[href^="/tvshow-"], a[href^="/tvshow-"]`).EachWithBreak(func(_ int, link *goquery.Selection) bool {
		href, _ := link.Attr("href")
		match := idPattern.FindStringSubmatch(href)
		if len(match) == 3 && strings.Contains(strings.ToLower(link.Text()), strings.ToLower(title)) {
			id, _ = strconv.ParseInt(match[2], 10, 64)
			return false
		}
		return true
	})
	if id == 0 {
		return 0, fmt.Errorf("%w: show not found", providercore.ErrProviderBrokenUpstream)
	}
	return id, nil
}

func findEpisodeID(ctx context.Context, svc providercore.Service, base string, showID int64, season, episode int32) (int64, error) {
	raw := fmt.Sprintf("%s/tvshow-%d-%d.html", base, showID, season)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return 0, err
	}
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return 0, err
	}
	var id int64
	doc.Find("table#table5 tr").EachWithBreak(func(_ int, row *goquery.Selection) bool {
		cells := row.Find("td")
		if cells.Length() < 2 || !strings.Contains(strings.TrimSpace(cells.Eq(0).Text()), fmt.Sprintf("x%02d", episode)) {
			return true
		}
		href, _ := cells.Eq(1).Find(`a[href^="episode-"]`).Attr("href")
		match := idPattern.FindStringSubmatch(href)
		if len(match) == 3 {
			id, _ = strconv.ParseInt(match[2], 10, 64)
			return false
		}
		return true
	})
	if id == 0 {
		return 0, fmt.Errorf("%w: episode not found", providercore.ErrProviderBrokenUpstream)
	}
	return id, nil
}

func findSubtitles(ctx context.Context, svc providercore.Service, base string, episodeID int64, fallbackLanguage string) ([]providercore.Candidate, error) {
	raw := fmt.Sprintf("%s/episode-%d.html", base, episodeID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return nil, err
	}
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return nil, err
	}
	var out []providercore.Candidate
	doc.Find(".subtitlen").Each(func(_ int, item *goquery.Selection) {
		parent := item.Parent()
		href, _ := parent.Attr("href")
		match := idPattern.FindStringSubmatch(href)
		if len(match) != 3 {
			return
		}
		id, _ := strconv.ParseInt(match[2], 10, 64)
		lang := fallbackLanguage
		if src, ok := item.Find("h5 img").Attr("src"); ok {
			lang = strings.TrimSuffix(path.Base(src), path.Ext(src))
		}
		release := strings.TrimSpace(item.Find("h5").Text())
		out = append(out, providercore.Candidate{ProviderName: key, LanguageID: lang, FileID: id, Format: "srt", ReleaseName: release, SourceURL: fmt.Sprintf("%s/download-%d.html", base, id), SourceRef: raw})
	})
	if len(out) == 0 {
		return nil, fmt.Errorf("%w: no subtitles found", providercore.ErrProviderBrokenUpstream)
	}
	return out, nil
}
