package subsynchro

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subsynchro",
	Name:       "Subsynchro",
	BaseURL:    "https://subsynchro.com",
	MediaTypes: []string{"movie"},
	Archive:    false,
})
