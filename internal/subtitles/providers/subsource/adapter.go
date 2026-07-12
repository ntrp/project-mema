package subsource

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
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/security"
)

const key = "subsource"
const defaultBaseURL = "https://api.subsource.net"
const maxBody = 50 << 20

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

type titleResponse struct {
	Data []titleRow `json:"data"`
}
type titleRow struct {
	MovieID               int64 `json:"movieId"`
	Title, AlternateTitle string
	ReleaseYear           any `json:"releaseYear"`
}
type subtitleResponse struct {
	Success bool          `json:"success"`
	Data    []subtitleRow `json:"data"`
}
type subtitleRow struct {
	SubtitleID  int64  `json:"subtitleId"`
	Language    string `json:"language"`
	ReleaseInfo any    `json:"releaseInfo"`
	Link        string `json:"link"`
}

func init() { providers.Register(key, Adapter) }

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	if _, ok := providercore.NewConfig(cfg).RequiredSecret("apiKey"); !ok {
		return fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	_, _, err := request(ctx, svc, cfg, "/api/v1/movies/search", url.Values{"searchType": {"text"}, "q": {"test"}}, false)
	return err
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if _, ok := providercore.NewConfig(cfg).RequiredSecret("apiKey"); !ok {
		return nil, fmt.Errorf("%w: apiKey is required", providercore.ErrProviderPrerequisiteMissing)
	}
	movieID, err := findTitle(ctx, svc, cfg, sr)
	if err != nil || movieID == 0 {
		return nil, err
	}
	query := url.Values{"language": {strings.ToLower(sr.LanguageID)}, "limit": {"100"}, "movieId": {strconv.FormatInt(movieID, 10)}}
	if sr.SeasonNumber != nil {
		query.Set("seasonNumber", strconv.Itoa(int(*sr.SeasonNumber)))
	}
	if sr.EpisodeNumber != nil {
		query.Set("episodeNumber", strconv.Itoa(int(*sr.EpisodeNumber)))
	}
	body, _, err := request(ctx, svc, cfg, "/api/v1/subtitles", query, false)
	if err != nil {
		return nil, err
	}
	var payload subtitleResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	out := make([]providercore.Candidate, 0, len(payload.Data))
	for _, row := range payload.Data {
		release := releaseName(row.ReleaseInfo)
		out = append(out, providercore.Candidate{ProviderName: key, LanguageID: firstNonEmpty(strings.ToLower(row.Language), sr.LanguageID), FileID: row.SubtitleID, Format: "srt", ReleaseName: release, SourceURL: endpoint(cfg, fmt.Sprintf("/api/v1/subtitles/%d/download", row.SubtitleID)), SourceRef: "https://subsource.net" + row.Link})
	}
	return out, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	raw := cand.SourceURL
	if raw == "" && cand.FileID != 0 {
		raw = endpoint(cfg, fmt.Sprintf("/api/v1/subtitles/%d/download", cand.FileID))
	}
	if raw == "" {
		return providercore.Download{}, fmt.Errorf("%w: missing download URL", providercore.ErrProviderPrerequisiteMissing)
	}
	body, _, err := requestURL(ctx, svc, cfg, raw, url.Values{}, true)
	if err != nil {
		return providercore.Download{}, err
	}
	member, err := providercore.ExtractSubtitle("subtitle.zip", body, security.ArchiveLimits{})
	if err != nil {
		return providercore.Download{}, err
	}
	return providercore.Download{Content: member.Content, URL: raw}, nil
}

func findTitle(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) (int64, error) {
	query := url.Values{}
	imdb := sr.MediaContext.ExternalIDs["imdb"]
	if imdb != "" {
		query.Set("searchType", "imdb")
		query.Set("imdb", imdb)
	} else {
		query.Set("searchType", "text")
		query.Set("q", strings.ToLower(sr.Title))
	}
	if sr.SeasonNumber != nil {
		query.Set("season", strconv.Itoa(int(*sr.SeasonNumber)))
	}
	rows, err := searchTitles(ctx, svc, cfg, query)
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 && imdb != "" {
		query.Del("imdb")
		query.Set("searchType", "text")
		query.Set("q", strings.ToLower(sr.Title))
		rows, err = searchTitles(ctx, svc, cfg, query)
	}
	if err != nil {
		return 0, err
	}
	for _, row := range rows {
		if titleMatches(sr.Title, row.Title, row.AlternateTitle) && yearMatches(sr.Year, row.ReleaseYear) {
			return row.MovieID, nil
		}
	}
	return 0, nil
}

func searchTitles(ctx context.Context, svc providercore.Service, cfg providercore.Config, query url.Values) ([]titleRow, error) {
	body, _, err := request(ctx, svc, cfg, "/api/v1/movies/search", query, false)
	if err != nil {
		return nil, err
	}
	var payload titleResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func request(ctx context.Context, svc providercore.Service, cfg providercore.Config, path string, query url.Values, download bool) ([]byte, *http.Response, error) {
	return requestURL(ctx, svc, cfg, endpoint(cfg, path), query, download)
}
func requestURL(ctx context.Context, svc providercore.Service, cfg providercore.Config, raw string, query url.Values, download bool) ([]byte, *http.Response, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, nil, err
	}
	q := u.Query()
	for k, values := range query {
		for _, value := range values {
			q.Add(k, value)
		}
	}
	q.Set("api_key", providercore.NewConfig(cfg).Secret("apiKey"))
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := svc.DoProviderRequest(req, key, download)
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

func endpoint(cfg providercore.Config, path string) string {
	return strings.TrimRight(providercore.NewConfig(cfg).BaseURL(defaultBaseURL), "/") + path
}
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
func titleMatches(want string, values ...string) bool {
	want = strings.ToLower(strings.TrimSpace(want))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value != "" && (strings.Contains(value, want) || strings.Contains(want, value)) {
			return true
		}
	}
	return false
}
func yearMatches(want *int32, raw any) bool {
	if want == nil {
		return true
	}
	switch value := raw.(type) {
	case float64:
		return int32(value) == *want
	case string:
		parsed, _ := strconv.Atoi(value)
		return int32(parsed) == *want
	}
	return false
}
func releaseName(raw any) string {
	switch value := raw.(type) {
	case string:
		return value
	case []any:
		var parts []string
		for _, item := range value {
			if text, ok := item.(string); ok {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, " ")
	}
	return ""
}
