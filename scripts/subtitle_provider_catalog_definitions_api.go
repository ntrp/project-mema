//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsAPI = []subtitleProviderDefinition{
	apiProvider("assrt", "Assrt", []string{"movie", "serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.assrt.net"}, AllowedDownloadHosts: []string{"api.assrt.net"}}, []catalog.Field{tokenField()}),
	apiProvider("betaseries", "Betaseries", []string{"serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.betaseries.com"}, AllowedDownloadHosts: []string{"api.betaseries.com"}}, []catalog.Field{tokenField()}),
	apiProvider("gestdown", "Gestdown", []string{"serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.gestdown.info"}, AllowedDownloadHosts: []string{"api.gestdown.info"}}, []catalog.Field{}),
	apiProvider("jimaku", "Jimaku", []string{"movie", "serie"}, catalog.Dependencies{Archive: true}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"jimaku.cc"}, AllowedDownloadHosts: []string{"jimaku.cc"}}, []catalog.Field{apiKeyField(), boolField("enableNameSearchFallback", "Enable name-search fallback"), boolField("enableArchivesDownload", "Enable archive downloads"), boolField("enableAiSubs", "Enable AI subtitles")}),
	apiProvider("regielive", "Regielive", []string{"movie", "serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.regielive.ro", "subtitrari.regielive.ro"}, AllowedDownloadHosts: []string{"api.regielive.ro", "subtitrari.regielive.ro"}}, []catalog.Field{}),
	apiProvider("subdl", "Subdl", []string{"movie", "serie"}, catalog.Dependencies{Archive: true}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.subdl.com", "subdl.com"}, AllowedDownloadHosts: []string{"dl.subdl.com", "subdl.com", "api.subdl.com"}}, []catalog.Field{apiKeyField()}),
	apiProvider("subsource", "Subsource", []string{"movie", "serie"}, catalog.Dependencies{Archive: true}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.subsource.net", "subsource.net"}, AllowedDownloadHosts: []string{"api.subsource.net", "subsource.net"}}, []catalog.Field{apiKeyField()}),
	apiProvider("subsarr", "Subsarr", []string{"movie", "serie"}, catalog.Dependencies{LocalHTTPEndpoint: true}, catalog.OutboundPolicy{AllowLocalHosts: true}, []catalog.Field{baseURLField("")}),
	apiProvider("subsro", "Subs.ro", []string{"movie", "serie"}, catalog.Dependencies{Archive: true}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.subs.ro"}, AllowedDownloadHosts: []string{"api.subs.ro"}}, []catalog.Field{apiKeyField()}),
	apiProvider("subtis", "Subt.is", []string{"movie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"api.subt.is"}, AllowedDownloadHosts: []string{"api.subt.is"}}, []catalog.Field{}),
	apiProvider("subx", "SubX", []string{"movie", "serie"}, catalog.Dependencies{}, catalog.OutboundPolicy{AllowedBaseHosts: []string{"subx-api.duckdns.org"}, AllowedDownloadHosts: []string{"subx-api.duckdns.org"}}, []catalog.Field{apiKeyField()}),
}

func apiProvider(key, displayName string, mediaTypes []string, deps catalog.Dependencies, policy catalog.OutboundPolicy, fields []catalog.Field) subtitleProviderDefinition {
	return subtitleProviderDefinition{Key: key, DisplayName: displayName, ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Runtime supported through the native Go structured HTTP provider adapter.", MediaTypes: mediaTypes, Dependencies: deps, OutboundPolicy: policy, Fields: fields}
}
