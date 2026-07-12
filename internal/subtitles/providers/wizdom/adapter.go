package wizdom

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/sitehtml"
)

const key = "wizdom"
const defaultBaseURL = "https://wizdom.xyz"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

type subtitle struct {
	ID      int64  `json:"id"`
	Version string `json:"version"`
}
type releaseResponse struct {
	Subs json.RawMessage `json:"subs"`
}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sitehtml.BaseURL(cfg, defaultBaseURL)+"/api/releases/tt0111161", nil)
	if err != nil {
		return err
	}
	return sitehtml.Test(req, svc, key)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if !sitehtml.Supports(sr.MediaType, "movie", "serie") {
		return nil, sitehtml.Unsupported(key, sr.MediaType)
	}
	imdb := firstNonEmpty(sr.MediaContext.ExternalIDs["imdb"], sr.MediaContext.ExternalIDs["imdb_id"])
	if imdb == "" {
		return nil, fmt.Errorf("%w: imdb id required", providercore.ErrProviderPrerequisiteMissing)
	}
	base := strings.TrimRight(sitehtml.BaseURL(cfg, defaultBaseURL), "/")
	endpoint := base + "/api/releases/" + imdb
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.DoProviderRequest(req, key, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, sitehtml.MaxHTMLBytes+1))
	if err != nil {
		return nil, err
	}
	rows, err := parseRows(body, sr)
	if err != nil {
		return nil, err
	}
	candidates := make([]providercore.Candidate, 0, len(rows))
	for _, row := range rows {
		candidates = append(candidates, providercore.Candidate{ProviderName: key, LanguageID: "heb", FileID: row.ID, Format: "srt", ReleaseName: row.Version, SourceURL: fmt.Sprintf("%s/api/files/sub/%d", base, row.ID), SourceRef: pageLink(base, imdb, sr.MediaType)})
	}
	return candidates, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	downloadURL := cand.SourceURL
	if downloadURL == "" && cand.FileID > 0 {
		downloadURL = fmt.Sprintf("%s/api/files/sub/%d", strings.TrimRight(sitehtml.BaseURL(cfg, defaultBaseURL), "/"), cand.FileID)
	}
	if strings.TrimSpace(downloadURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: candidate has no source URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	if cand.SourceRef != "" {
		req.Header.Set("Referer", cand.SourceRef)
	}
	return sitehtml.Download(req, svc, key, true, cand.ReleaseName+".zip")
}

func parseRows(body []byte, sr providercore.SearchRequest) ([]subtitle, error) {
	var payload releaseResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if sr.MediaType == "movie" {
		var rows []subtitle
		return rows, json.Unmarshal(payload.Subs, &rows)
	}
	if sr.SeasonNumber == nil || sr.EpisodeNumber == nil {
		return nil, fmt.Errorf("%w: season and episode required", providercore.ErrProviderPrerequisiteMissing)
	}
	var seasons any
	if err := json.Unmarshal(payload.Subs, &seasons); err != nil {
		return nil, err
	}
	season := indexed(seasons, int(*sr.SeasonNumber), strconv.Itoa(int(*sr.SeasonNumber)))
	episode := indexed(season, int(*sr.EpisodeNumber), strconv.Itoa(int(*sr.EpisodeNumber)))
	data, _ := json.Marshal(episode)
	var rows []subtitle
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, nil
	}
	return rows, nil
}

func indexed(value any, listIndex int, mapKey string) any {
	switch item := value.(type) {
	case []any:
		if listIndex >= 0 && listIndex < len(item) {
			return item[listIndex]
		}
	case map[string]any:
		return item[mapKey]
	}
	return nil
}
func pageLink(base, imdb, mediaType string) string {
	if mediaType == "movie" {
		return base + "/movies/" + imdb
	}
	return base + "/series/" + imdb
}
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
