package titrari

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "titrari",
	Name:       "Titrari",
	BaseURL:    "https://titrari.ro",
	MediaTypes: []string{"movie", "serie"},
	Archive:    true,
})
