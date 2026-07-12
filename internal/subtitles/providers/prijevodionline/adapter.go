package prijevodionline

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "prijevodionline",
	Name:       "Prijevodionline",
	BaseURL:    "https://prijevodi-online.org",
	MediaTypes: []string{"serie"},
	Archive:    false,
})
