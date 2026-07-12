//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsLocal = []subtitleProviderDefinition{
	{Key: "opensubtitlescom", DisplayName: "OpenSubtitles.com", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Runtime supported through the OpenSubtitles.com HTTP API.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.opensubtitles.com"}, AllowedDownloadHosts: []string{"www.opensubtitles.com", "dl.opensubtitles.com", "api.opensubtitles.com"}, AllowLocalHosts: false}, Fields: []catalog.Field{baseURLField("https://api.opensubtitles.com"), usernameField(), passwordField(), apiKeyField()}},
	{Key: "whisperai", DisplayName: "Whisper", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Runtime supported through a local Whisper-compatible transcription endpoint.", MediaTypes: []string{"movie", "serie", "anime", "audio"}, Dependencies: catalog.Dependencies{LocalHTTPEndpoint: true, FFmpeg: true, FFprobe: true}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{}, AllowedDownloadHosts: []string{}, AllowLocalHosts: true}, Fields: []catalog.Field{baseURLField(""), numericTextField("responseTimeoutSeconds", "Response timeout seconds"), numericTextField("transcriptionTimeoutSeconds", "Transcription timeout seconds"), numericTextField("logLevel", "Log level"), boolField("passVideoName", "Pass video name")}},
	{Key: "mock", DisplayName: "Mock", RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Fixture-backed local mock provider for tests and development.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{}, AllowedDownloadHosts: []string{}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
}
