package subx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/security"
)

const defaultBaseURL = "https://subx-api.duckdns.org"

var spainSpanishRE = regexp.MustCompile(`(?i)españa|ib[eé]rico|castellano|gallego|castilla|europ[ae]`)

type adapter struct{}

type searchResponse struct{ Items []item `json:"items"` }
type item struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PageURL     string  `json:"page_url"`
	Season      *int32  `json:"season"`
	Episode     *int32  `json:"episode"`
	Downloads   int     `json:"downloads"`
}

func init() { providers.Register("subx", adapter{}) }

func (adapter) Test(ctx context.Context, service providercore.Service, config providercore.Config) error {
	if _, ok := providercore.NewConfig(config).RequiredSecret("apiKey"); !ok {
		return fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := newRequest(ctx, config, http.MethodGet, "/api/subtitles/search?limit=1&video_type=movie&title=test", nil, false)
	if err != nil { return err }
	resp, err := service.DoProviderRequest(req, "subx", false)
	if resp != nil { io.Copy(io.Discard, resp.Body); resp.Body.Close() }
	if err != nil { return err }
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("%w: invalid SubX API key", providercore.ErrProviderPrerequisiteMissing)
	}
	if resp.StatusCode == http.StatusTooManyRequests { return fmt.Errorf("%w: SubX rate limit exceeded", providercore.ErrProviderBrokenUpstream) }
	if resp.StatusCode >= 500 { return fmt.Errorf("%w: SubX status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode) }
	return nil
}

func (adapter) Search(ctx context.Context, service providercore.Service, config providercore.Config, request providercore.SearchRequest) ([]providercore.Candidate, error) {
	if request.MediaType != "movie" && request.MediaType != "serie" { return nil, nil }
	reqURL := searchURL(config, request)
	req, err := newRequest(ctx, config, http.MethodGet, reqURL, nil, false)
	if err != nil { return nil, err }
	resp, err := service.DoProviderRequest(req, "subx", false)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest { return nil, nil }
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden { return nil, fmt.Errorf("%w: invalid SubX API key", providercore.ErrProviderPrerequisiteMissing) }
	if resp.StatusCode == http.StatusTooManyRequests { return nil, fmt.Errorf("%w: SubX rate limit exceeded", providercore.ErrProviderBrokenUpstream) }
	if resp.StatusCode < 200 || resp.StatusCode > 299 { return nil, fmt.Errorf("SubX status %d", resp.StatusCode) }
	var parsed searchResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 10<<20)).Decode(&parsed); err != nil { return nil, err }
	return candidates(config, request, parsed.Items), nil
}

func (adapter) Download(ctx context.Context, service providercore.Service, config providercore.Config, candidate providercore.Candidate) (providercore.Download, error) {
	raw := strings.TrimSpace(candidate.SourceURL)
	if raw == "" { return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing) }
	req, err := newRequest(ctx, config, http.MethodGet, raw, nil, true)
	if err != nil { return providercore.Download{}, err }
	resp, err := service.DoProviderRequest(req, "subx", true)
	if err != nil { return providercore.Download{}, err }
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 20<<20))
	if err != nil { return providercore.Download{}, err }
	member, err := providercore.ExtractSubtitle(path.Base(raw), data, security.ArchiveLimits{})
	if err != nil { return providercore.Download{}, err }
	return providercore.Download{Content: member.Content, URL: raw}, nil
}

func searchURL(config providercore.Config, request providercore.SearchRequest) string {
	u, _ := url.Parse(providercore.NewConfig(config).BaseURL(defaultBaseURL) + "/api/subtitles/search")
	q := u.Query(); q.Set("limit", "200")
	if request.MediaType == "serie" { q.Set("video_type", "episode") } else { q.Set("video_type", "movie") }
	if id := request.MediaContext.ExternalIDs["imdb"]; id != "" { q.Set("imdb_id", id) } else { q.Set("title", cleanTitle(request.Title)) }
	if request.Year != nil { q.Set("year", strconv.Itoa(int(*request.Year))) }
	u.RawQuery = q.Encode(); return u.String()
}

func candidates(config providercore.Config, request providercore.SearchRequest, items []item) []providercore.Candidate {
	out, packs := []providercore.Candidate{}, []providercore.Candidate{}
	for _, it := range items {
		if request.SeasonNumber != nil && (it.Season == nil || *it.Season != *request.SeasonNumber) { continue }
		c := candidate(config, it)
		if request.EpisodeNumber != nil {
			if it.Episode != nil && *it.Episode == *request.EpisodeNumber { out = append(out, c) } else if it.Episode == nil && it.Season != nil { packs = append(packs, c) }
			continue
		}
		out = append(out, c)
	}
	if len(out) == 0 && request.EpisodeNumber != nil { return packs }
	return out
}

func candidate(config providercore.Config, it item) providercore.Candidate {
	lang := "es-MX"; if spainSpanishRE.MatchString(it.Description) { lang = "es" }
	base := providercore.NewConfig(config).BaseURL(defaultBaseURL)
	page := it.PageURL; if page == "" { page = fmt.Sprintf("%s/api/subtitles/%d", base, it.ID) }
	return providercore.Candidate{ProviderName: "subx", LanguageID: lang, FileID: it.ID, Format: "srt", ReleaseName: strings.TrimSpace(it.Title + " | " + it.Description), DownloadCount: it.Downloads, SourceURL: fmt.Sprintf("%s/api/subtitles/%d/download", base, it.ID), SourceRef: page}
}

func newRequest(ctx context.Context, config providercore.Config, method, raw string, body io.Reader, _ bool) (*http.Request, error) {
	if !strings.HasPrefix(raw, "http") { raw = providercore.NewConfig(config).BaseURL(defaultBaseURL) + raw }
	req, err := http.NewRequestWithContext(ctx, method, raw, body)
	if err != nil { return nil, err }
	secret, _ := providercore.NewConfig(config).RequiredSecret("apiKey")
	req.Header.Set("Authorization", "Bearer "+secret); req.Header.Set("Accept", "application/json")
	return req, nil
}

func cleanTitle(s string) string { return strings.Join(strings.Fields(strings.NewReplacer(".", " ", "_", " ").Replace(s)), " ") }
