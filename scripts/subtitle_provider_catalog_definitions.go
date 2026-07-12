//go:build subtitlecatalog

package main

const bazarrProviderListCommit = "e54edd769b7062280118a14aa0fef3808829714d"

var subtitleProviderCatalogDefinitions = concatSubtitleProviderDefinitions(
	subtitleProviderCatalogDefinitionsAPI,
	subtitleProviderCatalogDefinitionsPublic,
	subtitleProviderCatalogDefinitionsAuth,
	subtitleProviderCatalogDefinitionsSpecial,
	subtitleProviderCatalogDefinitionsLocal,
)

func concatSubtitleProviderDefinitions(groups ...[]subtitleProviderDefinition) []subtitleProviderDefinition {
	total := 0
	for _, group := range groups {
		total += len(group)
	}
	entries := make([]subtitleProviderDefinition, 0, total)
	for _, group := range groups {
		entries = append(entries, group...)
	}
	return entries
}
