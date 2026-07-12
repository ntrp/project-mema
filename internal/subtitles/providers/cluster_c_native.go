package providers

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"

	"github.com/PuerkitoBio/goquery"
)

type nativeCProvider struct {
	key         string
	baseURL     string
	captcha     bool
	rawDownload bool
	search      func(providercore.SearchRequest) (string, url.Values, string)
	parse       func([]byte, string, string) ([]providercore.Candidate, error)
	download    func(providercore.Candidate) (string, url.Values, string)
}

var nativeCProviders = []nativeCProvider{
	karagargaAdapter(),
	ktuvitAdapter(),
	legendasdivxAdapter(),
	legendasnetAdapter(),
	napisy24Adapter(),
	pipocasAdapter(),
	subs4seriesAdapter(),
	subscenterAdapter(),
	titloviAdapter(),
}

func init() {
	for _, adapter := range nativeCProviders {
		Register(adapter.key, adapter)
	}
}

func (p nativeCProvider) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	if err := p.prereq(cfg); err != nil {
		return err
	}
	_, err := p.do(ctx, svc, cfg, http.MethodGet, "/", nil, "", false)
	return err
}

func (p nativeCProvider) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, req providercore.SearchRequest) ([]providercore.Candidate, error) {
	if err := p.prereq(cfg); err != nil {
		return nil, err
	}
	rawPath, form, method := p.search(req)
	data, err := p.do(ctx, svc, cfg, method, rawPath, form, "", false)
	if err != nil {
		return nil, err
	}
	return p.parse(data, p.absolute(cfg, rawPath), req.LanguageID)
}

func (p nativeCProvider) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	if err := p.prereq(cfg); err != nil {
		return providercore.Download{}, err
	}
	dlPath, form, method := p.download(cand)
	if strings.TrimSpace(dlPath) == "" {
		return providercore.Download{}, providercore.ErrProviderPrerequisiteMissing
	}
	data, err := p.do(ctx, svc, cfg, method, dlPath, form, acceptArchive(), true)
	if err != nil {
		return providercore.Download{}, err
	}
	if p.rawDownload {
		return providercore.Download{Content: data, URL: p.absolute(cfg, dlPath)}, nil
	}
	member, err := providercore.ExtractSubtitle(path.Base(dlPath), data, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: p.absolute(cfg, dlPath)}, nil
}

func (p nativeCProvider) prereq(cfg providercore.Config) error {
	if providercore.NewConfig(cfg).CookieString() == "" {
		if p.captcha {
			return providercore.ErrCaptchaRequired
		}
		return providercore.ErrPrivateMembershipRequired
	}
	return nil
}

func (p nativeCProvider) do(ctx context.Context, svc providercore.Service, cfg providercore.Config, method, rawPath string, form url.Values, accept string, download bool) ([]byte, error) {
	if method == "" {
		method = http.MethodGet
	}
	rawURL := p.absolute(cfg, rawPath)
	var body io.Reader
	if len(form) > 0 && method == http.MethodGet {
		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for key, values := range form {
			for _, value := range values {
				q.Add(key, value)
			}
		}
		u.RawQuery = q.Encode()
		rawURL = u.String()
	} else if len(form) > 0 {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	view := providercore.NewConfig(cfg)
	req.Header.Set("User-Agent", nativeUserAgent(view))
	req.Header.Set("Cookie", view.CookieString())
	req.Header.Set("Referer", view.BaseURL(p.baseURL)+"/")
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	resp, err := svc.DoProviderRequest(req, p.key, download)
	if err != nil {
		return nil, err
	}
	maxBytes := int64(4 << 20)
	if download {
		maxBytes = 50 << 20
	}
	data, err := readLimited(resp, maxBytes)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p nativeCProvider) absolute(cfg providercore.Config, raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	base, _ := url.Parse(strings.TrimRight(providercore.NewConfig(cfg).BaseURL(p.baseURL), "/") + "/")
	ref, _ := url.Parse(strings.TrimSpace(raw))
	return base.ResolveReference(ref).String()
}

func parseNativeHTML(provider, rowSelector, linkSelector string) func([]byte, string, string) ([]providercore.Candidate, error) {
	return func(data []byte, pageURL, fallback string) ([]providercore.Candidate, error) {
		doc, err := providercore.ParseHTML(data)
		if err != nil {
			return nil, err
		}
		out := []providercore.Candidate{}
		seen := map[string]bool{}
		doc.Find(rowSelector).Each(func(_ int, row *goquery.Selection) {
			link := row.Find(linkSelector).First()
			if link.Length() == 0 {
				link = row.Find("a[href]").First()
			}
			href, _ := link.Attr("href")
			if strings.TrimSpace(href) == "" {
				return
			}
			title := firstNonEmpty(attr(row, "data-release"), strings.TrimSpace(row.Find(".release, .title, .episode, .naziv, td:first-child").First().Text()), strings.TrimSpace(link.Text()))
			if title == "" {
				return
			}
			lang := firstNonEmpty(attr(row, "data-language"), strings.TrimSpace(row.Find(".language, .lang, .jezik").First().Text()), fallback)
			abs := resolveAgainst(pageURL, href)
			key := title + "\x00" + abs
			if seen[key] {
				return
			}
			seen[key] = true
			out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: lang, Format: formatFrom(href), ReleaseName: title, SourceURL: abs, SourceRef: href})
		})
		return out, nil
	}
}

func sourceDownload(c providercore.Candidate) (string, url.Values, string) {
	return firstNonEmpty(c.SourceURL, c.SourceRef), nil, http.MethodGet
}
func acceptArchive() string {
	return "application/zip,application/x-rar-compressed,application/octet-stream,text/plain,*/*"
}
func nativeUserAgent(c providercore.ConfigView) string {
	if v := c.StringSetting("userAgent"); v != "" {
		return v
	}
	return "Mozilla/5.0 (compatible; MemaSubtitleProvider/1.0)"
}
func nativeQuery(req providercore.SearchRequest) string { return searchQuery(req) }
func resolveAgainst(pageURL, href string) string {
	base, _ := url.Parse(pageURL)
	ref, _ := url.Parse(strings.TrimSpace(href))
	return base.ResolveReference(ref).String()
}
func formatFrom(raw string) string {
	ext := strings.TrimPrefix(path.Ext(raw), ".")
	if ext == "" || len(ext) > 5 {
		return "srt"
	}
	return strings.ToLower(ext)
}
