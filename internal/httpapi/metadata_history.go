package httpapi

import (
	"context"
	"log/slog"

	"media-manager/internal/storage"
)

func (s *Server) recordMetadataSearchHistory(
	ctx context.Context,
	provider storage.MetadataProvider,
	mediaType string,
	query string,
	year *int32,
	cacheHit bool,
	response any,
	searchErr error,
) {
	errorMessage := optionalErrorMessage(searchErr)
	entry, err := s.settings.RecordMetadataSearchHistory(ctx, storage.MetadataSearchHistoryInput{
		ProviderID:   provider.ID,
		ProviderName: provider.Name,
		ProviderType: provider.Type,
		MediaType:    mediaType,
		Query:        query,
		Year:         year,
		CacheHit:     cacheHit,
		Success:      searchErr == nil,
		ItemCount:    cacheItemCount(response),
		Error:        errorMessage,
		Response:     metadataHistoryResponse(response, errorMessage),
	})
	if err != nil {
		slog.Error("record metadata search history", "error", err)
		return
	}
	s.events.Publish("metadata.search.history.created", metadataSearchHistoryEntryResponse(entry))
}

func optionalErrorMessage(err error) *string {
	if err == nil {
		return nil
	}
	message := err.Error()
	return &message
}

func metadataHistoryResponse(response any, errorMessage *string) any {
	if errorMessage == nil {
		return response
	}
	return map[string]string{"error": *errorMessage}
}
