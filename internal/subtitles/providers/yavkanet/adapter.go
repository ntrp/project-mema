package yavkanet

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "yavkanet",
	Name:       "Yavka.net",
	BaseURL:    "https://yavka.net",
	MediaTypes: []string{"movie", "serie"},
	Archive:    false,
})
