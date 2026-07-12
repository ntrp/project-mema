package soustitreseu

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "soustitreseu",
	Name:       "Sous-titres.eu",
	BaseURL:    "https://sous-titres.eu",
	MediaTypes: []string{"movie", "serie"},
	Archive:    true,
})
