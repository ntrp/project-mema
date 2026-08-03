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
	providerType, externalID := mediaRefreshProviderLink(item)
	provider, ok := metadataProviderByType(providers, providerType)
	if !ok && providerType != strings.TrimSpace(*item.ExternalProvider) {
		providerType, externalID = strings.TrimSpace(*item.ExternalProvider), strings.TrimSpace(*item.ExternalID)
		provider, ok = metadataProviderByType(providers, providerType)
	}
	if !ok {
		writeError(w, http.StatusNotFound, "metadata_provider_not_found", "Metadata provider is not configured")
		return
	}

	primary, err := s.freshMetadataProviderDetails(r.Context(), provider, metadata.DetailsRequest{
		MediaType:  item.Type,
		ExternalID: externalID,
	})
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}

	resource := s.resourceWithMetadataFallback(
		r.Context(), providers, primary, s.freshMetadataProviderDetails, s.searchMetadataProvider,
	)
	input := refreshedMediaInput(item, resource)
	input, err = s.enrichAniDBFacts(r.Context(), item, providers, input, primary)
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}
	updated, err := s.settings.UpdateMediaItemMetadata(r.Context(), item.ID, input)
	if err != nil {
		writeSettingsError(w, err, "Could not refresh media metadata")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Media metadata refreshed", map[string]any{"mediaItemId": item.ID.String()})
	writeJSON(w, http.StatusOK, mediaItemResponse(updated))
}

func mediaRefreshProviderLink(item storage.MediaItem) (string, string) {
	for _, mapping := range item.ProviderMappings {
		if mapping.EntityType == "media_item" && mapping.ProviderName == "anilist" && strings.TrimSpace(mapping.ExternalID) != "" {
			return "anilist", strings.TrimSpace(mapping.ExternalID)
		}
	}
	return strings.TrimSpace(*item.ExternalProvider), strings.TrimSpace(*item.ExternalID)
}

func providerMappingInputs(values []storage.MediaProviderMapping) []storage.MediaProviderMappingInput {
	inputs := make([]storage.MediaProviderMappingInput, 0, len(values))
	for _, value := range values {
		inputs = append(inputs, storage.MediaProviderMappingInput{
			SeasonID: value.SeasonID, EpisodeID: value.EpisodeID, EntityType: value.EntityType,
			ProviderName: value.ProviderName, ProviderEntityType: value.ProviderEntityType,
			ExternalID: value.ExternalID, Canonical: value.Canonical, Confidence: value.Confidence, Source: value.Source,
		})
	}
	return inputs
}

func mediaAliasInputs(values []storage.MediaItemAlias) []storage.MediaItemAliasInput {
	inputs := make([]storage.MediaItemAliasInput, 0, len(values))
	for _, value := range values {
		inputs = append(inputs, storage.MediaItemAliasInput{Alias: value.Alias, Language: value.Language,
			Kind: value.Kind, ProviderName: value.ProviderName, ProviderMappingID: value.ProviderMappingID, Source: value.Source})
	}
	return inputs
}

func refreshedMediaInput(item storage.MediaItem, resource metadata.MediaResource) storage.MediaItemInput {
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
		ProviderMappings:      providerMappingInputs(item.ProviderMappings),
		Aliases:               mediaAliasInputs(item.Aliases),
		Tags:                  item.Tags,
	}
	return applyMediaResource(input, resource)
}
