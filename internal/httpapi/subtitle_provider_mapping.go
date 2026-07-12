package httpapi

import (
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
	"media-manager/internal/subtitles/catalog"
)

func subtitleProviderCatalogKey(providerType string) string {
	if providerType == "opensubtitles" {
		return "opensubtitlescom"
	}
	return providerType
}

func subtitleProviderRuntime(providerType string) (SubtitleProviderRuntimeStatus, string) {
	entry, ok := catalog.Lookup(subtitleProviderCatalogKey(providerType))
	if !ok {
		return Unsupported, "Subtitle provider is not in the catalog."
	}
	return SubtitleProviderRuntimeStatus(entry.RuntimeStatus), entry.RuntimeMessage
}

func apiSettingValues(settings storage.SubtitleProviderSettings) map[string]SubtitleProviderSettingValue {
	values := make(map[string]SubtitleProviderSettingValue, len(settings))
	for key, value := range settings {
		values[key] = SubtitleProviderSettingValue{
			StringValue:  value.StringValue,
			NumberValue:  value.NumberValue,
			BooleanValue: value.BooleanValue,
			StringValues: stringValuesPtr(value.StringValues),
		}
	}
	return values
}

func storageSettingValues(settings *map[string]SubtitleProviderSettingValue) storage.SubtitleProviderSettings {
	values := storage.SubtitleProviderSettings{}
	if settings == nil {
		return values
	}
	for key, value := range *settings {
		var stringValues []string
		if value.StringValues != nil {
			stringValues = append([]string{}, (*value.StringValues)...)
		}
		values[key] = storage.SubtitleProviderSettingValue{
			StringValue:  value.StringValue,
			NumberValue:  value.NumberValue,
			BooleanValue: value.BooleanValue,
			StringValues: stringValues,
		}
	}
	return values
}

func stringValuesPtr(values []string) *[]string {
	if values == nil {
		return nil
	}
	copyValues := append([]string{}, values...)
	return &copyValues
}

func subtitleProviderConfigFromInput(input storage.SubtitleProviderInput) subtitles.Config {
	return subtitles.Config{
		Name:          input.Name,
		Type:          input.Type,
		BaseURL:       input.BaseURL,
		Username:      input.Username,
		Password:      input.Password,
		APIKey:        input.APIKey,
		MockSubtitles: subtitleProviderMockConfigFromInput(input.MockSubtitles),
	}
}

func subtitleProviderMockConfigFromInput(rows []storage.MockSubtitleProviderRowInput) []subtitles.MockSubtitle {
	items := make([]subtitles.MockSubtitle, 0, len(rows))
	for _, row := range rows {
		items = append(items, subtitles.MockSubtitle{
			Title:      row.Title,
			LanguageID: row.LanguageID,
			Format:     row.Format,
		})
	}
	return items
}
