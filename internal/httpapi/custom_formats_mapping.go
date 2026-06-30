package httpapi

import (
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func customFormatInput(w http.ResponseWriter, request CustomFormatRequest) (storage.CustomFormatInput, bool) {
	name := strings.Join(strings.Fields(request.Name), " ")
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.CustomFormatInput{}, false
	}
	includeSpecs, ok := customFormatSpecsInput(w, request.IncludeSpecs)
	if !ok {
		return storage.CustomFormatInput{}, false
	}
	excludeSpecs, ok := customFormatSpecsInput(w, request.ExcludeSpecs)
	if !ok {
		return storage.CustomFormatInput{}, false
	}
	if len(includeSpecs) == 0 && len(excludeSpecs) == 0 {
		writeError(w, http.StatusBadRequest, "spec_required", "Add at least one condition")
		return storage.CustomFormatInput{}, false
	}
	return storage.CustomFormatInput{Name: name, IncludeSpecs: includeSpecs, ExcludeSpecs: excludeSpecs}, true
}

func customFormatSpecsInput(w http.ResponseWriter, specs []CustomFormatSpec) ([]storage.CustomFormatSpec, bool) {
	inputs := make([]storage.CustomFormatSpec, 0, len(specs))
	seen := map[string]struct{}{}
	for _, spec := range specs {
		input, ok := customFormatSpecInput(w, spec)
		if !ok {
			return nil, false
		}
		if _, exists := seen[input.ID]; exists {
			writeError(w, http.StatusBadRequest, "duplicate_spec", "Condition IDs must be unique")
			return nil, false
		}
		seen[input.ID] = struct{}{}
		inputs = append(inputs, input)
	}
	return inputs, true
}

func customFormatSpecInput(w http.ResponseWriter, spec CustomFormatSpec) (storage.CustomFormatSpec, bool) {
	input := storage.CustomFormatSpec{
		ID:       strings.TrimSpace(spec.Id),
		Name:     strings.Join(strings.Fields(spec.Name), " "),
		Type:     string(spec.Type),
		Value:    strings.TrimSpace(spec.Value),
		Required: spec.Required,
	}
	if input.ID == "" || input.Name == "" || input.Value == "" {
		writeError(w, http.StatusBadRequest, "invalid_spec", "Condition ID, name, and value are required")
		return storage.CustomFormatSpec{}, false
	}
	if !spec.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_spec_type", "Condition type is not supported")
		return storage.CustomFormatSpec{}, false
	}
	return input, true
}

func customFormatListResponse(formats []storage.CustomFormat) CustomFormatListResponse {
	response := CustomFormatListResponse{Formats: make([]CustomFormat, 0, len(formats))}
	for _, format := range formats {
		response.Formats = append(response.Formats, customFormatResponse(format))
	}
	return response
}

func customFormatResponse(format storage.CustomFormat) CustomFormat {
	return CustomFormat{
		Id:           openapi_types.UUID(format.ID),
		Name:         format.Name,
		IncludeSpecs: customFormatSpecResponses(format.IncludeSpecs),
		ExcludeSpecs: customFormatSpecResponses(format.ExcludeSpecs),
		CreatedAt:    format.CreatedAt,
		UpdatedAt:    format.UpdatedAt,
	}
}

func customFormatSpecResponses(specs []storage.CustomFormatSpec) []CustomFormatSpec {
	response := make([]CustomFormatSpec, 0, len(specs))
	for _, spec := range specs {
		response = append(response, CustomFormatSpec{
			Id:       spec.ID,
			Name:     spec.Name,
			Type:     CustomFormatSpecType(spec.Type),
			Value:    spec.Value,
			Required: spec.Required,
		})
	}
	return response
}
