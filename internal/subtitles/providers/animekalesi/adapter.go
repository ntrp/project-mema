package animekalesi

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "animekalesi",
	Name:       "Animekalesi",
	BaseURL:    "https://animekalesi.com",
	MediaTypes: []string{"serie"},
	Archive:    true,
})
