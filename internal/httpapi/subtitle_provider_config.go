package httpapi

import (
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

func subtitleProviderConfig(provider storage.SubtitleProvider) subtitles.Config {
	return subtitles.Config{
		Name:           provider.Name,
		Type:           provider.Type,
		BaseURL:        provider.BaseURL,
		Username:       provider.Username,
		Password:       provider.Password,
		APIKey:         provider.APIKey,
		Settings:       subtitleProviderConfigSettings(provider.Settings),
		SecretSettings: subtitleProviderConfigSecrets(provider.SecretSettings),
		MockSubtitles:  subtitleProviderMockConfig(provider.MockSubtitles),
	}
}

func subtitleProviderMockConfig(rows []storage.MockSubtitleProviderRow) []subtitles.MockSubtitle {
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
