package hosszupuska

import (
	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/publichtml"
)

var Adapter providercore.Adapter = publichtml.New(publichtml.Spec{
	Key:        "hosszupuska",
	Name:       "Hosszupuska",
	BaseURL:    "https://hosszupuskasub.com",
	MediaTypes: []string{"serie"},
	Archive:    false,
})
