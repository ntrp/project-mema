package subssabbz

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subssabbz",
	Name:       "Subs.sab.bz",
	BaseURL:    "http://subs.sab.bz",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
