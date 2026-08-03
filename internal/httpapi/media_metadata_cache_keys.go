package httpapi

import "strings"

func metadataDetailsCacheKey(providerType, externalID string) string {
	version := "v3"
	if providerType == "anidb" {
		version = "v7"
	}
	return "details:" + version + ":" + strings.ToLower(strings.TrimSpace(externalID))
}
