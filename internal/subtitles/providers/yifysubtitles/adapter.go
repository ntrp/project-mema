package yifysubtitles

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/sitehtml"
)

const key = "yifysubtitles"
const defaultBaseURL = "https://yifysubtitles.ch"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sitehtml.BaseURL(cfg, defaultBaseURL), nil)
	if err != nil {
		return err
	}
	return sitehtml.Test(req, svc, key)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if !sitehtml.Supports(sr.MediaType, "movie") {
		return nil, sitehtml.Unsupported(key, sr.MediaType)
	}
	imdb := sr.MediaContext.ExternalIDs["imdb"]
	if imdb == "" {
		imdb = sr.MediaContext.ExternalIDs["imdb_id"]
	}
	if imdb == "" {
		return nil, fmt.Errorf("%w: imdb id required", providercore.ErrProviderPrerequisiteMissing)
	}
	movieURL := sitehtml.BaseURL(cfg, defaultBaseURL) + "/movie-imdb/" + url.PathEscape(imdb)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, movieURL, nil)
	if err != nil {
		return nil, err
	}
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return nil, err
	}
	return parseMoviePage(ctx, svc, doc.Selection, movieURL, sr.LanguageID)
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(cand.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: candidate has no source URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cand.SourceURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	return sitehtml.Download(req, svc, key, true, cand.ReleaseName)
}

func parseMoviePage(ctx context.Context, svc providercore.Service, root *goquery.Selection, pageURL, requestedLanguage string) ([]providercore.Candidate, error) {
	var detailLinks []string
	root.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, _ := link.Attr("href")
		if strings.Contains(strings.ToLower(href), "/subtitle/") {
			detailLinks = append(detailLinks, sitehtml.Resolve(pageURL, href))
		}
	})
	seen := map[string]bool{}
	var candidates []providercore.Candidate
	for _, detailURL := range detailLinks {
		if seen[detailURL] {
			continue
		}
		seen[detailURL] = true
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, detailURL, nil)
		if err != nil {
			return nil, err
		}
		doc, _, err := sitehtml.DoHTML(svc, req, key)
		if err != nil {
			return nil, err
		}
		cand := yifyCandidateFromDetail(doc.Selection, detailURL, requestedLanguage)
		if cand.SourceURL != "" {
			candidates = append(candidates, cand)
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("%w: no subtitle links found", providercore.ErrProviderBrokenUpstream)
	}
	return candidates, nil
}

func yifyCandidateFromDetail(root *goquery.Selection, detailURL, requestedLanguage string) providercore.Candidate {
	downloadURL := ""
	root.Find("a[href]").EachWithBreak(func(_ int, link *goquery.Selection) bool {
		href, _ := link.Attr("href")
		lower := strings.ToLower(href + " " + link.Text())
		if strings.Contains(lower, "download") || strings.HasSuffix(lower, ".zip") {
			downloadURL = sitehtml.Resolve(detailURL, href)
			return false
		}
		return true
	})
	if downloadURL == "" {
		return providercore.Candidate{}
	}
	lang := sitehtml.Attr(root, "data-lang", "data-language")
	if lang == "" {
		lang = sitehtml.Text(root, ".sub-lang", ".language")
	}
	if lang == "" {
		lang = requestedLanguage
	}
	release := sitehtml.Attr(root, "data-release", "data-title")
	if release == "" {
		release = sitehtml.Text(root, ".release", ".subtitle-download", "h1")
	}
	return providercore.Candidate{ProviderName: key, LanguageID: lang, Format: "srt", ReleaseName: release, SourceURL: downloadURL, SourceRef: detailURL}
}
