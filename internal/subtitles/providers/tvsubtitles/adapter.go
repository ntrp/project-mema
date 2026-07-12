package tvsubtitles

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "tvsubtitles",
	Name:       "TVsubtitles",
	BaseURL:    "https://tvsubtitles.net",
	MediaTypes: []string{"serie"},
	Archive:    false,
})
