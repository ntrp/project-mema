package httpapi

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) RefreshMediaItemMetadata(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	item, err := s.settings.GetMediaItem(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	if item.ExternalProvider == nil || item.ExternalID == nil || strings.TrimSpace(*item.ExternalID) == "" {
		writeError(w, http.StatusBadRequest, "metadata_refresh_unavailable", "Media item has no metadata provider link")
		return
	}

	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	provider, ok := metadataProviderByType(providers, strings.TrimSpace(*item.ExternalProvider))
	if !ok {
		writeError(w, http.StatusNotFound, "metadata_provider_not_found", "Metadata provider is not configured")
		return
	}

	details, err := s.freshMetadataProviderDetails(r.Context(), provider, metadata.DetailsRequest{
		MediaType:  item.Type,
		ExternalID: strings.TrimSpace(*item.ExternalID),
	})
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}

	updated, err := s.settings.UpdateMediaItemMetadata(r.Context(), item.ID, refreshedMediaInput(item, details))
	if err != nil {
		writeSettingsError(w, err, "Could not refresh media metadata")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Media metadata refreshed", map[string]any{"mediaItemId": item.ID.String()})
	writeJSON(w, http.StatusOK, mediaItemResponse(updated))
}

func refreshedMediaInput(item storage.MediaItem, details metadata.Details) storage.MediaItemInput {
	input := storage.MediaItemInput{
		Type:                  item.Type,
		Title:                 item.Title,
		Year:                  item.Year,
		Monitored:             item.Monitored,
		ExternalProvider:      item.ExternalProvider,
		ExternalID:            item.ExternalID,
		Overview:              item.Overview,
		PosterPath:            item.PosterPath,
		MediaMetadataSnapshot: item.MediaMetadataSnapshot,
		MonitorMode:           item.MonitorMode,
		SeriesType:            item.SeriesType,
		MinimumAvailability:   item.MinimumAvailability,
		QualityProfileID:      item.QualityProfileID,
		LibraryFolderID:       item.LibraryFolderID,
		Tags:                  item.Tags,
	}
	return applyMediaDetails(input, details)
}
