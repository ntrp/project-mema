package subtitrarinoi

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subtitrarinoi",
	Name:       "Subtitrari Noi",
	BaseURL:    "https://subtitrari-noi.ro",
	MediaTypes: []string{"movie", "serie"},
	Archive:    true,
})
