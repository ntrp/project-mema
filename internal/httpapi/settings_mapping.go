package httpapi

import (
	"errors"
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func downloadClientInput(w http.ResponseWriter, request DownloadClientRequest) (storage.DownloadClientInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.DownloadClientInput{}, false
	}
	if !request.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Download client type is not supported")
		return storage.DownloadClientInput{}, false
	}
	if baseURL == "" {
		writeError(w, http.StatusBadRequest, "invalid_base_url", "Base URL is required")
		return storage.DownloadClientInput{}, false
	}
	if request.Priority < 0 || request.Priority > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_priority", "Priority must be between 0 and 1000")
		return storage.DownloadClientInput{}, false
	}

	return storage.DownloadClientInput{
		Name:     name,
		Type:     string(request.Type),
		BaseURL:  baseURL,
		Username: optionalTrimmedString(request.Username),
		Password: optionalTrimmedString(request.Password),
		APIKey:   optionalTrimmedString(request.ApiKey),
		Category: optionalTrimmedString(request.Category),
		Enabled:  request.Enabled,
		Priority: request.Priority,
	}, true
}

func indexerInput(w http.ResponseWriter, request IndexerRequest) (storage.IndexerInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.IndexerInput{}, false
	}
	if !request.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Indexer type is not supported")
		return storage.IndexerInput{}, false
	}
	if baseURL == "" {
		writeError(w, http.StatusBadRequest, "invalid_base_url", "Base URL is required")
		return storage.IndexerInput{}, false
	}
	if request.Priority < 0 || request.Priority > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_priority", "Priority must be between 0 and 1000")
		return storage.IndexerInput{}, false
	}

	categories := []int32{}
	if request.Categories != nil {
		categories = append(categories, (*request.Categories)...)
	}

	return storage.IndexerInput{
		Name:       name,
		Type:       string(request.Type),
		BaseURL:    baseURL,
		APIKey:     optionalTrimmedString(request.ApiKey),
		Categories: categories,
		Enabled:    request.Enabled,
		Priority:   request.Priority,
	}, true
}

func downloadClientResponse(client storage.DownloadClient) DownloadClient {
	return DownloadClient{
		Id:        openapi_types.UUID(client.ID),
		Name:      client.Name,
		Type:      DownloadClientType(client.Type),
		BaseUrl:   client.BaseURL,
		Username:  client.Username,
		Password:  client.Password,
		ApiKey:    client.APIKey,
		Category:  client.Category,
		Enabled:   client.Enabled,
		Priority:  client.Priority,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}

func indexerResponse(indexer storage.Indexer) Indexer {
	categories := append([]int32(nil), indexer.Categories...)
	if categories == nil {
		categories = []int32{}
	}
	return Indexer{
		Id:         openapi_types.UUID(indexer.ID),
		Name:       indexer.Name,
		Type:       IndexerType(indexer.Type),
		BaseUrl:    indexer.BaseURL,
		ApiKey:     indexer.APIKey,
		Categories: &categories,
		Enabled:    indexer.Enabled,
		Priority:   indexer.Priority,
		CreatedAt:  indexer.CreatedAt,
		UpdatedAt:  indexer.UpdatedAt,
	}
}

func optionalTrimmedString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func writeSettingsError(w http.ResponseWriter, err error, message string) {
	if errors.Is(err, storage.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", message)
		return
	}
	writeError(w, http.StatusInternalServerError, "settings_update_failed", message)
}
