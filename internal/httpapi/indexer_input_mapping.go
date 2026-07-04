package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func indexerInput(
	w http.ResponseWriter,
	request IndexerRequest,
	languages []storage.Language,
) (storage.IndexerInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	definitionID := strings.TrimSpace(request.DefinitionId)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.IndexerInput{}, false
	}
	definition, found := indexers.CatalogEntryByID(definitionID)
	if !found {
		writeError(w, http.StatusBadRequest, "invalid_definition", "Indexer definition is not supported")
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
	if definition.Protocol == "usenet" && definition.SupportsRedirect && request.Redirect != nil && !*request.Redirect {
		writeError(w, http.StatusBadRequest, "invalid_redirect", "Redirect must be enabled for Usenet indexers")
		return storage.IndexerInput{}, false
	}

	categories := []int32{}
	if request.Categories != nil {
		categories = append(categories, (*request.Categories)...)
	}
	fieldValues := []IndexerFieldValue{}
	if request.Fields != nil {
		fieldValues = append(fieldValues, (*request.Fields)...)
	}
	fields, err := json.Marshal(fieldValues)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_fields", "Indexer fields are invalid")
		return storage.IndexerInput{}, false
	}
	capabilities, err := json.Marshal(catalogCapabilities(definition.Capabilities))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_capabilities", "Indexer capabilities are invalid")
		return storage.IndexerInput{}, false
	}

	return storage.IndexerInput{
		DefinitionID:       definition.DefinitionID,
		Name:               name,
		Implementation:     firstNonEmptyString(request.Implementation, definition.Implementation),
		ImplementationName: firstNonEmptyString(request.ImplementationName, definition.ImplementationName),
		Protocol:           definition.Protocol,
		Privacy:            definition.Privacy,
		Language:           newCatalogLanguageMapper(languages).code(definition.Language),
		Encoding:           optionalCatalogString(definition.Encoding),
		Description:        optionalCatalogString(definition.Description),
		IndexerURLs:        append([]string{}, definition.IndexerURLs...),
		LegacyURLs:         append([]string{}, definition.LegacyURLs...),
		BaseURL:            baseURL,
		APIKey:             optionalTrimmedString(request.ApiKey),
		Categories:         categories,
		Fields:             fields,
		Capabilities:       capabilities,
		Redirect:           optionalBool(request.Redirect, true),
		AppProfileID:       appProfileID(request.AppProfileId),
		MinimumSeeders:     request.MinimumSeeders,
		SeedRatio:          request.SeedRatio,
		SeedTime:           request.SeedTime,
		PackSeedTime:       request.PackSeedTime,
		PreferMagnetURL:    optionalBool(request.PreferMagnetUrl, false),
		SupportsRSS:        definition.SupportsRSS,
		SupportsSearch:     definition.SupportsSearch,
		SupportsRedirect:   definition.SupportsRedirect,
		SupportsPagination: definition.SupportsPagination,
		Enabled:            request.Enabled,
		Priority:           request.Priority,
	}, true
}

func firstNonEmptyString(value *string, fallback string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return fallback
	}
	return strings.TrimSpace(*value)
}

func optionalBool(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func appProfileID(value *string) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return "default"
}
