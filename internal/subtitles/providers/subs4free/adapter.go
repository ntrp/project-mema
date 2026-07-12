package subs4free

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/htmlutil"
	"media-manager/internal/subtitles/security"
)

const (
	provider = "subs4free"
	baseURL  = "https://www.subs4free.info"
)

var yearRE = regexp.MustCompile(`^\((\d{4})\)$`)

type adapter struct{}

var Adapter providercore.Adapter = adapter{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return htmlutil.Test(ctx, svc, cfg, provider, baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && !strings.EqualFold(sr.MediaType, "movie") {
		return nil, fmt.Errorf("%w: subs4free does not support %s", providercore.ErrProviderPrerequisiteMissing, sr.MediaType)
	}
	root := htmlutil.BaseURL(cfg, baseURL)
	links, err := suggestions(ctx, svc, root, sr.Title, sr.Year)
	if err != nil {
		return nil, err
	}
	if len(links) == 0 {
		links = []string{fmt.Sprintf("/search_report.php?search=%s&searchType=1", url.QueryEscape(sr.Title))}
	}
	var out []providercore.Candidate
	for _, link := range links {
		cands, err := query(ctx, svc, root, htmlutil.Resolve(root, link), sr.LanguageID)
		if err != nil {
			return nil, err
		}
		out = append(out, cands...)
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, cand.SourceURL, nil, true)
	if err != nil {
		return providercore.Download{}, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "download page")
	if err != nil {
		return providercore.Download{}, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	id := attr(doc.Find(`input[name="id"]`).First(), "value")
	if id == "" {
		return providercore.Download{}, fmt.Errorf("%w: no subs4free download id", providercore.ErrProviderBrokenUpstream)
	}
	for _, u := range []string{"https://images.subs4free.info/favicon.ico", "https://www.subs4series.com/includes/anti-block-layover.php?launch=1", "https://www.subs4series.com/includes/anti-block.php"} {
		if r, e := htmlutil.Request(ctx, svc, provider, http.MethodGet, u, nil, false); e == nil {
			_, _ = htmlutil.ReadResponse(r, 1024, "anti-block")
		}
	}
	form := url.Values{"id": {id}, "x": {strconv.Itoa(rand.Intn(10))}, "y": {strconv.Itoa(rand.Intn(10))}}
	resp, err = htmlutil.Request(ctx, svc, provider, http.MethodPost, htmlutil.BaseURL(cfg, baseURL)+"/getSub.php", strings.NewReader(form.Encode()), true)
	if err != nil {
		return providercore.Download{}, err
	}
	archive, err := htmlutil.ReadResponse(resp, htmlutil.MaxDownloadBytes, "download")
	if err != nil {
		return providercore.Download{}, err
	}
	member, err := providercore.ExtractSubtitle(htmlutil.Filename(cand.SourceURL, cand.ReleaseName)+".zip", archive, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: cand.SourceURL}, nil
}

func suggestions(ctx context.Context, svc providercore.Service, root, title string, year *int32) ([]string, error) {
	raw := fmt.Sprintf("%s/search_report.php?search=%s&searchType=1", root, url.QueryEscape(title))
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, raw, nil, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "suggestions")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	want := normalize(title)
	var links []string
	doc.Find(`select[name="Mov_sel"] > option[value]`).Each(func(_ int, opt *goquery.Selection) {
		text := normalize(opt.Text())
		if text == want || (year != nil && text == fmt.Sprintf("%s%d", want, *year)) {
			links = append(links, attr(opt, "value"))
		}
	})
	return links, nil
}

func query(ctx context.Context, svc providercore.Service, root, pageURL, requested string) ([]providercore.Candidate, error) {
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, pageURL, nil, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "movie")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	title := strings.TrimSpace(doc.Find("td#dates_header > table u").First().Text())
	year := ""
	doc.Find("td#dates_header > table div").Contents().Each(func(_ int, s *goquery.Selection) {
		if m := yearRE.FindStringSubmatch(strings.TrimSpace(s.Text())); len(m) == 2 {
			year = m[1]
		}
	})
	var out []providercore.Candidate
	doc.Find(".movie-details").Each(func(_ int, item *goquery.Selection) {
		link := item.Find("a[href]").First()
		href := attr(link, "href")
		if href == "" {
			return
		}
		lang := strings.TrimSuffix(strings.TrimSuffix(classAt(item.Find(".sprite").First(), 1), "gif"), " ")
		lang = mapLang(lang)
		if requested != "" && requested != lang {
			return
		}
		version := strings.Join(strings.Fields(item.Find("span").First().Text()), " ")
		release := strings.TrimSpace(strings.Join([]string{title, year, version}, " "))
		out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: lang, Format: "srt", ReleaseName: release, SourceURL: htmlutil.Resolve(root, href), SourceRef: pageURL})
	})
	return out, nil
}

func attr(s *goquery.Selection, name string) string { v, _ := s.Attr(name); return v }
func normalize(s string) string                     { return strings.ToLower(strings.Join(strings.Fields(s), "")) }
func classAt(s *goquery.Selection, i int) string {
	c, _ := s.Attr("class")
	p := strings.Fields(c)
	if len(p) > i {
		return p[i]
	}
	return ""
}
func mapLang(s string) string {
	if s == "el" || s == "gr" {
		return "ell"
	}
	if s == "en" {
		return "eng"
	}
	return s
}
