package greeksubs

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "greeksubs",
	Name:       "Greek Subs",
	BaseURL:    "https://greeksubs.net",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
