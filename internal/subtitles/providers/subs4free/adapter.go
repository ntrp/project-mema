package subs4free

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "subs4free",
	Name:       "Subs4Free",
	BaseURL:    "https://subs4free.info",
	MediaTypes: []string{"movie"},
	Archive:    true,
})
