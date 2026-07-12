package animesubinfo

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "animesubinfo",
	Name:       "Animesubinfo",
	BaseURL:    "https://animesub.info",
	MediaTypes: []string{"movie", "serie"},
	Archive:    true,
})
