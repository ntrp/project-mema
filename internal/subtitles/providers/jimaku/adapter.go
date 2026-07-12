package jimaku

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

	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/security"
)

const providerKey = "jimaku"
const defaultBaseURL = "https://jimaku.cc/api"
const maxBody = 10 << 20
const corruptThreshold = 500

var adapter Adapter

type Adapter struct{}

type entry struct {
	ID   int64 `json:"id"`
	Name string `json:"name"`
	Flags struct{ Movie bool `json:"movie"` } `json:"flags"`
}

type fileItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Size int64  `json:"size"`
}

func init() { providers.Register(providerKey, adapter) }

func (Adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	if _, ok := providercore.NewConfig(config).RequiredSecret("apiKey"); !ok {
		return fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	_, err := do(ctx, service, config, http.MethodGet, "entries/search?query=test", false)
	return classify(err)
}

func (Adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	view := providercore.NewConfig(config)
	if _, ok := view.RequiredSecret("apiKey"); !ok {
		return nil, fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	entries, err := searchEntry(ctx, service, config, request)
	if err != nil || len(entries) == 0 {
		return nil, err
	}
	files, err := searchFiles(ctx, service, config, entries[0], request)
	if err != nil {
		return nil, err
	}
	return candidates(view, request, files), nil
}

func (Adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	data, err := do(ctx, service, config, http.MethodGet, candidate.SourceURL, true)
	if err != nil {
		return providercore.Download{}, classify(err)
	}
	member, err := providercore.ExtractSubtitle(candidate.ReleaseName, data, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: candidate.SourceURL}, nil
}

func searchEntry(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]entry, error) {
	params := url.Values{}
	if id := externalID(request, "anilist"); id != "" {
		params.Set("anilist_id", id)
	} else if providercore.NewConfig(config).BoolSetting("enableNameSearchFallback") || request.MediaType == "movie" {
		name := strings.ToLower(request.Title)
		if request.MediaType == "serie" && request.SeasonNumber != nil && *request.SeasonNumber > 1 {
			name += " " + strconv.Itoa(int(*request.SeasonNumber))
		}
		params.Set("query", name)
	} else {
		return nil, nil
	}
	var out []entry
	if err := getJSON(ctx, service, config, "entries/search?"+params.Encode(), &out); err != nil {
		return nil, err
	}
	if len(out) == 0 && params.Get("query") != "" {
		params.Set("anime", "false")
		if err := getJSON(ctx, service, config, "entries/search?"+params.Encode(), &out); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func searchFiles(ctx context.Context, service providercore.Service, config providercore.Config, e entry, request providercore.SearchRequest) ([]fileItem, error) {
	params := url.Values{}
	if request.MediaType == "serie" && !e.Flags.Movie && request.EpisodeNumber != nil {
		params.Set("episode", strconv.Itoa(int(*request.EpisodeNumber)))
	}
	var out []fileItem
	endpoint := fmt.Sprintf("entries/%d/files", e.ID)
	if err := getJSON(ctx, service, config, endpoint+query(params), &out); err != nil {
		return nil, err
	}
	if len(out) == 0 && params.Get("episode") != "" {
		if err := getJSON(ctx, service, config, endpoint, &out); err != nil {
			return nil, err
		}
		out = onlyArchives(out)
	}
	return out, nil
}

func candidates(view providercore.ConfigView, request providercore.SearchRequest, files []fileItem) []providercore.Candidate {
	hasSubtitle, hasArchive := false, false
	for _, f := range files {
		if isArchive(f.Name) { hasArchive = true } else { hasSubtitle = true }
	}
	out := []providercore.Candidate{}
	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name), ".7z") || (f.Size > 0 && f.Size < corruptThreshold) { continue }
		archive := isArchive(f.Name)
		if archive && hasSubtitle && !view.BoolSetting("enableArchivesDownload") { continue }
		if !view.BoolSetting("enableAiSubs") && looksAI(f.Name) { continue }
		if len(guessedLanguages(f.Name)) > 1 { continue }
		out = append(out, providercore.Candidate{ProviderName: providerKey, LanguageID: request.LanguageID, Format: strings.TrimPrefix(strings.ToLower(path.Ext(f.Name)), "."), ReleaseName: f.Name, SourceURL: f.URL})
	}
	_ = hasArchive
	return out
}

func getJSON(ctx context.Context, service providercore.Service, config providercore.Config, endpoint string, v any) error {
	data, err := do(ctx, service, config, http.MethodGet, endpoint, false)
	if err != nil { return classify(err) }
	if err := json.Unmarshal(data, v); err != nil { return err }
	return nil
}

func do(ctx context.Context, service providercore.Service, config providercore.Config, method, raw string, download bool) ([]byte, error) {
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") { raw = absolute(config, raw) }
	req, err := http.NewRequestWithContext(ctx, method, raw, nil)
	if err != nil { return nil, err }
	secret, _ := providercore.NewConfig(config).RequiredSecret("apiKey")
	req.Header.Set("Authorization", secret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := service.DoProviderRequest(req, providerKey, download)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxBody+1))
	if err != nil { return nil, err }
	if len(data) > maxBody { return nil, fmt.Errorf("response size limit exceeded") }
	if resp.StatusCode == http.StatusUnauthorized { return nil, fmt.Errorf("http status %d", resp.StatusCode) }
	if resp.StatusCode < 200 || resp.StatusCode > 299 { return nil, fmt.Errorf("http status %d", resp.StatusCode) }
	return data, nil
}

func absolute(config providercore.Config, endpoint string) string {
	base, _ := url.Parse(providercore.NewConfig(config).BaseURL(defaultBaseURL) + "/")
	ref, _ := url.Parse(endpoint)
	return base.ResolveReference(ref).String()
}
func query(v url.Values) string { if len(v) == 0 { return "" }; return "?" + v.Encode() }
func externalID(r providercore.SearchRequest, key string) string { if r.MediaContext.ExternalIDs == nil { return "" }; return r.MediaContext.ExternalIDs[key] }
func isArchive(name string) bool { l := strings.ToLower(name); return strings.HasSuffix(l, ".zip") || strings.HasSuffix(l, ".rar") }
func onlyArchives(files []fileItem) []fileItem { out := []fileItem{}; for _, f := range files { if isArchive(f.Name) { out = append(out, f) } }; return out }
func looksAI(name string) bool { l := strings.ToLower(name); return strings.Contains(l, "whisperai") || strings.Contains(l, "whisper") }
func classify(err error) error { if err == nil { return nil }; s := err.Error(); if strings.Contains(s, "401") || strings.Contains(s, "403") { return fmt.Errorf("%w: %v", providercore.ErrProviderPrerequisiteMissing, err) }; if strings.Contains(s, "429") { return fmt.Errorf("%w: %v", providercore.ErrProviderBrokenUpstream, err) }; return err }
func guessedLanguages(name string) []string { l := strings.ToLower(name); if strings.Contains(l, ".ja-en.") || strings.Contains(l, "[ja en]") { return []string{"ja", "en"} }; return []string{"ja"} }
