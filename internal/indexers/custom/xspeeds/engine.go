package xspeeds

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "XSpeeds",
		DefaultBaseURL: "https://www.xspeeds.eu/",
		SearchPath:     "/browse.php",
		QueryParam:     "search",
		LoginPath:      "/takelogin.php",
	}, clients...)
}
