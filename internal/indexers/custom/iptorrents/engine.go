package iptorrents

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "IPTorrents",
		DefaultBaseURL: "https://iptorrents.com/",
		SearchPath:     "/t",
		QueryParam:     "q",
		CategoryParam:  "cat",
	}, clients...)
}
