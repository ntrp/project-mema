package metadata

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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
	if request.Section == "trending" {
		return s.discoverTMDBTrending(ctx, config, request)
	}
	mediaPath, ok := tmdbMediaPath(request.MediaType)
	if !ok || !tmdbSectionValid(request.MediaType, request.Section) {
		return nil, nil
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), mediaPath, request.Section)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("page", strconv.Itoa(discoverPage(request.Page)))
	endpoint = endpoint + "?" + values.Encode()

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

func (s *Service) discoverTMDBTrending(ctx context.Context, config Config, request DiscoverRequest) ([]SearchResult, error) {
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "trending", "all", "day")
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("page", strconv.Itoa(discoverPage(request.Page)))
	endpoint = endpoint + "?" + values.Encode()
	var payload tmdbSearchResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	limit := request.Limit
	if limit <= 0 || limit > 40 {
		limit = 20
	}
	return tmdbResults(payload.Results, "mixed", limit), nil
}

func discoverPage(page int) int {
	if page < 1 {
		return 1
	}
	return page
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
	appends := []string{"credits", "external_ids", "keywords", "recommendations", "similar", "videos"}
	if request.MediaType == "movie" {
		appends = append(appends, "release_dates")
	} else {
		appends = append(appends, "content_ratings")
	}
	values.Set("append_to_response", strings.Join(appends, ","))
	endpoint = endpoint + "?" + values.Encode()

	var payload tmdbDetails
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return Details{}, err
	}
	if request.MediaType == "series" {
		if err := s.loadTMDBSeasonEpisodes(ctx, config, mediaPath, externalID, &payload); err != nil {
			return Details{}, err
		}
	}
	return tmdbDetailsResult(payload, request.MediaType, externalID), nil
}

func (s *Service) Collection(ctx context.Context, config Config, collectionID string) (Collection, error) {
	if config.Type != "tmdb" {
		return Collection{}, ErrUnsupportedProvider
	}
	collectionID = strings.TrimSpace(collectionID)
	if collectionID == "" {
		return Collection{}, errors.New("collection id is required")
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "collection", collectionID)
	if err != nil {
		return Collection{}, err
	}

	var payload tmdbCollection
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return Collection{}, err
	}
	return Collection{
		ID:           strconv.FormatInt(payload.ID, 10),
		Name:         strings.TrimSpace(payload.Name),
		Overview:     optionalString(payload.Overview),
		PosterPath:   optionalString(payload.PosterPath),
		BackdropPath: optionalString(payload.BackdropPath),
		Parts:        tmdbResults(payload.Parts, "movie", len(payload.Parts)),
	}, nil
}

func (s *Service) loadTMDBSeasonEpisodes(ctx context.Context, config Config, mediaPath string, externalID string, details *tmdbDetails) error {
	for seasonIndex := range details.Seasons {
		season := &details.Seasons[seasonIndex]
		if season.SeasonNumber < 0 || season.EpisodeCount <= 0 {
			continue
		}
		endpoint, err := url.JoinPath(
			strings.TrimRight(config.BaseURL, "/"),
			mediaPath,
			externalID,
			"season",
			strconv.Itoa(int(season.SeasonNumber)),
		)
		if err != nil {
			return err
		}
		var payload tmdbSeasonDetails
		if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
			return err
		}
		season.Episodes = payload.Episodes
	}
	return nil
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
		resultType := mediaType
		if mediaType == "mixed" {
			resultType = tmdbResultMediaType(item.MediaType)
		}
		if resultType == "series" {
			title = strings.TrimSpace(item.Name)
			date = item.FirstAirDate
		}
		if title == "" || resultType == "" {
			continue
		}
		results = append(results, SearchResult{
			Title:            title,
			Type:             resultType,
			Year:             yearFromDate(date),
			ExternalProvider: "tmdb",
			ExternalID:       strconv.FormatInt(item.ID, 10),
			Overview:         optionalString(item.Overview),
			PosterPath:       optionalString(item.PosterPath),
			Popularity:       optionalFloat64(item.Popularity),
		})
		if len(results) >= limit {
			break
		}
	}
	return results
}

func tmdbResultMediaType(value string) string {
	switch value {
	case "movie":
		return "movie"
	case "tv":
		return "series"
	default:
		return ""
	}
}
