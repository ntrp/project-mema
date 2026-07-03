package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func (s *Server) ListIndexerAppProfiles(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	profiles := storage.DefaultIndexerAppProfiles()
	response := IndexerAppProfileListResponse{Profiles: make([]IndexerAppProfile, 0, len(profiles))}
	for _, profile := range profiles {
		response.Profiles = append(response.Profiles, indexerAppProfileResponse(profile))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) ListIndexerProxies(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	proxies, err := s.settings.ListIndexerProxies(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "indexer_proxy_list_failed", "Could not list indexer proxies")
		return
	}
	response := IndexerProxyListResponse{Proxies: make([]IndexerProxy, 0, len(proxies))}
	for _, proxy := range proxies {
		response.Proxies = append(response.Proxies, indexerProxyResponse(proxy))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateIndexerProxy(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	input, ok := indexerProxyInput(w, r)
	if !ok {
		return
	}
	proxy, err := s.settings.CreateIndexerProxy(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not create indexer proxy")
		return
	}
	writeJSON(w, http.StatusCreated, indexerProxyResponse(proxy))
}

func (s *Server) UpdateIndexerProxy(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	input, ok := indexerProxyInput(w, r)
	if !ok {
		return
	}
	proxy, err := s.settings.UpdateIndexerProxy(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update indexer proxy")
		return
	}
	writeJSON(w, http.StatusOK, indexerProxyResponse(proxy))
}

func (s *Server) DeleteIndexerProxy(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if err := s.settings.DeleteIndexerProxy(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete indexer proxy")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) TestIndexerProxy(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	proxy, err := s.settings.GetIndexerProxy(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find indexer proxy")
		return
	}
	result := testProxyLink(r.Context(), proxy.Link)
	writeJSON(w, http.StatusOK, integrationTestResponse(s.now(), result.success, result.message, result.latency, result.details))
}

func indexerProxyInput(w http.ResponseWriter, r *http.Request) (storage.IndexerProxyInput, bool) {
	var body IndexerProxyRequest
	if !decodeJSON(w, r, &body) {
		return storage.IndexerProxyInput{}, false
	}
	fields, err := json.Marshal(body.Fields)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_fields", "Indexer proxy fields are invalid")
		return storage.IndexerProxyInput{}, false
	}
	return storage.IndexerProxyInput{
		Name:                  body.Name,
		Implementation:        body.Implementation,
		Link:                  body.Link,
		Enabled:               body.Enabled,
		OnHealthIssue:         body.OnHealthIssue,
		SupportsOnHealthIssue: true,
		IncludeHealthWarnings: body.IncludeHealthWarnings,
		TestCommand:           "test",
		Fields:                fields,
	}, true
}

type proxyTestResult struct {
	success bool
	message string
	latency time.Duration
	details map[string]interface{}
}

func testProxyLink(ctx context.Context, link string) proxyTestResult {
	started := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return proxyTestResult{success: false, message: err.Error(), latency: time.Since(started)}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return proxyTestResult{success: false, message: err.Error(), latency: time.Since(started)}
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return proxyTestResult{
			success: false,
			message: resp.Status,
			latency: time.Since(started),
			details: map[string]interface{}{"statusCode": resp.StatusCode},
		}
	}
	return proxyTestResult{
		success: true,
		message: "Proxy link OK",
		latency: time.Since(started),
		details: map[string]interface{}{"statusCode": resp.StatusCode},
	}
}

func indexerAppProfileResponse(profile storage.IndexerAppProfile) IndexerAppProfile {
	return IndexerAppProfile{
		Id:                      profile.ID,
		Name:                    profile.Name,
		EnableRss:               profile.EnableRSS,
		EnableAutomaticSearch:   profile.EnableAutomaticSearch,
		EnableInteractiveSearch: profile.EnableInteractiveSearch,
	}
}

func indexerProxyResponse(proxy storage.IndexerProxy) IndexerProxy {
	return IndexerProxy{
		Id:                    openapi_types.UUID(proxy.ID),
		Name:                  proxy.Name,
		Implementation:        proxy.Implementation,
		Link:                  proxy.Link,
		Enabled:               proxy.Enabled,
		OnHealthIssue:         proxy.OnHealthIssue,
		SupportsOnHealthIssue: proxy.SupportsOnHealthIssue,
		IncludeHealthWarnings: proxy.IncludeHealthWarnings,
		TestCommand:           proxy.TestCommand,
		Fields:                indexerFieldValues(proxy.Fields),
		CreatedAt:             proxy.CreatedAt,
		UpdatedAt:             proxy.UpdatedAt,
	}
}
