package soustitreseu

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/htmlutil"
)

const (
	provider = "soustitreseu"
	baseURL  = "https://www.sous-titres.eu"
)

type adapter struct{}

var Adapter providercore.Adapter = adapter{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return htmlutil.Test(ctx, svc, cfg, provider, baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	root := htmlutil.BaseURL(cfg, baseURL)
	searchURL := root + "/search.html?q=" + url.QueryEscape(sr.Title)
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, searchURL, nil, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "search")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	selector := ".film > h3 > a[href]"
	folder := "films"
	if sr.MediaType == "serie" || sr.SeasonNumber != nil {
		selector = ".serie > h3 > a[href]"
		folder = "series"
	}
	var pages []string
	doc.Find(selector).Each(func(_ int, a *goquery.Selection) {
		if strings.Contains(strings.ToLower(a.Text()), strings.ToLower(sr.Title)) {
			pages = append(pages, htmlutil.Resolve(root+"/", attr(a, "href")))
		}
	})
	var out []providercore.Candidate
	for _, page := range pages {
		c, err := pageCandidates(ctx, svc, root, page, folder, sr)
		if err != nil {
			return nil, err
		}
		out = append(out, c...)
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return htmlutil.Download(ctx, svc, provider, cand.SourceURL, cand.ReleaseName, true)
}

func pageCandidates(ctx context.Context, svc providercore.Service, root, page, folder string, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, page, nil, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "title page")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	var out []providercore.Candidate
	doc.Find("a.subList[href]").Each(func(_ int, a *goquery.Selection) {
		release := strings.TrimSuffix(pathBase(attr(a, "href")), ".zip")
		if sr.SeasonNumber != nil && !episodeOK(a, release, sr) {
			return
		}
		lang := "fre"
		lowerRelease := strings.ToLower(release)
		if strings.Contains(lowerRelease, ".en.") || strings.HasSuffix(lowerRelease, ".en") {
			lang = "eng"
		}
		if sr.LanguageID != "" && lang != sr.LanguageID {
			return
		}
		out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: lang, Format: "srt", ReleaseName: release, SourceURL: root + "/" + folder + "/" + attr(a, "href"), SourceRef: page})
	})
	return out, nil
}

func episodeOK(a *goquery.Selection, release string, sr providercore.SearchRequest) bool {
	txt := strings.TrimSpace(a.Find(".episodenum").Text())
	if strings.Contains(txt, "×") {
		parts := strings.Split(txt, "×")
		s, _ := strconv.Atoi(parts[0])
		e, _ := strconv.Atoi(parts[1])
		return sr.SeasonNumber != nil && int32(s) == *sr.SeasonNumber && (sr.EpisodeNumber == nil || int32(e) == *sr.EpisodeNumber)
	}
	if strings.HasPrefix(txt, "S") {
		s, _ := strconv.Atoi(strings.TrimPrefix(txt, "S"))
		return sr.SeasonNumber != nil && int32(s) == *sr.SeasonNumber
	}
	return strings.Contains(strings.ToLower(release), fmt.Sprintf("s%02d", *sr.SeasonNumber))
}

func pathBase(p string) string {
	parts := strings.Split(strings.TrimRight(p, "/"), "/")
	if len(parts) == 0 {
		return p
	}
	u, _ := url.QueryUnescape(parts[len(parts)-1])
	return u
}
func attr(s *goquery.Selection, name string) string { v, _ := s.Attr(name); return v }
