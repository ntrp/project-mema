package subsro

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

func init() { providers.Register("subsro", Adapter) }

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	if _, ok := providercore.NewConfig(cfg).RequiredSecret("apiKey"); !ok {
		return fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	_, _, err := do(ctx, svc, cfg, http.MethodGet, "/subtitles", false)
	return err
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if _, ok := providercore.NewConfig(cfg).RequiredSecret("apiKey"); !ok {
		return nil, fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	imdb := imdbID(sr)
	if imdb == "" {
		return nil, fmt.Errorf("%w: imdb id is required", providercore.ErrProviderPrerequisiteMissing)
	}
	u, _ := url.Parse(absolute(cfg, "/subtitles"))
	q := u.Query()
	q.Set("key", providercore.NewConfig(cfg).Secret("apiKey"))
	q.Set("imdb", imdb)
	q.Set("imdb_id", imdb)
	if sr.LanguageID != "" {
		q.Set("language", sr.LanguageID)
	}
	if sr.SeasonNumber != nil {
		q.Set("season", strconv.Itoa(int(*sr.SeasonNumber)))
	}
	if sr.EpisodeNumber != nil {
		q.Set("episode", strconv.Itoa(int(*sr.EpisodeNumber)))
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
		raw = "/download/" + strconv.FormatInt(cand.FileID, 10)
	}
	if raw == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	u := absolute(cfg, raw)
	parsed, _ := url.Parse(u)
	q := parsed.Query()
	if _, ok := q["key"]; !ok {
		q.Set("key", providercore.NewConfig(cfg).Secret("apiKey"))
	}
	parsed.RawQuery = q.Encode()
	data, resp, err := do(ctx, svc, cfg, http.MethodGet, parsed.String(), true)
	if err != nil {
		return providercore.Download{}, err
	}
	member, err := providercore.ExtractSubtitle(downloadName(raw, resp), data, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: parsed.String()}, nil
}

func do(ctx context.Context, svc providercore.Service, cfg providercore.Config, method, raw string, download bool) ([]byte, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, absolute(cfg, raw), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := svc.DoProviderRequest(req, "subsro", download)
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

func parseCandidates(data []byte, fallbackLang string) []providercore.Candidate {
	var v any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if dec.Decode(&v) != nil {
		return nil
	}
	objs := collect(v)
	out := make([]providercore.Candidate, 0, len(objs))
	for _, obj := range objs {
		dl := firstString(obj, "download", "download_url", "downloadUrl", "url", "link", "archive")
		if dl == "" {
			continue
		}
		lang := firstString(obj, "language", "lang", "language_id")
		if lang == "" {
			lang = fallbackLang
		}
		out = append(out, providercore.Candidate{ProviderName: "subsro", LanguageID: lang, FileID: firstInt(obj, "id", "file_id", "subtitle_id"), Format: format(dl), ReleaseName: firstString(obj, "release", "release_name", "filename", "name", "title"), DownloadCount: int(firstInt(obj, "downloads", "download_count")), SourceURL: dl})
	}
	return out
}

func absolute(cfg providercore.Config, raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	base, _ := url.Parse(providercore.NewConfig(cfg).BaseURL("https://api.subs.ro") + "/")
	ref, _ := url.Parse(strings.TrimLeft(raw, "/"))
	if strings.HasPrefix(raw, "/") {
		ref, _ = url.Parse(raw)
	}
	return base.ResolveReference(ref).String()
}
func imdbID(r providercore.SearchRequest) string {
	if r.MediaContext.ExternalIDs != nil {
		return r.MediaContext.ExternalIDs["imdb"]
	}
	return ""
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
