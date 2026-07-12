//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsSpecial = []subtitleProviderDefinition{
	{Key: "animetosho", DisplayName: "Animetosho", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeCatalogOnly, RuntimeMessage: "Catalog entry only; runtime support is not implemented yet.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{}, AllowedDownloadHosts: []string{}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
	{Key: "bsplayer", DisplayName: "BSPlayer", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeCatalogOnly, RuntimeMessage: "Catalog entry only; runtime support is not implemented yet.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{Archive: true}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{}, AllowedDownloadHosts: []string{}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
	{Key: "napiprojekt", DisplayName: "Napiprojekt", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeCatalogOnly, RuntimeMessage: "Catalog entry only; runtime support is not implemented yet.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{}, AllowedDownloadHosts: []string{}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
}
