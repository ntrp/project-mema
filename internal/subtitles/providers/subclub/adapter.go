package subclub

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/htmlutil"
)

const (
	provider = "subclub"
	baseURL  = "https://www.subclub.eu"
)

var (
	titleRE = regexp.MustCompile(`^(.+?)\s*\((\d{4})\)(?:\s*\[(\d+)x(\d+)\])?\s*$`)
	downRE  = regexp.MustCompile(`down\.php\?id=(\d+)`)
)

type adapter struct{}

var Adapter providercore.Adapter = adapter{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return htmlutil.Test(ctx, svc, cfg, provider, baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	root := htmlutil.BaseURL(cfg, baseURL)
	hits, err := search(ctx, svc, root, sr.Title)
	if err != nil {
		return nil, err
	}
	var out []providercore.Candidate
	for _, h := range hits {
		if !matches(h, sr) {
			continue
		}
		files, _ := archiveFiles(ctx, svc, root, h.id)
		if len(files) == 0 {
			out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: "est", Format: "srt", ReleaseName: h.title, SourceURL: root + "/down.php?id=" + h.id, SourceRef: root + "/down.php?id=" + h.id})
			continue
		}
		for _, f := range files {
			out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: "est", Format: "srt", ReleaseName: f.name, SourceURL: f.url, SourceRef: root + "/down.php?id=" + h.id})
		}
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return htmlutil.Download(ctx, svc, provider, cand.SourceURL, cand.ReleaseName, false)
}

type hit struct {
	id, title             string
	year, season, episode int
	rating                float64
}
type file struct{ name, url string }

func search(ctx context.Context, svc providercore.Service, root, title string) ([]hit, error) {
	url := fmt.Sprintf("%s/jutud.php?otsing=%s", root, title)
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, url, nil, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "search")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	var hits []hit
	doc.Find("table#tale_list tbody > tr").Each(func(_ int, row *goquery.Selection) {
		tds := row.Find("td")
		if tds.Length() < 9 {
			return
		}
		link := tds.Eq(1).Find("a.sc_link[href]").First()
		href := attr(link, "href")
		m := downRE.FindStringSubmatch(href)
		if len(m) != 2 {
			return
		}
		tm := titleRE.FindStringSubmatch(strings.Join(strings.Fields(link.Text()), " "))
		if len(tm) == 0 {
			return
		}
		y, _ := strconv.Atoi(tm[2])
		s, _ := strconv.Atoi(tm[3])
		e, _ := strconv.Atoi(tm[4])
		r, _ := strconv.ParseFloat(strings.ReplaceAll(tds.Eq(7).Text(), ",", "."), 64)
		hits = append(hits, hit{id: m[1], title: tm[1], year: y, season: s, episode: e, rating: r})
	})
	sort.SliceStable(hits, func(i, j int) bool { return hits[i].rating > hits[j].rating })
	return hits, nil
}

func archiveFiles(ctx context.Context, svc providercore.Service, root, id string) ([]file, error) {
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, root+"/subtitles_archivecontent.php?id="+id, nil, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "archive content")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	var files []file
	doc.Find(`a[href*="down.php"]`).Each(func(_ int, a *goquery.Selection) {
		href := attr(a, "href")
		name := strings.TrimSpace(a.Text())
		if strings.Contains(href, "filename=") && isSub(name) {
			files = append(files, file{name, htmlutil.Resolve(root+"/", href)})
		}
	})
	return files, nil
}

func matches(h hit, sr providercore.SearchRequest) bool {
	if sr.LanguageID != "" && sr.LanguageID != "est" {
		return false
	}
	if sr.Year != nil && h.year != int(*sr.Year) {
		return false
	}
	if sr.SeasonNumber != nil && h.season != int(*sr.SeasonNumber) {
		return false
	}
	if sr.EpisodeNumber != nil && h.episode != int(*sr.EpisodeNumber) {
		return false
	}
	return true
}
func isSub(n string) bool {
	l := strings.ToLower(n)
	return strings.HasSuffix(l, ".srt") || strings.HasSuffix(l, ".ass") || strings.HasSuffix(l, ".sub") || strings.HasSuffix(l, ".vtt")
}
func attr(s *goquery.Selection, name string) string { v, _ := s.Attr(name); return v }
