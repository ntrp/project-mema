package metadata

import (
	"regexp"
	"strconv"
	"strings"
)

const anidbTMDBResourceType = "44"

var anidbDescriptionLink = regexp.MustCompile(`https?://(?:www\.)?anidb\.net/\S+\s+\[([^\]]+)\]`)

type anidbResource struct {
	Type     string                `xml:"type,attr"`
	Entities []anidbExternalEntity `xml:"externalentity"`
}

type anidbExternalEntity struct {
	Identifiers []string `xml:"identifier"`
}

func anidbOverview(value string) string {
	return strings.TrimSpace(anidbDescriptionLink.ReplaceAllString(value, "$1"))
}

func anidbProviderMappings(resources []anidbResource) []ProviderMapping {
	for _, resource := range resources {
		if resource.Type != anidbTMDBResourceType {
			continue
		}
		for _, entity := range resource.Entities {
			for _, value := range entity.Identifiers {
				id := strings.TrimSpace(value)
				parsed, err := strconv.Atoi(id)
				if err != nil || parsed <= 0 {
					continue
				}
				return []ProviderMapping{{
					ProviderName: "tmdb", ProviderEntityType: "serie", ExternalID: id,
					Source: map[string]any{"source": "anidb_resource", "resourceType": anidbTMDBResourceType},
				}}
			}
		}
	}
	return nil
}
