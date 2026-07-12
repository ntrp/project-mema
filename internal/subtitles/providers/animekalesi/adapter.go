package animekalesi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const providerKey = "animekalesi"
const defaultBaseURL = "https://www.animekalesi.com"
const maxBytes = 50 << 20

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	resp, err := get(ctx, svc, baseURL(cfg)+"/", false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return nil
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && !strings.EqualFold(sr.MediaType, "serie") {
		return nil, fmt.Errorf("%w: animekalesi supports series", providercore.ErrProviderPrerequisiteMissing)
	}
	root := baseURL(cfg)
	resp, err := get(ctx, svc, root+"/tum-anime-serileri.html", false)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	seriesTitle, seriesHref := findSeries(body, sr.Title)
	if seriesHref == "" {
		return nil, nil
	}
	subPage := root + "/" + strings.TrimPrefix(strings.Replace(seriesHref, "bolumler-", "altyazib-", 1), "/")
	resp, err = get(ctx, svc, subPage, false)
	if err != nil {
		return nil, err
	}
	body, err = io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	episodeLink := findEpisode(body, sr.SeasonNumber, sr.EpisodeNumber)
	if episodeLink == "" {
		return nil, nil
	}
	episodeURL := root + "/" + strings.TrimPrefix(episodeLink, "/")
	resp, err = get(ctx, svc, episodeURL, false)
	if err != nil {
		return nil, err
	}
	body, err = io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	dl, uploader := findDownload(body)
	if dl == "" {
		return nil, nil
	}
	season, episode := int32(1), int32(0)
	if sr.SeasonNumber != nil {
		season = *sr.SeasonNumber
	}
	if sr.EpisodeNumber != nil {
		episode = *sr.EpisodeNumber
	}
	release := fmt.Sprintf("%s - S%02dE%02d", seriesTitle, season, episode)
	if uploader != "" {
		release += " by " + uploader
	}
	return []providercore.Candidate{{ProviderName: providerKey, LanguageID: "tur", Format: "srt", ReleaseName: release, SourceURL: root + "/" + strings.TrimPrefix(dl, "/"), SourceRef: episodeURL}}, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	resp, err := get(ctx, svc, cand.SourceURL, true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes+1))
	if err != nil {
		return providercore.Download{}, err
	}
	if len(body) > maxBytes {
		return providercore.Download{}, security.ErrUnsafeArchive
	}
	if bytes.HasPrefix(body, []byte("PK\x03\x04")) {
		member, err := providercore.ExtractSubtitle("animekalesi.zip", body, security.ArchiveLimits{})
		if err != nil {
			return providercore.Download{}, err
		}
		body = member.Content
	}
	return providercore.Download{Content: body, URL: cand.SourceURL}, nil
}

func findSeries(body []byte, title string) (string, string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	want := norm(title)
	bestTitle, bestHref := "", ""
	doc.Find("td#bolumler a[href]").EachWithBreak(func(_ int, a *goquery.Selection) bool {
		text := strings.TrimSpace(a.Text())
		href, _ := a.Attr("href")
		if !strings.Contains(href, "bolumler-") {
			return true
		}
		n := norm(text)
		if n == want {
			bestTitle, bestHref = text, href
			return false
		}
		if bestHref == "" && (strings.Contains(n, want) || strings.Contains(want, n)) {
			bestTitle, bestHref = text, href
		}
		return true
	})
	return bestTitle, bestHref
}
func findEpisode(body []byte, season, episode *int32) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	out := ""
	doc.Find("td#ayazi_indir a[href]").EachWithBreak(func(_ int, a *goquery.Selection) bool {
		href, _ := a.Attr("href")
		title, _ := a.Attr("title")
		if !strings.Contains(href, "indir_bolum-") || !strings.Contains(title, "Bölüm Türkçe Altyazısı") {
			return true
		}
		s, e := parseSeasonEpisode(title)
		if (season == nil || s == int(*season)) && (episode == nil || e == int(*episode)) {
			out = href
			return false
		}
		return true
	})
	return out
}
func findDownload(body []byte) (string, string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	href, _ := doc.Find("div#altyazi_indir a[href]").First().Attr("href")
	uploader := ""
	doc.Find("strong").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if strings.TrimSpace(s.Text()) == "Altyazı/Çeviri:" {
			uploader = strings.TrimSpace(s.Parent().Text())
			uploader = strings.TrimSpace(strings.TrimPrefix(uploader, "Altyazı/Çeviri:"))
			return false
		}
		return true
	})
	return href, uploader
}
func parseSeasonEpisode(title string) (int, int) {
	season := 1
	ep := 0
	re := regexp.MustCompile(`(\d+)\.\s*Bölüm`)
	if m := re.FindStringSubmatch(title); len(m) == 2 {
		fmt.Sscanf(m[1], "%d", &ep)
	}
	sr := regexp.MustCompile(`(\d+)\.\s*Sezon`)
	if m := sr.FindStringSubmatch(title); len(m) == 2 {
		fmt.Sscanf(m[1], "%d", &season)
	}
	return season, ep
}
func norm(s string) string {
	r := strings.NewReplacer("İ", "i", "I", "i", "Ğ", "g", "Ü", "u", "Ş", "s", "Ö", "o", "Ç", "c", "ı", "i", "ğ", "g", "ü", "u", "ş", "s", "ö", "o", "ç", "c")
	s = strings.ToLower(r.Replace(s))
	return strings.Join(strings.Fields(regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(s, "")), " ")
}
func get(ctx context.Context, svc providercore.Service, raw string, dl bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Language", "tr,en-US;q=0.7,en;q=0.3")
	resp, err := svc.DoProviderRequest(req, providerKey, dl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		resp.Body.Close()
		return nil, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	return resp, nil
}
func baseURL(cfg providercore.Config) string {
	if strings.TrimSpace(cfg.BaseURL) != "" {
		return strings.TrimRight(cfg.BaseURL, "/")
	}
	return defaultBaseURL
}
