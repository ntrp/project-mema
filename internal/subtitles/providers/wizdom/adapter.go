package wizdom

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/sitehtml"
)

const key = "wizdom"
const defaultBaseURL = "https://wizdom.xyz"

var Adapter providercore.Adapter = adapter{}

type adapter struct{}

type subtitle struct {
	ID          int64  `json:"id"`
	VersionName string `json:"versioname"`
	VersionAlt  string `json:"versionName"`
	Language    string `json:"language"`
}

func (adapter) Test(ctx context.Context, svc providercore.Service, cfg providercore.Config) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sitehtml.BaseURL(cfg, defaultBaseURL)+"/api/search?action=by_id&imdb=tt0111161", nil)
	if err != nil {
		return err
	}
	return sitehtml.Test(req, svc, key)
}

func (adapter) Search(ctx context.Context, svc providercore.Service, cfg providercore.Config, sr providercore.SearchRequest) ([]providercore.Candidate, error) {
	if !sitehtml.Supports(sr.MediaType, "movie", "serie") {
		return nil, sitehtml.Unsupported(key, sr.MediaType)
	}
	imdb := sr.MediaContext.ExternalIDs["imdb"]
	if imdb == "" {
		imdb = sr.MediaContext.ExternalIDs["imdb_id"]
	}
	if imdb == "" {
		return nil, fmt.Errorf("%w: imdb id required", providercore.ErrProviderPrerequisiteMissing)
	}
	endpoint := sitehtml.BaseURL(cfg, defaultBaseURL) + "/api/search?action=by_id&imdb=" + url.QueryEscape(imdb)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := svc.DoProviderRequest(req, key, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("%w: http status %d", providercore.ErrProviderBrokenUpstream, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, sitehtml.MaxHTMLBytes+1))
	if err != nil {
		return nil, err
	}
	var rows []subtitle
	if err := json.Unmarshal(body, &rows); err != nil {
		return nil, err
	}
	var candidates []providercore.Candidate
	for _, row := range rows {
		lang := row.Language
		if lang == "" {
			lang = sr.LanguageID
		}
		release := row.VersionName
		if release == "" {
			release = row.VersionAlt
		}
		candidates = append(candidates, providercore.Candidate{ProviderName: key, LanguageID: lang, FileID: row.ID, Format: "srt", ReleaseName: release, SourceURL: fmt.Sprintf("%s/api/files/sub/%d", sitehtml.BaseURL(cfg, defaultBaseURL), row.ID), SourceRef: endpoint})
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("%w: no subtitles found", providercore.ErrProviderBrokenUpstream)
	}
	return candidates, nil
}

func (adapter) Download(ctx context.Context, svc providercore.Service, cfg providercore.Config, cand providercore.Candidate) (providercore.Download, error) {
	downloadURL := cand.SourceURL
	if downloadURL == "" && cand.FileID > 0 {
		downloadURL = fmt.Sprintf("%s/api/files/sub/%d", sitehtml.BaseURL(cfg, defaultBaseURL), cand.FileID)
	}
	if strings.TrimSpace(downloadURL) == "" {
		return providercore.Download{}, fmt.Errorf("%w: candidate has no source URL", providercore.ErrProviderPrerequisiteMissing)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return providercore.Download{}, err
	}
	return sitehtml.Download(req, svc, key, false, cand.ReleaseName)
}
