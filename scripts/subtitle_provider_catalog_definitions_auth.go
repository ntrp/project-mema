//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsAuth = []subtitleProviderDefinition{
	clusterCProviderDefinition("addic7ed", "Addic7Ed", "https://www.addic7ed.com", []string{"www.addic7ed.com", "addic7ed.com"}, catalog.Dependencies{Captcha: true}, "Community subtitles for television episodes and movies, including current releases."),
	clusterCProviderDefinition("avistaz", "Avistaz", "https://avistaz.to", []string{"avistaz.to"}, catalog.Dependencies{}, "Private Asian media community covering films and television series."),
	clusterCProviderDefinition("cinemaz", "Cinemaz", "https://cinemaz.to", []string{"cinemaz.to"}, catalog.Dependencies{}, "Private movie community focused on international and specialized cinema."),
	clusterCProviderDefinition("hdbits", "HDBits", "https://hdbits.org", []string{"hdbits.org"}, catalog.Dependencies{}, "Private high-definition community for movie and television releases."),
	clusterCProviderDefinition("karagarga", "Karagarga", "https://karagarga.in", []string{"karagarga.in"}, catalog.Dependencies{}, "Private community specializing in arthouse, classic, and world cinema."),
	clusterCProviderDefinition("ktuvit", "Ktuvit", "https://www.ktuvit.me", []string{"www.ktuvit.me", "ktuvit.me"}, catalog.Dependencies{Captcha: true}, "Hebrew subtitle community for movies and television series."),
	clusterCProviderDefinition("legendasdivx", "Legendasdivx", "https://www.legendasdivx.pt", []string{"www.legendasdivx.pt", "legendasdivx.pt"}, catalog.Dependencies{}, "Portuguese subtitle community for movies and television series."),
	clusterCProviderDefinition("legendasnet", "Legendasnet", "https://legendas.net", []string{"legendas.net", "www.legendas.net"}, catalog.Dependencies{Captcha: true}, "Portuguese subtitle community and release index for films and series."),
	clusterCProviderDefinition("napisy24", "Napisy24", "https://napisy24.pl", []string{"napisy24.pl", "www.napisy24.pl"}, catalog.Dependencies{}, "Polish subtitle community for movies and television series."),
	clusterCProviderDefinition("pipocas", "Pipocas", "https://pipocas.tv", []string{"pipocas.tv", "www.pipocas.tv"}, catalog.Dependencies{}, "Portuguese subtitle community for movies and television series."),
	clusterCProviderDefinition("subs4series", "Subs4Series", "https://www.subs4series.com", []string{"www.subs4series.com", "subs4series.com"}, catalog.Dependencies{}, "Greek subtitle community specializing in television series."),
	clusterCProviderDefinition("subscenter", "Subscenter", "https://www.subscenter.org", []string{"www.subscenter.org", "subscenter.org"}, catalog.Dependencies{Captcha: true}, "Arabic subtitle community for movies and television series."),
	clusterCProviderDefinition("titlovi", "Titlovi", "https://titlovi.com", []string{"titlovi.com", "www.titlovi.com"}, catalog.Dependencies{}, "Balkan subtitle community for movies and television series."),
	clusterCProviderDefinition("titulky", "Titulky", "https://www.titulky.com", []string{"www.titulky.com", "titulky.com"}, catalog.Dependencies{}, "Czech and Slovak subtitle community for films and television series."),
	clusterCProviderDefinition("turkcealtyaziorg", "Turkcealtyaziorg", "https://turkcealtyazi.org", []string{"turkcealtyazi.org", "www.turkcealtyazi.org"}, catalog.Dependencies{}, "Turkish subtitle community for movies and television series."),
	clusterCProviderDefinition("xsubs", "X Subs", "https://xsubs.tv", []string{"xsubs.tv", "www.xsubs.tv"}, catalog.Dependencies{}, "Greek subtitle community for films and television series."),
	clusterCProviderDefinition("zimuku", "Zimuku", "https://zimuku.org", []string{"zimuku.org", "www.zimuku.org"}, catalog.Dependencies{Archive: true}, "Chinese subtitle library and community for movies and television shows."),
}

func clusterCProviderDefinition(key string, name string, baseURL string, hosts []string, deps catalog.Dependencies, description string) subtitleProviderDefinition {
	return subtitleProviderDefinition{
		Key: key, DisplayName: name, ProvenanceCommit: bazarrProviderListCommit,
		RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: description,
		MediaTypes: []string{"movie", "serie"}, Dependencies: deps,
		OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: hosts, AllowedDownloadHosts: hosts},
		Fields:         []catalog.Field{baseURLField(baseURL), cookiesField(), userAgentField()},
	}
}
