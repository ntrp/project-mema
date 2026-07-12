package subsunacs

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subsunacs",
	Name:       "Subsunacs",
	BaseURL:    "https://subsunacs.net",
	MediaTypes: []string{"movie", "serie"},
	Archive:    true,
})
