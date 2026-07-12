package subtitrarinoi

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

const baseURL = "https://www.subtitrari-noi.ro"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return nativeutil.Test(ctx, svc, cfg, "subtitrarinoi", baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && sr.MediaType != "movie" && sr.MediaType != "serie" {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	q := sr.Title
	if id := imdb(sr); id != "" {
		q = strings.TrimPrefix(id, "tt")
	}
	form := url.Values{"search_q": {"1"}, "tip": {"2"}, "an": {"Toti anii"}, "gen": {"Toate"}, "cautare": {q}, "query_q": {q}}
	data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subtitrarinoi", BaseURL: baseURL, Method: http.MethodPost, Path: "/paginare_filme.php", Form: form, Headers: map[string]string{"X-Requested-With": "XMLHttpRequest", "Referer": nativeutil.Absolute(cfg, baseURL, "/")}})
	if err != nil {
		return nil, err
	}
	return parse(data, cfg), nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return nativeutil.DownloadSubtitle(ctx, svc, cfg, "subtitrarinoi", baseURL, cand)
}

func parse(data []byte, cfg providercore.Config) []providercore.Candidate {
	doc, err := nativeutil.Document(data)
	if err != nil {
		return nil
	}
	comments := doc.Find("div:not([id]):not([class]):not([align])")
	out := []providercore.Candidate{}
	doc.Find("div#round").Each(func(i int, row *goquery.Selection) {
		link := nativeutil.Attr(row, ".buton a", "href")
		if link == "" {
			return
		}
		fullTitle := nativeutil.FirstText(row, "#content-main a")
		release := strings.TrimSpace(strings.Split(fullTitle, "(")[0])
		if note := strings.TrimSpace(comments.Eq(i).Text()); note != "" {
			release = release + " " + strings.Join(strings.Fields(note), " ")
		}
		downloads, _ := strconv.Atoi(regexp.MustCompile(`\D+`).ReplaceAllString(nativeutil.FirstText(row, "#content-right p"), ""))
		out = append(out, providercore.Candidate{ProviderName: "subtitrarinoi", LanguageID: "ron", Format: nativeutil.Format(link), ReleaseName: release, DownloadCount: downloads, SourceURL: nativeutil.Absolute(cfg, baseURL, link)})
	})
	return out
}

func imdb(sr providercore.SearchRequest) string {
	if sr.MediaContext.ExternalIDs != nil {
		return sr.MediaContext.ExternalIDs["imdb"]
	}
	return ""
}
