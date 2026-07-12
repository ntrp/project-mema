//go:build subtitlecatalog

package main

import "media-manager/internal/subtitles/catalog"

var subtitleProviderCatalogDefinitionsAuth = []subtitleProviderDefinition{
	clusterCProviderDefinition("addic7ed", "Addic7Ed", "https://www.addic7ed.com", []string{"www.addic7ed.com", "addic7ed.com"}, catalog.Dependencies{Captcha: true}, "Requires a logged-in browser cookie because interactive CAPTCHA cannot be solved by the runtime."),
	clusterCProviderDefinition("avistaz", "Avistaz", "https://avistaz.to", []string{"avistaz.to"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("cinemaz", "Cinemaz", "https://cinemaz.to", []string{"cinemaz.to"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("hdbits", "HDBits", "https://hdbits.org", []string{"hdbits.org"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("karagarga", "Karagarga", "https://karagarga.in", []string{"karagarga.in"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("ktuvit", "Ktuvit", "https://www.ktuvit.me", []string{"www.ktuvit.me", "ktuvit.me"}, catalog.Dependencies{Captcha: true}, captchaCookieRuntimeMessage),
	clusterCProviderDefinition("legendasdivx", "Legendasdivx", "https://www.legendasdivx.pt", []string{"www.legendasdivx.pt", "legendasdivx.pt"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("legendasnet", "Legendasnet", "https://legendas.net", []string{"legendas.net", "www.legendas.net"}, catalog.Dependencies{Captcha: true}, captchaCookieRuntimeMessage),
	clusterCProviderDefinition("napisy24", "Napisy24", "https://napisy24.pl", []string{"napisy24.pl", "www.napisy24.pl"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("pipocas", "Pipocas", "https://pipocas.tv", []string{"pipocas.tv", "www.pipocas.tv"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("subs4series", "Subs4Series", "https://www.subs4series.com", []string{"www.subs4series.com", "subs4series.com"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("subscenter", "Subscenter", "https://www.subscenter.org", []string{"www.subscenter.org", "subscenter.org"}, catalog.Dependencies{Captcha: true}, captchaCookieRuntimeMessage),
	clusterCProviderDefinition("titlovi", "Titlovi", "https://titlovi.com", []string{"titlovi.com", "www.titlovi.com"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("titulky", "Titulky", "https://www.titulky.com", []string{"www.titulky.com", "titulky.com"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("turkcealtyaziorg", "Turkcealtyaziorg", "https://turkcealtyazi.org", []string{"turkcealtyazi.org", "www.turkcealtyazi.org"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("xsubs", "X Subs", "https://xsubs.tv", []string{"xsubs.tv", "www.xsubs.tv"}, catalog.Dependencies{}, privateTrackerRuntimeMessage),
	clusterCProviderDefinition("zimuku", "Zimuku", "https://zimuku.org", []string{"zimuku.org", "www.zimuku.org"}, catalog.Dependencies{Archive: true}, privateTrackerRuntimeMessage),
}

const (
	privateTrackerRuntimeMessage = "Runtime supported for configured members with valid session cookies; unauthenticated use returns a typed private-membership error."
	captchaCookieRuntimeMessage = "Runtime supported only with valid session cookies from a user-solved CAPTCHA session; missing cookies return a typed CAPTCHA/private prerequisite error."
)

func clusterCProviderDefinition(key string, name string, baseURL string, hosts []string, deps catalog.Dependencies, message string) subtitleProviderDefinition {
	return subtitleProviderDefinition{
		Key: key, DisplayName: name, ProvenanceCommit: bazarrProviderListCommit,
		RuntimeStatus: catalog.RuntimeSupported, RuntimeMessage: message,
		MediaTypes: []string{"movie", "serie"}, Dependencies: deps,
		OutboundPolicy: catalog.OutboundPolicy{AllowedBaseHosts: hosts, AllowedDownloadHosts: hosts},
		Fields: []catalog.Field{baseURLField(baseURL), cookiesField(), userAgentField()},
	}
}
