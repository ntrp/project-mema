package publichtml

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const maxHTMLBytes = 2 << 20
const maxDownloadBytes = 50 << 20

// Spec describes a public subtitle site that exposes searchable HTML pages with download links.
type Spec struct {
	Key        string
	Name       string
	BaseURL    string
	MediaTypes []string
	Archive    bool
}

type Adapter struct{ spec Spec }

func New(spec Spec) Adapter { return Adapter{spec: spec} }

func (a Adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	base := a.baseURL(cfg)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base, nil)
	if err != nil {
		return err
	}
	resp, err := svc.DoProviderRequest(req, a.spec.Key, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return fmt.Errorf("%s test failed: http status %d", a.spec.Key, resp.StatusCode)
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return nil
}

func (a Adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if !a.supports(sr.MediaType) {
		return nil, fmt.Errorf("%w: %s does not support %s", providercore.ErrProviderPrerequisiteMissing, a.spec.Key, sr.MediaType)
	}
	searchURL, err := a.searchURL(cfg, sr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.DoProviderRequest(req, a.spec.Key, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("%w: rate limited", providercore.ErrProviderBrokenUpstream)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%w: search http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxHTMLBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxHTMLBytes {
		return nil, fmt.Errorf("%w: search response too large", providercore.ErrProviderBrokenUpstream)
	}
	return a.parseCandidates(searchURL, body, sr.LanguageID)
}

func (a Adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(cand.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: candidate has no source URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cand.SourceURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	resp, err := svc.DoProviderRequest(req, a.spec.Key, true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return providercore.Download{}, fmt.Errorf("%w: rate limited", providercore.ErrProviderBrokenUpstream)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return providercore.Download{}, fmt.Errorf("%w: download http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxDownloadBytes+1))
	if err != nil {
		return providercore.Download{}, err
	}
	if len(body) > maxDownloadBytes {
		return providercore.Download{}, fmt.Errorf("%w: download response too large", security.ErrUnsafeArchive)
	}
	content := body
	if a.spec.Archive || looksArchived(cand.SourceURL) || looksArchived(resp.Header.Get("Content-Type")) {
		member, err := providercore.ExtractSubtitle(filenameFor(cand.SourceURL, cand.ReleaseName), body, security.ArchiveLimits{})
		if err != nil {
			return providercore.Download{}, err
		}
		content = member.Content
	}
	return providercore.Download{Content: content, URL: cand.SourceURL}, nil
}

func (a Adapter) supports(mediaType string) bool {
	if mediaType == "" {
		return true
	}
	for _, allowed := range a.spec.MediaTypes {
		if strings.EqualFold(allowed, mediaType) {
			return true
		}
	}
	return false
}

func (a Adapter) baseURL(cfg providercore.Config) string {
	if strings.TrimSpace(cfg.BaseURL) != "" {
		return strings.TrimRight(cfg.BaseURL, "/")
	}
	return strings.TrimRight(a.spec.BaseURL, "/")
}

func (a Adapter) searchURL(cfg providercore.Config, sr providercore.SearchRequest) (string, error) {
	base, err := url.Parse(a.baseURL(cfg))
	if err != nil {
		return "", err
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/search"
	q := base.Query()
	q.Set("q", sr.Title)
	if sr.Year != nil {
		q.Set("year", fmt.Sprint(*sr.Year))
	}
	if sr.SeasonNumber != nil {
		q.Set("season", fmt.Sprint(*sr.SeasonNumber))
	}
	if sr.EpisodeNumber != nil {
		q.Set("episode", fmt.Sprint(*sr.EpisodeNumber))
	}
	if sr.LanguageID != "" {
		q.Set("language", sr.LanguageID)
	}
	for provider, id := range sr.MediaContext.ExternalIDs {
		if id != "" {
			q.Set(provider, id)
		}
	}
	base.RawQuery = q.Encode()
	return base.String(), nil
}

func (a Adapter) parseCandidates(pageURL string, body []byte, requestedLanguage string) ([]providercore.Candidate, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var candidates []providercore.Candidate
	doc.Find("[data-subtitle], .subtitle, .result, tr, li").Each(func(_ int, row *goquery.Selection) {
		link := firstDownloadLink(row)
		if link == "" {
			return
		}
		candidates = append(candidates, a.candidateFrom(row, pageURL, link, requestedLanguage))
	})
	if len(candidates) == 0 {
		doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
			href, _ := link.Attr("href")
			if !isDownloadHref(href) {
				return
			}
			candidates = append(candidates, a.candidateFrom(link, pageURL, href, requestedLanguage))
		})
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("%w: no subtitle links found", providercore.ErrProviderBrokenUpstream)
	}
	return candidates, nil
}

func firstDownloadLink(row *goquery.Selection) string {
	var href string
	row.Find("a[href]").EachWithBreak(func(_ int, link *goquery.Selection) bool {
		candidate, _ := link.Attr("href")
		if isDownloadHref(candidate) {
			href = candidate
			return false
		}
		return true
	})
	return href
}

func isDownloadHref(href string) bool {
	lower := strings.ToLower(strings.TrimSpace(href))
	return lower != "" && (strings.Contains(lower, "download") || strings.Contains(lower, "subtitle") || strings.Contains(lower, "subtitles") || strings.HasSuffix(lower, ".srt") || strings.HasSuffix(lower, ".zip"))
}

func (a Adapter) candidateFrom(sel *goquery.Selection, pageURL, href, requestedLanguage string) providercore.Candidate {
	abs := resolveURL(pageURL, href)
	lang := attrOrText(sel, requestedLanguage, "data-lang", "data-language", ".language", ".lang")
	format := attrOrText(sel, "srt", "data-format", ".format")
	release := attrOrText(sel, strings.TrimSpace(sel.Text()), "data-release", "data-title", ".release", ".title", "a")
	return providercore.Candidate{ProviderName: a.spec.Key, LanguageID: lang, Format: format, ReleaseName: strings.Join(strings.Fields(release), " "), SourceURL: abs, SourceRef: pageURL}
}

func attrOrText(sel *goquery.Selection, fallback string, names ...string) string {
	for _, name := range names {
		if strings.HasPrefix(name, "data-") {
			if value, ok := sel.Attr(name); ok && strings.TrimSpace(value) != "" {
				return strings.TrimSpace(value)
			}
			continue
		}
		if text := strings.TrimSpace(sel.Find(name).First().Text()); text != "" {
			return text
		}
	}
	return fallback
}

func resolveURL(pageURL, href string) string {
	base, err := url.Parse(pageURL)
	if err != nil {
		return href
	}
	ref, err := url.Parse(strings.TrimSpace(href))
	if err != nil {
		return href
	}
	return base.ResolveReference(ref).String()
}

func looksArchived(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "zip") || strings.Contains(lower, "rar") || strings.Contains(lower, "gzip") || strings.HasSuffix(lower, ".gz") || strings.HasSuffix(lower, ".xz")
}

func filenameFor(rawURL, fallback string) string {
	parsed, err := url.Parse(rawURL)
	if err == nil && path.Base(parsed.Path) != "." && path.Base(parsed.Path) != "/" {
		return path.Base(parsed.Path)
	}
	if fallback != "" {
		return fallback
	}
	return "subtitle.srt"
}
