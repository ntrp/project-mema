package subtitulamostv

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/nativeutil"
)

const baseURL = "https://www.subtitulamos.tv"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return nativeutil.Test(ctx, svc, cfg, "subtitulamostv", baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && sr.MediaType != "serie" {
		return nil, fmt.Errorf("%w: unsupported media type", providercore.ErrProviderPrerequisiteMissing)
	}
	if sr.SeasonNumber == nil || sr.EpisodeNumber == nil {
		return nil, fmt.Errorf("%w: season and episode are required", providercore.ErrProviderPrerequisiteMissing)
	}
	data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subtitulamostv", BaseURL: baseURL, Path: "/search/query", Form: url.Values{"q": {sr.Title}}})
	if err != nil {
		return nil, err
	}
	var shows []show
	if json.Unmarshal(data, &shows) != nil {
		return nil, nil
	}
	var out []providercore.Candidate
	for _, show := range shows {
		if !strings.EqualFold(show.Name, sr.Title) {
			continue
		}
		out = append(out, fetchEpisode(ctx, svc, cfg, show.ID, *sr.SeasonNumber, *sr.EpisodeNumber, sr.LanguageID)...)
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subtitulamostv", BaseURL: baseURL, Path: cand.SourceURL, Download: true})
	return providercore.Download{Content: data, URL: nativeutil.Absolute(cfg, baseURL, cand.SourceURL)}, err
}

type show struct {
	ID   int    `json:"show_id"`
	Name string `json:"show_name"`
}

func fetchEpisode(ctx context.Context, svc providercore.Service, cfg providercore.Config, showID int, season, episode int32, fallback string) []providercore.Candidate {
	path := "/shows/" + strconv.Itoa(showID)
	for i := 0; i < 3; i++ {
		data, _, err := nativeutil.Do(ctx, svc, cfg, nativeutil.RequestSpec{Provider: "subtitulamostv", BaseURL: baseURL, Path: path})
		if err != nil {
			return nil
		}
		doc, err := nativeutil.Document(data)
		if err != nil {
			return nil
		}
		if next := choice(doc.Selection, "#season-choices a", season); next != "" && next != path {
			path = next
			continue
		}
		if next := choice(doc.Selection, "#episode-choices a", episode); next != "" && next != path {
			path = next
			continue
		}
		return parseEpisode(doc.Selection, fallback)
	}
	return nil
}

func choice(root *goquery.Selection, selector string, want int32) string {
	var href string
	root.Find(selector).EachWithBreak(func(_ int, a *goquery.Selection) bool {
		if strings.TrimSpace(a.Text()) != strconv.Itoa(int(want)) {
			return true
		}
		if class, _ := a.Attr("class"); strings.Contains(class, "selected") {
			href = ""
			return false
		}
		href, _ = a.Attr("href")
		return false
	})
	return href
}

func parseEpisode(root *goquery.Selection, fallback string) []providercore.Candidate {
	out := []providercore.Candidate{}
	root.Find("div.language-container").Each(func(_ int, lang *goquery.Selection) {
		language := nativeutil.Lang(fallback, nativeutil.FirstText(lang, "div.language-name"))
		lang.Find("div.version-container").Each(func(_ int, rel *goquery.Selection) {
			link := nativeutil.Attr(rel, `a[href*="/download"]`, "href")
			if link == "" {
				return
			}
			release := strings.TrimSpace(rel.Find("p").Eq(1).Text())
			if release == "" {
				release = nativeutil.FirstText(rel, "p")
			}
			out = append(out, providercore.Candidate{ProviderName: "subtitulamostv", LanguageID: language, Format: "srt", ReleaseName: release, SourceURL: link})
		})
	})
	return out
}
