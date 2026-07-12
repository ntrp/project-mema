package nekur

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "nekur",
	Name:       "Nekur",
	BaseURL:    "http://subtitri.nekur.net",
	MediaTypes: []string{"movie"},
	Archive:    true,
})
