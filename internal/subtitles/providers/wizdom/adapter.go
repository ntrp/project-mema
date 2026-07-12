package wizdom

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "wizdom",
	Name:       "Wizdom",
	BaseURL:    "https://wizdom.xyz",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
