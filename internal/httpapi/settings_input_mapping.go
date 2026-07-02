package httpapi

import (
	"errors"
	"net/http"
	"strings"

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

func languageInput(w http.ResponseWriter, request LanguageRequest) (storage.LanguageInput, bool) {
	return validateLanguageInput(w, storage.LanguageInput{
		Code:        request.Code,
		DisplayName: request.DisplayName,
		Aliases:     request.Aliases,
	})
}

func languageUpdateInput(w http.ResponseWriter, request LanguageUpdateRequest) (storage.LanguageInput, bool) {
	return validateLanguageInput(w, storage.LanguageInput{
		DisplayName: request.DisplayName,
		Aliases:     request.Aliases,
	})
}

func validateLanguageInput(w http.ResponseWriter, input storage.LanguageInput) (storage.LanguageInput, bool) {
	if strings.TrimSpace(input.DisplayName) == "" {
		writeError(w, http.StatusBadRequest, "invalid_display_name", "Display name is required")
		return storage.LanguageInput{}, false
	}
	return input, true
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
