package subf2m

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subf2m",
	Name:       "Subf2M",
	BaseURL:    "https://subf2m.co",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
