package httpapi

import (
	"errors"
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func downloadClientInput(w http.ResponseWriter, request DownloadClientRequest) (storage.DownloadClientInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.DownloadClientInput{}, false
	}
	if !request.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Download client type is not supported")
		return storage.DownloadClientInput{}, false
	}
	if baseURL == "" {
		writeError(w, http.StatusBadRequest, "invalid_base_url", "Base URL is required")
		return storage.DownloadClientInput{}, false
	}
	if request.Priority < 0 || request.Priority > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_priority", "Priority must be between 0 and 1000")
		return storage.DownloadClientInput{}, false
	}

	return storage.DownloadClientInput{
		Name:     name,
		Type:     string(request.Type),
		BaseURL:  baseURL,
		Username: optionalTrimmedString(request.Username),
		Password: optionalTrimmedString(request.Password),
		APIKey:   optionalTrimmedString(request.ApiKey),
		Category: optionalTrimmedString(request.Category),
		Enabled:  request.Enabled,
		Priority: request.Priority,
	}, true
}

func indexerInput(w http.ResponseWriter, request IndexerRequest) (storage.IndexerInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.IndexerInput{}, false
	}
	if !request.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Indexer type is not supported")
		return storage.IndexerInput{}, false
	}
	if baseURL == "" {
		writeError(w, http.StatusBadRequest, "invalid_base_url", "Base URL is required")
		return storage.IndexerInput{}, false
	}
	if request.Priority < 0 || request.Priority > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_priority", "Priority must be between 0 and 1000")
		return storage.IndexerInput{}, false
	}

	categories := []int32{}
	if request.Categories != nil {
		categories = append(categories, (*request.Categories)...)
	}

	return storage.IndexerInput{
		Name:       name,
		Type:       string(request.Type),
		BaseURL:    baseURL,
		APIKey:     optionalTrimmedString(request.ApiKey),
		Categories: categories,
		Enabled:    request.Enabled,
		Priority:   request.Priority,
	}, true
}

func metadataProviderInput(w http.ResponseWriter, request MetadataProviderRequest) (storage.MetadataProviderInput, bool) {
	name := strings.TrimSpace(request.Name)
	baseURL := strings.TrimSpace(request.BaseUrl)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.MetadataProviderInput{}, false
	}
	if !request.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Metadata provider type is not supported")
		return storage.MetadataProviderInput{}, false
	}
	if baseURL == "" {
		writeError(w, http.StatusBadRequest, "invalid_base_url", "Base URL is required")
		return storage.MetadataProviderInput{}, false
	}
	if request.Priority < 0 || request.Priority > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_priority", "Priority must be between 0 and 1000")
		return storage.MetadataProviderInput{}, false
	}
	apiKey := optionalTrimmedString(request.ApiKey)
	accessToken := optionalTrimmedString(request.AccessToken)
	return storage.MetadataProviderInput{
		Name:        name,
		Type:        string(request.Type),
		BaseURL:     baseURL,
		APIKey:      apiKey,
		PIN:         optionalTrimmedString(request.Pin),
		AccessToken: accessToken,
		Enabled:     request.Enabled,
		Priority:    request.Priority,
	}, true
}

func userCreateInput(w http.ResponseWriter, request UserCreateRequest) (storage.UserInput, bool) {
	password := strings.TrimSpace(request.Password)
	if password == "" || len(password) < 8 {
		writeError(w, http.StatusBadRequest, "invalid_password", "Password must be at least 8 characters")
		return storage.UserInput{}, false
	}
	input, ok := userInput(w, request.Username, request.Role, &password)
	if !ok {
		return storage.UserInput{}, false
	}
	return input, true
}

func userUpdateInput(w http.ResponseWriter, request UserUpdateRequest) (storage.UserInput, bool) {
	password := optionalTrimmedString(request.Password)
	if password != nil && len(*password) < 8 {
		writeError(w, http.StatusBadRequest, "invalid_password", "Password must be at least 8 characters")
		return storage.UserInput{}, false
	}
	return userInput(w, request.Username, request.Role, password)
}

func userInput(w http.ResponseWriter, username string, role UserRole, password *string) (storage.UserInput, bool) {
	username = strings.TrimSpace(username)
	if username == "" {
		writeError(w, http.StatusBadRequest, "invalid_username", "Username is required")
		return storage.UserInput{}, false
	}
	if !role.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_role", "User role is not supported")
		return storage.UserInput{}, false
	}
	return storage.UserInput{
		Username: username,
		Password: password,
		Role:     string(role),
	}, true
}

func tagInput(w http.ResponseWriter, request TagRequest) (string, bool) {
	name := strings.Join(strings.Fields(request.Name), " ")
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return "", false
	}
	return name, true
}

func qualitySizeSettingsInput(
	w http.ResponseWriter,
	request QualitySizeSettingsUpdateRequest,
) ([]storage.QualitySizeSettingInput, bool) {
	inputs := make([]storage.QualitySizeSettingInput, 0, len(request.Qualities))
	for _, quality := range request.Qualities {
		qualityID := strings.TrimSpace(quality.QualityId)
		if qualityID == "" {
			writeError(w, http.StatusBadRequest, "invalid_quality", "Quality is required")
			return nil, false
		}
		if quality.MinimumSizeMbPerMinute < 0 {
			writeError(w, http.StatusBadRequest, "invalid_size", "Minimum size must be zero or greater")
			return nil, false
		}
		if quality.PreferredSizeMbPerMinute != nil && *quality.PreferredSizeMbPerMinute < quality.MinimumSizeMbPerMinute {
			writeError(w, http.StatusBadRequest, "invalid_size", "Preferred size must be greater than or equal to minimum size")
			return nil, false
		}
		if quality.MaximumSizeMbPerMinute != nil && *quality.MaximumSizeMbPerMinute < quality.MinimumSizeMbPerMinute {
			writeError(w, http.StatusBadRequest, "invalid_size", "Maximum size must be greater than or equal to minimum size")
			return nil, false
		}
		if quality.PreferredSizeMbPerMinute != nil && quality.MaximumSizeMbPerMinute != nil &&
			*quality.PreferredSizeMbPerMinute > *quality.MaximumSizeMbPerMinute {
			writeError(w, http.StatusBadRequest, "invalid_size", "Preferred size must be less than or equal to maximum size")
			return nil, false
		}

		inputs = append(inputs, storage.QualitySizeSettingInput{
			QualityID:                qualityID,
			MinimumSizeMBPerMinute:   quality.MinimumSizeMbPerMinute,
			PreferredSizeMBPerMinute: quality.PreferredSizeMbPerMinute,
			MaximumSizeMBPerMinute:   quality.MaximumSizeMbPerMinute,
		})
	}
	return inputs, true
}

func mediaProfileInput(w http.ResponseWriter, request MediaProfileRequest) (storage.MediaProfileInput, bool) {
	name := strings.Join(strings.Fields(request.Name), " ")
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_name", "Name is required")
		return storage.MediaProfileInput{}, false
	}

	qualityIDs := make([]string, 0, len(request.QualityIds))
	seen := map[string]struct{}{}
	for _, value := range request.QualityIds {
		qualityID := strings.TrimSpace(value)
		if qualityID == "" {
			continue
		}
		if _, ok := seen[qualityID]; ok {
			continue
		}
		seen[qualityID] = struct{}{}
		qualityIDs = append(qualityIDs, qualityID)
	}
	if len(qualityIDs) == 0 {
		writeError(w, http.StatusBadRequest, "quality_required", "Select at least one quality")
		return storage.MediaProfileInput{}, false
	}

	return storage.MediaProfileInput{
		Name:       name,
		QualityIDs: qualityIDs,
	}, true
}

func fileNamingSettingsInput(
	w http.ResponseWriter,
	request FileNamingSettingsRequest,
) (storage.FileNamingSettingsInput, bool) {
	input := storage.FileNamingSettingsInput{
		MovieFileFormat:      strings.TrimSpace(request.MovieFileFormat),
		MovieFolderFormat:    strings.TrimSpace(request.MovieFolderFormat),
		SeriesEpisodeFormat:  strings.TrimSpace(request.SeriesEpisodeFormat),
		DailyEpisodeFormat:   strings.TrimSpace(request.DailyEpisodeFormat),
		AnimeEpisodeFormat:   strings.TrimSpace(request.AnimeEpisodeFormat),
		SeriesFolderFormat:   strings.TrimSpace(request.SeriesFolderFormat),
		SeasonFolderFormat:   strings.TrimSpace(request.SeasonFolderFormat),
		SpecialsFolderFormat: strings.TrimSpace(request.SpecialsFolderFormat),
	}
	if input.MovieFileFormat == "" ||
		input.MovieFolderFormat == "" ||
		input.SeriesEpisodeFormat == "" ||
		input.DailyEpisodeFormat == "" ||
		input.AnimeEpisodeFormat == "" ||
		input.SeriesFolderFormat == "" ||
		input.SeasonFolderFormat == "" ||
		input.SpecialsFolderFormat == "" {
		writeError(w, http.StatusBadRequest, "invalid_template", "All file naming templates are required")
		return storage.FileNamingSettingsInput{}, false
	}
	return input, true
}

func fileNamingSettingsResponse(settings storage.FileNamingSettings) FileNamingSettings {
	return FileNamingSettings{
		MovieFileFormat:      settings.MovieFileFormat,
		MovieFolderFormat:    settings.MovieFolderFormat,
		SeriesEpisodeFormat:  settings.SeriesEpisodeFormat,
		DailyEpisodeFormat:   settings.DailyEpisodeFormat,
		AnimeEpisodeFormat:   settings.AnimeEpisodeFormat,
		SeriesFolderFormat:   settings.SeriesFolderFormat,
		SeasonFolderFormat:   settings.SeasonFolderFormat,
		SpecialsFolderFormat: settings.SpecialsFolderFormat,
		CreatedAt:            settings.CreatedAt,
		UpdatedAt:            settings.UpdatedAt,
	}
}

func downloadClientConfig(client storage.DownloadClient) downloadclients.Config {
	return downloadclients.Config{
		Name:     client.Name,
		Type:     client.Type,
		BaseURL:  client.BaseURL,
		Username: client.Username,
		Password: client.Password,
		APIKey:   client.APIKey,
		Category: client.Category,
	}
}

func downloadClientInputConfig(input storage.DownloadClientInput) downloadclients.Config {
	return downloadclients.Config{
		Name:     input.Name,
		Type:     input.Type,
		BaseURL:  input.BaseURL,
		Username: input.Username,
		Password: input.Password,
		APIKey:   input.APIKey,
		Category: input.Category,
	}
}

func indexerConfig(indexer storage.Indexer) indexers.Config {
	return indexers.Config{
		ID:         indexer.ID.String(),
		Name:       indexer.Name,
		Type:       indexer.Type,
		BaseURL:    indexer.BaseURL,
		APIKey:     indexer.APIKey,
		Categories: append([]int32(nil), indexer.Categories...),
	}
}

func metadataProviderConfig(provider storage.MetadataProvider) metadata.Config {
	return metadata.Config{
		ID:                    provider.ID,
		Name:                  provider.Name,
		Type:                  provider.Type,
		BaseURL:               provider.BaseURL,
		APIKey:                provider.APIKey,
		PIN:                   provider.PIN,
		AccessToken:           provider.AccessToken,
		SessionToken:          provider.SessionToken,
		SessionTokenExpiresAt: provider.SessionTokenExpiresAt,
	}
}

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
		Id:         openapi_types.UUID(indexer.ID),
		Name:       indexer.Name,
		Type:       IndexerType(indexer.Type),
		BaseUrl:    indexer.BaseURL,
		ApiKey:     indexer.APIKey,
		Categories: &categories,
		Enabled:    indexer.Enabled,
		Priority:   indexer.Priority,
		CreatedAt:  indexer.CreatedAt,
		UpdatedAt:  indexer.UpdatedAt,
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

func mediaProfileListResponse(profiles []storage.MediaProfile) MediaProfileListResponse {
	response := MediaProfileListResponse{Profiles: make([]MediaProfile, 0, len(profiles))}
	for _, profile := range profiles {
		response.Profiles = append(response.Profiles, mediaProfileResponse(profile))
	}
	return response
}

func mediaProfileResponse(profile storage.MediaProfile) MediaProfile {
	return MediaProfile{
		Id:         profile.ID,
		Name:       profile.Name,
		QualityIds: profile.QualityIDs,
		CreatedAt:  profile.CreatedAt,
		UpdatedAt:  profile.UpdatedAt,
	}
}

func optionalTrimmedString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func writeSettingsError(w http.ResponseWriter, err error, message string) {
	if errors.Is(err, storage.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", message)
		return
	}
	if errors.Is(err, storage.ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, "invalid_input", message)
		return
	}
	writeError(w, http.StatusInternalServerError, "settings_update_failed", message)
}
