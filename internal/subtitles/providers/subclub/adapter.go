package subclub

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subclub",
	Name:       "Subclub",
	BaseURL:    "https://subclub.eu",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
