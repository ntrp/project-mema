//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsSpecial = []subtitleProviderDefinition{
	{Key: "animetosho", DisplayName: "Animetosho", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Anime release index providing embedded subtitles and related files.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{Archive: true, AniDB: true}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{"feed.animetosho.org"}, AllowedDownloadHosts: []string{"animetosho.org"}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
	{Key: "bsplayer", DisplayName: "BSPlayer", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Community subtitle database used by BSPlayer for automatic subtitle matching.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{Archive: true}, Warning: "The upstream SOAP API is historically unreliable and may be unavailable intermittently.", OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{"*.api.bsplayer-subtitles.com"}, AllowedDownloadHosts: []string{"*.bsplayer-subtitles.com", "bsplayer-subtitles.com"}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
	{Key: "napiprojekt", DisplayName: "Napiprojekt", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Polish subtitle database with automatic matching for movies and television series.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{"napiprojekt.pl", "www.napiprojekt.pl"}, AllowedDownloadHosts: []string{"napiprojekt.pl", "www.napiprojekt.pl"}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
}
