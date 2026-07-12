package subtitriid

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subtitriid",
	Name:       "Subtitriid",
	BaseURL:    "https://subtitri.do.am",
	MediaTypes: []string{"movie"},
	Archive:    true,
})
