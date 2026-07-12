package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"

	"github.com/PuerkitoBio/goquery"
)

type clusterCProvider struct {
	key        string
	baseURL    string
	private    bool
	captcha    bool
	searchPath string
}

var clusterCProviders = []clusterCProvider{
	{key: "addic7ed", baseURL: "https://www.addic7ed.com", captcha: true, searchPath: "/search.php"},
	{key: "avistaz", baseURL: "https://avistaz.to", private: true, searchPath: "/subtitles"},
	{key: "cinemaz", baseURL: "https://cinemaz.to", private: true, searchPath: "/subtitles"},
	{key: "hdbits", baseURL: "https://hdbits.org", private: true, searchPath: "/browse.php"},
}

func init() {
	for _, provider := range clusterCProviders {
		Register(provider.key, provider)
	}
}

func (p clusterCProvider) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	if err := p.checkPrerequisites(config); err != nil {
		return err
	}
	request, err := p.newRequest(ctx, config, http.MethodGet, p.base(config), nil, false)
	if err != nil {
		return err
	}
	response, err := service.DoProviderRequest(request, p.key, false)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 399 {
		return fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, response.StatusCode)
	}
	return nil
}

func (p clusterCProvider) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	if err := p.checkPrerequisites(config); err != nil {
		return nil, err
	}
	searchURL, err := p.searchURL(config, request)
	if err != nil {
		return nil, err
	}
	httpRequest, err := p.newRequest(ctx, config, http.MethodGet, searchURL, nil, false)
	if err != nil {
		return nil, err
	}
	response, err := service.DoProviderRequest(httpRequest, p.key, false)
	if err != nil {
		return nil, err
	}
	data, err := readLimited(response, 4<<20)
	if err != nil {
		return nil, err
	}
	return p.parseCandidates(data, searchURL, request.LanguageID)
}

func (p clusterCProvider) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	if err := p.checkPrerequisites(config); err != nil {
		return providercore.Download{}, err
	}
	downloadURL := strings.TrimSpace(candidate.SourceURL)
	if downloadURL == "" {
		downloadURL = strings.TrimSpace(candidate.SourceRef)
	}
	if downloadURL == "" {
		return providercore.Download{}, providercore.ErrProviderPrerequisiteMissing
	}
	absolute, err := p.resolveURL(config, downloadURL)
	if err != nil {
		return providercore.Download{}, err
	}
	request, err := p.newRequest(ctx, config, http.MethodGet, absolute, nil, true)
	if err != nil {
		return providercore.Download{}, err
	}
	response, err := service.DoProviderRequest(request, p.key, true)
	if err != nil {
		return providercore.Download{}, err
	}
	data, err := readLimited(response, 25<<20)
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: data, URL: absolute}, nil
}

func (p clusterCProvider) checkPrerequisites(config providercore.Config) error {
	cookies := providercore.NewConfig(config).CookieString()
	if p.captcha && cookies == "" {
		return providercore.ErrCaptchaRequired
	}
	if p.private && cookies == "" {
		return providercore.ErrPrivateMembershipRequired
	}
	return nil
}

func (p clusterCProvider) newRequest(ctx context.Context, config providercore.Config, method, rawURL string, body io.Reader, download bool) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	view := providercore.NewConfig(config)
	request.Header.Set("User-Agent", userAgent(view))
	if cookies := view.CookieString(); cookies != "" {
		request.Header.Set("Cookie", cookies)
	}
	if download {
		request.Header.Set("Accept", "application/zip,application/x-rar-compressed,application/octet-stream,text/plain,*/*")
	}
	return request, nil
}

func (p clusterCProvider) searchURL(config providercore.Config, request providercore.SearchRequest) (string, error) {
	base, err := url.Parse(p.base(config))
	if err != nil {
		return "", err
	}
	base.Path = path.Join(base.Path, p.searchPath)
	query := base.Query()
	query.Set("q", searchQuery(request))
	if request.LanguageID != "" {
		query.Set("language", request.LanguageID)
	}
	if request.SeasonNumber != nil {
		query.Set("season", strconv.Itoa(int(*request.SeasonNumber)))
	}
	if request.EpisodeNumber != nil {
		query.Set("episode", strconv.Itoa(int(*request.EpisodeNumber)))
	}
	base.RawQuery = query.Encode()
	return base.String(), nil
}

func (p clusterCProvider) parseCandidates(data []byte, pageURL string, fallbackLanguage string) ([]providercore.Candidate, error) {
	document, err := providercore.ParseHTML(data)
	if err != nil {
		return nil, err
	}
	candidates := []providercore.Candidate{}
	seen := map[string]struct{}{}
	document.Find("[data-subtitle], a[href*='download'], a[href*='subtitle'], tr").Each(func(_ int, item *goquery.Selection) {
		if goquery.NodeName(item) == "a" && item.ParentsFiltered("tr,[data-subtitle]").Length() > 0 {
			return
		}
		candidate := p.candidateFromSelection(item, pageURL, fallbackLanguage)
		fingerprint := candidate.ReleaseName + "\x00" + candidate.SourceURL
		if candidate.ReleaseName != "" && candidate.SourceURL != "" {
			if _, ok := seen[fingerprint]; ok {
				return
			}
			seen[fingerprint] = struct{}{}
			candidates = append(candidates, candidate)
		}
	})
	return candidates, nil
}

func (p clusterCProvider) candidateFromSelection(item *goquery.Selection, pageURL string, fallbackLanguage string) providercore.Candidate {
	link := item
	if goquery.NodeName(item) != "a" {
		link = item.Find("a[href*='download'], a[href*='subtitle'], a[href]").First()
	}
	href, _ := link.Attr("href")
	absolute, _ := url.Parse(pageURL)
	resolved := absolute.ResolveReference(mustParseRef(href)).String()
	title := firstNonEmpty(attr(item, "data-release"), attr(link, "data-release"), strings.TrimSpace(link.Text()), strings.TrimSpace(item.Find("td").First().Text()))
	language := firstNonEmpty(attr(item, "data-language"), attr(link, "data-language"), fallbackLanguage)
	format := firstNonEmpty(attr(item, "data-format"), attr(link, "data-format"), "srt")
	return providercore.Candidate{ProviderName: p.key, LanguageID: language, Format: format, ReleaseName: title, SourceURL: resolved, SourceRef: href}
}

func (p clusterCProvider) resolveURL(config providercore.Config, ref string) (string, error) {
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		return ref, nil
	}
	base, err := url.Parse(p.base(config))
	if err != nil {
		return "", err
	}
	return base.ResolveReference(mustParseRef(ref)).String(), nil
}

func (p clusterCProvider) base(config providercore.Config) string {
	return providercore.NewConfig(config).BaseURL(p.baseURL)
}

func readLimited(response *http.Response, maxBytes int64) ([]byte, error) {
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("response size limit exceeded")
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, response.StatusCode)
	}
	return data, nil
}

func searchQuery(request providercore.SearchRequest) string {
	parts := []string{request.Title}
	if request.Year != nil {
		parts = append(parts, strconv.Itoa(int(*request.Year)))
	}
	if request.SeasonNumber != nil && request.EpisodeNumber != nil {
		parts = append(parts, fmt.Sprintf("S%02dE%02d", *request.SeasonNumber, *request.EpisodeNumber))
	}
	if request.FilePath != "" {
		parts = append(parts, path.Base(request.FilePath))
	}
	return strings.Join(parts, " ")
}

func mustParseRef(ref string) *url.URL { parsed, _ := url.Parse(strings.TrimSpace(ref)); return parsed }
func attr(item *goquery.Selection, name string) string {
	value, _ := item.Attr(name)
	return strings.TrimSpace(value)
}
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
func userAgent(config providercore.ConfigView) string {
	if value := config.StringSetting("userAgent"); value != "" {
		return value
	}
	return "Mozilla/5.0 (compatible; MemaSubtitleProvider/1.0)"
}
