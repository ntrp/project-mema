package yifysubtitles

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "yifysubtitles",
	Name:       "YIFY Subtitles",
	BaseURL:    "https://yifysubtitles.ch",
	MediaTypes: []string{"movie"},
	Archive:    true,
})
