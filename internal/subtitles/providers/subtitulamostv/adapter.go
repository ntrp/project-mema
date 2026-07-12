package subtitulamostv

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subtitulamostv",
	Name:       "Subtitulamos.tv",
	BaseURL:    "https://subtitulamos.tv",
	MediaTypes: []string{"serie"},
	Archive:    false,
})
