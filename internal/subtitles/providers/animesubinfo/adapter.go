package animesubinfo

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

const providerKey = "animesubinfo"
const defaultBaseURL = "http://animesub.info"
const maxBytes = 50 << 20

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	resp, err := do(ctx, svc, http.MethodGet, baseURL(cfg)+"/", nil, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return nil
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.LanguageID != "" && sr.LanguageID != "pol" {
		return nil, nil
	}
	titles := searchTitles(sr)
	seen := map[string]bool{}
	var out []providercore.Candidate
	for _, st := range titles {
		u, _ := url.Parse(baseURL(cfg) + "/szukaj.php")
		q := u.Query()
		q.Set("szukane", st.title)
		q.Set("pTitle", st.kind)
		q.Set("pSortuj", "pobrn")
		u.RawQuery = q.Encode()
		resp, err := do(ctx, svc, http.MethodGet, u.String(), nil, false)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		for _, c := range parse(body, u.String(), baseURL(cfg)+"/sciagnij.php") {
			if !seen[c.SourceURL] {
				out = append(out, c)
				seen[c.SourceURL] = true
			}
		}
		if len(out) > 0 && strings.Contains(st.title, " ep") {
			break
		}
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	vals := url.Values{}
	vals.Set("id", strconv.FormatInt(cand.FileID, 10))
	vals.Set("sh", cand.SourceRef)
	vals.Set("single_file", "Pobierz napisy")
	resp, err := do(ctx, svc, http.MethodPost, cand.SourceURL, strings.NewReader(vals.Encode()), true)
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
		member, err := providercore.ExtractSubtitle("animesubinfo.zip", body, security.ArchiveLimits{})
		if err != nil {
			return providercore.Download{}, err
		}
		body = member.Content
	}
	return providercore.Download{Content: body, URL: cand.SourceURL}, nil
}

type stitle struct{ kind, title string }

func searchTitles(sr providercore.SearchRequest) []stitle {
	title := strings.TrimSpace(sr.Title)
	epTitle := ""
	if sr.EpisodeNumber != nil {
		epTitle = fmt.Sprintf("%s ep%02d", title, *sr.EpisodeNumber)
	}
	var out []stitle
	if epTitle != "" {
		for _, k := range []string{"org", "en", "pl"} {
			out = append(out, stitle{k, epTitle})
		}
	} else {
		for _, k := range []string{"org", "en", "pl"} {
			out = append(out, stitle{k, title})
		}
	}
	for _, a := range sr.MediaContext.Aliases {
		if strings.TrimSpace(a.Value) != "" {
			out = append(out, stitle{"en", strings.TrimSpace(a.Value)})
		}
	}
	return out
}
func parse(body []byte, page, dl string) []providercore.Candidate {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	var out []providercore.Candidate
	doc.Find("table.Napisy").Each(func(_ int, table *goquery.Selection) {
		rows := table.Find("tr.KNap")
		if rows.Length() < 3 {
			return
		}
		row1 := rows.Eq(0).Find("td")
		row2 := rows.Eq(1).Find("td")
		row3 := rows.Eq(2).Find("td")
		titleOrg := textAt(row1, 0)
		titleEng := textAt(row2, 0)
		titleAlt := textAt(row3, 0)
		format := textAt(row1, 3)
		size := textAt(row2, row2.Length()-1)
		downloads := downloadCount(textAt(row3, 3))
		form := table.Find("tr.KKom form[method='POST']").First()
		idStr, _ := form.Find("input[name='id']").Attr("value")
		hash, _ := form.Find("input[name='sh']").Attr("value")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		if id == 0 || hash == "" {
			return
		}
		release := strings.TrimSpace(strings.Join([]string{titleOrg, titleEng, titleAlt, size}, " - "))
		out = append(out, providercore.Candidate{ProviderName: providerKey, LanguageID: "pol", FileID: id, Format: format, ReleaseName: release, DownloadCount: downloads, SourceURL: dl, SourceRef: hash})
	})
	return out
}
func textAt(s *goquery.Selection, i int) string {
	if i < 0 || i >= s.Length() {
		return ""
	}
	return strings.TrimSpace(s.Eq(i).Text())
}
func downloadCount(s string) int {
	re := regexp.MustCompile(`\d+`)
	m := re.FindString(s)
	n, _ := strconv.Atoi(m)
	return n
}
func do(ctx context.Context, svc providercore.Service, method, raw string, body io.Reader, dl bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, raw, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Charset", "ISO-8859-2,utf-8;q=0.7,*;q=0.3")
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
