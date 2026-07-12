package nativehtml

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/sitehtml"
	"media-manager/internal/subtitles/security"
)

type Spec struct {
	Key             string
	BaseURL         string
	SearchPath      func(providercore.SearchRequest) string
	Query           func(url.Values, providercore.SearchRequest)
	Candidate       func(*goquery.Selection, string, string) (providercore.Candidate, bool)
	NeedsCookie     bool
	ArchiveDownload bool
	UserAgent       string
}

type Adapter struct{ Spec Spec }

func New(spec Spec) Adapter { return Adapter{Spec: spec} }

func Register(spec Spec) { providers.Register(spec.Key, New(spec)) }

func (a Adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	if err := a.auth(cfg); err != nil {
		return err
	}
	req, err := a.request(ctx, cfg, http.MethodGet, a.base(cfg), nil, false)
	if err != nil {
		return err
	}
	return sitehtml.Test(req, svc, a.Spec.Key)
}

func (a Adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if err := a.auth(cfg); err != nil {
		return nil, err
	}
	searchURL, err := a.searchURL(cfg, sr)
	if err != nil {
		return nil, err
	}
	req, err := a.request(ctx, cfg, http.MethodGet, searchURL, nil, false)
	if err != nil {
		return nil, err
	}
	doc, _, err := sitehtml.DoHTML(svc, req, a.Spec.Key)
	if err != nil {
		return nil, err
	}
	return a.parse(doc, searchURL, sr.LanguageID), nil
}

func (a Adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	if err := a.auth(cfg); err != nil {
		return providercore.Download{}, err
	}
	raw := strings.TrimSpace(cand.SourceURL)
	if raw == "" {
		raw = strings.TrimSpace(cand.SourceRef)
	}
	if raw == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := a.request(ctx, cfg, http.MethodGet, a.resolve(cfg, raw), nil, true)
	if err != nil {
		return providercore.Download{}, err
	}
	return download(req, svc, a.Spec.Key, a.Spec.ArchiveDownload, cand.ReleaseName)
}

func (a Adapter) auth(cfg providercore.Config) error {
	if a.Spec.NeedsCookie && providercore.NewConfig(cfg).CookieString() == "" {
		return providercore.ErrPrivateMembershipRequired
	}
	return nil
}

func (a Adapter) request(ctx context.Context, cfg providercore.Config, method, raw string, body io.Reader, download bool) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, raw, body)
	if err != nil {
		return nil, err
	}
	ua := a.Spec.UserAgent
	if ua == "" {
		ua = "Mozilla/5.0 (compatible; MemaSubtitleProvider/1.0)"
	}
	req.Header.Set("User-Agent", ua)
	if cookies := providercore.NewConfig(cfg).CookieString(); cookies != "" {
		req.Header.Set("Cookie", cookies)
	}
	if download {
		req.Header.Set("Accept", "application/zip,application/x-rar-compressed,application/octet-stream,text/plain,*/*")
	}
	return req, nil
}

func (a Adapter) searchURL(cfg providercore.Config, sr providercore.SearchRequest) (string, error) {
	base, err := url.Parse(a.base(cfg))
	if err != nil {
		return "", err
	}
	searchPath := "/"
	if a.Spec.SearchPath != nil {
		searchPath = a.Spec.SearchPath(sr)
	}
	base.Path = joinURLPath(base.Path, searchPath)
	q := base.Query()
	if a.Spec.Query != nil {
		a.Spec.Query(q, sr)
	}
	base.RawQuery = q.Encode()
	return base.String(), nil
}

func (a Adapter) parse(doc *goquery.Document, pageURL, fallbackLang string) []providercore.Candidate {
	out := []providercore.Candidate{}
	seen := map[string]struct{}{}
	doc.Find("tr, li, article, div.subtitle, div.item, div.media, a[href]").Each(func(_ int, s *goquery.Selection) {
		if a.Spec.Candidate == nil {
			return
		}
		cand, ok := a.Spec.Candidate(s, pageURL, fallbackLang)
		if !ok || cand.SourceURL == "" || cand.ReleaseName == "" {
			return
		}
		cand.ProviderName = a.Spec.Key
		key := cand.SourceURL
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		out = append(out, cand)
	})
	return out
}

func (a Adapter) resolve(cfg providercore.Config, raw string) string {
	return sitehtml.Resolve(a.base(cfg)+"/", raw)
}

func (a Adapter) base(cfg providercore.Config) string {
	return providercore.NewConfig(cfg).BaseURL(a.Spec.BaseURL)
}

func QueryText(sr providercore.SearchRequest) string {
	parts := []string{sr.Title}
	if sr.Year != nil {
		parts = append(parts, strconv.Itoa(int(*sr.Year)))
	}
	if sr.SeasonNumber != nil && sr.EpisodeNumber != nil {
		parts = append(parts, fmt.Sprintf("S%02dE%02d", *sr.SeasonNumber, *sr.EpisodeNumber))
	}
	return strings.Join(parts, " ")
}

func LinkCandidate(sel *goquery.Selection, pageURL, fallbackLang string, selectors ...string) (providercore.Candidate, bool) {
	link := sel
	if goquery.NodeName(sel) != "a" {
		selector := strings.Join(selectors, ",")
		if selector == "" {
			selector = "a[href]"
		}
		link = sel.Find(selector).First()
	}
	href, ok := link.Attr("href")
	if !ok || strings.TrimSpace(href) == "" || strings.HasPrefix(strings.TrimSpace(href), "#") {
		return providercore.Candidate{}, false
	}
	release := first(sitehtml.Attr(sel, "data-release", "title"), sitehtml.Attr(link, "title", "data-release"), sitehtml.Text(sel, ".title", ".release", ".name", "td"), strings.TrimSpace(link.Text()))
	lang := first(sitehtml.Attr(sel, "data-language", "lang"), sitehtml.Attr(link, "data-language", "lang"), fallbackLang)
	return providercore.Candidate{LanguageID: lang, Format: format(href), ReleaseName: release, SourceURL: sitehtml.Resolve(pageURL, href), SourceRef: href}, true
}

func download(req *http.Request, svc providercore.Service, key string, archive bool, release string) (providercore.Download, error) {
	resp, err := svc.DoProviderRequest(req, key, true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return providercore.Download{}, fmt.Errorf("%w: download http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, sitehtml.MaxDownloadBytes+1))
	if err != nil {
		return providercore.Download{}, err
	}
	if len(body) > sitehtml.MaxDownloadBytes {
		return providercore.Download{}, fmt.Errorf("%w: download response too large", security.ErrUnsafeArchive)
	}
	content := body
	if archive || sitehtml.LooksArchived(req.URL.String()) || sitehtml.LooksArchived(resp.Header.Get("Content-Type")) {
		member, err := providercore.ExtractSubtitle(sitehtml.FilenameFor(req.URL.String(), release), body, security.ArchiveLimits{})
		if err != nil {
			return providercore.Download{}, err
		}
		content = member.Content
	}
	return providercore.Download{Content: content, URL: req.URL.String()}, nil
}

func joinURLPath(basePath, child string) string {
	if child == "" || child == "/" {
		return basePath
	}
	return path.Join(basePath, child)
}

func format(raw string) string {
	ext := strings.TrimPrefix(path.Ext(strings.Split(raw, "?")[0]), ".")
	if ext == "" || len(ext) > 4 {
		return "srt"
	}
	return ext
}

func first(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
