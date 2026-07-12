package bayflix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const (
	providerKey      = "bayflix"
	defaultBaseURL   = "https://bayflix.sb"
	maxResponseBytes = 50 << 20
)

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

type searchItem struct {
	ID           string `json:"_id"`
	Title        string `json:"title"`
	SubtitleLink string `json:"subtitle_link"`
	ReleaseName  any    `json:"release_name"`
	ReleaseDate  string `json:"release_date"`
	MediaType    string `json:"media_type"`
	Description  string `json:"description"`
	Language     string `json:"language"`
	Downloads    int    `json:"downloads"`
}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL(cfg)+"/", nil)
	if err != nil {
		return err
	}
	resp, err := svc.DoProviderRequest(req, providerKey, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return fmt.Errorf("%s test failed: http status %d", providerKey, resp.StatusCode)
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return nil
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if sr.MediaType != "" && !strings.EqualFold(sr.MediaType, "movie") && !strings.EqualFold(sr.MediaType, "serie") {
		return nil, fmt.Errorf("%w: bayflix supports movie and serie", providercore.ErrProviderPrerequisiteMissing)
	}
	u, _ := url.Parse(baseURL(cfg) + "/api/subtitles/search")
	q := u.Query()
	q.Set("title", strings.TrimSpace(sr.Title))
	u.RawQuery = q.Encode()
	resp, err := doJSON(ctx, svc, http.MethodGet, u.String(), nil, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var items []searchItem
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes)).Decode(&items); err != nil {
		return nil, err
	}
	out := make([]providercore.Candidate, 0, len(items))
	for _, item := range items {
		if item.ID == "" || !matchesYear(item.ReleaseDate, sr.Year) || !matchesEpisode(item, sr) {
			continue
		}
		link := item.SubtitleLink
		if link == "" {
			link = baseURL(cfg) + "/api/subtitles/download/" + url.PathEscape(item.ID)
		}
		lang := sr.LanguageID
		if item.Language != "" {
			lang = item.Language
		}
		if lang == "" {
			lang = "bul"
		}
		out = append(out, providercore.Candidate{ProviderName: providerKey, LanguageID: lang, Format: "srt", ReleaseName: strings.Join(releaseNames(item.ReleaseName), "\n"), DownloadCount: item.Downloads, SourceURL: link, SourceRef: u.String()})
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, _ providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	resp, err := doJSON(ctx, svc, http.MethodGet, cand.SourceURL, nil, true)
	if err != nil {
		return providercore.Download{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return providercore.Download{}, err
	}
	if len(body) > maxResponseBytes {
		return providercore.Download{}, security.ErrUnsafeArchive
	}
	member, err := providercore.ExtractSubtitle(filename(cand.SourceURL), body, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: cand.SourceURL}, nil
}

func doJSON(ctx context.Context, svc providercore.Service, method, rawURL string, body io.Reader, dl bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", defaultBaseURL+"/")
	resp, err := svc.DoProviderRequest(req, providerKey, dl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		resp.Body.Close()
		return nil, fmt.Errorf("%w: rate limited", providercore.ErrProviderBrokenUpstream)
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
func releaseNames(v any) []string {
	switch x := v.(type) {
	case string:
		if x != "" {
			return []string{x}
		}
	case []any:
		r := []string{}
		for _, e := range x {
			if s, ok := e.(string); ok && s != "" {
				r = append(r, s)
			}
		}
		return r
	case []string:
		return x
	}
	return nil
}
func matchesYear(date string, y *int32) bool {
	if y == nil || len(date) < 4 {
		return true
	}
	n, err := strconv.Atoi(date[:4])
	return err != nil || n == int(*y) || n == int(*y)-1 || n == int(*y)+1
}
func matchesEpisode(item searchItem, sr providercore.SearchRequest) bool {
	if sr.SeasonNumber == nil || sr.EpisodeNumber == nil {
		return true
	}
	want1, want2 := int(*sr.SeasonNumber), int(*sr.EpisodeNumber)
	for _, s := range append(releaseNames(item.ReleaseName), strings.Split(item.Description, "\n")...) {
		if epSeason(s) == want1 && epEpisode(s) == want2 {
			return true
		}
	}
	return false
}
func epSeason(s string) int  { a, b := episodeTuple(s); _ = b; return a }
func epEpisode(s string) int { _, b := episodeTuple(s); return b }
func episodeTuple(s string) (int, int) {
	lower := strings.ToLower(s)
	for _, sep := range []string{"e", "x"} {
		for i := 0; i < len(lower); i++ {
			if lower[i] == 's' || sep == "x" {
				rest := lower[i:]
				var a, b int
				var n int
				if sep == "e" {
					n, _ = fmt.Sscanf(rest, "s%de%d", &a, &b)
				} else {
					n, _ = fmt.Sscanf(rest, "%dx%d", &a, &b)
				}
				if n == 2 {
					return a, b
				}
			}
		}
	}
	return 0, 0
}
func filename(raw string) string {
	u, err := url.Parse(raw)
	if err == nil {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1] + ".zip"
		}
	}
	return "subtitle.zip"
}
