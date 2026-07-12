package sitehtml

import (
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

const MaxHTMLBytes = 2 << 20
const MaxDownloadBytes = 50 << 20

func BaseURL(cfg providercore.Config, fallback string) string {
	if strings.TrimSpace(cfg.BaseURL) != "" {
		return strings.TrimRight(cfg.BaseURL, "/")
	}
	return strings.TrimRight(fallback, "/")
}

func Test(ctxReq *http.Request, svc providercore.Service, key string) error {
	resp, err := svc.DoProviderRequest(ctxReq, key, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return fmt.Errorf("%s test failed: http status %d", key, resp.StatusCode)
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return nil
}

func DoHTML(svc providercore.Service, req *http.Request, key string) (*goquery.Document, []byte, error) {
	resp, err := svc.DoProviderRequest(req, key, false)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, nil, fmt.Errorf("%w: rate limited", providercore.ErrProviderBrokenUpstream)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxHTMLBytes+1))
	if err != nil {
		return nil, nil, err
	}
	if len(body) > MaxHTMLBytes {
		return nil, nil, fmt.Errorf("%w: response too large", providercore.ErrProviderBrokenUpstream)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	return doc, body, err
}

func Download(ctxReq *http.Request, svc providercore.Service, key string, archive bool, release string) (providercore.Download, error) {
	resp, err := svc.DoProviderRequest(ctxReq, key, true)
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
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxDownloadBytes+1))
	if err != nil {
		return providercore.Download{}, err
	}
	if len(body) > MaxDownloadBytes {
		return providercore.Download{}, fmt.Errorf("%w: download response too large", security.ErrUnsafeArchive)
	}
	content := body
	if archive || LooksArchived(ctxReq.URL.String()) || LooksArchived(resp.Header.Get("Content-Type")) {
		member, err := providercore.ExtractSubtitle(FilenameFor(ctxReq.URL.String(), release), body, security.ArchiveLimits{})
		if err != nil {
			return providercore.Download{}, err
		}
		content = member.Content
	}
	return providercore.Download{Content: content, URL: ctxReq.URL.String()}, nil
}

func Resolve(pageURL, href string) string {
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

func Text(sel *goquery.Selection, selectors ...string) string {
	for _, selector := range selectors {
		if text := strings.Join(strings.Fields(sel.Find(selector).First().Text()), " "); text != "" {
			return text
		}
	}
	return strings.Join(strings.Fields(sel.Text()), " ")
}

func Attr(sel *goquery.Selection, names ...string) string {
	for _, name := range names {
		if value, ok := sel.Attr(name); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func Supports(mediaType string, allowed ...string) bool {
	if mediaType == "" {
		return true
	}
	for _, item := range allowed {
		if strings.EqualFold(mediaType, item) {
			return true
		}
	}
	return false
}

func Unsupported(key, mediaType string) error {
	return fmt.Errorf("%w: %s does not support %s", providercore.ErrProviderPrerequisiteMissing, key, mediaType)
}

func LooksArchived(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "zip") || strings.Contains(lower, "rar") || strings.Contains(lower, "gzip") || strings.HasSuffix(lower, ".gz") || strings.HasSuffix(lower, ".xz")
}

func FilenameFor(rawURL, fallback string) string {
	parsed, err := url.Parse(rawURL)
	if err == nil && path.Base(parsed.Path) != "." && path.Base(parsed.Path) != "/" {
		return path.Base(parsed.Path)
	}
	if fallback != "" {
		return fallback
	}
	return "subtitle.srt"
}
