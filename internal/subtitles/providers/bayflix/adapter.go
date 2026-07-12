package bayflix

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "bayflix",
	Name:       "Bayflix",
	BaseURL:    "https://bayflix.sb",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
