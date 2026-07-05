package httpapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func (s *Server) ListDiscoverBlacklist(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	items, err := s.settings.ListDiscoverBlacklist(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "discover_blacklist_list_failed", "Could not list discover blacklist")
		return
	}
	writeJSON(w, http.StatusOK, discoverBlacklistResponse(items))
}

func (s *Server) AddDiscoverBlacklistItem(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body DiscoverBlacklistRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := discoverBlacklistInput(w, body)
	if !ok {
		return
	}
	item, err := s.settings.SaveDiscoverBlacklistItem(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not add media to discover blacklist")
		return
	}
	writeJSON(w, http.StatusCreated, discoverBlacklistItemResponse(item))
}

func (s *Server) DeleteDiscoverBlacklistItem(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if err := s.settings.DeleteDiscoverBlacklistItem(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not remove media from discover blacklist")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func discoverBlacklistInput(
	w http.ResponseWriter,
	request DiscoverBlacklistRequest,
) (storage.DiscoverBlacklistInput, bool) {
	title := strings.Join(strings.Fields(request.Title), " ")
	if title == "" {
		writeError(w, http.StatusBadRequest, "invalid_title", "Title is required")
		return storage.DiscoverBlacklistInput{}, false
	}
	if request.Type != MediaTypeMovie && request.Type != MediaTypeSerie {
		writeError(w, http.StatusBadRequest, "invalid_type", "Media type is required")
		return storage.DiscoverBlacklistInput{}, false
	}
	return storage.DiscoverBlacklistInput{
		Type:             string(request.Type),
		Title:            title,
		Year:             request.Year,
		ExternalProvider: request.ExternalProvider,
		ExternalID:       request.ExternalId,
		Overview:         request.Overview,
		PosterPath:       request.PosterPath,
	}, true
}

func discoverBlacklistResponse(items []storage.DiscoverBlacklistItem) DiscoverBlacklistResponse {
	response := DiscoverBlacklistResponse{Items: make([]DiscoverBlacklistItem, 0, len(items))}
	for _, item := range items {
		response.Items = append(response.Items, discoverBlacklistItemResponse(item))
	}
	return response
}

func discoverBlacklistItemResponse(item storage.DiscoverBlacklistItem) DiscoverBlacklistItem {
	return DiscoverBlacklistItem{
		Id:               openapi_types.UUID(item.ID),
		Type:             MediaType(item.Type),
		Title:            item.Title,
		Year:             item.Year,
		ExternalProvider: item.ExternalProvider,
		ExternalId:       item.ExternalID,
		Overview:         item.Overview,
		PosterPath:       item.PosterPath,
		CreatedAt:        item.CreatedAt,
	}
}

func filterDiscoverBlacklist(
	results []MediaSearchResult,
	blacklist []storage.DiscoverBlacklistItem,
) []MediaSearchResult {
	if len(results) == 0 || len(blacklist) == 0 {
		return results
	}
	keys := discoverBlacklistKeys(blacklist)
	filtered := make([]MediaSearchResult, 0, len(results))
	for _, result := range results {
		if key := discoverResultExternalKey(result); key != "" {
			if _, ok := keys[key]; ok {
				continue
			}
		}
		if _, ok := keys[discoverResultTitleKey(result)]; ok {
			continue
		}
		filtered = append(filtered, result)
	}
	return filtered
}

func discoverBlacklistKeys(items []storage.DiscoverBlacklistItem) map[string]struct{} {
	keys := map[string]struct{}{}
	for _, item := range items {
		if item.ExternalProvider != nil && item.ExternalID != nil {
			keys[discoverExternalKey(item.Type, *item.ExternalProvider, *item.ExternalID)] = struct{}{}
		}
		keys[discoverTitleKey(item.Type, item.Title, item.Year)] = struct{}{}
	}
	return keys
}

func discoverResultExternalKey(result MediaSearchResult) string {
	if result.ExternalProvider == nil || result.ExternalId == nil {
		return ""
	}
	return discoverExternalKey(string(result.Type), *result.ExternalProvider, *result.ExternalId)
}

func discoverResultTitleKey(result MediaSearchResult) string {
	return discoverTitleKey(string(result.Type), result.Title, result.Year)
}

func discoverExternalKey(mediaType string, provider string, id string) string {
	return strings.ToLower(strings.TrimSpace(mediaType)) + ":external:" +
		strings.ToLower(strings.TrimSpace(provider)) + ":" + strings.ToLower(strings.TrimSpace(id))
}

func discoverTitleKey(mediaType string, title string, year *int32) string {
	yearValue := int32(0)
	if year != nil {
		yearValue = *year
	}
	return strings.ToLower(strings.TrimSpace(mediaType)) + ":title:" +
		strings.ToLower(strings.Join(strings.Fields(title), " ")) + ":" + fmt.Sprint(yearValue)
}
