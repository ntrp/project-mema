package subf2m

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/htmlutil"
)

const (
	provider = "subf2m"
	baseURL  = "https://subf2m.co"
)

var (
	movieTitleRE = regexp.MustCompile(`^(.+?)(\s+\((\d{4})\))?$`)
	tvTitleRE    = regexp.MustCompile(`^(.+?)\s+[-\(]\s?(.*?)\s+(season|series)\)?(\s+\((\d{4})\))?$`)
	episodeRE    = regexp.MustCompile(`(?i)(?:season|s)\s*?(\d{1,2})\s?[-−]\s?(\d{1,2})|s(\d{1,2})e(\d{1,2})`)
	languages    = map[string]string{"english": "eng", "french": "fre", "german": "ger", "spanish": "spa", "brazillian-portuguese": "por", "farsi_persian": "per", "arabic": "ara", "greek": "gre", "turkish": "tur"}
)

type adapter struct{}

var Adapter providercore.Adapter = adapter{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return htmlutil.Test(ctx, svc, cfg, provider, baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	root := htmlutil.BaseURL(cfg, baseURL)
	paths, err := searchPaths(ctx, svc, root, sr)
	if err != nil || len(paths) == 0 {
		return nil, err
	}
	langPaths := wantedLangPaths(sr.LanguageID)
	var out []providercore.Candidate
	for _, p := range paths {
		for _, lp := range langPaths {
			c, err := pageCandidates(ctx, svc, root, p, lp, sr)
			if err != nil {
				return nil, err
			}
			out = append(out, c...)
		}
		if len(out) > 0 {
			break
		}
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, cand.SourceURL, nil, false)
	if err != nil {
		return providercore.Download{}, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "download page")
	if err != nil {
		return providercore.Download{}, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	href := attr(doc.Find(`#downloadButton[href]`).First(), "href")
	if href == "" {
		return providercore.Download{}, fmt.Errorf("%w: couldn't get download url", providercore.ErrProviderBrokenUpstream)
	}
	return htmlutil.Download(ctx, svc, provider, htmlutil.Resolve(htmlutil.BaseURL(cfg, baseURL), href), cand.ReleaseName, true)
}

func searchPaths(ctx context.Context, svc providercore.Service, root string, sr providercore.SearchRequest) ([]string, error) {
	query := url.QueryEscape(strings.ToLower(sr.Title))
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, root+"/subtitles/searchbytitle?query="+query+"&l=", nil, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "search")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	type hit struct {
		href  string
		score int
	}
	var hits []hit
	doc.Find(`li div.title a[href]`).Each(func(_ int, a *goquery.Selection) {
		text := strings.ToLower(strings.Join(strings.Fields(a.Text()), " "))
		score := 0
		if sr.MediaType == "serie" || sr.SeasonNumber != nil {
			m := tvTitleRE.FindStringSubmatch(text)
			if len(m) == 0 {
				return
			}
			if sr.SeasonNumber != nil && (m[2] == strconv.Itoa(int(*sr.SeasonNumber)) || strings.Contains(m[2], "complete")) {
				score += 10
			}
		} else {
			m := movieTitleRE.FindStringSubmatch(text)
			if len(m) == 0 {
				return
			}
			if sr.Year != nil && m[3] == strconv.Itoa(int(*sr.Year)) {
				score += 10
			}
		}
		hits = append(hits, hit{attr(a, "href"), score})
	})
	sort.SliceStable(hits, func(i, j int) bool { return hits[i].score > hits[j].score })
	var paths []string
	for i, h := range hits {
		if i == 3 {
			break
		}
		paths = append(paths, h.href)
	}
	return paths, nil
}

func pageCandidates(ctx context.Context, svc providercore.Service, root, path, langPath string, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	lang := languages[langPath]
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, root+path+"/"+langPath, nil, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "subtitle page")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if imdb := sr.MediaContext.ExternalIDs["imdb"]; imdb != "" && !strings.Contains(doc.Text(), imdb) {
		return nil, nil
	}
	var out []providercore.Candidate
	doc.Find("li.item").Each(func(_ int, item *goquery.Selection) {
		release := strings.Join(strings.Fields(item.Text()), " ")
		if sr.SeasonNumber != nil && sr.EpisodeNumber != nil && !episodeMatches(release, *sr.SeasonNumber, *sr.EpisodeNumber) {
			return
		}
		href := attr(item.Find("a.download.icon-download[href]").First(), "href")
		if href == "" {
			return
		}
		out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: lang, Format: "srt", ReleaseName: release, SourceURL: htmlutil.Resolve(root, href), SourceRef: root + path + "/" + langPath})
	})
	return out, nil
}

func wantedLangPaths(lang string) []string {
	if lang == "" {
		out := make([]string, 0, len(languages))
		for k := range languages {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}
	for path, id := range languages {
		if id == lang {
			return []string{path}
		}
	}
	return nil
}
func episodeMatches(release string, season, episode int32) bool {
	m := episodeRE.FindStringSubmatch(release)
	if len(m) == 0 {
		return strings.Contains(strings.ToLower(release), "complete")
	}
	vals := m[1:]
	s, e := 0, 0
	if vals[0] != "" {
		s, _ = strconv.Atoi(vals[0])
		e, _ = strconv.Atoi(vals[1])
	} else {
		s, _ = strconv.Atoi(vals[2])
		e, _ = strconv.Atoi(vals[3])
	}
	return int32(s) == season && int32(e) == episode
}
func attr(s *goquery.Selection, name string) string { v, _ := s.Attr(name); return v }
