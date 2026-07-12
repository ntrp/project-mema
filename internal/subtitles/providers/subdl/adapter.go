package subdl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const providerKey = "subdl"
const defaultBaseURL = "https://api.subdl.com/api/v1"
const defaultDownloadBase = "https://dl.subdl.com"
const maxBody = 10 << 20

var adapter Adapter

type Adapter struct{}

type searchResponse struct {
	Success   *bool       `json:"success"`
	Status    *bool       `json:"status"`
	Error     string      `json:"error"`
	Subtitles []subdlItem `json:"subtitles"`
}

type subdlItem struct {
	Name         string       `json:"name"`
	URL          string       `json:"url"`
	SubtitlePage string       `json:"subtitlePage"`
	Language     string       `json:"language"`
	Author       string       `json:"author"`
	Comment      string       `json:"comment"`
	HI           bool         `json:"hi"`
	Season       *int32       `json:"season"`
	Episode      *int32       `json:"episode"`
	EpisodeFrom  *int32       `json:"episode_from"`
	EpisodeEnd   *int32       `json:"episode_end"`
	Releases     []string     `json:"releases"`
	UnpackFiles  []unpackItem `json:"unpack_files"`
}

type unpackItem struct {
	URL     string `json:"url"`
	Name    string `json:"name"`
	FileID  string `json:"file_n_id"`
	Season  *int32 `json:"season"`
	Episode *int32 `json:"episode"`
}

type sourceRef struct{ IsPack bool `json:"isPack"` }

func init() { providers.Register(providerKey, adapter) }

func (Adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	if _, ok := providercore.NewConfig(config).RequiredSecret("apiKey"); !ok {
		return fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	_, err := get(ctx, service, config, "subtitles", url.Values{"api_key": {providercore.NewConfig(config).Secret("apiKey")}}, false)
	return classify(err)
}

func (Adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	if _, ok := providercore.NewConfig(config).RequiredSecret("apiKey"); !ok {
		return nil, fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	items, err := query(ctx, service, config, request)
	if err != nil { return nil, err }
	return buildCandidates(request, items), nil
}

func (Adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	raw := resolveDownload(candidate.SourceURL)
	data, err := getRaw(ctx, service, config, raw, true)
	if err != nil { return providercore.Download{}, classify(err) }
	name := path.Base(raw)
	if !isSubtitleName(name) && candidate.ReleaseName != "" { name = candidate.ReleaseName + ".zip" }
	member, err := providercore.ExtractSubtitle(name, data, security.ArchiveLimits{})
	if err != nil { return providercore.Download{}, err }
	return providercore.Download{Content: member.Content, URL: raw}, nil
}

func query(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]subdlItem, error) {
	base := params(config, request)
	var all []subdlItem
	seen := map[string]bool{}
	for _, p := range searchParams(base, request) {
		items, err := fetch(ctx, service, config, p)
		if err != nil { return nil, err }
		for _, item := range items {
			if !seen[item.Name] { all = append(all, item); seen[item.Name] = true }
		}
	}
	if len(all) == 0 && request.MediaType == "serie" {
		fallback := cloneValues(base)
		fallback.Del("season_number")
		fallback.Del("episode_number")
		items, err := fetch(ctx, service, config, fallback)
		if err != nil { return nil, err }
		all = append(all, items...)
	}
	return all, nil
}

func searchParams(base url.Values, request providercore.SearchRequest) []url.Values {
	params := []url.Values{cloneValues(base)}
	if request.MediaType != "serie" { return params }
	seasonOnly := cloneValues(base)
	seasonOnly.Del("episode_number")
	params = append(params, seasonOnly)
	if abs := absoluteEpisode(request); abs != 0 && request.EpisodeNumber != nil && abs != *request.EpisodeNumber {
		absolute := cloneValues(base)
		absolute.Set("episode_number", strconv.Itoa(int(abs)))
		absolute.Del("season_number")
		params = append(params, absolute)
	}
	return params
}

func params(config providercore.Config, request providercore.SearchRequest) url.Values {
	apiKey := providercore.NewConfig(config).Secret("apiKey")
	p := url.Values{"api_key": {apiKey}, "languages": {request.LanguageID}, "subs_per_page": {"30"}, "comment": {"1"}, "releases": {"1"}}
	if request.MediaType == "serie" {
		p.Set("type", "tv"); p.Set("bazarr", "1"); p.Set("unpack", "1")
		if request.SeasonNumber != nil { p.Set("season_number", strconv.Itoa(int(*request.SeasonNumber))) }
		if request.EpisodeNumber != nil { p.Set("episode_number", strconv.Itoa(int(*request.EpisodeNumber))) }
		if id := imdbID(request); id != "" { p.Set("imdb_id", id) } else { p.Set("film_name", request.Title) }
		return p
	}
	p.Set("type", "movie")
	if id := imdbID(request); id != "" { p.Set("imdb_id", id) } else { p.Set("film_name", request.Title) }
	if id := externalID(request, "tmdb"); id != "" { p.Set("tmdb_id", id) }
	return p
}

func fetch(ctx context.Context, service providercore.Service, config providercore.Config, params url.Values) ([]subdlItem, error) {
	data, err := get(ctx, service, config, "subtitles", params, false)
	if err != nil { return nil, classify(err) }
	var parsed searchResponse
	if err := json.Unmarshal(data, &parsed); err != nil { return nil, err }
	if (parsed.Success != nil && !*parsed.Success) || (parsed.Status != nil && !*parsed.Status) { return nil, nil }
	return parsed.Subtitles, nil
}

func buildCandidates(request providercore.SearchRequest, items []subdlItem) []providercore.Candidate {
	out := []providercore.Candidate{}
	for _, item := range items {
		itemURL, release, pack := item.URL, strings.Join(item.Releases, ", "), false
		if release == "" { release = item.Name }
		if request.MediaType == "serie" && item.EpisodeFrom != nil && item.EpisodeEnd != nil && *item.EpisodeFrom != *item.EpisodeEnd {
			target := episode(request); abs := absoluteEpisode(request)
			if !inRange(target, abs, *item.EpisodeFrom, *item.EpisodeEnd) { continue }
			pack = true
			if unpack := matchingUnpack(item.UnpackFiles, target, abs); unpack != nil { itemURL = unpack.URL; pack = false }
		}
		ref, _ := json.Marshal(sourceRef{IsPack: pack})
		out = append(out, providercore.Candidate{ProviderName: providerKey, LanguageID: request.LanguageID, Format: "zip", ReleaseName: release, SourceURL: itemURL, SourceRef: string(ref)})
	}
	return out
}

func get(ctx context.Context, service providercore.Service, config providercore.Config, endpoint string, params url.Values, download bool) ([]byte, error) {
	u, _ := url.Parse(providercore.NewConfig(config).BaseURL(defaultBaseURL) + "/")
	ref, _ := url.Parse(endpoint)
	resolved := u.ResolveReference(ref)
	resolved.RawQuery = params.Encode()
	return getRaw(ctx, service, config, resolved.String(), download)
}

func getRaw(ctx context.Context, service providercore.Service, _ providercore.Config, raw string, download bool) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil { return nil, err }
	req.Header.Set("User-Agent", "Sub-Zero/2")
	resp, err := service.DoProviderRequest(req, providerKey, download)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil { return nil, err }
	if len(data) > maxBody { return nil, fmt.Errorf("response size limit exceeded") }
	if resp.StatusCode == http.StatusTooManyRequests { return nil, rateLimitError(resp, data) }
	if resp.StatusCode == http.StatusForbidden { return nil, fmt.Errorf("http status 403") }
	if resp.StatusCode == http.StatusNotFound { return nil, fmt.Errorf("resource not found") }
	if resp.StatusCode != http.StatusOK { return nil, fmt.Errorf("http status %d", resp.StatusCode) }
	return data, nil
}

func rateLimitError(resp *http.Response, data []byte) error {
	var payload struct{ Error string `json:"error"` }
	_ = json.Unmarshal(data, &payload)
	if payload.Error == "daily_limit" || payload.Error == "api_download_limit_exceeded" { return fmt.Errorf("download limit exceeded") }
	if payload.Error == "rate_limit" { if d, _ := strconv.Atoi(resp.Header.Get("Retry-After")); d > 0 { time.Sleep(time.Duration(d) * time.Second) }; return fmt.Errorf("rate limit") }
	if payload.Error == "service_busy" { return fmt.Errorf("service busy") }
	return fmt.Errorf("http status 429")
}

func resolveDownload(raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") { return raw }
	base, _ := url.Parse(defaultDownloadBase + "/")
	ref, _ := url.Parse(strings.TrimPrefix(raw, "/"))
	return base.ResolveReference(ref).String()
}
func cloneValues(in url.Values) url.Values { out := url.Values{}; for k, v := range in { out[k] = append([]string(nil), v...) }; return out }
func imdbID(r providercore.SearchRequest) string { if r.MediaContext.ExternalIDs == nil { return "" }; return r.MediaContext.ExternalIDs["imdb"] }
func externalID(r providercore.SearchRequest, key string) string { if r.MediaContext.ExternalIDs == nil { return "" }; return r.MediaContext.ExternalIDs[key] }
func absoluteEpisode(r providercore.SearchRequest) int32 { for _, n := range r.MediaContext.EpisodeNumbering { if n.AbsoluteNumber != nil { return *n.AbsoluteNumber } }; return 0 }
func episode(r providercore.SearchRequest) int32 { if r.EpisodeNumber == nil { return 0 }; return *r.EpisodeNumber }
func inRange(ep, abs, from, end int32) bool { return (ep >= from && ep <= end) || (abs != 0 && abs >= from && abs <= end) }
func matchingUnpack(files []unpackItem, ep, abs int32) *unpackItem { for i := range files { if files[i].Episode != nil && (*files[i].Episode == ep || (abs != 0 && *files[i].Episode == abs)) { return &files[i] } }; return nil }
func isSubtitleName(name string) bool { lower := strings.ToLower(name); return strings.HasSuffix(lower, ".srt") || strings.HasSuffix(lower, ".ass") || strings.HasSuffix(lower, ".ssa") || strings.HasSuffix(lower, ".vtt") || strings.HasSuffix(lower, ".sub") }
func classify(err error) error { if err == nil { return nil }; s := err.Error(); if strings.Contains(s, "403") || strings.Contains(s, "Invalid API key") { return fmt.Errorf("%w: %v", providercore.ErrProviderPrerequisiteMissing, err) }; if strings.Contains(s, "429") || strings.Contains(s, "rate limit") || strings.Contains(s, "service busy") { return fmt.Errorf("%w: %v", providercore.ErrProviderBrokenUpstream, err) }; return err }
