package httpapi

import (
	"media-manager/internal/dlna"
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

func fileDeleteSettingsResponse(settings storage.FileDeleteSettings) FileDeleteSettings {
	return FileDeleteSettings{
		Mode:          FileDeleteMode(settings.Mode),
		RecycleFolder: settings.RecycleFolder,
		CreatedAt:     settings.CreatedAt,
		UpdatedAt:     settings.UpdatedAt,
	}
}

func (s *Server) dlnaSettingsResponse(settings storage.DLNASettings) DLNASettings {
	return DLNASettings{
		Enabled:                 settings.Enabled,
		FriendlyName:            settings.FriendlyName,
		Interfaces:              append([]string{}, settings.Interfaces...),
		AllowedCidrs:            append([]string{}, settings.AllowedCIDRs...),
		AnnounceIntervalSeconds: settings.AnnounceIntervalSeconds,
		TranscodeEnabled:        settings.TranscodeEnabled,
		ThumbnailsEnabled:       settings.ThumbnailsEnabled,
		SubtitlesEnabled:        settings.SubtitlesEnabled,
		DefaultRendererProfile:  settings.DefaultRendererProfile,
		CreatedAt:               settings.CreatedAt,
		UpdatedAt:               settings.UpdatedAt,
		Status:                  dlnaStatusResponse(s.currentDLNAStatus()),
	}
}

func (s *Server) currentDLNAStatus() dlna.Status {
	if s.dlna == nil {
		return dlna.Status{}
	}
	return s.dlna.Status()
}

func dlnaStatusResponse(status dlna.Status) DLNAStatus {
	return DLNAStatus{
		Running:          status.Running,
		BoundInterfaces:  append([]string{}, status.BoundInterfaces...),
		AdvertisedUrls:   append([]string{}, status.AdvertisedURLs...),
		LastError:        status.LastError,
		LastSsdpEvent:    status.LastSSDPEvent,
		LastSoapAction:   status.LastSOAPAction,
		RecentClients:    dlnaClientDiagnostics(status.RecentClients),
		ActiveStreams:    dlnaStreamDiagnostics(status.ActiveStreams),
		ActiveTranscodes: dlnaStreamDiagnostics(status.ActiveTranscodes),
	}
}

func dlnaClientDiagnostics(clients []dlna.ClientStatus) []DLNAClientDiagnostic {
	values := make([]DLNAClientDiagnostic, 0, len(clients))
	for _, client := range clients {
		values = append(values, DLNAClientDiagnostic{
			Ip:             client.IP,
			UserAgent:      client.UserAgent,
			ProfileId:      client.ProfileID,
			LastSoapAction: client.LastSOAPAction,
			LastError:      client.LastError,
			LastSeen:       client.LastSeen,
		})
	}
	return values
}

func dlnaStreamDiagnostics(streams []dlna.StreamStatus) []DLNAStreamDiagnostic {
	values := make([]DLNAStreamDiagnostic, 0, len(streams))
	for _, stream := range streams {
		values = append(values, DLNAStreamDiagnostic{
			Id:        stream.ID,
			ClientIp:  stream.ClientIP,
			Path:      stream.Path,
			ProfileId: stream.ProfileID,
			StartedAt: stream.StartedAt,
		})
	}
	return values
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
		ID:             indexer.ID.String(),
		DefinitionID:   indexer.DefinitionID,
		Name:           indexer.Name,
		Implementation: indexer.Implementation,
		Protocol:       indexer.Protocol,
		BaseURL:        indexer.BaseURL,
		APIKey:         indexer.APIKey,
		Categories:     append([]int32(nil), indexer.Categories...),
		Fields:         append([]byte(nil), indexer.Fields...),
		Redirect:       indexer.Redirect,
	}
}

func indexerInputConfig(input storage.IndexerInput) indexers.Config {
	return indexers.Config{
		DefinitionID:   input.DefinitionID,
		Name:           input.Name,
		Implementation: input.Implementation,
		Protocol:       input.Protocol,
		BaseURL:        input.BaseURL,
		APIKey:         input.APIKey,
		Categories:     append([]int32(nil), input.Categories...),
		Fields:         append([]byte(nil), input.Fields...),
		Redirect:       input.Redirect,
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
