package metadata

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
	"time"
)

var (
	ErrCredentialsRequired = errors.New("metadata provider credentials are required")
	ErrUnsupportedProvider = errors.New("unsupported metadata provider")
)

type ProviderHTTPError struct {
	StatusCode int
}

func (e ProviderHTTPError) Error() string {
	return fmt.Sprintf("metadata provider returned HTTP %d", e.StatusCode)
}

func (s *Service) Search(ctx context.Context, config Config, request SearchRequest) ([]SearchResult, error) {
	switch config.Type {
	case "tmdb":
		if request.MediaType != "movie" && request.MediaType != "series" {
			return nil, nil
		}
		return s.searchTMDB(ctx, config, request)
	case "tvdb":
		return s.searchTVDB(ctx, config, request)
	default:
		return nil, ErrUnsupportedProvider
	}
}

func (s *Service) Discover(ctx context.Context, config Config, request DiscoverRequest) ([]SearchResult, error) {
	if config.Type != "tmdb" {
		return nil, ErrUnsupportedProvider
	}
	return s.discoverTMDB(ctx, config, request)
}

func (s *Service) Details(ctx context.Context, config Config, request DetailsRequest) (Details, error) {
	if config.Type != "tmdb" {
		return Details{}, ErrUnsupportedProvider
	}
	return s.detailsTMDB(ctx, config, request)
}

func (s *Service) Test(ctx context.Context, config Config) TestResult {
	start := time.Now()
	results, err := s.Search(ctx, config, SearchRequest{Query: "test", MediaType: "movie"})
	latency := time.Since(start)
	if err != nil {
		return TestResult{
			Success: false,
			Message: err.Error(),
			Latency: latency,
			Details: map[string]interface{}{"provider": config.Type},
		}
	}
	return TestResult{
		Success: true,
		Message: "Metadata provider connection OK",
		Latency: latency,
		Details: map[string]interface{}{
			"provider": config.Type,
			"results":  len(results),
		},
	}
}

func (s *Service) searchTMDB(ctx context.Context, config Config, request SearchRequest) ([]SearchResult, error) {
	mediaPath, ok := tmdbMediaPath(request.MediaType)
	if !ok {
		return nil, nil
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "search", mediaPath)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("query", request.Query)
	if request.Year != nil {
		if request.MediaType == "series" {
			values.Set("first_air_date_year", strconv.Itoa(int(*request.Year)))
		} else {
			values.Set("year", strconv.Itoa(int(*request.Year)))
		}
	}
	endpoint = endpoint + "?" + values.Encode()

	var payload tmdbSearchResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	return tmdbResults(payload.Results, request.MediaType, 20), nil
}

func (s *Service) discoverTMDB(ctx context.Context, config Config, request DiscoverRequest) ([]SearchResult, error) {
	mediaPath, ok := tmdbMediaPath(request.MediaType)
	if !ok || !tmdbSectionValid(request.MediaType, request.Section) {
		return nil, nil
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), mediaPath, request.Section)
	if err != nil {
		return nil, err
	}

	var payload tmdbSearchResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	limit := request.Limit
	if limit <= 0 || limit > 20 {
		limit = 20
	}
	return tmdbResults(payload.Results, request.MediaType, limit), nil
}

func (s *Service) detailsTMDB(ctx context.Context, config Config, request DetailsRequest) (Details, error) {
	mediaPath, ok := tmdbMediaPath(request.MediaType)
	if !ok {
		return Details{}, ErrUnsupportedProvider
	}
	externalID := strings.TrimSpace(request.ExternalID)
	if externalID == "" {
		return Details{}, errors.New("external id is required")
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), mediaPath, externalID)
	if err != nil {
		return Details{}, err
	}
	values := url.Values{}
	values.Set("append_to_response", "credits")
	endpoint = endpoint + "?" + values.Encode()

	var payload tmdbDetails
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return Details{}, err
	}
	return tmdbDetailsResult(payload, request.MediaType, externalID), nil
}

func (s *Service) searchTVDB(ctx context.Context, config Config, request SearchRequest) ([]SearchResult, error) {
	token, err := s.tvdbToken(ctx, config)
	if err != nil {
		return nil, err
	}
	config.AccessToken = &token

	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "search")
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("query", request.Query)
	switch request.MediaType {
	case "movie":
		values.Set("type", "movie")
	case "series":
		values.Set("type", "series")
	}
	if request.Year != nil {
		values.Set("year", strconv.Itoa(int(*request.Year)))
	}
	values.Set("limit", "10")
	endpoint = endpoint + "?" + values.Encode()

	var payload tvdbSearchResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	results := make([]SearchResult, 0, len(payload.Data))
	for _, item := range payload.Data {
		mediaType := tvdbMediaType(item.Type)
		if mediaType == "" || mediaType != request.MediaType {
			continue
		}
		title := strings.TrimSpace(item.Name)
		if title == "" {
			title = strings.TrimSpace(item.Title)
		}
		if title == "" {
			continue
		}
		results = append(results, SearchResult{
			Title:            title,
			Type:             mediaType,
			Year:             yearFromString(item.Year),
			ExternalProvider: "tvdb",
			ExternalID:       firstNonEmpty(item.TVDBID, item.ID, item.ObjectID),
			Overview:         optionalString(firstNonEmpty(item.Overview, firstString(item.OverviewTranslated))),
			PosterPath:       optionalString(firstNonEmpty(item.ImageURL, item.Poster, item.Thumbnail)),
		})
	}
	return results, nil
}

func (s *Service) tvdbToken(ctx context.Context, config Config) (string, error) {
	if config.AccessToken != nil && strings.TrimSpace(*config.AccessToken) != "" {
		return strings.TrimSpace(*config.AccessToken), nil
	}
	if config.SessionToken != nil && config.SessionTokenExpiresAt != nil && time.Until(*config.SessionTokenExpiresAt) > 24*time.Hour {
		return *config.SessionToken, nil
	}
	if config.APIKey == nil || strings.TrimSpace(*config.APIKey) == "" {
		return "", errors.New("TVDB API key is required")
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "login")
	if err != nil {
		return "", err
	}
	body := map[string]string{"apikey": strings.TrimSpace(*config.APIKey)}
	if config.PIN != nil && strings.TrimSpace(*config.PIN) != "" {
		body["pin"] = strings.TrimSpace(*config.PIN)
	}
	var payload tvdbLoginResponse
	if err := s.doJSON(ctx, Config{ID: config.ID, Type: config.Type, BaseURL: config.BaseURL}, http.MethodPost, endpoint, body, &payload); err != nil {
		return "", err
	}
	token := strings.TrimSpace(payload.Data.Token)
	if token == "" {
		return "", errors.New("TVDB login did not return a token")
	}
	expiresAt := time.Now().Add(29 * 24 * time.Hour)
	if s.tokenStore != nil {
		_ = s.tokenStore.UpdateMetadataProviderSessionToken(ctx, config.ID, token, expiresAt)
	}
	return token, nil
}

func (s *Service) doJSON(ctx context.Context, config Config, method string, endpoint string, body any, target any) error {
	if err := s.wait(ctx, config); err != nil {
		return err
	}
	err := s.doJSONOnce(ctx, config, method, endpoint, body, target)
	if retry, wait := retryAfter(err); retry {
		if wait > 5*time.Second {
			return err
		}
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
		if waitErr := s.wait(ctx, config); waitErr != nil {
			return waitErr
		}
		return s.doJSONOnce(ctx, config, method, endpoint, body, target)
	}
	return err
}

func (s *Service) doJSONOnce(ctx context.Context, config Config, method string, endpoint string, body any, target any) error {
	if config.Type == "tmdb" &&
		(config.APIKey == nil || strings.TrimSpace(*config.APIKey) == "") &&
		(config.AccessToken == nil || strings.TrimSpace(*config.AccessToken) == "") {
		return ErrCredentialsRequired
	}

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if config.AccessToken != nil && strings.TrimSpace(*config.AccessToken) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(*config.AccessToken))
	} else if config.Type == "tmdb" && config.APIKey != nil && strings.TrimSpace(*config.APIKey) != "" {
		values := req.URL.Query()
		values.Set("api_key", strings.TrimSpace(*config.APIKey))
		req.URL.RawQuery = values.Encode()
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return rateLimitedError{retryAfter: parseRetryAfter(resp.Header.Get("Retry-After"))}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ProviderHTTPError{StatusCode: resp.StatusCode}
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func (s *Service) wait(ctx context.Context, config Config) error {
	interval := 500 * time.Millisecond
	if config.Type == "tvdb" {
		interval = time.Second
	}
	key := config.ID.String()
	s.mu.Lock()
	last := s.lastByID[key]
	wait := time.Until(last.Add(interval))
	if wait <= 0 {
		s.lastByID[key] = time.Now()
		s.mu.Unlock()
		return nil
	}
	s.mu.Unlock()

	timer := time.NewTimer(wait)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}

	s.mu.Lock()
	s.lastByID[key] = time.Now()
	s.mu.Unlock()
	return nil
}

type rateLimitedError struct {
	retryAfter time.Duration
}

func (e rateLimitedError) Error() string {
	return "metadata provider rate limit reached"
}

func retryAfter(err error) (bool, time.Duration) {
	var rateErr rateLimitedError
	if errors.As(err, &rateErr) {
		if rateErr.retryAfter <= 0 {
			return true, time.Second
		}
		return true, rateErr.retryAfter
	}
	return false, 0
}

func IsRateLimited(err error) bool {
	var rateErr rateLimitedError
	return errors.As(err, &rateErr)
}

func ProviderStatusCode(err error) (int, bool) {
	var providerErr ProviderHTTPError
	if errors.As(err, &providerErr) {
		return providerErr.StatusCode, true
	}
	return 0, false
}

func parseRetryAfter(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	seconds, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return time.Duration(seconds) * time.Second
	}
	when, err := http.ParseTime(value)
	if err != nil {
		return 0
	}
	return time.Until(when)
}

func yearFromDate(value string) *int32 {
	if len(value) < 4 {
		return nil
	}
	return yearFromString(value[:4])
}

func yearFromString(value string) *int32 {
	if len(value) < 4 {
		return nil
	}
	year, err := strconv.ParseInt(value[:4], 10, 32)
	if err != nil {
		return nil
	}
	result := int32(year)
	return &result
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstString(values []string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func tvdbMediaType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "movie":
		return "movie"
	case "series":
		return "series"
	default:
		return ""
	}
}

func tmdbMediaPath(mediaType string) (string, bool) {
	switch mediaType {
	case "movie":
		return "movie", true
	case "series":
		return "tv", true
	default:
		return "", false
	}
}

func tmdbSectionValid(mediaType string, section string) bool {
	switch mediaType {
	case "movie":
		return section == "popular" || section == "upcoming" || section == "top_rated"
	case "series":
		return section == "popular" || section == "on_the_air" || section == "top_rated"
	default:
		return false
	}
}

func tmdbResults(items []tmdbMedia, mediaType string, limit int) []SearchResult {
	results := make([]SearchResult, 0, min(len(items), limit))
	for _, item := range items {
		title := strings.TrimSpace(item.Title)
		date := item.ReleaseDate
		if mediaType == "series" {
			title = strings.TrimSpace(item.Name)
			date = item.FirstAirDate
		}
		if title == "" {
			continue
		}
		results = append(results, SearchResult{
			Title:            title,
			Type:             mediaType,
			Year:             yearFromDate(date),
			ExternalProvider: "tmdb",
			ExternalID:       strconv.FormatInt(item.ID, 10),
			Overview:         optionalString(item.Overview),
			PosterPath:       optionalString(item.PosterPath),
		})
		if len(results) >= limit {
			break
		}
	}
	return results
}

func tmdbDetailsResult(item tmdbDetails, mediaType string, externalID string) Details {
	title := strings.TrimSpace(item.Title)
	date := item.ReleaseDate
	if mediaType == "series" {
		title = strings.TrimSpace(item.Name)
		date = item.FirstAirDate
	}
	details := Details{
		Title:            title,
		Type:             mediaType,
		Year:             yearFromDate(date),
		ExternalProvider: "tmdb",
		ExternalID:       externalID,
		Overview:         optionalString(item.Overview),
		PosterPath:       optionalString(item.PosterPath),
		BackdropPath:     optionalString(item.BackdropPath),
		Status:           optionalString(item.Status),
		OriginalLanguage: optionalString(item.OriginalLanguage),
		Genres:           tmdbNames(item.Genres),
		Facts:            []Fact{},
		Seasons:          []Season{},
		Cast:             []Person{},
	}
	if mediaType == "movie" {
		details.ReleaseDate = optionalString(item.ReleaseDate)
		if item.Runtime > 0 {
			details.RuntimeMinutes = &item.Runtime
		}
	} else {
		details.FirstAirDate = optionalString(item.FirstAirDate)
		if item.NumberOfSeasons > 0 {
			details.SeasonCount = &item.NumberOfSeasons
		}
		if item.NumberOfEpisodes > 0 {
			details.EpisodeCount = &item.NumberOfEpisodes
		}
		if len(item.EpisodeRunTime) > 0 && item.EpisodeRunTime[0] > 0 {
			value := item.EpisodeRunTime[0]
			details.RuntimeMinutes = &value
		}
		for _, season := range item.Seasons {
			name := strings.TrimSpace(season.Name)
			if name == "" {
				continue
			}
			mapped := Season{
				Name:       name,
				AirDate:    optionalString(season.AirDate),
				PosterPath: optionalString(season.PosterPath),
			}
			if season.EpisodeCount > 0 {
				mapped.EpisodeCount = &season.EpisodeCount
			}
			details.Seasons = append(details.Seasons, mapped)
		}
	}
	if item.VoteAverage > 0 {
		details.VoteAverage = &item.VoteAverage
	}
	if len(item.CreatedBy) > 0 {
		details.Facts = append(details.Facts, Fact{Label: "Creator", Value: strings.Join(tmdbNames(item.CreatedBy), ", ")})
	}
	if len(item.Networks) > 0 {
		details.Facts = append(details.Facts, Fact{Label: "Networks", Value: strings.Join(tmdbNames(item.Networks), ", ")})
	}
	for _, cast := range item.Credits.Cast {
		name := strings.TrimSpace(cast.Name)
		if name == "" {
			continue
		}
		details.Cast = append(details.Cast, Person{
			Name:        name,
			Role:        optionalString(cast.Character),
			ProfilePath: optionalString(cast.ProfilePath),
		})
		if len(details.Cast) >= 14 {
			break
		}
	}
	return details
}

func tmdbNames(items []tmdbName) []string {
	names := []string{}
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}
