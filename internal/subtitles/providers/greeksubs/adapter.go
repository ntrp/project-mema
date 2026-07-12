package greeksubs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
)

const providerKey = "greeksubs"
const defaultBaseURL = "https://greeksubs.net"
const maxBytes = 2 << 20

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	resp, err := get(ctx, svc, baseURL(cfg)+"/", "", false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return nil
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	imdb := sr.MediaContext.ExternalIDs["imdb"]
	if imdb == "" {
		imdb = sr.MediaContext.ExternalIDs["imdb_id"]
	}
	if imdb == "" {
		return nil, fmt.Errorf("%w: greeksubs requires imdb id", providercore.ErrProviderPrerequisiteMissing)
	}
	search := baseURL(cfg) + "/en/view/" + strings.TrimPrefix(imdb, "tt")
	resp, err := get(ctx, svc, search, "", false)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if sr.SeasonNumber != nil && sr.EpisodeNumber != nil {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
		next := ""
		doc.Find("div.col-lg-offset-2.col-md-8.text-center.top30.bottom10 > a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
			text := a.Text()
			s, e, ok := seasonEpisode(text)
			if ok && s == int(*sr.SeasonNumber) && e == int(*sr.EpisodeNumber) {
				href, _ := a.Attr("href")
				next = href
				return false
			}
			return true
		})
		if next == "" {
			return nil, nil
		}
		resp, err = get(ctx, svc, next, "", false)
		if err != nil {
			return nil, err
		}
		body, err = io.ReadAll(io.LimitReader(resp.Body, maxBytes))
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
	}
	return parse(body, search, sr.LanguageID, baseURL(cfg)), nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	resp, err := get(ctx, svc, cand.SourceURL, cand.SourceRef, true)
	if err != nil {
		return providercore.Download{}, err
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	resp.Body.Close()
	if err != nil {
		return providercore.Download{}, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	data := strings.NewReader(formData(doc))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cand.SourceURL, data)
	if err != nil {
		return providercore.Download{}, err
	}
	req.Header.Set("Referer", cand.SourceURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post, err := svc.DoProviderRequest(req, providerKey, true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer post.Body.Close()
	if post.StatusCode < 200 || post.StatusCode > 299 {
		return providercore.Download{}, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, post.StatusCode)
	}
	content, err := io.ReadAll(io.LimitReader(post.Body, maxBytes))
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: content, URL: cand.SourceURL}, nil
}

func parse(body []byte, ref, lang, base string) []providercore.Candidate {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	sec, _ := doc.Find("input#secCode").Attr("value")
	out := []providercore.Candidate{}
	doc.Find("#elSub > tbody > tr").Each(func(_ int, tr *goquery.Selection) {
		onclick := ""
		tr.Find("a[onclick]").EachWithBreak(func(_ int, a *goquery.Selection) bool {
			onclick, _ = a.Attr("onclick")
			return !strings.Contains(onclick, "downloadMe")
		})
		id := downloadID(onclick)
		if id == "" {
			return
		}
		rowLang := lang
		if alt, ok := tr.Find("img").First().Attr("alt"); ok {
			rowLang = alpha2ToID(alt)
		}
		version := strings.Join(strings.Fields(tr.Text()), " ")
		out = append(out, providercore.Candidate{ProviderName: providerKey, LanguageID: rowLang, Format: "srt", ReleaseName: version, SourceURL: base + "/dll/" + id + "/0/" + sec, SourceRef: ref})
	})
	return out
}
func formData(doc *goquery.Document) string {
	vals := []string{}
	for _, name := range []string{"langcode", "uid", "output", "dll"} {
		if v, ok := doc.Find("input[name='" + name + "']").Attr("value"); ok {
			vals = append(vals, name+"="+v)
		}
	}
	return strings.Join(vals, "&")
}
func get(ctx context.Context, svc providercore.Service, raw, ref string, dl bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return nil, err
	}
	if ref != "" {
		req.Header.Set("Referer", ref)
	}
	resp, err := svc.DoProviderRequest(req, providerKey, dl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 {
		resp.Body.Close()
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil
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
func seasonEpisode(s string) (int, int, bool) {
	re := regexp.MustCompile(`Season (\d+) Episode (\d+)`)
	m := re.FindStringSubmatch(s)
	if len(m) != 3 {
		return 0, 0, false
	}
	var a, b int
	fmt.Sscanf(m[1], "%d", &a)
	fmt.Sscanf(m[2], "%d", &b)
	return a, b, true
}
func downloadID(s string) string {
	re := regexp.MustCompile(`downloadMe\('([^']+)'\)`)
	m := re.FindStringSubmatch(s)
	if len(m) == 2 {
		return m[1]
	}
	return ""
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
