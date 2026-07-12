package vladoonmooo

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "vladoonmooo",
	Name:       "Vladoonmooo",
	BaseURL:    "https://vladoon.mooo.com",
	MediaTypes: []string{"movie", "serie"},
	Archive:    true,
})
