package yifysubtitles

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/sitehtml"
)

const key = "yifysubtitles"
const defaultBaseURL = "https://yifysubtitles.ch"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

var yifyLanguages = map[string]string{"english": "eng", "spanish": "spa", "french": "fra", "german": "deu", "italian": "ita", "portuguese": "por", "dutch": "nld", "polish": "pol", "romanian": "ron", "russian": "rus", "turkish": "tur", "arabic": "ara", "greek": "ell", "hebrew": "heb", "chinese": "zho", "japanese": "jpn", "korean": "kor"}

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
	imdb := firstNonEmpty(sr.MediaContext.ExternalIDs["imdb"], sr.MediaContext.ExternalIDs["imdb_id"])
	if imdb == "" {
		return nil, fmt.Errorf("%w: imdb id required", providercore.ErrProviderPrerequisiteMissing)
	}
	movieURL := strings.TrimRight(sitehtml.BaseURL(cfg, defaultBaseURL), "/") + "/movie-imdb/" + url.PathEscape(imdb)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, movieURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", sitehtml.BaseURL(cfg, defaultBaseURL))
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return nil, err
	}
	return parseMovieRows(doc.Selection, movieURL, sr.LanguageID), nil
}
func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	if strings.TrimSpace(cand.SourceURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: candidate has no detail URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cand.SourceURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	req.Header.Set("Referer", cand.SourceURL)
	doc, _, err := sitehtml.DoHTML(svc, req, key)
	if err != nil {
		return providercore.Download{}, err
	}
	href, ok := doc.Find("a.download-subtitle").First().Attr("href")
	if !ok {
		return providercore.Download{}, fmt.Errorf("%w: download button missing", providercore.ErrProviderBrokenUpstream)
	}
	downloadURL := sitehtml.Resolve(cand.SourceURL, href)
	downloadReq, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	downloadReq.Header.Set("Referer", cand.SourceURL)
	return sitehtml.Download(downloadReq, svc, key, true, cand.ReleaseName+".zip")
}

func parseMovieRows(root *goquery.Selection, pageURL, requested string) []providercore.Candidate {
	var out []providercore.Candidate
	root.Find("table.other-subs tbody tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 5 {
			return
		}
		rating, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(0).Text()))
		language := normalizeLanguage(cells.Eq(1).Text())
		if requested != "" && language != requested {
			return
		}
		link := cells.Eq(2).Find("a[href]").First()
		href, ok := link.Attr("href")
		if !ok {
			return
		}
		release := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(cells.Eq(2).Text()), "subtitle "))
		hi := cells.Eq(3).Find("span.hi-subtitle").Length() > 0
		ref := sitehtml.Resolve(pageURL, href)
		out = append(out, providercore.Candidate{ProviderName: key, LanguageID: language, Format: "srt", ReleaseName: release, DownloadCount: rating, SourceURL: ref, SourceRef: fmt.Sprintf("hi=%t", hi)})
	})
	sort.SliceStable(out, func(i, j int) bool { return out[i].DownloadCount > out[j].DownloadCount })
	return out
}
func normalizeLanguage(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	if mapped := yifyLanguages[value]; mapped != "" {
		return mapped
	}
	return value
}
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
