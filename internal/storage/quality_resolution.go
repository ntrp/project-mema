package storage

import (
	"strings"
)

type QualityResolutionBounds struct {
	MinWidth  int32
	MinHeight int32
}

func QualityResolutionForID(qualityID string) (QualityResolutionBounds, bool) {
	definition, ok := QualitySizeDefinitionMap()[strings.TrimSpace(qualityID)]
	if !ok {
		return QualityResolutionBounds{}, false
	}
	return qualityResolutionForText(definition.ID + " " + definition.Name)
}

func QualityIDFromPath(path string) string {
	return mediaRenameQualityDefinition(path).ID
}

func qualityResolutionForText(value string) (QualityResolutionBounds, bool) {
	switch {
	case strings.Contains(value, "2160"):
		return QualityResolutionBounds{MinWidth: 3840, MinHeight: 2160}, true
	case strings.Contains(value, "1080"):
		return QualityResolutionBounds{MinWidth: 1920, MinHeight: 1080}, true
	case strings.Contains(value, "720"):
		return QualityResolutionBounds{MinWidth: 1280, MinHeight: 720}, true
	case strings.Contains(value, "576"):
		return QualityResolutionBounds{MinWidth: 1024, MinHeight: 576}, true
	case strings.Contains(value, "480"):
		return QualityResolutionBounds{MinWidth: 854, MinHeight: 480}, true
	default:
		return QualityResolutionBounds{}, false
	}
}
