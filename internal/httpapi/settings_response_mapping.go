package httpapi

import (
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func downloadClientResponse(client storage.DownloadClient) DownloadClient {
	return DownloadClient{
		Id:          openapi_types.UUID(client.ID),
		Name:        client.Name,
		Type:        DownloadClientType(client.Type),
		Protocol:    IndexerProtocol(client.Protocol),
		BaseUrl:     client.BaseURL,
		Username:    client.Username,
		Password:    client.Password,
		ApiKey:      client.APIKey,
		PasswordSet: client.Password != nil,
		ApiKeySet:   client.APIKey != nil,
		Category:    client.Category,
		Enabled:     client.Enabled,
		Priority:    client.Priority,
		CreatedAt:   client.CreatedAt,
		UpdatedAt:   client.UpdatedAt,
	}
}

func indexerResponse(indexer storage.Indexer, languages catalogLanguageMapper) Indexer {
	categories := append([]int32(nil), indexer.Categories...)
	if categories == nil {
		categories = []int32{}
	}
	mediaTypeScopes := indexerMediaTypeScopeResponses(indexer.MediaTypeScopes)
	tagScopes := append([]string(nil), indexer.TagScopes...)
	if tagScopes == nil {
		tagScopes = []string{}
	}
	return Indexer{
		Id:                 openapi_types.UUID(indexer.ID),
		DefinitionId:       indexer.DefinitionID,
		Name:               indexer.Name,
		Implementation:     &indexer.Implementation,
		ImplementationName: &indexer.ImplementationName,
		Protocol:           IndexerProtocol(indexer.Protocol),
		Privacy:            IndexerPrivacy(indexer.Privacy),
		Language:           languages.code(indexer.Language),
		Encoding:           indexer.Encoding,
		Description:        indexer.Description,
		IndexerUrls:        &indexer.IndexerURLs,
		LegacyUrls:         &indexer.LegacyURLs,
		BaseUrl:            indexer.BaseURL,
		ApiKey:             indexer.APIKey,
		ApiKeySet:          indexer.APIKey != nil,
		Categories:         &categories,
		MediaTypeScopes:    &mediaTypeScopes,
		TagScopes:          &tagScopes,
		Fields:             indexerFieldValues(indexer.Fields),
		Capabilities:       indexerCapabilities(indexer.Capabilities),
		Redirect:           &indexer.Redirect,
		AppProfileId:       &indexer.AppProfileID,
		MinimumSeeders:     indexer.MinimumSeeders,
		SeedRatio:          indexer.SeedRatio,
		SeedTime:           indexer.SeedTime,
		PackSeedTime:       indexer.PackSeedTime,
		PreferMagnetUrl:    &indexer.PreferMagnetURL,
		SupportsRss:        indexer.SupportsRSS,
		SupportsSearch:     indexer.SupportsSearch,
		SupportsRedirect:   indexer.SupportsRedirect,
		SupportsPagination: indexer.SupportsPagination,
		Enabled:            indexer.Enabled,
		Priority:           indexer.Priority,
		HealthStatus:       IndexerHealthStatus(indexer.HealthStatus),
		LastQueryAt:        indexer.LastQueryAt,
		LastSuccessAt:      indexer.LastSuccessAt,
		LastFailureAt:      indexer.LastFailureAt,
		NextCheckAt:        indexer.NextCheckAt,
		LastStatusCode:     indexer.LastStatusCode,
		LastError:          indexer.LastError,
		FailureCount:       indexer.FailureCount,
		CreatedAt:          indexer.CreatedAt,
		UpdatedAt:          indexer.UpdatedAt,
	}
}

func indexerMediaTypeScopeResponses(values []string) []IndexerMediaType {
	scopes := make([]IndexerMediaType, 0, len(values))
	for _, value := range values {
		scopes = append(scopes, IndexerMediaType(value))
	}
	return scopes
}

func indexerSearchResponse(
	settings storage.IndexerSearchSettings,
	stats storage.IndexerSearchCacheStats,
	cacheEntries []storage.IndexerSearchCacheEntry,
	historyEntries []storage.IndexerSearchHistoryEntry,
	historyStats storage.QueryHistoryStats,
) IndexerSearchResponse {
	cache := make([]IndexerSearchCacheEntry, 0, len(cacheEntries))
	for _, entry := range cacheEntries {
		cache = append(cache, indexerSearchCacheEntryResponse(entry))
	}
	history := make([]IndexerSearchHistoryEntry, 0, len(historyEntries))
	for _, entry := range historyEntries {
		history = append(history, indexerSearchHistoryEntryResponse(entry))
	}
	return IndexerSearchResponse{
		Settings: IndexerSearchSettings{
			CacheDurationMinutes:         settings.CacheDurationMinutes,
			HistoryRetentionDays:         settings.HistoryRetentionDays,
			AutomaticBlocklistExpiryDays: settings.AutomaticBlocklistExpiryDays,
		},
		Stats: IndexerSearchCacheStats{
			TotalEntries:   stats.TotalEntries,
			ActiveEntries:  stats.ActiveEntries,
			ExpiredEntries: stats.ExpiredEntries,
			IndexerCount:   stats.IndexerCount,
		},
		CacheEntries:        cache,
		HistoryEntries:      history,
		HistoryTotalEntries: historyStats.TotalEntries,
		HistoryStats:        queryHistoryStatsResponse(historyStats),
	}
}

func queryHistoryStatsResponse(stats storage.QueryHistoryStats) QueryHistoryStats {
	return QueryHistoryStats{
		TotalEntries: stats.TotalEntries,
		CacheHits:    stats.CacheHits,
		CacheMisses:  stats.CacheMisses,
		Failures:     stats.Failures,
	}
}

func indexerSearchCacheEntryResponse(entry storage.IndexerSearchCacheEntry) IndexerSearchCacheEntry {
	return IndexerSearchCacheEntry{
		IndexerId:       openapi_types.UUID(entry.IndexerID),
		IndexerName:     entry.IndexerName,
		IndexerProtocol: IndexerProtocol(entry.IndexerProtocol),
		MediaType:       MediaType(entry.MediaType),
		Query:           entry.Query,
		ResultCount:     entry.ResultCount,
		ExpiresAt:       entry.ExpiresAt,
		CreatedAt:       entry.CreatedAt,
		UpdatedAt:       entry.UpdatedAt,
		Expired:         entry.Expired,
	}
}

func indexerSearchHistoryEntryResponse(entry storage.IndexerSearchHistoryEntry) IndexerSearchHistoryEntry {
	return IndexerSearchHistoryEntry{
		IndexerName:     entry.IndexerName,
		IndexerProtocol: IndexerProtocol(entry.IndexerProtocol),
		MediaType:       MediaType(entry.MediaType),
		Query:           entry.Query,
		CacheHit:        entry.CacheHit,
		Success:         entry.Success,
		ResultCount:     entry.ResultCount,
		Error:           entry.Error,
		Response:        entry.Response,
		CreatedAt:       entry.CreatedAt,
	}
}

func metadataProviderResponse(provider storage.MetadataProvider) MetadataProvider {
	return MetadataProvider{
		Id:             openapi_types.UUID(provider.ID),
		Name:           provider.Name,
		Type:           MetadataProviderType(provider.Type),
		BaseUrl:        provider.BaseURL,
		ApiKey:         provider.APIKey,
		Pin:            provider.PIN,
		AccessToken:    provider.AccessToken,
		ApiKeySet:      provider.APIKey != nil,
		PinSet:         provider.PIN != nil,
		AccessTokenSet: provider.AccessToken != nil,
		Enabled:        provider.Enabled,
		Priority:       provider.Priority,
		CreatedAt:      provider.CreatedAt,
		UpdatedAt:      provider.UpdatedAt,
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
		ProviderId:   openapi_types.UUID(entry.ProviderID),
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

func metadataSearchHistoryEntryResponse(entry storage.MetadataSearchHistoryEntry) MetadataSearchHistoryEntry {
	return MetadataSearchHistoryEntry{
		ProviderName: entry.ProviderName,
		ProviderType: MetadataProviderType(entry.ProviderType),
		MediaType:    MediaType(entry.MediaType),
		Query:        entry.Query,
		CacheKind:    MetadataSearchHistoryEntryCacheKind(cacheKind(entry.Query)),
		Year:         entry.Year,
		CacheHit:     entry.CacheHit,
		Success:      entry.Success,
		ItemCount:    entry.ItemCount,
		Error:        entry.Error,
		Response:     entry.Response,
		CreatedAt:    entry.CreatedAt,
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

func languageResponse(language storage.Language) Language {
	return Language{
		Code:        language.Code,
		DisplayName: language.DisplayName,
		Aliases:     append([]string{}, language.Aliases...),
		CreatedAt:   language.CreatedAt,
		UpdatedAt:   language.UpdatedAt,
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
