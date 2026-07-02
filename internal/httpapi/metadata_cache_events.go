package httpapi

import (
	"context"
	"reflect"
	"time"

	"media-manager/internal/storage"
)

type metadataCacheUpdatedEvent struct {
	Entry MetadataCacheEntry `json:"entry"`
	Stats MetadataCacheStats `json:"stats"`
}

func (s *Server) publishMetadataCacheUpdated(
	ctx context.Context,
	provider storage.MetadataProvider,
	mediaType string,
	query string,
	year *int32,
	value any,
	expiresAt time.Time,
) {
	stats, err := s.settings.MetadataCacheStats(ctx)
	if err != nil {
		return
	}
	entry := storage.MetadataCacheEntry{
		ProviderID:   provider.ID,
		ProviderName: provider.Name,
		ProviderType: provider.Type,
		MediaType:    mediaType,
		Query:        query,
		Year:         cacheEntryYear(year),
		ItemCount:    cacheItemCount(value),
		ExpiresAt:    expiresAt,
		CreatedAt:    s.now(),
		UpdatedAt:    s.now(),
		Expired:      false,
	}
	s.events.Publish("metadata.cache.updated", metadataCacheUpdatedEvent{
		Entry: metadataCacheEntryResponse(entry),
		Stats: metadataCacheStatsResponse(stats),
	})
}

func cacheEntryYear(year *int32) int32 {
	if year == nil {
		return 0
	}
	return *year
}

func cacheItemCount(value any) int32 {
	if value == nil {
		return 0
	}
	reflected := reflect.ValueOf(value)
	if reflected.Kind() == reflect.Slice || reflected.Kind() == reflect.Array {
		return int32(reflected.Len())
	}
	return 1
}
