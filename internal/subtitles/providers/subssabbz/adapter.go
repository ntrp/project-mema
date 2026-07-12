package subssabbz

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

const baseURL = "http://subs.sab.bz"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return nativeutil.Test(ctx, svc, cfg, "subssabbz", baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && sr.MediaType != "movie" && sr.MediaType != "serie" {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	form := url.Values{"act": {"search"}, "movie": {title(sr)}, "select-language": {"2"}, "upldr": {""}, "yr": {""}, "release": {""}}
	if sr.Year != nil {
		form.Set("yr", strconv.Itoa(int(*sr.Year)))
	}
	if sr.LanguageID == "eng" {
		form.Set("select-language", "1")
	}
	data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subssabbz", BaseURL: baseURL, Method: http.MethodPost, Path: "/index.php?", Form: form, Headers: map[string]string{"Referer": nativeutil.Absolute(cfg, baseURL, "/")}})
	if err != nil {
		return nil, err
	}
	return parse(data, sr.LanguageID, cfg), nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return nativeutil.DownloadSubtitle(ctx, svc, cfg, "subssabbz", baseURL, cand)
}

func parse(data []byte, fallback string, cfg providercore.Config) []providercore.Candidate {
	doc, err := nativeutil.Document(data)
	if err != nil {
		return nil
	}
	out := []providercore.Candidate{}
	doc.Find("tr.subs-row").EachWithBreak(func(i int, row *goquery.Selection) bool {
		if i >= 25 {
			return false
		}
		link := nativeutil.Attr(row, "td.c2field a", "href")
		if link == "" {
			return true
		}
		name := nativeutil.FirstText(row, "td.c2field a")
		count, _ := strconv.Atoi(strings.TrimSpace(row.Find("td").Eq(6).Text()))
		out = append(out, providercore.Candidate{ProviderName: "subssabbz", LanguageID: nativeutil.Lang(fallback, ""), Format: nativeutil.Format(link), ReleaseName: name, DownloadCount: count, SourceURL: nativeutil.Absolute(cfg, baseURL, link)})
		return true
	})
	return out
}

func title(sr providercore.SearchRequest) string {
	value := sr.Title
	if sr.SeasonNumber != nil && sr.EpisodeNumber != nil {
		value = fmt.Sprintf("%s %02d %02d", sr.Title, *sr.SeasonNumber, *sr.EpisodeNumber)
	}
	return regexp.MustCompile(`[^\pL\pN' ]+`).ReplaceAllString(value, " ")
}
