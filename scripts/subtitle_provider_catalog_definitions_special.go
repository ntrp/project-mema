//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsSpecial = []subtitleProviderDefinition{
	{Key: "animetosho", DisplayName: "Animetosho", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Runtime support implemented for AniDB episode-id based anime subtitle lookup and xz downloads.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{Archive: true, AniDB: true}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{"feed.animetosho.org"}, AllowedDownloadHosts: []string{"animetosho.org"}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
	{Key: "bsplayer", DisplayName: "BSPlayer", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Runtime support reports the Bazarr-documented broken SOAP search upstream honestly; gzip downloads are supported for existing candidates.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{Archive: true}, Warning: "BSPlayer search is disabled because the upstream SOAP API is unreliable.", OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{"*.api.bsplayer-subtitles.com"}, AllowedDownloadHosts: []string{"*.bsplayer-subtitles.com", "bsplayer-subtitles.com"}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
	{Key: "napiprojekt", DisplayName: "Napiprojekt", ProvenanceCommit: bazarrProviderListCommit, RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: "Runtime support implemented for Polish hash-based subtitle lookup and download.", MediaTypes: []string{"movie", "serie"}, Dependencies: catalog.Dependencies{}, OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: []string{"napiprojekt.pl", "www.napiprojekt.pl"}, AllowedDownloadHosts: []string{"napiprojekt.pl", "www.napiprojekt.pl"}, AllowLocalHosts: false}, Fields: []catalog.Field{}},
}
