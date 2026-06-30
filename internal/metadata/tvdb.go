package metadata

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
