package indexers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (s *Service) searchCardigann(ctx context.Context, config Config, query string, mediaType string) ([]Release, error) {
	def, err := s.loader.load(ctx, config.DefinitionID)
	if err != nil {
		return nil, err
	}
	templateCtx := newCardigannContext(def, config, query, mediaType)
	templateCtx.Keywords = applyCardigannFilters(templateCtx.Keywords, def.Search.KeywordFilters, templateCtx)
	requests, err := cardigannSearchRequests(def, config, templateCtx)
	if err != nil {
		return nil, err
	}
	releases := []Release{}
	for _, request := range requests {
		body, err := s.executeCardigannRequest(ctx, request)
		if err != nil {
			return nil, err
		}
		found, err := s.parseCardigannResponse(ctx, def, config, templateCtx, request, body)
		if err != nil {
			return nil, err
		}
		releases = append(releases, found...)
	}
	return releases, nil
}

func (s *Service) testCardigann(ctx context.Context, config Config) TestResult {
	def, err := s.loader.load(ctx, config.DefinitionID)
	if err != nil {
		return failedResult("Indexer definition could not be loaded", "error", err.Error())
	}
	endpoint := config.BaseURL
	if def.Login != nil && def.Login.Test != nil && def.Login.Test.Path != "" {
		if resolved, err := resolveCardigannURL(config.BaseURL, def.Login.Test.Path); err == nil {
			endpoint = resolved
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return failedResult("Invalid indexer request", "error", err.Error())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return requestFailedResult(err)
	}
	defer closeBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return statusFailedResult(resp.StatusCode)
	}
	return successResult("Cardigann indexer reachable", "definitionId", def.ID)
}

func (s *Service) parseCardigannResponse(
	ctx context.Context,
	def cardigannDefinition,
	config Config,
	templateCtx cardigannContext,
	request cardigannSearchRequest,
	body []byte,
) ([]Release, error) {
	if isCardigannJSONResponse(def, request) {
		return s.parseCardigannJSON(ctx, def, config, templateCtx, request, body)
	}
	return s.parseCardigannHTML(ctx, def, config, templateCtx, request, body)
}

func isCardigannJSONResponse(def cardigannDefinition, request cardigannSearchRequest) bool {
	if request.SearchPath.Response != nil && strings.EqualFold(request.SearchPath.Response.Type, "json") {
		return true
	}
	return strings.HasPrefix(strings.TrimSpace(def.Search.Rows.Selector), "$")
}

func (s *Service) parseCardigannHTML(
	ctx context.Context,
	def cardigannDefinition,
	config Config,
	templateCtx cardigannContext,
	request cardigannSearchRequest,
	body []byte,
) ([]Release, error) {
	selector, err := renderCardigannTemplate(def.Search.Rows.Selector, templateCtx)
	if err != nil {
		return nil, err
	}
	rows, err := cardigannHTMLRows(body, selector)
	if err != nil {
		return nil, err
	}
	if def.Search.Rows.After > 0 && len(rows) > def.Search.Rows.After {
		rows = rows[def.Search.Rows.After:]
	}
	releases := make([]Release, 0, len(rows))
	for _, row := range rows {
		release, ok, err := s.releaseFromHTMLRow(ctx, def, config, templateCtx, request.URL, row)
		if err != nil {
			return nil, err
		}
		if ok {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func (s *Service) parseCardigannJSON(
	ctx context.Context,
	def cardigannDefinition,
	config Config,
	templateCtx cardigannContext,
	request cardigannSearchRequest,
	body []byte,
) ([]Release, error) {
	selector, err := renderCardigannTemplate(def.Search.Rows.Selector, templateCtx)
	if err != nil {
		return nil, err
	}
	rows := cardigannJSONRows(body, selector, def.Search.Rows)
	releases := make([]Release, 0, len(rows))
	for _, row := range rows {
		release, ok, err := s.releaseFromJSONRow(ctx, def, config, templateCtx, request.URL, row)
		if err != nil {
			return nil, err
		}
		if ok {
			releases = append(releases, release)
		}
	}
	return releases, nil
}

func baseCardigannRelease(config Config) Release {
	return Release{IndexerID: config.ID, IndexerName: config.Name, IndexerProtocol: config.Protocol}
}

func isOptionalCardigannField(name string, selector cardigannSelector) bool {
	if selector.Optional || strings.Contains(name, "_optional") || strings.HasPrefix(name, "_") {
		return true
	}
	switch canonicalCardigannField(name) {
	case "imdb", "imdbid", "tmdbid", "rageid", "tvdbid", "poster", "banner", "description", "genre":
		return true
	default:
		return false
	}
}

func cardigannFieldError(name string, err error) error {
	return fmt.Errorf("parse field %s: %w", name, err)
}
