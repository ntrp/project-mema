package httpapi

import (
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func downloadClientResponse(client storage.DownloadClient) DownloadClient {
	return DownloadClient{
		Id:        openapi_types.UUID(client.ID),
		Name:      client.Name,
		Type:      DownloadClientType(client.Type),
		BaseUrl:   client.BaseURL,
		Username:  client.Username,
		Password:  client.Password,
		ApiKey:    client.APIKey,
		Category:  client.Category,
		Enabled:   client.Enabled,
		Priority:  client.Priority,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}

func indexerResponse(indexer storage.Indexer) Indexer {
	categories := append([]int32(nil), indexer.Categories...)
	if categories == nil {
		categories = []int32{}
	}
	return Indexer{
		Id:             openapi_types.UUID(indexer.ID),
		Name:           indexer.Name,
		Type:           IndexerType(indexer.Type),
		BaseUrl:        indexer.BaseURL,
		ApiKey:         indexer.APIKey,
		Categories:     &categories,
		Enabled:        indexer.Enabled,
		Priority:       indexer.Priority,
		HealthStatus:   IndexerHealthStatus(indexer.HealthStatus),
		LastQueryAt:    indexer.LastQueryAt,
		LastSuccessAt:  indexer.LastSuccessAt,
		LastFailureAt:  indexer.LastFailureAt,
		NextCheckAt:    indexer.NextCheckAt,
		LastStatusCode: indexer.LastStatusCode,
		LastError:      indexer.LastError,
		FailureCount:   indexer.FailureCount,
		CreatedAt:      indexer.CreatedAt,
		UpdatedAt:      indexer.UpdatedAt,
	}
}

func metadataProviderResponse(provider storage.MetadataProvider) MetadataProvider {
	return MetadataProvider{
		Id:          openapi_types.UUID(provider.ID),
		Name:        provider.Name,
		Type:        MetadataProviderType(provider.Type),
		BaseUrl:     provider.BaseURL,
		ApiKey:      provider.APIKey,
		Pin:         provider.PIN,
		AccessToken: provider.AccessToken,
		Enabled:     provider.Enabled,
		Priority:    provider.Priority,
		CreatedAt:   provider.CreatedAt,
		UpdatedAt:   provider.UpdatedAt,
	}
}

func metadataCacheStatsResponse(stats storage.MetadataCacheStats) MetadataCacheStats {
	return MetadataCacheStats{
		TotalEntries:   stats.TotalEntries,
		ActiveEntries:  stats.ActiveEntries,
		ExpiredEntries: stats.ExpiredEntries,
		ProviderCount:  stats.ProviderCount,
	}
}

func metadataCacheEntryResponse(entry storage.MetadataCacheEntry) MetadataCacheEntry {
	return MetadataCacheEntry{
		ProviderName: entry.ProviderName,
		ProviderType: MetadataProviderType(entry.ProviderType),
		MediaType:    MediaType(entry.MediaType),
		Query:        entry.Query,
		CacheKind:    MetadataCacheEntryCacheKind(cacheKind(entry.Query)),
		Year:         entry.Year,
		ItemCount:    entry.ItemCount,
		ExpiresAt:    entry.ExpiresAt,
		CreatedAt:    entry.CreatedAt,
		UpdatedAt:    entry.UpdatedAt,
		Expired:      entry.Expired,
	}
}

func cacheKind(query string) string {
	switch {
	case strings.HasPrefix(query, "discover:"):
		return "discover"
	case strings.HasPrefix(query, "details:"):
		return "details"
	default:
		return "search"
	}
}

func managedUserResponse(user storage.User) ManagedUser {
	return ManagedUser{
		Id:        openapi_types.UUID(user.ID),
		Username:  user.Username,
		Role:      UserRole(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func tagResponse(tag storage.Tag) Tag {
	return Tag{
		Id:        openapi_types.UUID(tag.ID),
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}

func qualitySizeSettingsResponse(settings []storage.QualitySizeSetting) QualitySizeSettingsResponse {
	response := QualitySizeSettingsResponse{Qualities: make([]QualitySizeSetting, 0, len(settings))}
	for _, setting := range settings {
		response.Qualities = append(response.Qualities, qualitySizeSettingResponse(setting))
	}
	return response
}

func qualitySizeSettingResponse(setting storage.QualitySizeSetting) QualitySizeSetting {
	return QualitySizeSetting{
		QualityId:                setting.ID,
		Name:                     setting.Name,
		SortOrder:                setting.SortOrder,
		MinimumSizeMbPerMinute:   setting.MinimumSizeMBPerMinute,
		PreferredSizeMbPerMinute: setting.PreferredSizeMBPerMinute,
		MaximumSizeMbPerMinute:   setting.MaximumSizeMBPerMinute,
		CreatedAt:                setting.CreatedAt,
		UpdatedAt:                setting.UpdatedAt,
	}
}
