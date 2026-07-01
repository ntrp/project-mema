package httpapi

import (
	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

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
