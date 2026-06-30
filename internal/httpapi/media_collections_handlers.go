package httpapi

import (
	"net/http"
	"strings"

	"media-manager/internal/metadata"
)

func (s *Server) GetMediaCollection(
	w http.ResponseWriter,
	r *http.Request,
	providerType MetadataProviderType,
	collectionID string,
) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	if providerType != Tmdb {
		writeError(w, http.StatusNotFound, "metadata_provider_not_found", "Could not find metadata provider")
		return
	}
	provider, ok, err := s.tmdbProvider(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "metadata_provider_not_found", "Could not find metadata provider")
		return
	}

	collection, err := s.metadata.Collection(r.Context(), metadataProviderConfig(provider), collectionID)
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, mediaCollectionResponse(providerType, collection))
}

func mediaCollectionResponse(providerType MetadataProviderType, collection metadata.Collection) MediaCollection {
	results := make([]MediaSearchResult, 0, len(collection.Parts))
	for _, part := range collection.Parts {
		results = append(results, metadataSearchResultResponse(part))
	}
	name := strings.TrimSpace(collection.Name)
	if name == "" {
		name = "Collection"
	}
	return MediaCollection{
		Id:           collection.ID,
		Name:         name,
		Provider:     providerType,
		Overview:     collection.Overview,
		PosterPath:   collection.PosterPath,
		BackdropPath: collection.BackdropPath,
		Results:      results,
	}
}
