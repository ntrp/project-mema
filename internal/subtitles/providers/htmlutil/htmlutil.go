package htmlutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const MaxHTMLBytes = 2 << 20
const MaxDownloadBytes = 50 << 20

func BaseURL(cfg providercore.Config, fallback string) string {
	return providercore.NewConfig(cfg).BaseURL(fallback)
}

func ReadResponse(resp *http.Response, max int64, what string) ([]byte, error) {
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("%w: rate limited", providercore.ErrProviderBrokenUpstream)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%w: %s http status %d", providercore.ErrProviderBrokenUpstream, what, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, max+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > max {
		return nil, fmt.Errorf("%w: %s response too large", providercore.ErrProviderBrokenUpstream, what)
	}
	return body, nil
}

func Request(ctx context.Context, svc providercore.Service, provider, method, rawURL string, body io.Reader, isDownload bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	return svc.DoProviderRequest(req, provider, isDownload)
}

func Test(ctx context.Context, svc providercore.Service, cfg providercore.Config, provider, fallback string) error {
	resp, err := Request(ctx, svc, provider, http.MethodGet, BaseURL(cfg, fallback), nil, false)
	if err != nil {
		return err
	}
	_, err = ReadResponse(resp, 1024, "test")
	return err
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

func WithQuery(baseURL, page string, values map[string]string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + page
	q := parsed.Query()
	for key, value := range values {
		if strings.TrimSpace(value) != "" {
			q.Set(key, value)
		}
	}
	parsed.RawQuery = q.Encode()
	return parsed.String(), nil
}

func Download(ctx context.Context, svc providercore.Service, provider, rawURL, release string, archiveRequired bool) (providercore.Download, error) {
	resp, err := Request(ctx, svc, provider, http.MethodGet, rawURL, nil, true)
	if err != nil {
		return providercore.Download{}, err
	}
	body, err := ReadResponse(resp, MaxDownloadBytes, "download")
	if err != nil {
		return providercore.Download{}, err
	}
	content := body
	if archiveRequired || LooksArchive(rawURL) || LooksArchive(resp.Header.Get("Content-Type")) {
		member, err := providercore.ExtractSubtitle(Filename(rawURL, release), body, security.ArchiveLimits{})
		if err != nil {
			return providercore.Download{}, err
		}
		content = member.Content
	}
	return providercore.Download{Content: content, URL: rawURL}, nil
}

func LooksArchive(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "zip") || strings.Contains(lower, "rar") || strings.Contains(lower, "gzip") || strings.HasSuffix(lower, ".gz") || strings.HasSuffix(lower, ".xz")
}

func Filename(rawURL, fallback string) string {
	parsed, err := url.Parse(rawURL)
	if err == nil {
		base := path.Base(parsed.Path)
		if base != "." && base != "/" && base != "" {
			return base
		}
	}
	if fallback != "" {
		return fallback
	}
	return "subtitle.srt"
}
