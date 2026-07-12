package subtitles

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type openSubtitlesAdapter struct {
	providerKey string
}

func (a openSubtitlesAdapter) Test(ctx context.Context, service *Service, config Config) error {
	return service.testOpenSubtitles(ctx, a.normalizedConfig(config))
}

func (a openSubtitlesAdapter) Search(
	ctx context.Context,
	service *Service,
	config Config,
	request SearchRequest,
) ([]Candidate, error) {
	return service.searchOpenSubtitles(ctx, a.normalizedConfig(config), request)
}

func (a openSubtitlesAdapter) Download(
	ctx context.Context,
	service *Service,
	config Config,
	candidate Candidate,
) (Download, error) {
	return service.downloadOpenSubtitles(ctx, a.normalizedConfig(config), candidate)
}

func (a openSubtitlesAdapter) normalizedConfig(config Config) Config {
	config.Type = a.providerKey
	return config
}

func (s *Service) testOpenSubtitles(ctx context.Context, config Config) error {
	if config.APIKey == nil || strings.TrimSpace(*config.APIKey) == "" {
		return ErrCredentialsRequired
	}
	base, err := url.Parse(strings.TrimSpace(config.BaseURL))
	if err != nil || base.Scheme == "" || base.Host == "" {
		return errors.New("subtitle provider base URL is invalid")
	}
	if err := validateProviderURL(config.Type, base.String(), false); err != nil {
		return err
	}
	endpoint := base.JoinPath("api", "v1", "infos", "languages")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Api-Key", strings.TrimSpace(*config.APIKey))
	req.Header.Set("User-Agent", "project-mema")
	resp, err := s.doProviderRequest(req, config.Type, false)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("subtitle provider returned HTTP %d", resp.StatusCode)
	}
	return nil
}

func (s *Service) searchOpenSubtitles(
	ctx context.Context,
	config Config,
	request SearchRequest,
) ([]Candidate, error) {
	endpoint, err := openSubtitlesEndpoint(config.BaseURL, "subtitles")
	if err != nil {
		return nil, err
	}
	values := endpoint.Query()
	values.Set("query", request.Title)
	values.Set("languages", openSubtitlesLanguage(request.LanguageID))
	if request.Year != nil {
		values.Set("year", strconv.Itoa(int(*request.Year)))
	}
	if request.SeasonNumber != nil {
		values.Set("season_number", strconv.Itoa(int(*request.SeasonNumber)))
	}
	if request.EpisodeNumber != nil {
		values.Set("episode_number", strconv.Itoa(int(*request.EpisodeNumber)))
	}
	endpoint.RawQuery = values.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	openSubtitlesHeaders(req, config)
	resp, err := s.doProviderRequest(req, config.Type, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("subtitle search returned HTTP %d", resp.StatusCode)
	}
	var payload openSubtitlesSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return openSubtitlesCandidates(config.Name, request.LanguageID, payload), nil
}

func (s *Service) downloadOpenSubtitles(
	ctx context.Context,
	config Config,
	candidate Candidate,
) (Download, error) {
	endpoint, err := openSubtitlesEndpoint(config.BaseURL, "download")
	if err != nil {
		return Download{}, err
	}
	token, err := s.openSubtitlesToken(ctx, config)
	if err != nil {
		return Download{}, err
	}
	body, _ := json.Marshal(map[string]any{"file_id": candidate.FileID})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return Download{}, err
	}
	openSubtitlesHeaders(req, config)
	openSubtitlesAuth(req, token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.doProviderRequest(req, config.Type, false)
	if err != nil {
		return Download{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Download{}, fmt.Errorf("subtitle download returned HTTP %d", resp.StatusCode)
	}
	var payload openSubtitlesDownloadResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return Download{}, err
	}
	if strings.TrimSpace(payload.Link) == "" {
		return Download{}, errors.New("subtitle provider returned no download link")
	}
	if err := validateProviderURL(config.Type, payload.Link, true); err != nil {
		return Download{}, err
	}
	return s.fetchSubtitle(ctx, config, payload.Link)
}

func (s *Service) openSubtitlesToken(ctx context.Context, config Config) (string, error) {
	if config.Username == nil || strings.TrimSpace(*config.Username) == "" ||
		config.Password == nil || strings.TrimSpace(*config.Password) == "" {
		return "", nil
	}
	endpoint, err := openSubtitlesEndpoint(config.BaseURL, "login")
	if err != nil {
		return "", err
	}
	body, _ := json.Marshal(map[string]string{
		"username": strings.TrimSpace(*config.Username),
		"password": strings.TrimSpace(*config.Password),
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	openSubtitlesHeaders(req, config)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.doProviderRequest(req, config.Type, false)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("subtitle provider login returned HTTP %d", resp.StatusCode)
	}
	var payload openSubtitlesLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if strings.TrimSpace(payload.Token) == "" {
		return "", errors.New("subtitle provider returned no login token")
	}
	return strings.TrimSpace(payload.Token), nil
}

func (s *Service) fetchSubtitle(ctx context.Context, config Config, link string) (Download, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return Download{}, err
	}
	openSubtitlesHeaders(req, config)
	resp, err := s.doProviderRequest(req, config.Type, true)
	if err != nil {
		return Download{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Download{}, fmt.Errorf("subtitle file returned HTTP %d", resp.StatusCode)
	}
	const subtitleDownloadLimit = 10 << 20
	content, err := io.ReadAll(io.LimitReader(resp.Body, subtitleDownloadLimit+1))
	if err != nil {
		return Download{}, err
	}
	if len(content) > subtitleDownloadLimit {
		return Download{}, errors.New("subtitle file is too large")
	}
	return Download{Content: content, URL: link}, nil
}
