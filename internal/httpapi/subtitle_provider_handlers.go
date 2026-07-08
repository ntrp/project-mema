package httpapi

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

func (s *Server) ListSubtitleProviders(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	providers, err := s.settings.ListSubtitleProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list subtitle providers")
		return
	}
	response := SubtitleProviderListResponse{Providers: make([]SubtitleProvider, 0, len(providers))}
	for _, provider := range providers {
		response.Providers = append(response.Providers, subtitleProviderResponse(provider))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateSubtitleProvider(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body SubtitleProviderRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := subtitleProviderInput(w, body)
	if !ok {
		return
	}
	provider, err := s.settings.CreateSubtitleProvider(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_create_failed", "Could not create subtitle provider")
		return
	}
	writeJSON(w, http.StatusCreated, subtitleProviderResponse(provider))
}

func (s *Server) UpdateSubtitleProvider(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body SubtitleProviderRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := subtitleProviderInput(w, body)
	if !ok {
		return
	}
	current, err := s.settings.GetSubtitleProvider(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not update subtitle provider")
		return
	}
	input = preserveSubtitleProviderSecrets(input, body, current)
	provider, err := s.settings.UpdateSubtitleProvider(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update subtitle provider")
		return
	}
	writeJSON(w, http.StatusOK, subtitleProviderResponse(provider))
}

func preserveSubtitleProviderSecrets(
	input storage.SubtitleProviderInput,
	request SubtitleProviderRequest,
	current storage.SubtitleProvider,
) storage.SubtitleProviderInput {
	if request.ApiKey == nil {
		input.APIKey = current.APIKey
	}
	if request.Password == nil {
		input.Password = current.Password
	}
	return input
}

func (s *Server) DeleteSubtitleProvider(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if err := s.settings.DeleteSubtitleProvider(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete subtitle provider")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) TestSubtitleProvider(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	provider, err := s.settings.GetSubtitleProvider(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find subtitle provider")
		return
	}
	service := s.subtitles
	if service == nil {
		service = subtitles.NewService(nil)
	}
	result := service.Test(r.Context(), subtitleProviderConfig(provider))
	writeJSON(w, http.StatusOK, IntegrationTestResponse{
		Success:   result.Success,
		Message:   result.Message,
		CheckedAt: s.now(),
		LatencyMs: int32(result.Latency.Milliseconds()),
		Details:   result.Details,
	})
}

func subtitleProviderInput(w http.ResponseWriter, request SubtitleProviderRequest) (storage.SubtitleProviderInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.SubtitleProviderInput{}, false
	}
	if !request.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Subtitle provider type is not supported")
		return storage.SubtitleProviderInput{}, false
	}
	if baseURL == "" {
		writeError(w, http.StatusBadRequest, "invalid_base_url", "Base URL is required")
		return storage.SubtitleProviderInput{}, false
	}
	if request.Priority < 0 || request.Priority > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_priority", "Priority must be between 0 and 1000")
		return storage.SubtitleProviderInput{}, false
	}
	mockRows, ok := subtitleProviderMockRowsInput(w, request)
	if !ok {
		return storage.SubtitleProviderInput{}, false
	}
	return storage.SubtitleProviderInput{
		Name:          name,
		Type:          string(request.Type),
		BaseURL:       baseURL,
		Username:      optionalTrimmedString(request.Username),
		Password:      optionalTrimmedString(request.Password),
		APIKey:        optionalTrimmedString(request.ApiKey),
		Enabled:       request.Enabled,
		Priority:      request.Priority,
		MockSubtitles: mockRows,
	}, true
}

func subtitleProviderResponse(provider storage.SubtitleProvider) SubtitleProvider {
	return SubtitleProvider{
		Id:          openapi_types.UUID(provider.ID),
		Name:        provider.Name,
		Type:        SubtitleProviderType(provider.Type),
		BaseUrl:     provider.BaseURL,
		Username:    provider.Username,
		Password:    provider.Password,
		ApiKey:      provider.APIKey,
		Enabled:     provider.Enabled,
		Priority:    provider.Priority,
		ApiKeySet:   provider.APIKey != nil,
		PasswordSet: provider.Password != nil,
		MockSubtitles: subtitleProviderMockRowsResponse(
			provider.MockSubtitles,
		),
		CreatedAt: provider.CreatedAt,
		UpdatedAt: provider.UpdatedAt,
	}
}

func subtitleProviderConfig(provider storage.SubtitleProvider) subtitles.Config {
	return subtitles.Config{
		Name:          provider.Name,
		Type:          provider.Type,
		BaseURL:       provider.BaseURL,
		Username:      provider.Username,
		Password:      provider.Password,
		APIKey:        provider.APIKey,
		MockSubtitles: subtitleProviderMockConfig(provider.MockSubtitles),
	}
}

func subtitleProviderMockConfig(rows []storage.MockSubtitleProviderRow) []subtitles.MockSubtitle {
	items := make([]subtitles.MockSubtitle, 0, len(rows))
	for _, row := range rows {
		items = append(items, subtitles.MockSubtitle{
			Title:      row.Title,
			LanguageID: row.LanguageID,
			Format:     row.Format,
		})
	}
	return items
}

func subtitleProviderMockRowsInput(
	w http.ResponseWriter,
	request SubtitleProviderRequest,
) ([]storage.MockSubtitleProviderRowInput, bool) {
	if request.MockSubtitles == nil {
		return nil, true
	}
	if request.Type != Mock {
		writeError(w, http.StatusBadRequest, "invalid_mock_subtitles", "Mock subtitles are only supported by mock providers")
		return nil, false
	}
	rows := make([]storage.MockSubtitleProviderRowInput, 0, len(*request.MockSubtitles))
	for _, row := range *request.MockSubtitles {
		title := strings.TrimSpace(row.Title)
		languageID := strings.TrimSpace(row.LanguageId)
		format := subtitleProviderMockFormat(row.Format)
		if title == "" || languageID == "" {
			writeError(w, http.StatusBadRequest, "invalid_mock_subtitle", "Mock subtitle title and language are required")
			return nil, false
		}
		if format == "" {
			writeError(w, http.StatusBadRequest, "invalid_mock_subtitle_format", "Mock subtitle format is not supported")
			return nil, false
		}
		rows = append(rows, storage.MockSubtitleProviderRowInput{
			Title:      title,
			LanguageID: languageID,
			Format:     format,
		})
	}
	return rows, true
}

func subtitleProviderMockRowsResponse(rows []storage.MockSubtitleProviderRow) []MockSubtitleProviderRow {
	response := make([]MockSubtitleProviderRow, 0, len(rows))
	for _, row := range rows {
		response = append(response, MockSubtitleProviderRow{
			Id:         openapi_types.UUID(row.ID),
			Title:      row.Title,
			LanguageId: row.LanguageID,
			Format:     row.Format,
		})
	}
	return response
}

func subtitleProviderMockFormat(value string) string {
	switch strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".") {
	case "", "srt", "subrip":
		return "subrip"
	case "vtt", "webvtt":
		return "vtt"
	case "ass", "ssa":
		return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".")
	default:
		return ""
	}
}
