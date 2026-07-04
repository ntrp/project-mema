package bithdtv

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "BitHDTV",
		DefaultBaseURL: "https://www.bit-hdtv.com/",
		SearchPath:     "/torrents.php",
		QueryParam:     "search",
		CategoryParam:  "cat",
		LoginPath:      "/takelogin.php",
	}, clients...)
}
