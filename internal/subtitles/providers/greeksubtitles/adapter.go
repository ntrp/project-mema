package greeksubtitles

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const providerKey = "greeksubtitles"
const defaultBaseURL = "http://gr.greek-subtitles.com"
const downloadBaseURL = "http://www.greeksubtitles.info/getp.php?id="
const maxBytes = 50 << 20

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	resp, err := request(ctx, svc, http.MethodGet, baseURL(cfg)+"/", "", false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return nil
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && !strings.EqualFold(sr.MediaType, "movie") && !strings.EqualFold(sr.MediaType, "serie") {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	query := strings.TrimSpace(sr.Title)
	if sr.SeasonNumber != nil && sr.EpisodeNumber != nil {
		query += fmt.Sprintf(" S%02dE%02d", *sr.SeasonNumber, *sr.EpisodeNumber)
	} else if sr.Year != nil {
		query += fmt.Sprintf(" %04d", *sr.Year)
	}
	page := baseURL(cfg) + "/search.php?name=" + url.QueryEscape(query)
	var out []providercore.Candidate
	for page != "" {
		resp, err := request(ctx, svc, http.MethodGet, page, "", false)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		doc.Find("td.latest_name > a:nth-of-type(1)").Each(func(_ int, a *goquery.Selection) {
			href, _ := a.Attr("href")
			id := idFromHref(href)
			if id == 0 {
				return
			}
			lang := sr.LanguageID
			if img := a.Parent().Find("img").First(); img.Size() > 0 {
				if src, ok := img.Attr("src"); ok {
					lang = alpha2ToID(strings.TrimSuffix(pathBase(src), ".gif"))
				}
			}
			out = append(out, providercore.Candidate{ProviderName: providerKey, LanguageID: lang, Format: "srt", ReleaseName: strings.TrimSpace(a.Text()), FileID: int64(id), SourceURL: downloadBaseURL + strconv.Itoa(id), SourceRef: resolve(page, href)})
		})
		next := ""
		doc.Find("a[href]").EachWithBreak(func(_ int, a *goquery.Selection) bool {
			href, _ := a.Attr("href")
			if strings.Contains(a.Text(), "Next") && strings.Contains(href, "search.php") {
				next = resolve(baseURL(cfg)+"/", href)
				return false
			}
			return true
		})
		page = next
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	resp, err := request(ctx, svc, http.MethodGet, cand.SourceURL, cand.SourceRef, true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return providercore.Download{}, err
	}
	if len(body) > maxBytes {
		return providercore.Download{}, security.ErrUnsafeArchive
	}
	member, err := providercore.ExtractSubtitle("greeksubtitles.zip", body, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: cand.SourceURL}, nil
}

func request(ctx context.Context, svc providercore.Service, method, rawURL, referer string, dl bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return nil, err
	}
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	resp, err := svc.DoProviderRequest(req, providerKey, dl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		resp.Body.Close()
		return nil, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	return resp, nil
}
func baseURL(cfg providercore.Config) string {
	if strings.TrimSpace(cfg.BaseURL) != "" {
		return strings.TrimRight(cfg.BaseURL, "/")
	}
	return defaultBaseURL
}
func idFromHref(h string) int {
	parts := strings.Split(strings.Trim(h, "/"), "/")
	for _, p := range parts {
		if n, err := strconv.Atoi(p); err == nil {
			return n
		}
	}
	return 0
}
func pathBase(s string) string {
	parts := strings.Split(strings.Trim(s, "/"), "/")
	if len(parts) == 0 {
		return s
	}
	return parts[len(parts)-1]
}
func alpha2ToID(a string) string {
	switch strings.ToLower(a) {
	case "gr", "el":
		return "ell"
	case "en":
		return "eng"
	}
	return a
}
func resolve(base, ref string) string {
	b, err := url.Parse(base)
	if err != nil {
		return ref
	}
	r, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return b.ResolveReference(r).String()
}
