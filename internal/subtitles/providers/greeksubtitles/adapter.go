package greeksubtitles

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "greeksubtitles",
	Name:       "Greek Subtitles",
	BaseURL:    "https://gr.greek-subtitles.com",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
