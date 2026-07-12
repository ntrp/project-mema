//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsLocal = []subtitleProviderDefinition{
	{Key: "opensubtitlescom", DisplayName: "OpenSubtitles.com", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Runtime supported through the OpenSubtitles.com HTTP API.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.opensubtitles.com"}, AllowedDownloadHosts: []string{"www.opensubtitles.com", "dl.opensubtitles.com", "api.opensubtitles.com"}, AllowLocalHosts: false}, Fields: []catalog.Field{baseURLField("https://api.opensubtitles.com"), usernameField(), passwordField(), apiKeyField()}},
	{Key: "whisperai", DisplayName: "Whisper", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeUnsupported, RuntimeMessage: "Requires a reviewed local Whisper service integration before runtime use.", MediaTypes: []string{"movie", "serie", "anime", "audio"}, Dependencies: catalog.Dependencies{LocalHTTPEndpoint: true}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{}, AllowedDownloadHosts: []string{}, AllowLocalHosts: true}, Fields: []catalog.Field{baseURLField("")}},
	{Key: "mock", DisplayName: "Mock", RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Fixture-backed local mock provider for tests and development.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{}, AllowedDownloadHosts: []string{}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
}
