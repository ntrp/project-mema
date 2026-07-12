package subsunacs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/nativeutil"
)

const baseURL = "https://subsunacs.net"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return nativeutil.Test(ctx, svc, cfg, "subsunacs", baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && sr.MediaType != "movie" && sr.MediaType != "serie" {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	form := url.Values{"m": {queryTitle(sr)}, "l": {"0"}, "c": {""}, "y": {""}, "action": {"   Търси   "}, "a": {""}, "d": {""}, "u": {""}, "g": {""}, "t": {""}, "imdbcheck": {"1"}}
	if sr.Year != nil {
		form.Set("y", strconv.Itoa(int(*sr.Year)))
	}
	if sr.LanguageID == "eng" {
		form.Set("l", "1")
	}
	data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subsunacs", BaseURL: baseURL, Method: http.MethodPost, Path: "/search.php", Form: form, Headers: map[string]string{"Referer": nativeutil.Absolute(cfg, baseURL, "/index.php")}})
	if err != nil {
		return nil, err
	}
	return parse(data, sr.LanguageID, cfg), nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return nativeutil.DownloadSubtitle(ctx, svc, cfg, "subsunacs", baseURL, cand)
}

func parse(data []byte, fallback string, cfg providercore.Config) []providercore.Candidate {
	doc, err := nativeutil.Document(data)
	if err != nil {
		return nil
	}
	out := []providercore.Candidate{}
	doc.Find("tr[onmouseover]").EachWithBreak(func(i int, row *goquery.Selection) bool {
		if i >= 20 {
			return false
		}
		link := nativeutil.Attr(row, "td.tdMovie a.tooltip", "href")
		if link == "" {
			return true
		}
		name := nativeutil.FirstText(row, "td.tdMovie a.tooltip")
		downloads, _ := strconv.Atoi(strings.TrimSpace(row.Find("td").Eq(4).Text()))
		out = append(out, providercore.Candidate{ProviderName: "subsunacs", LanguageID: nativeutil.Lang(fallback, ""), Format: nativeutil.Format(link), ReleaseName: name, DownloadCount: downloads, SourceURL: nativeutil.Absolute(cfg, baseURL, link)})
		return true
	})
	return out
}

func queryTitle(sr providercore.SearchRequest) string {
	if sr.SeasonNumber != nil && sr.EpisodeNumber != nil {
		return fmt.Sprintf("%s %02d %02d", clean(sr.Title), *sr.SeasonNumber, *sr.EpisodeNumber)
	}
	return clean(sr.Title)
}

func clean(value string) string {
	return regexp.MustCompile(`[^\pL\pN' ]+`).ReplaceAllString(value, " ")
}
