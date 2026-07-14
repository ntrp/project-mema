//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsAPI = []subtitleProviderDefinition{
	apiProvider("assrt", "Assrt", "Chinese subtitle catalog for movies and television series.", []string{"movie", "serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.assrt.net"}, AllowedDownloadHosts: []string{"api.assrt.net"}}, []catalog.Field{tokenField()}),
	apiProvider("betaseries", "Betaseries", "Community subtitles for television series tracked on BetaSeries.", []string{"serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.betaseries.com"}, AllowedDownloadHosts: []string{"api.betaseries.com"}}, []catalog.Field{tokenField()}),
	apiProvider("gestdown", "Gestdown", "Subtitle catalog for television series and episode releases.", []string{"serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.gestdown.info"}, AllowedDownloadHosts: []string{"api.gestdown.info"}}, []catalog.Field{}),
	apiProvider("jimaku", "Jimaku", "Japanese subtitles for anime, films, and television series.", []string{"movie", "serie"}, catalog.Dependencies{Archive: true}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"jimaku.cc"}, AllowedDownloadHosts: []string{"jimaku.cc"}}, []catalog.Field{apiKeyField(), boolField("enableNameSearchFallback", "Enable name-search fallback"), boolField("enableArchivesDownload", "Enable archive downloads"), boolField("enableAiSubs", "Enable AI subtitles")}),
	apiProvider("regielive", "Regielive", "Romanian subtitles for movies and television series.", []string{"movie", "serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.regielive.ro", "subtitrari.regielive.ro"}, AllowedDownloadHosts: []string{"api.regielive.ro", "subtitrari.regielive.ro"}}, []catalog.Field{}),
	apiProvider("subdl", "Subdl", "Multilingual subtitle catalog for movies and television series.", []string{"movie", "serie"}, catalog.Dependencies{Archive: true}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.subdl.com", "subdl.com"}, AllowedDownloadHosts: []string{"dl.subdl.com", "subdl.com", "api.subdl.com"}}, []catalog.Field{apiKeyField()}),
	apiProvider("subsource", "Subsource", "Multilingual subtitle archive for movies and television series.", []string{"movie", "serie"}, catalog.Dependencies{Archive: true}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.subsource.net", "subsource.net"}, AllowedDownloadHosts: []string{"api.subsource.net", "subsource.net"}}, []catalog.Field{apiKeyField()}),
	apiProvider("subsarr", "Subsarr", "Self-hosted subtitle search service for movies and television series.", []string{"movie", "serie"}, catalog.Dependencies{LocalHTTPEndpoint: true}, catalog.OutboundPolicy{AllowLocalHosts: true}, []catalog.Field{baseURLField("")}),
	apiProvider("subsro", "Subs.ro", "Romanian subtitle catalog for movies and television series.", []string{"movie", "serie"}, catalog.Dependencies{Archive: true}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.subs.ro"}, AllowedDownloadHosts: []string{"api.subs.ro"}}, []catalog.Field{apiKeyField()}),
	apiProvider("subtis", "Subt.is", "Movie subtitle catalog spanning multiple languages and releases.", []string{"movie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.subt.is"}, AllowedDownloadHosts: []string{"api.subt.is"}}, []catalog.Field{}),
	apiProvider("subx", "SubX", "Subtitle search service for movies and television series.", []string{"movie", "serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"subx-api.duckdns.org"}, AllowedDownloadHosts: []string{"subx-api.duckdns.org"}}, []catalog.Field{apiKeyField()}),
}

func apiProvider(key, displayName string, description string, mediaTypes []string, deps catalog.Dependencies, policy catalog.OutboundPolicy, fields []catalog.Field) subtitleProviderDefinition {
	return subtitleProviderDefinition{Key: key, DisplayName: displayName, ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: description, MediaTypes: mediaTypes, Dependencies: deps, OutboundPolicy: policy, Fields: fields}
}
