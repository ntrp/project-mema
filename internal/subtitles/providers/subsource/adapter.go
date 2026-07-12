package subsource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/security"
)

const maxBody = 50 << 20

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

func init() { providers.Register("subsource", Adapter) }

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	if _, ok := providercore.NewConfig(cfg).RequiredSecret("apiKey"); !ok {
		return fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	_, _, err := do(ctx, svc, cfg, http.MethodGet, "/api/search", nil, false)
	return err
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if _, ok := providercore.NewConfig(cfg).RequiredSecret("apiKey"); !ok {
		return nil, fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	searchURL := searchEndpoint(cfg, sr)
	data, _, err := do(ctx, svc, cfg, http.MethodGet, searchURL, nil, false)
	if err != nil {
		return nil, err
	}
	matches := parseMatches(data)
	var out []providercore.Candidate
	for _, match := range matches {
		detailURL := detailEndpoint(cfg, match, sr)
		detail, _, err := do(ctx, svc, cfg, http.MethodGet, detailURL, nil, false)
		if err != nil {
			return nil, err
		}
		out = append(out, parseSubtitles(detail, sr.LanguageID)...)
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	raw := strings.TrimSpace(cand.SourceURL)
	if raw == "" && cand.FileID != 0 {
		raw = "/api/downloadSubtitle?subtitleId=" + strconv.FormatInt(cand.FileID, 10)
	}
	if raw == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	data, resp, err := do(ctx, svc, cfg, http.MethodGet, raw, nil, true)
	if err != nil {
		return providercore.Download{}, err
	}
	member, err := providercore.ExtractSubtitle(downloadName(raw, resp), data, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: absolute(cfg, raw)}, nil
}

func do(ctx context.Context, svc providercore.Service, cfg providercore.Config, method, raw string, body io.Reader, download bool) ([]byte, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, absolute(cfg, raw), body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token, ok := providercore.NewConfig(cfg).RequiredSecret("apiKey"); ok {
		req.Header.Set("Authorization", token)
	}
	resp, err := svc.DoProviderRequest(req, "subsource", download)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return nil, resp, err
	}
	if len(data) > maxBody {
		return nil, resp, fmt.Errorf("provider response size limit exceeded")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp, fmt.Errorf("provider returned HTTP %d", resp.StatusCode)
	}
	return data, resp, nil
}

func searchEndpoint(cfg providercore.Config, sr providercore.SearchRequest) string {
	u, _ := url.Parse(absolute(cfg, "/api/search"))
	q := u.Query()
	q.Set("query", sr.Title)
	q.Set("q", sr.Title)
	if sr.Year != nil {
		q.Set("year", strconv.Itoa(int(*sr.Year)))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func detailEndpoint(cfg providercore.Config, m match, sr providercore.SearchRequest) string {
	if m.URL != "" {
		return m.URL
	}
	kind := "movie"
	if sr.MediaType == "serie" {
		kind = "tv"
	}
	u, _ := url.Parse(absolute(cfg, "/api/"+kind))
	q := u.Query()
	q.Set(kind+"Name", m.LinkName)
	if sr.SeasonNumber != nil {
		q.Set("season", strconv.Itoa(int(*sr.SeasonNumber)))
	}
	if sr.EpisodeNumber != nil {
		q.Set("episode", strconv.Itoa(int(*sr.EpisodeNumber)))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

type match struct {
	LinkName string
	URL      string
}

func parseMatches(data []byte) []match {
	var v any
	if json.Unmarshal(data, &v) != nil {
		return nil
	}
	objs := collect(v)
	out := make([]match, 0, len(objs))
	for _, obj := range objs {
		link := firstString(obj, "linkName", "link_name", "name", "title", "slug")
		u := firstString(obj, "details", "details_url", "url", "link")
		if link != "" || u != "" {
			out = append(out, match{LinkName: link, URL: u})
		}
	}
	return out
}

func parseSubtitles(data []byte, fallbackLang string) []providercore.Candidate {
	var v any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if dec.Decode(&v) != nil {
		return nil
	}
	objs := collect(v)
	out := make([]providercore.Candidate, 0, len(objs))
	for _, obj := range objs {
		dl := firstString(obj, "download", "download_url", "downloadUrl", "url", "link")
		if dl == "" {
			continue
		}
		lang := firstString(obj, "language", "lang", "language_id")
		if lang == "" {
			lang = fallbackLang
		}
		out = append(out, providercore.Candidate{ProviderName: "subsource", LanguageID: lang, FileID: firstInt(obj, "id", "subtitleId", "subtitle_id"), Format: format(dl), ReleaseName: firstString(obj, "release", "release_name", "filename", "name"), DownloadCount: int(firstInt(obj, "downloads", "download_count")), SourceURL: dl})
	}
	return out
}

func absolute(cfg providercore.Config, raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	base, _ := url.Parse(providercore.NewConfig(cfg).BaseURL("https://subsource.net") + "/")
	ref, _ := url.Parse(strings.TrimLeft(raw, "/"))
	if strings.HasPrefix(raw, "/") {
		ref, _ = url.Parse(raw)
	}
	return base.ResolveReference(ref).String()
}
func collect(v any) []map[string]any {
	var out []map[string]any
	switch x := v.(type) {
	case map[string]any:
		out = append(out, x)
		for _, c := range x {
			out = append(out, collect(c)...)
		}
	case []any:
		for _, c := range x {
			out = append(out, collect(c)...)
		}
	}
	return out
}
func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch x := v.(type) {
			case string:
				return strings.TrimSpace(x)
			case json.Number:
				return x.String()
			}
		}
	}
	return ""
}
func firstInt(m map[string]any, keys ...string) int64 {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch x := v.(type) {
			case json.Number:
				n, _ := x.Int64()
				return n
			case float64:
				return int64(x)
			}
		}
	}
	return 0
}
func format(raw string) string {
	ext := strings.TrimPrefix(path.Ext(raw), ".")
	if ext == "" || len(ext) > 4 {
		return "srt"
	}
	return ext
}
func downloadName(raw string, resp *http.Response) string {
	if resp != nil {
		_, p, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if p["filename"] != "" {
			return p["filename"]
		}
	}
	u, _ := url.Parse(raw)
	if b := path.Base(u.Path); b != "." && b != "/" {
		return b
	}
	return "subtitle.srt"
}
