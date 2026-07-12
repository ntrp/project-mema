package supersubtitles

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "supersubtitles",
	Name:       "Super Subtitles",
	BaseURL:    "https://feliratok.eu",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
