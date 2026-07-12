package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"

	"github.com/PuerkitoBio/goquery"
)

type nativeCProvider struct {
	key      string
	baseURL  string
	captcha  bool
	search   func(providercore.SearchRequest) (string, url.Values, string)
	parse    func([]byte, string, string) ([]providercore.Candidate, error)
	download func(providercore.Candidate) (string, url.Values, string)
}

func init() {
	for _, adapter := range []nativeCProvider{karagargaAdapter(), ktuvitAdapter(), legendasdivxAdapter(), legendasnetAdapter()} {
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
	path, form, method := p.search(req)
	data, err := p.do(ctx, svc, cfg, method, path, form, "", false)
	if err != nil {
		return nil, err
	}
	return p.parse(data, p.absolute(cfg, path), req.LanguageID)
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
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	resp, err := svc.DoProviderRequest(req, p.key, download)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 50<<20 {
		return nil, fmt.Errorf("%w: response too large", security.ErrUnsafeArchive)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	return data, nil
}

func (p nativeCProvider) absolute(cfg providercore.Config, raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	base, _ := url.Parse(strings.TrimRight(providercore.NewConfig(cfg).BaseURL(p.baseURL), "/") + "/")
	ref, _ := url.Parse(strings.TrimLeft(raw, "/"))
	if strings.HasPrefix(raw, "/") {
		ref, _ = url.Parse(raw)
	}
	return base.ResolveReference(ref).String()
}

func karagargaAdapter() nativeCProvider {
	return nativeCProvider{key: "karagarga", baseURL: "https://karagarga.in", search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/browse.php", url.Values{"search": {nativeQuery(req)}, "search_type": {"title"}}, http.MethodGet
	}, parse: parseNativeHTML("karagarga", "tr", "a[href*='download.php'], a[href*='details.php']"), download: sourceDownload}
}

func ktuvitAdapter() nativeCProvider {
	return nativeCProvider{key: "ktuvit", baseURL: "https://www.ktuvit.me", captcha: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/Services/GetModuleAjax.ashx", url.Values{"moduleName": {"SubtitlesList"}, "SeriesName": {nativeQuery(req)}, "FilmName": {nativeQuery(req)}, "lang": {req.LanguageID}}, http.MethodPost
	}, parse: parseKtuvit, download: func(c providercore.Candidate) (string, url.Values, string) {
		return firstNonEmpty(c.SourceURL, "/Services/DownloadFile.ashx"), url.Values{"subtitleID": {strconv.FormatInt(c.FileID, 10)}}, http.MethodPost
	}}
}

func legendasdivxAdapter() nativeCProvider {
	return nativeCProvider{key: "legendasdivx", baseURL: "https://www.legendasdivx.pt", search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/modules.php", url.Values{"name": {"Downloads"}, "op": {"search"}, "query": {nativeQuery(req)}}, http.MethodGet
	}, parse: parseNativeHTML("legendasdivx", "tr, .download, .subtitle", "a[href*='d_op=getit'], a[href*='download']"), download: sourceDownload}
}

func legendasnetAdapter() nativeCProvider {
	return nativeCProvider{key: "legendasnet", baseURL: "https://legendas.net", captcha: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/search", url.Values{"q": {nativeQuery(req)}, "language": {req.LanguageID}}, http.MethodGet
	}, parse: parseNativeHTML("legendasnet", "tr, .subtitle, .item", "a[href*='download'], a[href*='/download/']"), download: sourceDownload}
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
			href, _ := link.Attr("href")
			if strings.TrimSpace(href) == "" {
				return
			}
			title := firstNonEmpty(attr(row, "data-release"), strings.TrimSpace(link.Text()), strings.TrimSpace(row.Find("td").First().Text()))
			lang := firstNonEmpty(attr(row, "data-language"), fallback)
			abs := resolveAgainst(pageURL, href)
			key := title + "\x00" + abs
			if seen[key] {
				return
			}
			seen[key] = true
			out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: lang, Format: "srt", ReleaseName: title, SourceURL: abs, SourceRef: href})
		})
		return out, nil
	}
}

func parseKtuvit(data []byte, pageURL, fallback string) ([]providercore.Candidate, error) {
	var payload struct {
		Subtitles []struct {
			ID                                    int64
			Name, FileName, Language, DownloadURL string
		}
	}
	if json.Unmarshal(data, &payload) == nil && len(payload.Subtitles) > 0 {
		out := make([]providercore.Candidate, 0, len(payload.Subtitles))
		for _, sub := range payload.Subtitles {
			out = append(out, providercore.Candidate{ProviderName: "ktuvit", FileID: sub.ID, LanguageID: firstNonEmpty(sub.Language, fallback), Format: "srt", ReleaseName: firstNonEmpty(sub.FileName, sub.Name), SourceURL: firstNonEmpty(sub.DownloadURL, "/Services/DownloadFile.ashx")})
		}
		return out, nil
	}
	return parseNativeHTML("ktuvit", "tr, .subtitle", "a[href*='DownloadFile'], a[href*='download']")(data, pageURL, fallback)
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
