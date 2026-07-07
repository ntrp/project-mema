package metadata

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func (s *Service) detailsTVDB(ctx context.Context, config Config, request DetailsRequest) (Details, error) {
	token, err := s.tvdbToken(ctx, config)
	if err != nil {
		return Details{}, err
	}
	config.AccessToken = &token

	path, ok := tvdbDetailsPath(request.MediaType)
	if !ok {
		return Details{}, ErrUnsupportedProvider
	}
	externalID := strings.TrimSpace(request.ExternalID)
	if externalID == "" {
		return Details{}, errors.New("external id is required")
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), path, externalID, "extended")
	if err != nil {
		return Details{}, err
	}
	values := url.Values{}
	values.Set("meta", "translations")
	endpoint = endpoint + "?" + values.Encode()

	var payload tvdbDetailsResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return Details{}, err
	}
	if strings.TrimSpace(tvdbOverview(payload.Data)) == "" {
		if err := s.loadTVDBTranslation(ctx, config, path, externalID, &payload.Data); err != nil {
			return Details{}, err
		}
	}
	return tvdbDetailsResult(payload.Data, request.MediaType, externalID), nil
}

func (s *Service) loadTVDBTranslation(ctx context.Context, config Config, path string, externalID string, details *tvdbDetails) error {
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), path, externalID, "translations", "eng")
	if err != nil {
		return err
	}
	var payload tvdbTranslationResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		var providerErr ProviderHTTPError
		if errors.As(err, &providerErr) && providerErr.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	if details.Overview == "" {
		details.Overview = payload.Data.Overview
	}
	if details.Name == "" {
		details.Name = payload.Data.Name
	}
	return nil
}

func tvdbDetailsPath(mediaType string) (string, bool) {
	switch mediaType {
	case "movie":
		return "movies", true
	case "serie":
		return "series", true
	default:
		return "", false
	}
}
