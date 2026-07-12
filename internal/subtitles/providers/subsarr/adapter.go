package subsarr

import (
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

func init() { providers.Register("subsarr", Adapter) }

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	if strings.TrimSpace(providercore.NewConfig(cfg).BaseURL("")) == "" {
		return fmt.Errorf("%w: baseUrl is required", providercore.ErrProviderPrerequisiteMissing)
	}
	_, _, err := do(ctx, svc, cfg, http.MethodGet, "/api/health", false)
	return err
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if strings.TrimSpace(providercore.NewConfig(cfg).BaseURL("")) == "" {
		return nil, fmt.Errorf("%w: baseUrl is required", providercore.ErrProviderPrerequisiteMissing)
	}
	u, err := endpoint(cfg, "/api/subtitles/search")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("title", sr.Title)
	q.Set("q", sr.Title)
	if sr.LanguageID != "" {
		q.Set("language", sr.LanguageID)
	}
	if sr.MediaType != "" {
		q.Set("type", sr.MediaType)
	}
	if sr.Year != nil {
		q.Set("year", strconv.Itoa(int(*sr.Year)))
	}
	if sr.SeasonNumber != nil {
		q.Set("season", strconv.Itoa(int(*sr.SeasonNumber)))
	}
	if sr.EpisodeNumber != nil {
		q.Set("episode", strconv.Itoa(int(*sr.EpisodeNumber)))
	}
	if sr.FilePath != "" {
		q.Set("path", sr.FilePath)
	}
	u.RawQuery = q.Encode()
	data, _, err := do(ctx, svc, cfg, http.MethodGet, u.String(), false)
	if err != nil {
		return nil, err
	}
	return parseCandidates(data, sr.LanguageID), nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	raw := strings.TrimSpace(cand.SourceURL)
	if raw == "" && cand.FileID != 0 {
		raw = "/api/subtitles/download/" + strconv.FormatInt(cand.FileID, 10)
	}
	if raw == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	data, resp, err := do(ctx, svc, cfg, http.MethodGet, raw, false)
	if err != nil {
		return providercore.Download{}, err
	}
	member, err := providercore.ExtractSubtitle(downloadName(raw, resp), data, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: absoluteString(cfg, raw)}, nil
}

func do(ctx context.Context, svc providercore.Service, cfg providercore.Config, method, raw string, download bool) ([]byte, *http.Response, error) {
	u := absoluteString(cfg, raw)
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := svc.DoProviderRequest(req, "subsarr", download)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil {
		return nil, resp, err
	}
	if len(body) > maxBody {
		return nil, resp, fmt.Errorf("provider response size limit exceeded")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, resp, fmt.Errorf("provider returned HTTP %d", resp.StatusCode)
	}
	return body, resp, nil
}

func endpoint(cfg providercore.Config, p string) (*url.URL, error) {
	return url.Parse(absoluteString(cfg, p))
}

func absoluteString(cfg providercore.Config, raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	base, _ := url.Parse(providercore.NewConfig(cfg).BaseURL("") + "/")
	ref, _ := url.Parse(strings.TrimLeft(raw, "/"))
	if strings.HasPrefix(raw, "/") {
		ref, _ = url.Parse(raw)
	}
	return base.ResolveReference(ref).String()
}

func parseCandidates(data []byte, fallbackLang string) []providercore.Candidate {
	var v any
	if json.Unmarshal(data, &v) != nil {
		return nil
	}
	objects := collect(v)
	out := make([]providercore.Candidate, 0, len(objects))
	for _, obj := range objects {
		dl := firstString(obj, "download_url", "downloadUrl", "download", "url", "link", "file")
		if dl == "" {
			continue
		}
		lang := firstString(obj, "language", "lang", "language_id", "locale")
		if lang == "" {
			lang = fallbackLang
		}
		out = append(out, providercore.Candidate{ProviderName: "subsarr", LanguageID: lang, FileID: firstInt(obj, "id", "file_id", "subtitle_id"), Format: format(dl), ReleaseName: firstString(obj, "release", "release_name", "filename", "name", "title"), DownloadCount: int(firstInt(obj, "downloads", "download_count")), SourceURL: dl})
	}
	return out
}

func collect(v any) []map[string]any {
	var out []map[string]any
	switch x := v.(type) {
	case map[string]any:
		out = append(out, x)
		for _, child := range x {
			out = append(out, collect(child)...)
		}
	case []any:
		for _, child := range x {
			out = append(out, collect(child)...)
		}
	}
	return out
}
func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok {
				return strings.TrimSpace(s)
			}
		}
	}
	return ""
}
func firstInt(m map[string]any, keys ...string) int64 {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			switch x := v.(type) {
			case float64:
				return int64(x)
			case int64:
				return x
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
