package hosszupuska

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const providerKey = "hosszupuska"
const defaultBaseURL = "http://hosszupuskasub.com"
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
		return nil, fmt.Errorf("%w: hosszupuska supports series", providercore.ErrProviderPrerequisiteMissing)
	}
	if sr.SeasonNumber == nil || sr.EpisodeNumber == nil {
		return nil, fmt.Errorf("%w: season and episode required", providercore.ErrProviderPrerequisiteMissing)
	}
	u, _ := url.Parse(baseURL(cfg) + "/sorozatok.php")
	q := u.Query()
	q.Set("cim", strings.ReplaceAll(sr.Title, " ", "+"))
	q.Set("evad", fmt.Sprintf("%02d", *sr.SeasonNumber))
	q.Set("resz", fmt.Sprintf("%02d", *sr.EpisodeNumber))
	q.Set("nyelvtipus", "%")
	q.Set("x", "24")
	q.Set("y", "8")
	u.RawQuery = strings.ReplaceAll(q.Encode(), "%2B", "+")
	resp, err := get(ctx, svc, u.String(), false)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return parse(body, int(*sr.SeasonNumber), int(*sr.EpisodeNumber), u.String(), sr.LanguageID), nil
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
	member, err := providercore.ExtractSubtitle("hosszupuska.zip", body, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: cand.SourceURL}, nil
}

func parse(body []byte, season, episode int, ref, requested string) []providercore.Candidate {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	out := []providercore.Candidate{}
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		html, _ := goquery.OuterHtml(table)
		if !strings.Contains(html, "over2.jpg") || !strings.Contains(html, "css/infooldal.png") {
			return
		}
		table.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			rowhtml, _ := goquery.OuterHtml(tr)
			if !strings.Contains(rowhtml, "over2.jpg") {
				return
			}
			cells := tr.Find("td")
			title := strings.TrimSpace(cells.Eq(1).Text())
			s, e := se(title)
			if s != season || e != episode {
				return
			}
			lang := requested
			if src, ok := cells.Eq(2).Find("img").First().Attr("src"); ok {
				lang = langFromGif(pathBase(src))
			}
			href, _ := cells.Eq(6).Find("a").First().Attr("href")
			if href == "" {
				return
			}
			version := version(title)
			id := fileID(href)
			out = append(out, providercore.Candidate{ProviderName: providerKey, LanguageID: lang, FileID: id, Format: "srt", ReleaseName: strings.Join(splitReleases(version), ", "), SourceURL: resolve(defaultBaseURL+"/", href), SourceRef: ref})
		})
	})
	return out
}
func se(s string) (int, int) {
	re := regexp.MustCompile(`s(\d{1,2})e(\d{1,2})`)
	m := re.FindStringSubmatch(strings.ToLower(s))
	if len(m) != 3 {
		return 0, 0
	}
	a, _ := strconv.Atoi(strings.TrimLeft(m[1], "0"))
	b, _ := strconv.Atoi(strings.TrimLeft(m[2], "0"))
	return a, b
}
func version(s string) string {
	parts := regexp.MustCompile(`\(([^)]*)\)`).FindAllStringSubmatch(s, -1)
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1][1]
}
func splitReleases(s string) []string {
	out := []string{}
	for _, p := range strings.Split(s, ",") {
		if x := strings.TrimSpace(p); x != "" {
			out = append(out, x)
		}
	}
	return out
}
func langFromGif(g string) string {
	if g == "1.gif" {
		return "hun"
	}
	if g == "2.gif" {
		return "eng"
	}
	return ""
}
func fileID(h string) int64 {
	u, err := url.Parse(h)
	if err == nil {
		h = u.Query().Get("file")
	}
	re := regexp.MustCompile(`\d+`)
	n, _ := strconv.ParseInt(re.FindString(h), 10, 64)
	return n
}
func pathBase(s string) string {
	parts := strings.Split(strings.Trim(s, "/"), "/")
	return parts[len(parts)-1]
}
func resolve(base, ref string) string {
	b, _ := url.Parse(base)
	r, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return b.ResolveReference(r).String()
}
func get(ctx context.Context, svc providercore.Service, raw string, dl bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return nil, err
	}
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
