package prijevodionline

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/htmlutil"
)

const (
	provider = "prijevodionline"
	baseURL  = "https://www.prijevodi-online.org"
)

var keyRE = regexp.MustCompile(`epizode\.key\s*=\s*['"]([0-9a-f]{32})['"]`)

type adapter struct{}

var Adapter providercore.Adapter = adapter{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	return htmlutil.Test(ctx, svc, cfg, provider, baseURL)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && !strings.EqualFold(sr.MediaType, "serie") {
		return nil, fmt.Errorf("%w: prijevodionline does not support %s", providercore.ErrProviderPrerequisiteMissing, sr.MediaType)
	}
	if sr.SeasonNumber == nil || sr.EpisodeNumber == nil {
		return nil, fmt.Errorf("%w: season and episode required", providercore.ErrProviderPrerequisiteMissing)
	}
	root := htmlutil.BaseURL(cfg, baseURL)
	seriesID, slug, err := findSeries(ctx, svc, root, sr.Title)
	if err != nil || seriesID == 0 {
		return nil, err
	}
	episodeID, key, seriesURL, err := findEpisode(ctx, svc, root, seriesID, slug, *sr.SeasonNumber, *sr.EpisodeNumber)
	if err != nil || episodeID == 0 {
		return nil, err
	}
	return fetch(ctx, svc, root, seriesURL, episodeID, key, sr)
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	return htmlutil.Download(ctx, svc, provider, cand.SourceURL, cand.ReleaseName, true)
}

func findSeries(ctx context.Context, svc providercore.Service, root, title string) (int, string, error) {
	letter := "0"
	if title != "" && ((title[0] >= 'A' && title[0] <= 'Z') || (title[0] >= 'a' && title[0] <= 'z')) {
		letter = strings.ToLower(title[:1])
	}
	url := root + "/serije/index/" + letter
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, url, nil, false)
	if err != nil {
		return 0, "", err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "search")
	if err != nil {
		return 0, "", err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	want := sanitize(title)
	var id int
	var slug string
	doc.Find(`tr[id^="serija-"] td.naziv > a[href]`).EachWithBreak(func(_ int, a *goquery.Selection) bool {
		if sanitize(a.Text()) != want {
			return true
		}
		parts := strings.Split(attr(a, "href"), "/")
		if len(parts) >= 5 {
			id, _ = strconv.Atoi(parts[3])
			slug = parts[4]
		}
		return false
	})
	return id, slug, nil
}

func findEpisode(ctx context.Context, svc providercore.Service, root string, id int, slug string, season, episode int32) (int, string, string, error) {
	seriesURL := fmt.Sprintf("%s/serije/view/%d/%s", root, id, slug)
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodGet, seriesURL, nil, false)
	if err != nil {
		return 0, "", "", err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "series")
	if err != nil {
		return 0, "", "", err
	}
	key := ""
	if m := keyRE.FindSubmatch(body); len(m) == 2 {
		key = string(m[1])
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	seasonHead := doc.Find(fmt.Sprintf(`#epizode h3#sezona-%d`, season)).First()
	var epID int
	seasonHead.NextAll().EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if goquery.NodeName(s) == "h3" {
			return false
		}
		if goquery.NodeName(s) != "div" || !strings.HasPrefix(attr(s, "id"), "epizoda-") {
			return true
		}
		num, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(s.Find("li.broj").Text()), "."))
		if int32(num) == episode {
			epID, _ = strconv.Atoi(strings.TrimPrefix(attr(s, "id"), "epizoda-"))
			return false
		}
		return true
	})
	return epID, key, seriesURL, nil
}

func fetch(ctx context.Context, svc providercore.Service, root, ref string, episodeID int, key string, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	form := strings.NewReader("key=" + key)
	resp, err := htmlutil.Request(ctx, svc, provider, http.MethodPost, fmt.Sprintf("%s/prijevod/get/%d", root, episodeID), form, false)
	if err != nil {
		return nil, err
	}
	body, err := htmlutil.ReadResponse(resp, htmlutil.MaxHTMLBytes, "subtitles")
	if err != nil {
		return nil, err
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	var out []providercore.Candidate
	doc.Find(`tr[id^="prijevod-"]`).Each(func(_ int, row *goquery.Selection) {
		id := attr(row, "id")
		if strings.Contains(id, "opis") {
			return
		}
		link := row.Find("td.naziv a[href]").First()
		href := attr(link, "href")
		lang := langFromHref(href)
		if lang == "" || (sr.LanguageID != "" && sr.LanguageID != "hbs" && lang != sr.LanguageID) {
			return
		}
		if sr.LanguageID == "hbs" {
			lang = "hbs"
		}
		release := strings.Join(strings.Fields(doc.Find("#prijevod-opis-"+strings.TrimPrefix(id, "prijevod-")).Text()), " ")
		out = append(out, providercore.Candidate{ProviderName: provider, LanguageID: lang, Format: "srt", ReleaseName: release, SourceURL: htmlutil.Resolve(root, href), SourceRef: ref})
	})
	return out, nil
}

func attr(s *goquery.Selection, name string) string { v, _ := s.Attr(name); return v }
func sanitize(s string) string                      { return strings.ToLower(strings.Join(strings.Fields(s), "")) }
func langFromHref(h string) string {
	l := strings.ToLower(h)
	if strings.HasSuffix(l, "-hr") {
		return "hrv"
	}
	if strings.HasSuffix(l, "-sr") {
		return "srp"
	}
	if strings.HasSuffix(l, "-cg") {
		return "mne"
	}
	return ""
}
