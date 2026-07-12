package httpapi

import (
	"net/http"

	"media-manager/internal/subtitles"
	"media-manager/internal/subtitles/catalog"
)

func (s *Server) ListSubtitleProviderCatalog(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	entries, err := catalog.All()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "subtitle_catalog_failed", "Could not load subtitle provider catalog")
		return
	}
	response := SubtitleProviderCatalogResponse{Providers: make([]SubtitleProviderCatalogEntry, 0, len(entries))}
	for _, entry := range entries {
		response.Providers = append(response.Providers, subtitleProviderCatalogEntry(entry))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) TestDraftSubtitleProvider(w http.ResponseWriter, r *http.Request) {
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
	service := s.subtitles
	if service == nil {
		service = subtitles.NewService(nil)
	}
	result := service.Test(r.Context(), subtitleProviderConfigFromInput(input))
	writeJSON(w, http.StatusOK, IntegrationTestResponse{
		Success:   result.Success,
		Message:   result.Message,
		CheckedAt: s.now(),
		LatencyMs: int32(result.Latency.Milliseconds()),
		Details:   result.Details,
	})
}
